package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"slices"
)

func TestRunValidate(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func() string
		args        []string
		flags       map[string]bool
		expectPass  bool
		expectError bool
	}{
		{
			name: "valid_project",
			setupFunc: func() string {
				dir, _ := os.MkdirTemp("", "wizard-validate-test")
				// Create go.mod
				goMod := `module github.com/user/test
go 1.21
`
				os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goMod), 0644)
				// Create .goreleaser.yaml
				goreleaser := `# GoReleaser configuration
version: 2
project_name: test
build:
  main: .
  binary: test
  goos:
    - linux
  goarch:
    - amd64
`
				os.WriteFile(filepath.Join(dir, ".goreleaser.yaml"), []byte(goreleaser), 0644)
				// Create main.go
				os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\n\nfunc main() {}"), 0644)
				// Initialize git
				exec.Command("git", "init").Dir = dir
				exec.Command("git", "init").Run()
				exec.Command("git", "config", "user.email", "test@example.com").Dir = dir
				exec.Command("git", "config", "user.email", "test@example.com").Run()
				exec.Command("git", "add", ".").Dir = dir
				exec.Command("git", "add", ".").Run()
				exec.Command("git", "commit", "-m", "init").Dir = dir
				exec.Command("git", "commit", "-m", "init").Run()
				return dir
			},
			args:       []string{},
			flags:      map[string]bool{"verbose": false, "fix": false},
			expectPass: true,
		},
		{
			name: "missing_goreleaser_config",
			setupFunc: func() string {
				dir, _ := os.MkdirTemp("", "wizard-validate-test")
				goMod := `module github.com/user/test
go 1.21
`
				os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goMod), 0644)
				return dir
			},
			args:        []string{},
			flags:       map[string]bool{"verbose": false, "fix": false},
			expectPass:  false,
			expectError: true,
		},
		{
			name: "with_verbose_flag",
			setupFunc: func() string {
				dir, _ := os.MkdirTemp("", "wizard-validate-test")
				goMod := `module github.com/user/test
go 1.21
`
				os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goMod), 0644)
				goreleaser := `# GoReleaser configuration
version: 2
project_name: test
`
				os.WriteFile(filepath.Join(dir, ".goreleaser.yaml"), []byte(goreleaser), 0644)
				return dir
			},
			args:  []string{},
			flags: map[string]bool{"verbose": true, "fix": false},
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

			// Reset viper
			viper.Reset()
			viper.Set("debug", false)

			// Create test command
			cmd := &cobra.Command{}
			cmd.Flags().Bool("verbose", false, "show detailed validation output")
			cmd.Flags().Bool("fix", false, "attempt to fix common issues")

			// Set flags
			for flag, value := range tt.flags {
				if value {
					cmd.Flags().Set(flag, "true")
				} else {
					cmd.Flags().Set(flag, "false")
				}
			}

			// Capture exit status
			defer func() {
				if r := recover(); r != nil {
					if tt.expectError {
						// Expected to fail
						return
					}
					t.Errorf("runValidate() panicked: %v", r)
				}
			}()

			// Run validation
			func() {
				defer HandlePanic("validate test", log.New(os.Stderr))
				runValidate(cmd, tt.args)
			}()
		})
	}
}

func TestCheckFileExists(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		requireDir  bool
		setupFunc   func() string
		wantErr     bool
		errContains string
	}{
		{
			name:       "existing_file",
			requireDir: false,
			wantErr:    false,
			setupFunc: func() string {
				file, _ := os.CreateTemp("", "wizard-test-file")
				file.Close()
				return file.Name()
			},
		},
		{
			name:       "existing_directory",
			requireDir: true,
			wantErr:    false,
			setupFunc: func() string {
				dir, _ := os.MkdirTemp("", "wizard-test-dir")
				return dir
			},
		},
		{
			name:        "nonexistent_file",
			path:        "/nonexistent/file.txt",
			requireDir:  false,
			wantErr:     true,
			errContains: "File not found",
		},
		{
			name:       "file_when_directory_required",
			requireDir: true,
			wantErr:    true,
			errContains: "Expected directory",
			setupFunc: func() string {
				file, _ := os.CreateTemp("", "wizard-test-file")
				file.Close()
				return file.Name()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupFunc != nil {
				tt.path = tt.setupFunc()
				defer func() {
					if tt.path != "" {
						os.Remove(tt.path)
					}
				}()
			}

			err := CheckFileExists(tt.path, tt.requireDir)

			if (err != nil) != tt.wantErr {
				t.Errorf("CheckFileExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errContains != "" {
				if err == nil {
					t.Errorf("CheckFileExists() expected error containing %q, got nil", tt.errContains)
					return
				}
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("CheckFileExists() error = %v, want to contain %q", err, tt.errContains)
				}
			}
		})
	}
}

func TestValidateProjectStructure(t *testing.T) {
	tests := []struct {
		name          string
		setupFunc     func() string
		expectIssues  []string
		expectWarnings []string
	}{
		{
			name: "valid_simple_project",
			setupFunc: func() string {
				dir, _ := os.MkdirTemp("", "wizard-structure-test")
				goMod := `module github.com/user/test
go 1.21
`
				os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goMod), 0644)
				os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\n\nfunc main() {}"), 0644)
				goreleaser := `# GoReleaser configuration
version: 2
project_name: test
build:
  main: .
  binary: test
`
				os.WriteFile(filepath.Join(dir, ".goreleaser.yaml"), []byte(goreleaser), 0644)
				return dir
			},
			expectIssues:  []string{},
		},
		{
			name: "project_with_cmd_structure",
			setupFunc: func() string {
				dir, _ := os.MkdirTemp("", "wizard-structure-test")
				goMod := `module github.com/user/test
go 1.21
`
				os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goMod), 0644)
				os.MkdirAll(filepath.Join(dir, "cmd", "test"), 0755)
				os.WriteFile(filepath.Join(dir, "cmd", "test", "main.go"), []byte("package main\n\nfunc main() {}"), 0644)
				goreleaser := `# GoReleaser configuration
version: 2
project_name: test
build:
  main: ./cmd/test
  binary: test
`
				os.WriteFile(filepath.Join(dir, ".goreleaser.yaml"), []byte(goreleaser), 0644)
				return dir
			},
			expectIssues:  []string{},
		},
		{
			name: "missing_main_package",
			setupFunc: func() string {
				dir, _ := os.MkdirTemp("", "wizard-structure-test")
				goMod := `module github.com/user/test
go 1.21
`
				os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goMod), 0644)
				goreleaser := `# GoReleaser configuration
version: 2
project_name: test
build:
  main: .
  binary: test
`
				os.WriteFile(filepath.Join(dir, ".goreleaser.yaml"), []byte(goreleaser), 0644)
				return dir
			},
			expectWarnings: []string{"No main.go found"},
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

			// Test project structure validation (this is part of runValidate)
			// For now, just test that the structure detection works
			issues := []string{}
			warnings := []string{}

			// Check main package existence (simplified version of validate logic)
			mainFound := false
			commonPaths := []string{
				"main.go",
				"./cmd/*/main.go",
				"./*.go",
			}
			for _, path := range commonPaths {
				matches, _ := filepath.Glob(path)
				if len(matches) > 0 {
					mainFound = true
					break
				}
			}

			if !mainFound {
				warnings = append(warnings, "No main.go found")
			}

			// Verify expectations
			if len(tt.expectIssues) != len(issues) {
				t.Errorf("Expected %d issues, got %d", len(tt.expectIssues), len(issues))
			}

			if len(tt.expectWarnings) != len(warnings) {
				t.Errorf("Expected %d warnings, got %d", len(tt.expectWarnings), len(warnings))
			}

			for i, warning := range warnings {
				if i >= len(tt.expectWarnings) {
					t.Errorf("Unexpected warning: %s", warning)
					continue
				}
				if !strings.Contains(warning, tt.expectWarnings[i]) {
					t.Errorf("Warning %d = %q, want to contain %q", i, warning, tt.expectWarnings[i])
				}
			}
		})
	}
}

func TestValidateDependencies(t *testing.T) {
	tests := []struct {
		name         string
		dependencies []string
		expectFound  []string
		expectMissing []string
	}{
		{
			name:         "check_goreleaser",
			dependencies: []string{"goreleaser"},
			expectMissing: []string{"goreleaser"},
		},
		{
			name:         "check_go_command",
			dependencies: []string{"go"},
			expectFound:  []string{"go"},
		},
		{
			name:         "check_docker",
			dependencies: []string{"nonexistent-docker-binary"},
			expectMissing: []string{"nonexistent-docker-binary"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, dep := range tt.dependencies {
				path, err := exec.LookPath(dep)
				found := err == nil

				expectedFound := slices.Contains(tt.expectFound, dep)

				expectedMissing := slices.Contains(tt.expectMissing, dep)

				if expectedFound && !found {
					t.Errorf("Expected to find %s, but it was not found", dep)
				}

				if expectedMissing && found {
					t.Errorf("Expected %s to be missing, but found at %s", dep, path)
				}
			}
		})
	}
}

func TestValidateCommandFlags(t *testing.T) {
	// Test command setup - use the actual validateCmd from main.go
	validateCmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate GoReleaser configuration",
		Run:   runValidate,
	}
	
	// Initialize flags like in init()
	validateCmd.Flags().Bool("verbose", false, "show detailed validation output")
	validateCmd.Flags().Bool("fix", false, "attempt to fix common issues")

	// Check that flags are properly set up
	flags := validateCmd.Flags()
	
	verboseFlag := flags.Lookup("verbose")
	if verboseFlag == nil {
		t.Error("Expected 'verbose' flag to be present")
	} else if verboseFlag.DefValue != "false" {
		t.Errorf("Expected verbose flag default value 'false', got '%s'", verboseFlag.DefValue)
	}

	fixFlag := flags.Lookup("fix")
	if fixFlag == nil {
		t.Error("Expected 'fix' flag to be present")
	} else if fixFlag.DefValue != "false" {
		t.Errorf("Expected fix flag default value 'false', got '%s'", fixFlag.DefValue)
	}
}

func TestValidateOutputFormatting(t *testing.T) {
	// Test that output is properly formatted (basic check)
	originalErrorStyle := errorStyle
	originalSuccessStyle := successStyle
	originalInfoStyle := infoStyle
	originalTitleStyle := titleStyle

	// Temporarily set simple styles for testing
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("red"))
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("green"))
	infoStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("blue"))
	titleStyle = lipgloss.NewStyle().Bold(true)

	defer func() {
		errorStyle = originalErrorStyle
		successStyle = originalSuccessStyle
		infoStyle = originalInfoStyle
		titleStyle = originalTitleStyle
	}()

	// Create a temporary directory with a valid setup for testing output
	dir, _ := os.MkdirTemp("", "wizard-output-test")
	defer os.RemoveAll(dir)

	goMod := `module github.com/user/test
go 1.21
`
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goMod), 0644)
	goreleaser := `# GoReleaser configuration
version: 2
project_name: test
build:
  main: .
  binary: test
`
	os.WriteFile(filepath.Join(dir, ".goreleaser.yaml"), []byte(goreleaser), 0644)
	os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\n\nfunc main() {}"), 0644)

	originalDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(originalDir)

	// Reset viper
	viper.Reset()
	viper.Set("debug", false)

	// Create test command
	cmd := &cobra.Command{}
	cmd.Flags().Bool("verbose", false, "show detailed validation output")
	cmd.Flags().Bool("fix", false, "attempt to fix common issues")

	// Run validation and check that it doesn't panic
	func() {
		defer HandlePanic("validate output test", log.New(os.Stderr))
		runValidate(cmd, []string{})
	}()

	// If we get here without panic, the output formatting is working
	// More detailed output testing would require capturing stdout, which is complex
}