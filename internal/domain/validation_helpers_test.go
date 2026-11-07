package domain

import (
	"errors"
	"testing"
)

func TestContainsSecurityHelpers(t *testing.T) {
	tests := []struct {
		name string
		got  bool
		want bool
	}{
		{"traversal_parent", containsPathTraversal("../escape"), true},
		{"traversal_backslash", containsPathTraversal(`cmd\evil`), true},
		{"traversal_clean", containsPathTraversal("cmd/app"), false},
		{"shell_pipe", containsShellMetacharacters("a|b"), true},
		{"shell_substitution", containsShellMetacharacters("$(x)"), true},
		{"shell_clean", containsShellMetacharacters("cmd/app"), false},
		{"url_javascript", containsURLInjection("javascript:x"), true},
		{"url_data", containsURLInjection("data:y"), true},
		{"url_clean", containsURLInjection("https://ok.example"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("%s = %v, want %v", tt.name, tt.got, tt.want)
			}
		})
	}
}

func TestValidateEnumHelpers(t *testing.T) {
	if err := validateEnum("project type", "cli", true); err != nil {
		t.Errorf("validateEnum(valid) = %v, want nil", err)
	}

	err := validateEnum("project type", "bogus", false)
	if err == nil {
		t.Fatal("validateEnum(invalid) = nil, want error")
	}

	if !contains(err.Error(), "Invalid project type") {
		t.Errorf("validateEnum error %q does not name the item", err.Error())
	}
}

func TestValidateEnumSlice(t *testing.T) {
	validPlatforms := []Platform{PlatformLinux, PlatformDarwin}

	if err := ValidatePlatforms(validPlatforms); err != nil {
		t.Errorf("ValidatePlatforms(valid) = %v, want nil", err)
	}

	if err := ValidatePlatforms([]Platform{Platform("holodeck")}); err == nil {
		t.Error("ValidatePlatforms(invalid) = nil, want error")
	}
}

func TestValidatePlatformArchCompatibility(t *testing.T) {
	valid := []Platform{PlatformLinux}

	archs := []Architecture{ArchitectureAMD64}

	if err := ValidatePlatformArchCompatibility(valid, archs); err != nil {
		t.Errorf("compatible combo rejected: %v", err)
	}

	err := ValidatePlatformArchCompatibility([]Platform{PlatformDarwin}, []Architecture{Architecture386})
	if err == nil {
		t.Error("darwin/386 (incompatible combo) accepted")
	}
}

func TestErrorConstructorsUnwrapCause(t *testing.T) {
	cause := errors.New("disk full")

	err := FileWriteFailedError("/tmp/x", cause)

	domainErr, ok := errors.AsType[*DomainError](err)
	if !ok {
		t.Fatalf("FileWriteFailedError returned %T, want *DomainError", err)
	}

	if !errors.Is(err, cause) {
		t.Error("FileWriteFailedError does not unwrap to its cause")
	}

	if domainErr.Code != ErrFileWriteFailed {
		t.Errorf("code = %q, want %q", domainErr.Code, ErrFileWriteFailed)
	}
}

func TestWithSuggestionAndFieldChain(t *testing.T) {
	err := NewValidationError(ErrInvalidProjectName, "bad name", "123!!").WithField("project_name")

	domainErr, ok := errors.AsType[*DomainError](err)
	if !ok {
		t.Fatalf("got %T, want *DomainError", err)
	}

	suggestion := domainErr.WithSuggestion("use lowercase-hyphen names")

	if suggestion.Suggestion == "" {
		t.Error("WithSuggestion did not attach a suggestion")
	}

	if suggestion.Field != "project_name" {
		t.Errorf("Field = %q, want project_name", suggestion.Field)
	}
}

func TestDomainErrorFormatting(t *testing.T) {
	withContext := NewValidationError(ErrInvalidConfig, "broken", "detail").WithContext("generate")
	if got := withContext.Error(); got == "" || !contains(got, string(ErrInvalidConfig)) {
		t.Errorf("Error() = %q, want it to carry the code", got)
	}

	noContext := NewValidationError(ErrInvalidConfig, "broken", "detail")
	if got := noContext.Error(); !contains(got, string(ErrInvalidConfig)) {
		t.Errorf("Error() = %q, want it to carry the code", got)
	}
}

func TestGetAllPlatformsAndRecommendedType(t *testing.T) {
	platforms := GetAllPlatforms()
	if len(platforms) == 0 {
		t.Error("GetAllPlatforms() = empty")
	}

	// The recommendation path shells out to git; inside the checkout it must
	// not panic and must return a valid type.
	recommended := GetRecommendedProjectType()
	if !recommended.IsValid() {
		t.Errorf("GetRecommendedProjectType() = %q, want a valid project type", recommended)
	}
}

func TestDockerRegistryGetURL(t *testing.T) {
	tests := []struct {
		registry DockerRegistry
		nonEmpty bool
	}{
		{DockerRegistryGitHub, true},
		{DockerRegistryDockerHub, true},
		{DockerRegistry("bogus"), false},
	}

	for _, tt := range tests {
		t.Run("registry_"+string(tt.registry), func(t *testing.T) {
			url := tt.registry.GetURL()
			if tt.nonEmpty && url == "" {
				t.Errorf("GetURL(%q) = empty, want a URL", tt.registry)
			}
		})
	}
}

func TestGitProviderSupportsDocker(t *testing.T) {
	if !GitProviderGitHub.SupportsDocker() {
		t.Error("github must support docker")
	}
}
