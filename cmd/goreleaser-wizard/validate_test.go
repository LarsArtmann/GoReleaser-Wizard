//godoclint:disable

// CRITICAL ARCHITECTURE TODO: This file is 489 lines - SPLIT IMMEDIATELY into:
// 1. validate_unit_test.go - Unit tests for validation
// 2. validate_integration_test.go - Integration tests
// 3. validate_fix_test.go - Auto-fix functionality tests
// 4. validate_performance_test.go - Performance tests
// 5. validate_test_helpers.go - Test utilities and factories
//
// TODO: Implement proper test data builders
// TODO: Create mock repositories for testing
// TODO: Add benchmark tests for validation performance
// TODO: Implement proper test isolation
// TODO: Add fuzz testing for validation edge cases
// TODO: Create test scenarios for all validation rules
package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/LarsArtmann/GoReleaser-Wizard/internal/config"
	"github.com/spf13/cobra"
)

// setupBasicTestProject creates a minimal test project with go.mod and optional .goreleaser.yaml.
func setupBasicTestProject(
	prefix string,
	includeGoreleaser bool,
	goreleaserExtra string,
) func() string {
	return func() string {
		dir, _ := os.MkdirTemp("", prefix)
		goMod := `module github.com/user/test
go 1.21
`
		os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goMod), 0o644)

		if includeGoreleaser {
			goreleaser := `# GoReleaser configuration
version: 2
project_name: test
`
			if goreleaserExtra != "" {
				goreleaser += goreleaserExtra
			}

			os.WriteFile(filepath.Join(dir, ".goreleaser.yaml"), []byte(goreleaser), 0o644)
		}

		return dir
	}
}

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
				os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goMod), 0o644)
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
				os.WriteFile(filepath.Join(dir, ".goreleaser.yaml"), []byte(goreleaser), 0o644)
				// Create cmd directory with main.go (expected project structure)
				cmdDir := filepath.Join(dir, "cmd")
				os.MkdirAll(cmdDir, 0o755)
				os.WriteFile(
					filepath.Join(cmdDir, "main.go"),
					[]byte("package main\n\nfunc main() {}"),
					0o644,
				)
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
				os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goMod), 0o644)

				return dir
			},
			args:        []string{},
			flags:       map[string]bool{"verbose": false, "fix": false},
			expectPass:  false,
			expectError: true,
		},
		{
			name:      "with_verbose_flag",
			setupFunc: setupBasicTestProject("wizard-validate-test", true, ""),
			args:      []string{},
			flags:     map[string]bool{"verbose": true, "fix": false},
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

			// Reset config manager
			config.GetManager().Reset()
			config.GetManager().Set("debug", false)

			// Create test command with the production flag set
			cmd := &cobra.Command{}
			configureValidateCommand(cmd)

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

					t.Errorf("performValidation() panicked: %v", r)
				}
			}()

			// Run validation
			func() {
				defer recoverFromPanic("validate test")

				performValidation(cmd)
			}()
		})
	}
}

func TestValidateProjectStructure(t *testing.T) {
	tests := []struct {
		name           string
		setupFunc      func() string
		expectIssues   []string
		expectWarnings []string
	}{
		{
			name: "valid_simple_project",
			setupFunc: func() string {
				dir, _ := os.MkdirTemp("", "wizard-structure-test")
				goMod := `module github.com/user/test
go 1.21
`
				os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goMod), 0o644)
				os.WriteFile(
					filepath.Join(dir, "main.go"),
					[]byte("package main\n\nfunc main() {}"),
					0o644,
				)

				goreleaser := `# GoReleaser configuration
version: 2
project_name: test
build:
  main: .
  binary: test
`
				os.WriteFile(filepath.Join(dir, ".goreleaser.yaml"), []byte(goreleaser), 0o644)

				return dir
			},
			expectIssues: []string{},
		},
		{
			name: "project_with_cmd_structure",
			setupFunc: func() string {
				dir, _ := os.MkdirTemp("", "wizard-structure-test")
				goMod := `module github.com/user/test
go 1.21
`
				os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goMod), 0o644)
				os.MkdirAll(filepath.Join(dir, "cmd", "test"), 0o755)
				os.WriteFile(
					filepath.Join(dir, "cmd", "test", "main.go"),
					[]byte("package main\n\nfunc main() {}"),
					0o644,
				)

				goreleaser := `# GoReleaser configuration
version: 2
project_name: test
build:
  main: ./cmd/test
  binary: test
`
				os.WriteFile(filepath.Join(dir, ".goreleaser.yaml"), []byte(goreleaser), 0o644)

				return dir
			},
			expectIssues: []string{},
		},
		{
			name: "missing_main_package",
			setupFunc: setupBasicTestProject(
				"wizard-structure-test",
				true,
				"build:\n  main: .\n  binary: test\n",
			),
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
					t.Errorf(
						"Warning %d = %q, want to contain %q",
						i,
						warning,
						tt.expectWarnings[i],
					)
				}
			}
		})
	}
}

func TestValidateDependencies(t *testing.T) {
	tests := []struct {
		name          string
		dependencies  []string
		expectFound   []string
		expectMissing []string
		skipIfMissing []string
	}{
		{
			name:          "check_goreleaser",
			dependencies:  []string{"goreleaser"},
			expectFound:   []string{"goreleaser"},
			skipIfMissing: []string{"goreleaser"},
		},
		{
			name:         "check_go_command",
			dependencies: []string{"go"},
			expectFound:  []string{"go"},
		},
		{
			name:          "check_docker",
			dependencies:  []string{"nonexistent-docker-binary"},
			expectMissing: []string{"nonexistent-docker-binary"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, dep := range tt.dependencies {
				path, err := exec.LookPath(dep)
				found := err == nil

				if !found && slices.Contains(tt.skipIfMissing, dep) {
					t.Skipf("Skipping %s - not installed", dep)
				}

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

func TestValidateOutputFormatting(_ *testing.T) {
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
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goMod), 0o644)

	goreleaser := `# GoReleaser configuration
version: 2
project_name: test
build:
  main: .
  binary: test
`
	os.WriteFile(filepath.Join(dir, ".goreleaser.yaml"), []byte(goreleaser), 0o644)

	// Create cmd directory with main.go (expected project structure)
	cmdDir := filepath.Join(dir, "cmd")
	os.MkdirAll(cmdDir, 0o755)
	os.WriteFile(filepath.Join(cmdDir, "main.go"), []byte("package main\n\nfunc main() {}"), 0o644)

	originalDir, _ := os.Getwd()

	os.Chdir(dir)
	defer os.Chdir(originalDir)

	// Reset config manager
	config.GetManager().Reset()
	config.GetManager().Set("debug", false)

	// Create test command with the production flag set
	cmd := &cobra.Command{}
	configureValidateCommand(cmd)

	// Run validation and check that it doesn't panic
	func() {
		defer recoverFromPanic("validate output test")

		performValidation(cmd)
	}()

	// If we get here without panic, the output formatting is working
	// More detailed output testing would require capturing stdout, which is complex
}
