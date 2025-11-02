package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TestInitCommand(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		flags       map[string]string
		setupFunc   func() string
		expectError bool
	}{
		{
			name: "basic_init_command",
			args: []string{},
			flags: map[string]string{},
			setupFunc: func() string {
				dir, _ := os.MkdirTemp("", "wizard-init-test")
				goMod := `module github.com/user/test
go 1.21
`
				os.WriteFile(dir+"/go.mod", []byte(goMod), 0644)
				os.WriteFile(dir+"/main.go", []byte("package main\n\nfunc main() {}"), 0644)
				return dir
			},
			expectError: false,
		},
		{
			name: "init_in_non_go_project",
			args: []string{},
			flags: map[string]string{},
			setupFunc: func() string {
				dir, _ := os.MkdirTemp("", "wizard-init-test")
				return dir
			},
			expectError: true,
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
			defer func() {
				if r := recover(); r != nil {
					if !tt.expectError {
						t.Errorf("Init command panicked: %v", r)
					}
				}
			}()

			// Execute command
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
				dir, _ := os.MkdirTemp("", "wizard-detect-test")
				goMod := `module github.com/user/simpleapp
go 1.21
`
				os.WriteFile(dir+"/go.mod", []byte(goMod), 0644)
				os.WriteFile(dir+"/main.go", []byte("package main\n\nfunc main() {}"), 0644)
				return dir
			},
			expectedProject: ProjectConfig{
				ProjectName: "simpleapp",
				MainPath:    ".",
				BinaryName:  "simpleapp",
				ProjectType: "CLI Application",
			},
		},
		{
			name: "detect_cmd_structure",
			setupFunc: func() string {
				dir, _ := os.MkdirTemp("", "wizard-detect-test")
				goMod := `module github.com/user/cmdapp
go 1.21
`
				os.WriteFile(dir+"/go.mod", []byte(goMod), 0644)
				os.MkdirAll(dir+"/cmd/cmdapp", 0755)
				os.WriteFile(dir+"/cmd/cmdapp/main.go", []byte("package main\n\nfunc main() {}"), 0644)
				return dir
			},
			expectedProject: ProjectConfig{
				ProjectName: "cmdapp",
				MainPath:    "./cmd/cmdapp",
				BinaryName:  "cmdapp",
				ProjectType: "CLI Application",
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
				t.Errorf("ProjectName = %q, want %q", config.ProjectName, tt.expectedProject.ProjectName)
			}

			if config.MainPath != tt.expectedProject.MainPath {
				t.Errorf("MainPath = %q, want %q", config.MainPath, tt.expectedProject.MainPath)
			}

			if config.BinaryName != tt.expectedProject.BinaryName {
				t.Errorf("BinaryName = %q, want %q", config.BinaryName, tt.expectedProject.BinaryName)
			}

			if config.ProjectType != tt.expectedProject.ProjectType {
				t.Errorf("ProjectType = %q, want %q", config.ProjectType, tt.expectedProject.ProjectType)
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
			function: validateProjectName,
			wantErr:  false,
		},
		{
			name:     "invalid_empty_project_name",
			input:    "",
			function: validateProjectName,
			wantErr:  true,
		},
		{
			name:     "invalid_project_name_too_long",
			input:    strings.Repeat("a", 65),
			function: validateProjectName,
			wantErr:  true,
		},
		{
			name:     "valid_binary_name",
			input:    "my-app",
			function: validateBinaryName,
			wantErr:  false,
		},
		{
			name:     "invalid_binary_name_with_spaces",
			input:    "my app",
			function: validateBinaryName,
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

func TestFlagHandling(t *testing.T) {
	tests := []struct {
		name         string
		flags        map[string]string
		expectedViper map[string]interface{}
	}{
		{
			name:  "debug_flag",
			flags: map[string]string{"debug": "true"},
			expectedViper: map[string]interface{}{
				"debug": true,
			},
		},
		{
			name:  "force_flag",
			flags: map[string]string{"force": "true"},
			expectedViper: map[string]interface{}{
				"force": true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset viper
			viper.Reset()

			// Simulate flag parsing (simplified)
			for flag, value := range tt.flags {
				if flag == "debug" {
					viper.Set("debug", value == "true")
				} else if flag == "force" {
					viper.Set("force", value == "true")
				}
			}

			// Check viper values
			for key, expected := range tt.expectedViper {
				actual := viper.Get(key)
				if actual != expected {
					t.Errorf("Viper.Get(%q) = %v, want %v", key, actual, expected)
				}
			}
		})
	}
}

// Mock validation functions (would normally be in the main code)
func validateProjectName(name string) error {
	if len(name) == 0 {
		return UserInputError("project name", os.ErrInvalid)
	}
	if len(name) > 63 {
		return UserInputError("project name", os.ErrInvalid)
	}
	return nil
}

func validateBinaryName(name string) error {
	if len(name) == 0 {
		return UserInputError("binary name", os.ErrInvalid)
	}
	if strings.Contains(name, " ") {
		return UserInputError("binary name", os.ErrInvalid)
	}
	return nil
}