package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

// Custom error types for better error handling
var (
	ErrConfigExists      = errors.New("configuration already exists")
	ErrProjectNotFound   = errors.New("project not found")
	ErrInvalidInput      = errors.New("invalid input")
	ErrTemplateExecution = errors.New("template execution failed")
	ErrFileWrite         = errors.New("file write failed")
	ErrFileRead          = errors.New("file read failed")
	ErrPermission        = errors.New("permission denied")
	ErrDependency        = errors.New("missing dependency")
	ErrConfiguration     = errors.New("configuration error")
	ErrUserInput         = errors.New("user input error")
	ErrFileOperation     = errors.New("file operation error")
)

// WizardError provides detailed error information with recovery suggestions
type WizardError struct {
	Type       error
	Message    string
	Details    string
	Suggestion string
	Err        error
}

func (e *WizardError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func (e *WizardError) Unwrap() error {
	return e.Err
}

// NewWizardError creates a detailed error with recovery suggestions
func NewWizardError(errType error, message, details, suggestion string, err error) *WizardError {
	return &WizardError{
		Type:       errType,
		Message:    message,
		Details:    details,
		Suggestion: suggestion,
		Err:        err,
	}
}

// HandleError provides user-friendly error display and recovery suggestions
func HandleError(err error, logger *log.Logger) {
	if err == nil {
		return
	}

	// Check if it's a WizardError with details
	var wizErr *WizardError
	if errors.As(err, &wizErr) {
		// Display structured error information
		fmt.Println()
		fmt.Println(errorStyle.Render("❌ Error: " + wizErr.Message))

		if wizErr.Details != "" {
			fmt.Println(infoStyle.Render("Details: " + wizErr.Details))
		}

		if wizErr.Suggestion != "" {
			suggestStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("220")).
				Bold(true)
			fmt.Println(suggestStyle.Render("💡 Suggestion: " + wizErr.Suggestion))
		}

		// Log the full error for debugging
		if logger != nil {
			logger.Error("Wizard error",
				"type", wizErr.Type,
				"message", wizErr.Message,
				"details", wizErr.Details,
				"original", wizErr.Err)
		}
	} else {
		// Generic error handling
		fmt.Println()
		fmt.Println(errorStyle.Render("❌ Error: " + err.Error()))

		// Provide generic suggestions based on error content
		suggestion := getSuggestionForError(err)
		if suggestion != "" {
			suggestStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("220")).
				Bold(true)
			fmt.Println(suggestStyle.Render("💡 Suggestion: " + suggestion))
		}

		if logger != nil {
			logger.Error("Unexpected error", "error", err)
		}
	}
}

// getSuggestionForError provides suggestions for common errors
func getSuggestionForError(err error) string {
	errStr := strings.ToLower(err.Error())

	switch {
	case strings.Contains(errStr, "permission"):
		return "Try running with appropriate permissions or check file ownership"
	case strings.Contains(errStr, "not found"):
		return "Make sure you're in a Go project directory with go.mod"
	case strings.Contains(errStr, "already exists"):
		return "Use --force to overwrite or check existing configuration"
	case strings.Contains(errStr, "template"):
		return "This might be a bug. Please report it at https://github.com/LarsArtmann/template-GoReleaser/issues"
	case strings.Contains(errStr, "invalid"):
		return "Check your input and try again with valid values"
	case strings.Contains(errStr, "connection"):
		return "Check your internet connection and try again"
	default:
		return ""
	}
}

// RecoverFromPanic provides graceful panic recovery
func RecoverFromPanic(logger *log.Logger) {
	if r := recover(); r != nil {
		fmt.Println()
		fmt.Println(errorStyle.Render("💥 Unexpected error occurred!"))
		fmt.Println(infoStyle.Render("The wizard encountered an unexpected problem and had to stop."))

		suggestStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("220")).
			Bold(true)
		fmt.Println(suggestStyle.Render("💡 Please report this issue at:"))
		fmt.Println("   https://github.com/LarsArtmann/template-GoReleaser/issues")
		fmt.Println()
		fmt.Println("Include this information:")
		fmt.Printf("   Error: %v\n", r)

		if logger != nil {
			logger.Fatal("Panic recovered", "panic", r)
		}

		os.Exit(1)
	}
}

// ValidateFilePermissions checks if we can write to a directory
func ValidateFilePermissions(path string) error {
	// Check if directory exists and is writable
	if info, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			// Try to create the directory
			if err := os.MkdirAll(path, 0755); err != nil {
				return NewWizardError(
					ErrPermission,
					"Cannot create directory",
					fmt.Sprintf("Failed to create %s", path),
					"Check that you have write permissions in the parent directory",
					err,
				)
			}
		} else {
			return NewWizardError(
				ErrFileRead,
				"Cannot access directory",
				fmt.Sprintf("Failed to access %s", path),
				"Check that the path exists and you have read permissions",
				err,
			)
		}
	} else if !info.IsDir() {
		return NewWizardError(
			ErrInvalidInput,
			"Path is not a directory",
			fmt.Sprintf("%s exists but is not a directory", path),
			"Please specify a valid directory path",
			nil,
		)
	}

	// Test write permissions by creating a temporary file
	testFile := fmt.Sprintf("%s/.wizard_test_%d", path, os.Getpid())
	if f, err := os.Create(testFile); err != nil {
		return NewWizardError(
			ErrPermission,
			"No write permission",
			fmt.Sprintf("Cannot write to %s", path),
			"Check that you have write permissions in this directory",
			err,
		)
	} else {
		f.Close()
		os.Remove(testFile)
	}

	return nil
}

// SafeFileWrite writes a file with proper error handling and recovery
func SafeFileWrite(path string, content []byte, perm os.FileMode) error {
	// Create backup if file exists
	if fileExists(path) {
		backupPath := path + ".backup"
		if data, err := os.ReadFile(path); err == nil {
			if err := os.WriteFile(backupPath, data, perm); err != nil {
				log.Debug("Failed to create backup", "file", path, "error", err)
			}
		}
	}

	// Write the file
	if err := os.WriteFile(path, content, perm); err != nil {
		return NewWizardError(
			ErrFileWrite,
			"Failed to write file",
			fmt.Sprintf("Could not write to %s", path),
			"Check file permissions and disk space",
			err,
		)
	}

	return nil
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// HandlePanic provides standardized panic recovery
func HandlePanic(context string, logger *log.Logger) {
	if r := recover(); r != nil {
		RecoverFromPanic(logger)
	}
}

// CheckFileExists checks if a file exists and is accessible
func CheckFileExists(path string, requireDir bool) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return NewWizardError(
				ErrProjectNotFound,
				"File not found",
				fmt.Sprintf("File %s does not exist", path),
				"Check the file path and try again",
				err,
			)
		}
		return NewWizardError(
			ErrFileRead,
			"File access error",
			fmt.Sprintf("Cannot access %s", path),
			"Check file permissions and path",
			err,
		)
	}

	if requireDir && !info.IsDir() {
		return NewWizardError(
			ErrInvalidInput,
			"Expected directory",
			fmt.Sprintf("%s is not a directory", path),
			"Provide a valid directory path",
			nil,
		)
	}

	return nil
}

// LogAndDisplayError logs and displays an error with proper formatting
func LogAndDisplayError(err error, logger *log.Logger) {
	if err != nil {
		HandleError(err, logger)
	}
}

// UserInputError creates a user input validation error
func UserInputError(field string, err error) *WizardError {
	return NewWizardError(
		ErrInvalidInput,
		fmt.Sprintf("Invalid %s", field),
		err.Error(),
		"Please provide valid input and try again",
		err,
	)
}

// TemplateError creates a template processing error
func TemplateError(templateName string, err error) *WizardError {
	return NewWizardError(
		ErrTemplateExecution,
		fmt.Sprintf("Template error in %s", templateName),
		err.Error(),
		"This may be a bug, please report it",
		err,
	)
}

// SafeReadFile reads a file safely with error handling
func SafeReadFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, NewWizardError(
			ErrFileRead,
			"Failed to read file",
			fmt.Sprintf("Could not read %s", path),
			"Check file exists and permissions",
			err,
		)
	}
	return data, nil
}

// SafeCreateFile creates a file safely with error handling
func SafeCreateFile(path string) (*os.File, error) {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, NewWizardError(
			ErrPermission,
			"Cannot create directory",
			fmt.Sprintf("Failed to create directory %s", dir),
			"Check directory permissions",
			err,
		)
	}

	file, err := os.Create(path)
	if err != nil {
		return nil, NewWizardError(
			ErrFileWrite,
			"Cannot create file",
			fmt.Sprintf("Failed to create %s", path),
			"Check file permissions and disk space",
			err,
		)
	}
	return file, nil
}

// WrapFileError wraps file operation errors
func WrapFileError(operation, path string, err error) *WizardError {
	return NewWizardError(
		ErrFileWrite,
		fmt.Sprintf("Failed to %s", operation),
		fmt.Sprintf("Error with %s: %s", path, err.Error()),
		"Check file permissions and disk space",
		err,
	)
}
