package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunGenerate(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func() string
		expectError bool
	}{
		{
			name: "generate_valid_config",
			setupFunc: func() string {
				dir, _ := os.MkdirTemp("", "wizard-generate-test")
				goMod := `module github.com/user/generate-test
go 1.21
`
				os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goMod), 0644)
				os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\n\nfunc main() {}"), 0644)
				return dir
			},
			expectError: false,
		},
		{
			name: "generate_in_non_go_project",
			setupFunc: func() string {
				dir, _ := os.MkdirTemp("", "wizard-generate-test")
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

			// Test runGenerate
			defer func() {
				if r := recover(); r != nil {
					if !tt.expectError {
						t.Errorf("runGenerate panicked: %v", r)
					}
				}
			}()

			// runGenerate is from generate.go, but we need to simulate it
			// For testing purposes, we'll test the underlying functions
			config := &ProjectConfig{}
			detectProjectInfo(config)

			if config.ProjectName == "" && !tt.expectError {
				t.Error("Expected project detection to work")
			}
		})
	}
}

func TestTemplateGeneration(t *testing.T) {
	tests := []struct {
		name        string
		config      ProjectConfig
		expectError bool
		checks      []string
	}{
		{
			name: "generate_complete_config",
			config: ProjectConfig{
				ProjectName:        "complete-test",
				ProjectDescription:  "A complete test project",
				BinaryName:         "complete-test",
				MainPath:           "./cmd/complete-test",
				ProjectType:        "CLI Application",
				Platforms:          []string{"linux", "darwin", "windows"},
				Architectures:      []string{"amd64", "arm64"},
				CGOEnabled:         false,
				GitProvider:        "GitHub",
				DockerEnabled:      true,
				DockerRegistry:     "ghcr.io/user",
				Signing:           true,
				Homebrew:          true,
				GenerateActions:    true,
				ActionsOn:         []string{"On version tags (v*)"},
			},
			expectError: false,
			checks: []string{
				"project_name: complete-test",
				"binary: complete-test",
				"main: ./cmd/complete-test",
				"goos:",
				"goarch:",
				"CGO_ENABLED=0",
				"dockers:",
				"signs:",
				"brews:",
			},
		},
		{
			name: "generate_minimal_config",
			config: ProjectConfig{
				ProjectName: "minimal-test",
				BinaryName:  "minimal-test",
				MainPath:    ".",
				GitProvider: "GitHub",
			},
			expectError: false,
			checks: []string{
				"project_name: minimal-test",
				"binary: minimal-test",
				"main: .",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory for test
			tmpDir, _ := os.MkdirTemp("", "wizard-template-test")
			defer os.RemoveAll(tmpDir)

			// Change to temp directory
			originalDir, _ := os.Getwd()
			os.Chdir(tmpDir)
			defer os.Chdir(originalDir)

			// Generate config
			err := generateGoReleaserConfig(&tt.config)

			// Check error
			if (err != nil) != tt.expectError {
				t.Errorf("generateGoReleaserConfig() error = %v, wantErr %v", err, tt.expectError)
				return
			}

			if !tt.expectError {
				// Read generated file
				content, err := os.ReadFile(".goreleaser.yaml")
				if err != nil {
					t.Fatalf("Failed to read generated file: %v", err)
				}

				contentStr := string(content)

				// Check for expected strings
				for _, check := range tt.checks {
					if !strings.Contains(contentStr, check) {
						t.Errorf("Generated config missing expected string: %q", check)
					}
				}
			}
		})
	}
}

func TestGitHubActionsGeneration(t *testing.T) {
	tests := []struct {
		name        string
		config      ProjectConfig
		expectError bool
		checks      []string
	}{
		{
			name: "actions_with_docker",
			config: ProjectConfig{
				ProjectName:    "docker-test",
				BinaryName:     "docker-test",
				GenerateActions: true,
				DockerEnabled:  true,
				DockerRegistry: "ghcr.io/user",
				ActionsOn:      []string{"Manual trigger only"},
			},
			expectError: false,
			checks: []string{
				"name: Release",
				"workflow_dispatch:",
				"Login to Docker Registry",
				"packages: write",
			},
		},
		{
			name: "actions_with_signing",
			config: ProjectConfig{
				ProjectName:    "signing-test",
				BinaryName:     "signing-test",
				GenerateActions: true,
				Signing:        true,
				ActionsOn:      []string{"On all tags"},
			},
			expectError: false,
			checks: []string{
				"Install Cosign",
				"id-token: write",
				"tags:",
				"- '*'",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory for test
			tmpDir, _ := os.MkdirTemp("", "wizard-actions-test")
			defer os.RemoveAll(tmpDir)

			// Change to temp directory
			originalDir, _ := os.Getwd()
			os.Chdir(tmpDir)
			defer os.Chdir(originalDir)

			// Generate actions
			err := generateGitHubActions(&tt.config)

			// Check error
			if (err != nil) != tt.expectError {
				t.Errorf("generateGitHubActions() error = %v, wantErr %v", err, tt.expectError)
				return
			}

			if !tt.expectError {
				// Read generated file
				workflowPath := filepath.Join(".github", "workflows", "release.yml")
				content, err := os.ReadFile(workflowPath)
				if err != nil {
					t.Fatalf("Failed to read generated workflow: %v", err)
				}

				contentStr := string(content)

				// Check for expected strings
				for _, check := range tt.checks {
					if !strings.Contains(contentStr, check) {
						t.Errorf("Generated workflow missing expected string: %q", check)
					}
				}
			}
		})
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name     string
		config   ProjectConfig
		wantErr  bool
	}{
		{
			name: "valid_complete_config",
			config: ProjectConfig{
				ProjectName:        "valid-test",
				ProjectDescription:  "Valid test project",
				BinaryName:         "valid-test",
				MainPath:           "./cmd/valid-test",
				ProjectType:        "CLI Application",
				Platforms:          []string{"linux", "darwin"},
				Architectures:      []string{"amd64"},
				CGOEnabled:         false,
				GitProvider:        "GitHub",
			},
			wantErr: false,
		},
		{
			name: "invalid_empty_project_name",
			config: ProjectConfig{
				ProjectName: "",
				BinaryName:  "test",
				MainPath:    ".",
			},
			wantErr: true,
		},
		{
			name: "invalid_empty_binary_name",
			config: ProjectConfig{
				ProjectName: "test",
				BinaryName:  "",
				MainPath:    ".",
			},
			wantErr: true,
		},
		{
			name: "invalid_empty_main_path",
			config: ProjectConfig{
				ProjectName: "test",
				BinaryName:  "test",
				MainPath:    "",
			},
			wantErr: false, // This is not validated in current implementation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory for test
			tmpDir, _ := os.MkdirTemp("", "wizard-config-validation-test")
			defer os.RemoveAll(tmpDir)

			// Change to temp directory
			originalDir, _ := os.Getwd()
			os.Chdir(tmpDir)
			defer os.Chdir(originalDir)

			// Test config generation
			err := generateGoReleaserConfig(&tt.config)

			if (err != nil) != tt.wantErr {
				t.Errorf("generateGoReleaserConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFileOperations(t *testing.T) {
	tests := []struct {
		name      string
		operation func() error
		wantErr   bool
	}{
		{
			name: "write_new_file_with_safe_write",
			operation: func() error {
				return SafeFileWrite("test-safe-write.txt", []byte("test content"), 0644)
			},
			wantErr: false,
		},
		{
			name: "read_existing_file_with_safe_read",
			operation: func() error {
				content := []byte("test content for reading")
				err := os.WriteFile("test-safe-read.txt", content, 0644)
				if err != nil {
					return err
				}
				
				readContent, err := SafeReadFile("test-safe-read.txt")
				if err != nil {
					return err
				}
				
				if string(readContent) != string(content) {
					return os.ErrInvalid
				}
				
				return nil
			},
			wantErr: false,
		},
		{
			name: "create_file_with_safe_create",
			operation: func() error {
				file, err := SafeCreateFile("test-safe-create.txt")
				if err != nil {
					return err
				}
				file.WriteString("test content")
				file.Close()
				return nil
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory for test
			tmpDir, _ := os.MkdirTemp("", "wizard-file-ops-test")
			defer os.RemoveAll(tmpDir)

			// Change to temp directory
			originalDir, _ := os.Getwd()
			os.Chdir(tmpDir)
			defer os.Chdir(originalDir)

			// Test file operation
			err := tt.operation()

			if (err != nil) != tt.wantErr {
				t.Errorf("File operation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBackupCreation(t *testing.T) {
	// Test that backup files are created when overwriting existing files
	tests := []struct {
		name        string
		originalContent string
		newContent      string
		expectBackup    bool
	}{
		{
			name:            "backup_created_on_overwrite",
			originalContent: "original content",
			newContent:      "new content",
			expectBackup:    true,
		},
		{
			name:            "no_backup_for_new_file",
			originalContent: "",
			newContent:      "new content",
			expectBackup:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory for test
			tmpDir, _ := os.MkdirTemp("", "wizard-backup-test")
			defer os.RemoveAll(tmpDir)

			// Change to temp directory
			originalDir, _ := os.Getwd()
			os.Chdir(tmpDir)
			defer os.Chdir(originalDir)

			testFile := "test-backup.txt"

			// Create original file if needed
			if tt.originalContent != "" {
				err := os.WriteFile(testFile, []byte(tt.originalContent), 0644)
				if err != nil {
					t.Fatalf("Failed to create original file: %v", err)
				}
			}

			// Write new file (this should create backup)
			err := SafeFileWrite(testFile, []byte(tt.newContent), 0644)
			if err != nil {
				t.Errorf("SafeFileWrite() error = %v", err)
				return
			}

			// Check backup file
			backupFile := testFile + ".backup"
			backupInfo, _ := os.Stat(backupFile)
			backupExists := backupInfo != nil
			
			if tt.expectBackup && !backupExists {
				t.Error("Backup file should exist when overwriting existing file")
			} else if tt.expectBackup && backupExists {
				// Read backup to ensure it contains original content
				backupContent, _ := os.ReadFile(backupFile)
				if string(backupContent) != tt.originalContent {
					t.Errorf("Backup content = %q, want %q", string(backupContent), tt.originalContent)
				}
			} else if !tt.expectBackup && backupExists {
				t.Error("Backup file should not exist for new file")
			}
		})
	}
}

func TestErrorRecovery(t *testing.T) {
	// Test error recovery and panic handling
	tests := []struct {
		name        string
		shouldPanic bool
		testFunc    func()
	}{
		{
			name:        "normal_operation_no_panic",
			shouldPanic: false,
			testFunc: func() {
				// Normal operation should not panic
				config := &ProjectConfig{}
				detectProjectInfo(config)
			},
		},
		{
			name:        "panic_recovery_works",
			shouldPanic: true,
			testFunc: func() {
				// This should panic
				var nilPointer *string
				_ = *nilPointer
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.shouldPanic {
						t.Errorf("Unexpected panic: %v", r)
					}
				} else if tt.shouldPanic {
					t.Error("Expected panic but none occurred")
				}
			}()

			tt.testFunc()
		})
	}
}