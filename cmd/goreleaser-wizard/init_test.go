package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/LarsArtmann/GoReleaser-Wizard/internal/config"
	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
	"github.com/LarsArtmann/GoReleaser-Wizard/internal/validation"
	"github.com/spf13/cobra"
)

// setupBasicGoProject creates a basic Go project with go.mod and main.go.
func setupBasicGoProject(t *testing.T, moduleName, pattern string) string {
	t.Helper()

	dir, err := os.MkdirTemp("", pattern)
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	goMod := "module " + moduleName + "\ngo 1.21\n"
	os.WriteFile(dir+"/go.mod", []byte(goMod), 0o644)
	os.WriteFile(dir+"/main.go", []byte("package main\n\nfunc main() {}\n"), 0o644)

	return dir
}

// expectNoPanic sets up a deferred panic recovery that fails the test if a panic occurs.
func expectNoPanic(t *testing.T, expectError bool) {
	t.Helper()

	if r := recover(); r != nil {
		if !expectError {
			t.Errorf("Command panicked: %v", r)
		}
	}
}

func TestInitCommand(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		flags       map[string]string
		setupFunc   func() string
		expectError bool
	}{
		{
			name:  "basic_init_command",
			args:  []string{},
			flags: map[string]string{},
			setupFunc: func() string {
				return setupBasicGoProject(t, "github.com/user/myapp", "wizard-init-test")
			},
			expectError: false,
		},
		{
			name:  "init_in_non_go_project",
			args:  []string{},
			flags: map[string]string{},
			setupFunc: func() string {
				dir := CreateTempDir("wizard-init-test")

				return dir
			},
			// Note: The command uses Run (not RunE), so it displays errors but doesn't return them
			// The error is logged via displayError() and the command returns nil
			expectError: false,
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

			// Create command
			cmd := &cobra.Command{
				Use:   "init",
				Short: "Initialize GoReleaser configuration",
				Run:   runInitWizard,
			}

			// Add flags (simplified version)
			cmd.Flags().Bool("force", false, "force overwrite existing configuration")

			// Set flag values
			for flag, value := range tt.flags {
				cmd.Flags().Set(flag, value)
			}

			// Execute command (with panic recovery)
			defer expectNoPanic(t, tt.expectError)

			// Execute command
			// Explicit empty args: without SetArgs, Cobra parses
			// os.Args[1:], so `go test -update-golden` leaked the flag
			// into the CLI and failed as unknown-flag (gotcha 14b).
			cmd.SetArgs([]string{})
			err := cmd.Execute()

			if (err != nil) != tt.expectError {
				t.Errorf("Init command error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestProjectDetection(t *testing.T) {
	tests := []struct {
		name            string
		setupFunc       func() string
		expectedProject ProjectConfig
	}{
		{
			name: "detect_simple_project",
			setupFunc: func() string {
				return setupBasicGoProject(t, "github.com/user/simpleapp", "wizard-detect-test")
			},
			expectedProject: ProjectConfig{
				ProjectName: "simpleapp",
				MainPath:    ".",
				BinaryName:  "simpleapp",
				ProjectType: domain.ProjectTypeCLI,
			},
		},
		{
			name: "detect_cmd_structure",
			setupFunc: func() string {
				dir, _ := os.MkdirTemp("", "wizard-detect-test")
				goMod := `module github.com/user/cmdapp
go 1.21
`
				os.WriteFile(dir+"/go.mod", []byte(goMod), 0o644)
				os.MkdirAll(dir+"/cmd/cmdapp", 0o755)
				os.WriteFile(
					dir+"/cmd/cmdapp/main.go",
					[]byte("package main\n\nfunc main() {}"),
					0o644,
				)

				return dir
			},
			expectedProject: ProjectConfig{
				ProjectName: "cmdapp",
				MainPath:    "./cmd/cmdapp",
				BinaryName:  "cmdapp",
				ProjectType: domain.ProjectTypeCLI,
			},
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

			// Check detected information
			if config.ProjectName != tt.expectedProject.ProjectName {
				t.Errorf(
					"ProjectName = %q, want %q",
					config.ProjectName,
					tt.expectedProject.ProjectName,
				)
			}

			if config.MainPath != tt.expectedProject.MainPath {
				t.Errorf("MainPath = %q, want %q", config.MainPath, tt.expectedProject.MainPath)
			}

			if config.BinaryName != tt.expectedProject.BinaryName {
				t.Errorf(
					"BinaryName = %q, want %q",
					config.BinaryName,
					tt.expectedProject.BinaryName,
				)
			}

			if config.ProjectType != tt.expectedProject.ProjectType {
				t.Errorf(
					"ProjectType = %q, want %q",
					config.ProjectType,
					tt.expectedProject.ProjectType,
				)
			}
		})
	}
}

func TestFormValidation(t *testing.T) {
	// Test form field validation functions
	tests := []struct {
		name     string
		input    string
		function func(string) error
		wantErr  bool
	}{
		{
			name:     "valid_project_name",
			input:    "my-awesome-project",
			function: validation.ValidateProjectName,
			wantErr:  false,
		},
		{
			name:     "invalid_empty_project_name",
			input:    "",
			function: validation.ValidateProjectName,
			wantErr:  true,
		},
		{
			name:     "invalid_project_name_too_long",
			input:    strings.Repeat("a", 65),
			function: validation.ValidateProjectName,
			wantErr:  true,
		},
		{
			name:     "valid_binary_name",
			input:    "my-app",
			function: validation.ValidateBinaryName,
			wantErr:  false,
		},
		{
			name:     "invalid_binary_name_with_spaces",
			input:    "my app",
			function: validation.ValidateBinaryName,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.function(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("%s(%q) error = %v, wantErr %v", t.Name(), tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestCommandHelp(t *testing.T) {
	tests := []struct {
		name        string
		command     *cobra.Command
		args        []string
		expectUsage bool
	}{
		{
			name: "init_help",
			command: &cobra.Command{
				Use:   "init",
				Short: "Initialize GoReleaser configuration",
				Long:  "Interactive wizard to create GoReleaser configuration",
				Run:   runInitWizard,
			},
			args:        []string{"--help"},
			expectUsage: true,
		},
		{
			name: "validate_help",
			command: &cobra.Command{
				Use:   "validate",
				Short: "Validate GoReleaser configuration",
				Long:  "Validate your GoReleaser configuration and check for common issues",
				Run:   runValidate,
			},
			args:        []string{"--help"},
			expectUsage: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture output
			var buf bytes.Buffer
			tt.command.SetOut(&buf)
			tt.command.SetErr(&buf)

			// Execute with help flag
			tt.command.SetArgs(tt.args)
			err := tt.command.Execute()
			// Help should not return an error
			if err != nil {
				t.Errorf("Help command returned error: %v", err)
			}

			output := buf.String()
			if tt.expectUsage && !strings.Contains(output, "Usage:") {
				t.Errorf("Expected help output to contain 'Usage:', got: %s", output)
			}
		})
	}
}

// makeFlagTestCase creates a flag test case with the given flag name and expected boolean value.
func makeFlagTestCase(name, flagName string, flagValue bool) struct {
	name          string
	flags         map[string]string
	expectedViper map[string]any
} {
	return struct {
		name          string
		flags         map[string]string
		expectedViper map[string]any
	}{
		name:  name,
		flags: map[string]string{flagName: "true"},
		expectedViper: map[string]any{
			flagName: flagValue,
		},
	}
}

func TestFlagHandling(t *testing.T) {
	tests := []struct {
		name          string
		flags         map[string]string
		expectedViper map[string]any
	}{
		makeFlagTestCase("debug_flag", "debug", true),
		makeFlagTestCase("force_flag", "force", true),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset config manager
			config.GetManager().Reset()

			// Simulate flag parsing (simplified)
			for flag, value := range tt.flags {
				switch flag {
				case "debug":
					config.GetManager().Set("debug", value == "true")
				case "force":
					config.GetManager().Set("force", value == "true")
				default:
					config.GetManager().Set(flag, value)
				}
			}

			// Check config values
			for key, expected := range tt.expectedViper {
				actual := config.GetManager().GetRaw(key)
				if actual != expected {
					t.Errorf("Config.Get(%q) = %v, want %v", key, actual, expected)
				}
			}
		})
	}
}
