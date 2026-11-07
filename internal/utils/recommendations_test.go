package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestGetRecommendedPlatforms(t *testing.T) {
	platforms := GetRecommendedPlatforms()
	assert.Contains(t, platforms, domain.PlatformLinux, "should always include Linux")

	switch runtime.GOOS {
	case "darwin":
		assert.Contains(t, platforms, domain.PlatformDarwin)
	case "windows":
		assert.Contains(t, platforms, domain.PlatformWindows)
	}
}

func TestGetRecommendedArchitectures(t *testing.T) {
	archs := GetRecommendedArchitectures()
	assert.Contains(t, archs, domain.ArchitectureAMD64, "should always include AMD64")

	switch runtime.GOARCH {
	case "arm64":
		assert.Contains(t, archs, domain.ArchitectureARM64)
	case "386":
		assert.Contains(t, archs, domain.Architecture386)
	}
}

func TestGetRecommendedDockerRegistry(t *testing.T) {
	assert.Equal(t, domain.DockerRegistryGitHub, GetRecommendedDockerRegistry(domain.GitProviderGitHub))
	assert.Equal(t, domain.DockerRegistryGitLab, GetRecommendedDockerRegistry(domain.GitProviderGitLab))
	assert.Equal(t, domain.DockerRegistryCustom, GetRecommendedDockerRegistry(domain.GitProviderGitea))
}

func TestIsDevelopmentEnvironment(t *testing.T) {
	t.Cleanup(func() {
		os.Unsetenv("GO_ENV")
		os.Unsetenv("ENV")
		os.Unsetenv("NODE_ENV")
	})

	os.Setenv("GO_ENV", "development")
	assert.True(t, IsDevelopmentEnvironment())

	os.Setenv("GO_ENV", "production")
	assert.False(t, IsDevelopmentEnvironment())
}

func TestIsProductionEnvironment(t *testing.T) {
	t.Cleanup(func() {
		os.Unsetenv("GO_ENV")
		os.Unsetenv("ENV")
		os.Unsetenv("NODE_ENV")
	})

	os.Setenv("GO_ENV", "production")
	assert.True(t, IsProductionEnvironment())

	os.Setenv("GO_ENV", "development")
	assert.False(t, IsProductionEnvironment())
}

func TestGetEnvironment(t *testing.T) {
	t.Cleanup(func() {
		os.Unsetenv("GO_ENV")
		os.Unsetenv("ENV")
		os.Unsetenv("NODE_ENV")
	})

	os.Setenv("GO_ENV", "staging")
	assert.Equal(t, "staging", GetEnvironment())

	os.Unsetenv("GO_ENV")
	os.Setenv("ENV", "test")
	assert.Equal(t, "test", GetEnvironment())

	os.Unsetenv("ENV")
	os.Setenv("NODE_ENV", "ci")
	assert.Equal(t, "ci", GetEnvironment())

	os.Unsetenv("NODE_ENV")
	assert.Equal(t, "development", GetEnvironment())
}

func TestHasFile(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	tests := []struct {
		name  string
		setup func()
		file  string
		want  bool
	}{
		{
			name: "root_hit",
			setup: func() {
				os.WriteFile("go.mod", []byte("module t\n"), 0o644)
			},
			file: "go.mod",
			want: true,
		},
		{
			name:  "miss",
			setup: func() {},
			file:  "does-not-exist.txt",
			want:  false,
		},
		{
			name: "subdir_hit",
			setup: func() {
				os.MkdirAll("cmd", 0o755)
				os.WriteFile("cmd/main.go", []byte("package main\n"), 0o644)
			},
			file: "main.go",
			want: true,
		},
		{
			name:  "dir_slash_form_hit",
			setup: func() {},
			file:  "api/",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			if got := hasFile(tt.file); got != tt.want {
				t.Errorf("hasFile(%q) = %v, want %v", tt.file, got, tt.want)
			}
		})
	}
}

func TestGetRecommendedProjectType(t *testing.T) {
	tests := []struct {
		name  string
		setup func(dir string)
		want  domain.ProjectType
	}{
		{
			name:  "no_go_mod_defaults_to_cli",
			setup: func(dir string) {},
			want:  domain.ProjectTypeCLI,
		},
		{
			name: "go_mod_with_root_main_is_cli",
			setup: func(dir string) {
				os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module t\n"), 0o644)
				os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\n"), 0o644)
			},
			want: domain.ProjectTypeCLI,
		},
		{
			name: "go_mod_with_cmd_main_is_cli",
			setup: func(dir string) {
				os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module t\n"), 0o644)
				os.MkdirAll(filepath.Join(dir, "cmd"), 0o755)
				os.WriteFile(filepath.Join(dir, "cmd", "main.go"), []byte("package main\n"), 0o644)
			},
			want: domain.ProjectTypeCLI,
		},
		{
			name: "go_mod_with_api_dir_is_webapi",
			setup: func(dir string) {
				os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module t\n"), 0o644)
				os.MkdirAll(filepath.Join(dir, "api"), 0o755)
			},
			want: domain.ProjectTypeWebAPI,
		},
		{
			name: "go_mod_with_grpc_dir_is_grpc",
			setup: func(dir string) {
				os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module t\n"), 0o644)
				os.MkdirAll(filepath.Join(dir, "grpc"), 0o755)
			},
			want: domain.ProjectTypeGRPCService,
		},
		{
			name: "go_mod_alone_is_library",
			setup: func(dir string) {
				os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module t\n"), 0o644)
			},
			want: domain.ProjectTypeLibrary,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			t.Chdir(dir)

			tt.setup(dir)

			if got := GetRecommendedProjectType(); got != tt.want {
				t.Errorf("GetRecommendedProjectType() = %v, want %v", got, tt.want)
			}
		})
	}
}
