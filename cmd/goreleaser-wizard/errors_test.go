package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

func TestWizardError(t *testing.T) {
	tests := []struct {
		name       string
		errType    error
		message    string
		details    string
		suggestion string
		err        error
		want       string
	}{
		{
			name:       "basic_wizard_error",
			errType:    ErrConfigExists,
			message:    "Config already exists",
			details:    ".goreleaser.yaml found",
			suggestion: "Use --force to overwrite",
			err:        os.ErrExist,
			want:       "configuration already exists: Config already exists",
		},
		{
			name:       "minimal_error",
			errType:    ErrInvalidInput,
			message:    "Invalid input",
			details:    "",
			suggestion: "",
			err:        nil,
			want:       "invalid input: Invalid input",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wizErr := NewWizardError(tt.errType, tt.message, tt.details, tt.suggestion, tt.err)

			if wizErr.Error() != tt.want {
				t.Errorf("WizardError.Error() = %q, want %q", wizErr.Error(), tt.want)
			}

			if wizErr.Type != tt.errType {
				t.Errorf("WizardError.Type = %v, want %v", wizErr.Type, tt.errType)
			}

			if wizErr.Message != tt.message {
				t.Errorf("WizardError.Message = %q, want %q", wizErr.Message, tt.message)
			}

			if wizErr.Details != tt.details {
				t.Errorf("WizardError.Details = %q, want %q", wizErr.Details, tt.details)
			}

			if wizErr.Suggestion != tt.suggestion {
				t.Errorf("WizardError.Suggestion = %q, want %q", wizErr.Suggestion, tt.suggestion)
			}

			if wizErr.Unwrap() != tt.err {
				t.Errorf("WizardError.Unwrap() = %v, want %v", wizErr.Unwrap(), tt.err)
			}
		})
	}
}

func TestHandleError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		logger     *log.Logger
		wantOutput string
	}{
		{
			name:   "nil_error",
			err:    nil,
			logger: nil,
		},
		{
			name:       "wizard_error_with_details",
			err:        NewWizardError(ErrConfigExists, "Config exists", "Details here", "Use --force", os.ErrExist),
			logger:     nil,
			wantOutput: "❌ Error: Config exists\nDetails: Details here\n💡 Suggestion: Use --force",
		},
		{
			name:       "wizard_error_minimal",
			err:        NewWizardError(ErrInvalidInput, "Invalid input", "", "", nil),
			logger:     nil,
			wantOutput: "❌ Error: Invalid input",
		},
		{
			name:       "generic_error",
			err:        os.ErrPermission,
			logger:     nil,
			wantOutput: "❌ Error: permission denied\n💡 Suggestion: Try running with appropriate permissions or check file ownership",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			// Temporarily replace stdout
			originalStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Reset error styles to avoid ANSI codes in test output
			originalErrorStyle := errorStyle
			originalInfoStyle := infoStyle
			errorStyle = lipgloss.NewStyle()
			infoStyle = lipgloss.NewStyle()

			HandleError(tt.err, tt.logger)

			// Restore stdout
			w.Close()
			os.Stdout = originalStdout

			// Read captured output
			outputBuf := new(bytes.Buffer)
			outputBuf.ReadFrom(r)
			output := outputBuf.String()

			// Restore styles
			errorStyle = originalErrorStyle
			infoStyle = originalInfoStyle

			if len(tt.wantOutput) == 0 {
				if strings.TrimSpace(output) != "" {
					t.Errorf("HandleError() output = %q, want empty", output)
				}
			} else {
				if !strings.Contains(output, tt.wantOutput) {
					t.Errorf("HandleError() output = %q, want to contain %q", output, tt.wantOutput)
				}
			}
		})
	}
}

func TestGetSuggestionForError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "permission_error",
			err:      os.ErrPermission,
			expected: "Try running with appropriate permissions or check file ownership",
		},
		{
			name:     "not_found_error",
			err:      os.ErrNotExist,
			expected: "", // os.ErrNotExist.Error() is "file does not exist", not "not found"
		},
		{
			name:     "template_error",
			err:      NewWizardError(ErrTemplateExecution, "Template failed", "Details", "", nil),
			expected: "This might be a bug. Please report it at https://github.com/LarsArtmann/template-GoReleaser/issues",
		},
		{
			name:     "invalid_error",
			err:      NewWizardError(ErrInvalidInput, "Invalid", "Details", "", nil),
			expected: "Check your input and try again with valid values",
		},
		{
			name:     "connection_error",
			err:      &WizardError{Type: ErrDependency, Message: "Connection failed"},
			expected: "Check your internet connection and try again",
		},
		{
			name:     "unknown_error",
			err:      &WizardError{Type: ErrConfiguration, Message: "Unknown issue"},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getSuggestionForError(tt.err)
			if result != tt.expected {
				t.Errorf("getSuggestionForError() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestValidateFilePermissions(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func() string
		wantErr     bool
		errContains string
	}{
		{
			name: "valid_existing_directory",
			setupFunc: func() string {
				dir, _ := os.MkdirTemp("", "wizard-test")
				return dir
			},
			wantErr: false,
		},
		{
			name: "nonexistent_directory_creatable",
			setupFunc: func() string {
				return filepath.Join(os.TempDir(), "wizard-new-dir")
			},
			wantErr: false,
		},
		{
			name: "path_is_file",
			setupFunc: func() string {
				file, _ := os.CreateTemp("", "wizard-test-file")
				file.Close()
				return file.Name()
			},
			wantErr:     true,
			errContains: "not a directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testPath := tt.setupFunc()
			defer func() {
				if path := testPath; path != "" {
					os.RemoveAll(path)
				}
			}()

			err := ValidateFilePermissions(testPath)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFilePermissions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errContains != "" {
				if err == nil {
					t.Errorf("ValidateFilePermissions() expected error containing %q, got nil", tt.errContains)
					return
				}
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("ValidateFilePermissions() error = %v, want to contain %q", err, tt.errContains)
				}
			}
		})
	}
}

func TestSafeFileWrite(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		content     []byte
		perm        os.FileMode
		wantErr     bool
		errContains string
		setupFunc   func() string
	}{
		{
			name:    "write_new_file",
			path:    "test-new.txt",
			content: []byte("test content"),
			perm:    0644,
			wantErr: false,
		},
		{
			name:    "overwrite_existing_file",
			path:    "test-existing.txt",
			content: []byte("new content"),
			perm:    0644,
			wantErr: false,
			setupFunc: func() string {
				file, _ := os.Create("test-existing.txt")
				file.Write([]byte("original content"))
				file.Close()
				return "test-existing.txt"
			},
		},
		{
			name:        "write_to_invalid_path",
			path:        "/invalid/path/test.txt",
			content:     []byte("test"),
			perm:        0644,
			wantErr:     true,
			errContains: "Failed to write file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			if tt.setupFunc != nil {
				tt.path = tt.setupFunc()
			}

			defer func() {
				if tt.path != "" {
					os.Remove(tt.path)
					os.Remove(tt.path + ".backup")
				}
			}()

			// Test
			err := SafeFileWrite(tt.path, tt.content, tt.perm)

			if (err != nil) != tt.wantErr {
				t.Errorf("SafeFileWrite() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify file was written correctly
				data, err := os.ReadFile(tt.path)
				if err != nil {
					t.Errorf("Failed to read written file: %v", err)
					return
				}
				if string(data) != string(tt.content) {
					t.Errorf("File content = %q, want %q", string(data), string(tt.content))
				}

				// Verify permissions
				info, err := os.Stat(tt.path)
				if err != nil {
					t.Errorf("Failed to stat file: %v", err)
					return
				}
				// Note: permission mask on some systems, so we check for execute bits not being set
				if info.Mode().Perm()&0111 != 0 {
					t.Errorf("File has execute permissions, expected none")
				}
			} else if tt.errContains != "" {
				if err == nil {
					t.Errorf("SafeFileWrite() expected error containing %q, got nil", tt.errContains)
					return
				}
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("SafeFileWrite() error = %v, want to contain %q", err, tt.errContains)
				}
			}
		})
	}
}

func TestSafeReadFile(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		setupFunc   func() string
		wantErr     bool
		errContains string
		expected    string
	}{
		{
			name:     "read_existing_file",
			expected: "test content",
			wantErr:  false,
			setupFunc: func() string {
				file, _ := os.CreateTemp("", "wizard-read-test")
				file.Write([]byte("test content"))
				file.Close()
				return file.Name()
			},
		},
		{
			name:        "read_nonexistent_file",
			path:        "/nonexistent/file.txt",
			wantErr:     true,
			errContains: "Failed to read file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupFunc != nil {
				tt.path = tt.setupFunc()
				defer os.Remove(tt.path)
			}

			data, err := SafeReadFile(tt.path)

			if (err != nil) != tt.wantErr {
				t.Errorf("SafeReadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if string(data) != tt.expected {
					t.Errorf("SafeReadFile() = %q, want %q", string(data), tt.expected)
				}
			} else if tt.errContains != "" {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("SafeReadFile() error = %v, want to contain %q", err, tt.errContains)
				}
			}
		})
	}
}

func TestSafeCreateFile(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		wantErr     bool
		errContains string
	}{
		{
			name:    "create_file_in_existing_dir",
			path:    "test-create.txt",
			wantErr: false,
		},
		{
			name:    "create_file_with_subdirs",
			path:    "subdir/nested/file.txt",
			wantErr: false,
		},
		{
			name:        "create_file_in_invalid_path",
			path:        "/invalid/path/file.txt",
			wantErr:     true,
			errContains: "Cannot create directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if tt.path != "" {
					os.RemoveAll(filepath.Dir(tt.path))
				}
			}()

			file, err := SafeCreateFile(tt.path)

			if (err != nil) != tt.wantErr {
				t.Errorf("SafeCreateFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if file == nil {
					t.Error("SafeCreateFile() returned nil file")
					return
				}
				file.Close()

				// Verify file was created
				if _, err := os.Stat(tt.path); os.IsNotExist(err) {
					t.Errorf("File was not created at %s", tt.path)
				}
			} else if tt.errContains != "" {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("SafeCreateFile() error = %v, want to contain %q", err, tt.errContains)
				}
			}
		})
	}
}

func TestUserInputError(t *testing.T) {
	field := "project name"
	originalErr := os.ErrInvalid

	wizErr := UserInputError(field, originalErr)

	if wizErr.Type != ErrInvalidInput {
		t.Errorf("UserInputError() Type = %v, want %v", wizErr.Type, ErrInvalidInput)
	}

	expectedMsg := "Invalid project name"
	if wizErr.Message != expectedMsg {
		t.Errorf("UserInputError() Message = %q, want %q", wizErr.Message, expectedMsg)
	}

	if wizErr.Details != originalErr.Error() {
		t.Errorf("UserInputError() Details = %q, want %q", wizErr.Details, originalErr.Error())
	}

	expectedSuggestion := "Please provide valid input and try again"
	if wizErr.Suggestion != expectedSuggestion {
		t.Errorf("UserInputError() Suggestion = %q, want %q", wizErr.Suggestion, expectedSuggestion)
	}

	if wizErr.Err != originalErr {
		t.Errorf("UserInputError() Err = %v, want %v", wizErr.Err, originalErr)
	}
}

func TestTemplateError(t *testing.T) {
	templateName := "goreleaser.yaml"
	originalErr := os.ErrNotExist

	wizErr := TemplateError(templateName, originalErr)

	if wizErr.Type != ErrTemplateExecution {
		t.Errorf("TemplateError() Type = %v, want %v", wizErr.Type, ErrTemplateExecution)
	}

	expectedMsg := "Template error in goreleaser.yaml"
	if wizErr.Message != expectedMsg {
		t.Errorf("TemplateError() Message = %q, want %q", wizErr.Message, expectedMsg)
	}

	if wizErr.Err != originalErr {
		t.Errorf("TemplateError() Err = %v, want %v", wizErr.Err, originalErr)
	}
}

func TestWrapFileError(t *testing.T) {
	operation := "read"
	path := "/test/file.txt"
	originalErr := os.ErrPermission

	wizErr := WrapFileError(operation, path, originalErr)

	if wizErr.Type != ErrFileWrite {
		t.Errorf("WrapFileError() Type = %v, want %v", wizErr.Type, ErrFileWrite)
	}

	expectedMsg := "Failed to read"
	if wizErr.Message != expectedMsg {
		t.Errorf("WrapFileError() Message = %q, want %q", wizErr.Message, expectedMsg)
	}

	expectedDetails := "Error with /test/file.txt: permission denied"
	if wizErr.Details != expectedDetails {
		t.Errorf("WrapFileError() Details = %q, want %q", wizErr.Details, expectedDetails)
	}

	if wizErr.Err != originalErr {
		t.Errorf("WrapFileError() Err = %v, want %v", wizErr.Err, originalErr)
	}
}
