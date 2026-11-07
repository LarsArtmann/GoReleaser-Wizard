package domain

import (
	"testing"
)

func TestGitProviderProperties(t *testing.T) {
	tests := []struct {
		provider      GitProvider
		registry      DockerRegistry
		actions       bool
		pat           bool
		validProvider bool
	}{
		{GitProviderGitHub, DockerRegistryGitHub, true, false, true},
		{GitProviderGitLab, DockerRegistryGitLab, true, true, true},
		{GitProvider("forgejo"), DockerRegistry("ghcr.io"), false, true, false},
	}

	for _, tt := range tests {
		t.Run("provider_"+string(tt.provider), func(t *testing.T) {
			if got := tt.provider.IsValid(); got != tt.validProvider {
				t.Errorf("IsValid(%s) = %v, want %v", tt.provider, got, tt.validProvider)
			}

			if !tt.validProvider {
				return
			}

			if got := tt.provider.DefaultRegistry(); got != tt.registry {
				t.Errorf("DefaultRegistry(%s) = %q, want %q", tt.provider, got, tt.registry)
			}

			if got := tt.provider.ActionsSupported(); got != tt.actions {
				t.Errorf("ActionsSupported(%s) = %v, want %v", tt.provider, got, tt.actions)
			}

			if got := tt.provider.RequiresPersonalAccessToken(); got != tt.pat {
				t.Errorf("RequiresPersonalAccessToken(%s) = %v, want %v", tt.provider, got, tt.pat)
			}

			if tt.provider.APIURL() == "" {
				t.Errorf("APIURL(%s) is empty", tt.provider)
			}

			if tt.provider.WebURL() == "" {
				t.Errorf("WebURL(%s) is empty", tt.provider)
			}
		})
	}
}

func TestGetAllGitProvidersIncludeGitHub(t *testing.T) {
	providers := GetAllGitProviders()

	found := false

	for _, provider := range providers {
		if provider == GitProviderGitHub {
			found = true
		}
	}

	if !found {
		t.Error("GetAllGitProviders() does not include github")
	}
}

func TestConvertToDockerRegistry(t *testing.T) {
	tests := []struct {
		displayName string
		want        DockerRegistry
	}{
		{"github", DockerRegistryGitHub},
		{"ghcr.io", DockerRegistryGitHub},
		{"docker hub", DockerRegistryDockerHub},
		{"quay", DockerRegistryQuay},
		// Unknown names fall back to Docker Hub, not an invalid registry.
		{"unknown-display-name", DockerRegistryDockerHub},
	}

	for _, tt := range tests {
		t.Run("name_"+tt.displayName, func(t *testing.T) {
			if got := ConvertToDockerRegistry(tt.displayName); got != tt.want {
				t.Errorf("ConvertToDockerRegistry(%q) = %q, want %q", tt.displayName, got, tt.want)
			}
		})
	}
}

func TestDockerRegistryProperties(t *testing.T) {
	gh := DockerRegistryGitHub

	if !gh.IsValid() {
		t.Error("github registry must be valid")
	}

	if gh.URLPattern() == "" {
		t.Error("github registry URLPattern is empty")
	}

	if !gh.RequiresAuthentication() {
		t.Error("ghcr.io requires authentication for pushes")
	}

	if DockerRegistryDockerHub.DefaultNamespace() == "" {
		t.Error("docker hub DefaultNamespace is empty")
	}
}

func TestValidateDockerImageName(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"mytool", false},
		{"octo/mytool", false},
		// Empty is allowed: the generator defaults it to the project name.
		{"", false},
		{"UPPER/INVALID", true},
	}

	for _, tt := range tests {
		t.Run("image_"+tt.name, func(t *testing.T) {
			err := ValidateDockerImageName(tt.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDockerImageName(%q) = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}

func TestArchitectureCompatibility(t *testing.T) {
	tests := []struct {
		arch       Architecture
		platform   Platform
		compatible bool
	}{
		{ArchitectureAMD64, PlatformLinux, true},
		{ArchitectureARM64, PlatformLinux, true},
		{ArchitectureAMD64, PlatformDarwin, true},
		{Architecture386, PlatformDarwin, false},
		{Architecture386, PlatformWindows, true},
	}

	for _, tt := range tests {
		t.Run(string(tt.arch)+"_"+string(tt.platform), func(t *testing.T) {
			if got := tt.arch.IsCompatibleWith(tt.platform); got != tt.compatible {
				t.Errorf("IsCompatibleWith(%s, %s) = %v, want %v", tt.arch, tt.platform, got, tt.compatible)
			}
		})
	}
}

func TestBuildTagValidity(t *testing.T) {
	tests := []struct {
		name string
		tag  BuildTag
		want bool
	}{
		{"valid_name", BuildTag{Name: "json"}, true},
		{"empty_name", BuildTag{}, false},
		{"invalid_chars", BuildTag{Name: "json mode"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tag.IsValid(); got != tt.want {
				t.Errorf("IsValid(%q) = %v, want %v", tt.tag.Name, got, tt.want)
			}
		})
	}
}

func TestGetSummary(t *testing.T) {
	config := validTestConfig()
	config.DockerSupport = DockerSupportBuild
	config.SigningLevel = SigningLevelAdvanced

	summary := config.GetSummary()

	if summary.ProjectName != config.ProjectName {
		t.Errorf("summary ProjectName = %q, want %q", summary.ProjectName, config.ProjectName)
	}

	if !summary.DockerEnabled {
		t.Error("summary DockerEnabled = false for build support")
	}

	if !summary.SigningEnabled {
		t.Error("summary SigningEnabled = false for keyless signing")
	}

	if summary.ArchCount != len(config.Architectures) {
		t.Errorf("summary ArchCount = %d, want %d", summary.ArchCount, len(config.Architectures))
	}
}
