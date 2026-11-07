package domain

import (
	"errors"
	"strings"
	"testing"
)

func validTestConfig() *SafeProjectConfig {
	return &SafeProjectConfig{
		ProjectName:        "valid-app",
		BinaryName:         "valid-app",
		MainPath:           ".",
		ProjectType:        ProjectTypeCLI,
		GitProvider:        GitProviderGitHub,
		CGOStatus:          CGOStatusDisabled,
		DockerSupport:      DockerSupportNone,
		DockerRegistry:     DockerRegistryGitHub,
		SigningLevel:       SigningLevelNone,
		ActionLevel:        ActionLevelBasic,
		FeatureLevel:       FeatureLevelStandard,
		State:              ConfigStateValid,
		ActionsOn:          []ActionTrigger{ActionTriggerVersionTags},
		Platforms:          []Platform{PlatformLinux},
		Architectures:      []Architecture{ArchitectureAMD64},
		ProjectDescription: "A valid test application",
	}
}

func TestSafeProjectConfigValidateValid(t *testing.T) {
	if err := validTestConfig().Validate(); err != nil {
		t.Errorf("valid config failed Validate(): %v", err)
	}
}

func TestSafeProjectConfigValidateErrorPaths(t *testing.T) {
	tests := []struct {
		name       string
		mutate     func(*SafeProjectConfig)
		wantCode   ErrorCode
		wantField  string
		wantStatus bool
	}{
		{
			name:      "empty_project_name",
			mutate:    func(c *SafeProjectConfig) { c.ProjectName = "" },
			wantField: "project_name",
		},
		{
			name:      "reserved_project_name",
			mutate:    func(c *SafeProjectConfig) { c.ProjectName = "test" },
			wantField: "project_name",
		},
		{
			name:      "empty_binary_name",
			mutate:    func(c *SafeProjectConfig) { c.BinaryName = "" },
			wantField: "binary_name",
		},
		{
			name:      "empty_main_path",
			mutate:    func(c *SafeProjectConfig) { c.MainPath = "" },
			wantField: "main_path",
		},
		{
			name:       "path_traversal_in_main_path",
			mutate:     func(c *SafeProjectConfig) { c.MainPath = "../escape" },
			wantStatus: true,
		},
		{
			name:      "invalid_project_type",
			mutate:    func(c *SafeProjectConfig) { c.ProjectType = ProjectType("bogus") },
			wantField: "project_type",
		},
		{
			name:      "invalid_git_provider",
			mutate:    func(c *SafeProjectConfig) { c.GitProvider = GitProvider("sourcehut") },
			wantField: "git_provider",
		},
		{
			name:      "invalid_cgo_status",
			mutate:    func(c *SafeProjectConfig) { c.CGOStatus = CGOStatus("maybe") },
			wantField: "cgo_status",
		},
		{
			name:      "invalid_docker_support",
			mutate:    func(c *SafeProjectConfig) { c.DockerSupport = DockerSupport("sometimes") },
			wantField: "docker_support",
		},
		{
			name:      "docker_enabled_without_registry",
			mutate:    func(c *SafeProjectConfig) { c.DockerSupport = DockerSupportBuild; c.DockerRegistry = "" },
			wantField: "docker_registry",
		},
		{
			name:      "invalid_signing_level",
			mutate:    func(c *SafeProjectConfig) { c.SigningLevel = SigningLevel("paranoid") },
			wantField: "signing_level",
		},
		{
			name:      "invalid_action_level",
			mutate:    func(c *SafeProjectConfig) { c.ActionLevel = ActionLevel("everything") },
			wantField: "action_level",
		},
		{
			name:      "invalid_ci_provider",
			mutate:    func(c *SafeProjectConfig) { c.CIProvider = CIProvider("sourcehut") },
			wantField: "ci_provider",
		},
		{
			name:      "invalid_feature_level",
			mutate:    func(c *SafeProjectConfig) { c.FeatureLevel = FeatureLevel("deluxe") },
			wantField: "feature_level",
		},
		{
			name:      "invalid_state",
			mutate:    func(c *SafeProjectConfig) { c.State = ConfigState("vanished") },
			wantField: "state",
		},
		{
			name:      "invalid_project_description",
			mutate:    func(c *SafeProjectConfig) { c.ProjectDescription = "has\nnewline" },
			wantField: "project_description",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := validTestConfig()
			tt.mutate(config)

			err := config.Validate()
			if err == nil {
				t.Fatalf("Validate() = nil, want error for %s", tt.name)
			}

			if tt.wantStatus {
				return
			}

			domainErr, ok := errors.AsType[*DomainError](err)
			if !ok {
				t.Fatalf("error is %T, want *DomainError", err)
			}

			if domainErr.Field != tt.wantField {
				t.Errorf("error field = %q, want %q (error: %v)", domainErr.Field, tt.wantField, err)
			}
		})
	}
}

func TestSafeProjectConfigCloneDeepIndependence(t *testing.T) {
	original := validTestConfig()

	clone := original.Clone()

	if clone == original {
		t.Fatal("Clone() returned the same pointer")
	}

	if clone.ProjectName != original.ProjectName {
		t.Errorf("clone ProjectName = %q, want %q", clone.ProjectName, original.ProjectName)
	}

	clone.ProjectName = "mutated"
	clone.Platforms[0] = PlatformWindows

	if original.ProjectName != "valid-app" {
		t.Errorf("mutating clone changed original ProjectName to %q", original.ProjectName)
	}

	if original.Platforms[0] != PlatformLinux {
		t.Errorf("mutating clone slice changed original Platforms[0] to %q (shared backing array)", original.Platforms[0])
	}
}

func TestSafeProjectConfigApplyDefaults(t *testing.T) {
	config := &SafeProjectConfig{
		ProjectName: "defaults-app",
		BinaryName:  "defaults-app",
		MainPath:    ".",
		ProjectType: ProjectTypeCLI,
	}

	config.ApplyDefaults()

	if config.FeatureLevel != FeatureLevelStandard {
		t.Errorf("FeatureLevel = %q, want standard default", config.FeatureLevel)
	}

	if config.CGOStatus == "" {
		t.Error("CGOStatus left empty by ApplyDefaults for CLI project")
	}

	if len(config.Platforms) == 0 {
		t.Error("Platforms left empty by ApplyDefaults")
	}

	if len(config.Architectures) == 0 {
		t.Error("Architectures left empty by ApplyDefaults")
	}

	if config.DockerSupport == "" {
		t.Error("DockerSupport left empty by ApplyDefaults")
	}
}

func TestSafeProjectConfigJSONRoundTrip(t *testing.T) {
	original := validTestConfig()

	data, err := original.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON() = %v", err)
	}

	restored := &SafeProjectConfig{}

	if err := restored.FromJSON(data); err != nil {
		t.Fatalf("FromJSON() = %v", err)
	}

	if restored.ProjectName != original.ProjectName {
		t.Errorf("round-trip ProjectName = %q, want %q", restored.ProjectName, original.ProjectName)
	}

	if restored.ProjectType != original.ProjectType {
		t.Errorf("round-trip ProjectType = %q, want %q", restored.ProjectType, original.ProjectType)
	}

	if len(restored.Architectures) != len(original.Architectures) {
		t.Errorf("round-trip Architectures = %v, want %v", restored.Architectures, original.Architectures)
	}
}

func TestSafeProjectConfigYAMLRoundTrip(t *testing.T) {
	original := validTestConfig()

	data, err := original.ToYAML()
	if err != nil {
		t.Fatalf("ToYAML() = %v", err)
	}

	restored := &SafeProjectConfig{}

	if err := restored.FromYAML(data); err != nil {
		t.Fatalf("FromYAML() = %v", err)
	}

	if restored.ProjectName != original.ProjectName {
		t.Errorf("round-trip ProjectName = %q, want %q", restored.ProjectName, original.ProjectName)
	}
}

func TestEnumFromBoolTables(t *testing.T) {
	tests := []struct {
		name    string
		got     any
		wantStr string
	}{
		{"action_level_true", ActionLevelFromBool(true), string(ActionLevelBasic)},
		{"action_level_false", ActionLevelFromBool(false), string(ActionLevelNone)},
		{"cgo_status_true", CGOStatusFromBool(true), string(CGOStatusEnabled)},
		{"cgo_status_false", CGOStatusFromBool(false), string(CGOStatusDisabled)},
		{"docker_support_true", DockerSupportFromBool(true), string(DockerSupportBuild)},
		{"docker_support_false", DockerSupportFromBool(false), string(DockerSupportNone)},
		{"feature_level_true", FeatureLevelFromBool(true), string(FeatureLevelStandard)},
		{"feature_level_false", FeatureLevelFromBool(false), string(FeatureLevelBasic)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, ok := tt.got.(string); !ok || got != tt.wantStr {
				t.Errorf("%s = %v, want %q", tt.name, tt.got, tt.wantStr)
			}
		})
	}
}

func TestSecurityValidatorsRejectHostilePaths(t *testing.T) {
	tests := []struct {
		name     string
		mainPath string
	}{
		{"path_traversal", "../escape"},
		{"absolute_path_traversal", "cmd/../../etc"},
		{"shell_metacharacters", "cmd; rm -rf /"},
		{"command_substitution", "$(dangerous)"},
		{"backtick_substitution", "`dangerous`"},
		{"url_injection", "https://evil.example/payload"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateMainPath(tt.mainPath); err == nil {
				t.Errorf("ValidateMainPath(%q) = nil, want hostile path rejected", tt.mainPath)
			}
		})
	}
}

func TestValidateMainPathAcceptsLegitimatePaths(t *testing.T) {
	for _, path := range []string{".", "./cmd/app", "cmd/app", "internal/tool"} {
		t.Run("path_"+path, func(t *testing.T) {
			if err := ValidateMainPath(path); err != nil {
				t.Errorf("ValidateMainPath(%q) = %v, want nil", path, err)
			}
		})
	}
}

func TestValidateProjectNameReservedNames(t *testing.T) {
	for _, name := range []string{"test", "example"} {
		t.Run("name_"+name, func(t *testing.T) {
			if err := ValidateProjectName(name); err == nil {
				t.Errorf("ValidateProjectName(%q) = nil, want reserved name rejected", name)
			}
		})
	}
}

func TestDefaultCGOEnabledByProjectType(t *testing.T) {
	tests := []struct {
		projectType ProjectType
		want        bool
	}{
		{ProjectTypeCLI, false},
		{ProjectTypeLibrary, false},
		{ProjectTypeMobile, false},
		{ProjectTypeDesktop, true},
		{ProjectTypeWebAPI, true},
		{ProjectTypeGRPCService, true},
	}

	for _, tt := range tests {
		t.Run("type_"+string(tt.projectType), func(t *testing.T) {
			if got := tt.projectType.DefaultCGOEnabled(); got != tt.want {
				t.Errorf("DefaultCGOEnabled(%s) = %v, want %v", tt.projectType, got, tt.want)
			}
		})
	}
}

func TestErrorConstructersCarryContext(t *testing.T) {
	dockerErr := DockerNotSupportedError(ProjectTypeCLI)
	if dockerErr == nil {
		t.Fatal("DockerNotSupportedError = nil")
	}

	fileErr := FileNotFoundError("go.mod", nil)
	if fileErr == nil {
		t.Fatal("FileNotFoundError = nil")
	}

	if !strings.Contains(fileErr.Details, "go.mod") {
		t.Errorf("FileNotFoundError details %q do not name the file", fileErr.Details)
	}

	writeErr := FileWriteFailedError("/tmp/x", nil)
	if writeErr == nil {
		t.Fatal("FileWriteFailedError = nil")
	}
}

func contains(haystack, needle string) bool {
	return strings.Contains(haystack, needle)
}
