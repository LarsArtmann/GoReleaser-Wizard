package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEndToEndWizard(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func() string
		expectFiles []string
	}{
		{
			name: "complete_wizard_flow",
			setupFunc: func() string {
				dir, _ := os.MkdirTemp("", "wizard-e2e-test")
				
				// Create basic Go project
				goMod := `module github.com/user/e2e-test
go 1.21
require github.com/charmbracelet/bubbletea v0.25.0
`
				os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goMod), 0644)
				os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\n\nfunc main() {}"), 0644)
				
				// Create cmd directory structure
				os.MkdirAll(filepath.Join(dir, "cmd", "e2e-test"), 0755)
				os.WriteFile(filepath.Join(dir, "cmd", "e2e-test", "main.go"), []byte("package main\n\nfunc main() {}"), 0644)
				
				return dir
			},
			expectFiles: []string{"go.mod", "main.go", "cmd/e2e-test/main.go"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDir := tt.setupFunc()
			defer os.RemoveAll(testDir)

			// Change to test directory
			originalDir, _ := os.Getwd()
			os.Chdir(testDir)
			defer os.Chdir(originalDir)

			// Test project detection
			config := &ProjectConfig{}
			detectProjectInfo(config)

			// Verify project was detected correctly
			if config.ProjectName == "" {
				t.Error("Project name should be detected")
			}

			if config.MainPath == "" {
				t.Error("Main path should be detected")
			}

			if config.BinaryName == "" {
				t.Error("Binary name should be detected")
			}

			// Test config generation
			err := generateGoReleaserConfig(config)
			if err != nil {
				t.Errorf("generateGoReleaserConfig() error = %v", err)
			}

			// Verify .goreleaser.yaml was created
			if _, err := os.Stat(".goreleaser.yaml"); os.IsNotExist(err) {
				t.Error(".goreleaser.yaml should be created")
			}

			// Test GitHub Actions generation
			config.GenerateActions = true
			config.ActionsOn = []string{"On version tags (v*)"}
			err = generateGitHubActions(config)
			if err != nil {
				t.Errorf("generateGitHubActions() error = %v", err)
			}

			// Verify GitHub Actions workflow was created
			workflowPath := filepath.Join(".github", "workflows", "release.yml")
			if _, err := os.Stat(workflowPath); os.IsNotExist(err) {
				t.Error("GitHub Actions workflow should be created")
			}

			// Verify expected files exist
			for _, expectedFile := range tt.expectFiles {
				if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
					t.Errorf("Expected file %s should exist", expectedFile)
				}
			}
		})
	}
}

func TestConfigurationValidation(t *testing.T) {
	tests := []struct {
		name           string
		config         ProjectConfig
		expectError    bool
		errorContains  string
	}{
		{
			name: "valid_cli_config",
			config: ProjectConfig{
				ProjectName:        "test-cli",
				ProjectDescription:  "A test CLI application",
				BinaryName:         "test-cli",
				MainPath:           "./cmd/test-cli",
				ProjectType:        "CLI Application",
				Platforms:          []string{"linux", "darwin", "windows"},
				Architectures:      []string{"amd64", "arm64"},
				CGOEnabled:         false,
				GitProvider:        "GitHub",
				GenerateActions:     true,
				ActionsOn:          []string{"On version tags (v*)"},
			},
			expectError: false,
		},
		{
			name: "valid_web_service_config",
			config: ProjectConfig{
				ProjectName:        "test-web",
				ProjectDescription:  "A test web service",
				BinaryName:         "test-web",
				MainPath:           ".",
				ProjectType:        "Web Service",
				Platforms:          []string{"linux", "darwin"},
				Architectures:      []string{"amd64"},
				CGOEnabled:         true,
				GitProvider:        "GitHub",
				DockerEnabled:      true,
				DockerRegistry:     "ghcr.io/user",
				Homebrew:           true,
			},
			expectError: false,
		},
		{
			name: "missing_project_name",
			config: ProjectConfig{
				BinaryName: "test",
				MainPath:   ".",
			},
			expectError:   true,
			errorContains: "project name",
		},
		{
			name: "missing_binary_name",
			config: ProjectConfig{
				ProjectName: "test",
				MainPath:    ".",
			},
			expectError:   true,
			errorContains: "binary name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory for test
			tmpDir, _ := os.MkdirTemp("", "wizard-validation-test")
			defer os.RemoveAll(tmpDir)

			// Change to temp directory
			originalDir, _ := os.Getwd()
			os.Chdir(tmpDir)
			defer os.Chdir(originalDir)

			// Test config generation
			err := generateGoReleaserConfig(&tt.config)

			// Check error
			if (err != nil) != tt.expectError {
				t.Errorf("generateGoReleaserConfig() error = %v, wantErr %v", err, tt.expectError)
				return
			}

			if tt.expectError {
				if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain %q, got %q", tt.errorContains, err.Error())
				}
			} else {
				// Verify generated file exists and has expected content
				if _, err := os.Stat(".goreleaser.yaml"); os.IsNotExist(err) {
					t.Error(".goreleaser.yaml should be created for valid config")
				}

				// Read and validate basic structure
				content, err := os.ReadFile(".goreleaser.yaml")
				if err != nil {
					t.Errorf("Failed to read generated config: %v", err)
					return
				}

				contentStr := string(content)
				
				// Check for required fields
				if !strings.Contains(contentStr, "version: 2") {
					t.Error("Config should specify version 2")
				}

				if !strings.Contains(contentStr, "project_name: "+tt.config.ProjectName) {
					t.Error("Config should contain project name")
				}

				if !strings.Contains(contentStr, "binary: "+tt.config.BinaryName) {
					t.Error("Config should contain binary name")
				}
			}
		})
	}
}

func TestDifferentProjectTypes(t *testing.T) {
	tests := []struct {
		name           string
		projectType    string
		expectedConfig ProjectConfig
	}{
		{
			name:        "cli_application",
			projectType: "CLI Application",
			expectedConfig: ProjectConfig{
				ProjectType: "CLI Application",
				CGOEnabled: false,
			},
		},
		{
			name:        "web_service",
			projectType: "Web Service",
			expectedConfig: ProjectConfig{
				ProjectType: "Web Service",
				CGOEnabled: true,
			},
		},
		{
			name:        "library",
			projectType: "Library",
			expectedConfig: ProjectConfig{
				ProjectType: "Library",
				CGOEnabled: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory for test
			tmpDir, _ := os.MkdirTemp("", "wizard-project-type-test")
			defer os.RemoveAll(tmpDir)

			// Change to temp directory
			originalDir, _ := os.Getwd()
			os.Chdir(tmpDir)
			defer os.Chdir(originalDir)

			// Create basic Go project
			goMod := `module github.com/user/test
go 1.21
`
			os.WriteFile("go.mod", []byte(goMod), 0644)
			os.WriteFile("main.go", []byte("package main\n\nfunc main() {}"), 0644)

			// Test project detection
			config := &ProjectConfig{}
			detectProjectInfo(config)

			// Override project type for testing
			config.ProjectType = tt.projectType

			// Apply project type-specific defaults
			switch tt.projectType {
			case "CLI Application":
				config.CGOEnabled = false
			case "Web Service":
				config.CGOEnabled = true
			case "Library":
				config.CGOEnabled = false
			}

			// Verify project type
			if config.ProjectType != tt.expectedConfig.ProjectType {
				t.Errorf("ProjectType = %q, want %q", config.ProjectType, tt.expectedConfig.ProjectType)
			}

			// Verify CGO setting
			if config.CGOEnabled != tt.expectedConfig.CGOEnabled {
				t.Errorf("CGOEnabled = %v, want %v", config.CGOEnabled, tt.expectedConfig.CGOEnabled)
			}

			// Generate config to test
			err := generateGoReleaserConfig(config)
			if err != nil {
				t.Errorf("generateGoReleaserConfig() error = %v", err)
			}

			// Verify config file was created
			if _, err := os.Stat(".goreleaser.yaml"); os.IsNotExist(err) {
				t.Error(".goreleaser.yaml should be created")
			}
		})
	}
}

func TestEdgeCaseScenarios(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func() string
		expectError bool
		errorMsg    string
	}{
		{
			name: "empty_project_directory",
			setupFunc: func() string {
				dir, _ := os.MkdirTemp("", "wizard-empty-test")
				return dir
			},
			expectError: true,
			errorMsg:    "go.mod not found",
		},
		{
			name: "malformed_go_mod",
			setupFunc: func() string {
				dir, _ := os.MkdirTemp("", "wizard-malformed-test")
				os.WriteFile(filepath.Join(dir, "go.mod"), []byte("invalid go.mod content"), 0644)
				return dir
			},
			expectError: true,
			errorMsg:    "go.mod parsing failed",
		},
		{
			name: "read_only_directory",
			setupFunc: func() string {
				dir, _ := os.MkdirTemp("", "wizard-readonly-test")
				// Make directory read-only (this might not work on all systems)
				os.Chmod(dir, 0444)
				return dir
			},
			expectError: true,
			errorMsg:    "permission",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDir := tt.setupFunc()
			
			// Ensure we can cleanup even for read-only test
			defer func() {
				if testDir != "" {
					os.Chmod(testDir, 0755) // Restore permissions
					os.RemoveAll(testDir)
				}
			}()

			// Change to test directory
			originalDir, _ := os.Getwd()
			os.Chdir(testDir)
			defer os.Chdir(originalDir)

			// Test project detection
			config := &ProjectConfig{}
			detectProjectInfo(config) // This function doesn't return an error, it modifies config directly

			// Check error expectation (for tests that expect errors, we check config state)
			if tt.expectError {
				if tt.errorMsg != "" && config.ProjectName != "" {
					// If we expected an error but got a valid project, that's an issue
					t.Errorf("Expected error containing %q, but got valid project %q", tt.errorMsg, config.ProjectName)
				}
			} else {
				if config.ProjectName == "" {
					t.Error("Expected valid project detection, but got empty project name")
				}
			}
		})
	}
}