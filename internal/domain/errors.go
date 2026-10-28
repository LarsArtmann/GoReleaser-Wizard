package domain

import (
	"errors"
	"fmt"
)

// DomainError represents all domain-specific errors
// Provides type safety and detailed error information.
type DomainError struct {
	Code       ErrorCode `json:"code"`
	Message    string    `json:"message"`
	Details    string    `json:"details,omitempty"`
	Context    string    `json:"context,omitempty"`
	Field      string    `json:"field,omitempty"`
	Suggestion string    `json:"suggestion,omitempty"`
	Cause      error     `json:"cause,omitempty"`
}

// ErrorCode represents all possible error codes in the domain.
type ErrorCode string

const (
	// Validation Errors.

	ErrInvalidProjectName        ErrorCode = "INVALID_PROJECT_NAME"
	ErrInvalidBinaryName         ErrorCode = "INVALID_BINARY_NAME"
	ErrInvalidMainPath           ErrorCode = "INVALID_MAIN_PATH"
	ErrInvalidProjectDescription ErrorCode = "INVALID_PROJECT_DESCRIPTION"
	ErrInvalidPlatform           ErrorCode = "INVALID_PLATFORM"
	ErrInvalidArchitecture       ErrorCode = "INVALID_ARCHITECTURE"
	ErrInvalidGitProvider        ErrorCode = "INVALID_GIT_PROVIDER"
	ErrInvalidDockerRegistry     ErrorCode = "INVALID_DOCKER_REGISTRY"
	ErrInvalidActionTrigger      ErrorCode = "INVALID_ACTION_TRIGGER"
	ErrInvalidBuildTag           ErrorCode = "INVALID_BUILD_TAG"
	ErrInvalidConfigState        ErrorCode = "INVALID_CONFIG_STATE"

	// Configuration Errors.

	ErrDockerNotSupported     ErrorCode = "DOCKER_NOT_SUPPORTED"
	ErrPlatformArchMismatch   ErrorCode = "PLATFORM_ARCH_MISMATCH"
	ErrMainPathRequired       ErrorCode = "MAIN_PATH_REQUIRED"
	ErrInvalidStateTransition ErrorCode = "INVALID_STATE_TRANSITION"
	ErrMissingRequiredField   ErrorCode = "MISSING_REQUIRED_FIELD"
	ErrFieldTooLong           ErrorCode = "FIELD_TOO_LONG"
	ErrFieldTooShort          ErrorCode = "FIELD_TOO_SHORT"

	// Business Rule Errors.

	ErrDuplicateBuildTag ErrorCode = "DUPLICATE_BUILD_TAG"
	ErrTooManyBuildTags  ErrorCode = "TOO_MANY_BUILD_TAGS"
	ErrInvalidURLPattern ErrorCode = "INVALID_URL_PATTERN"
	ErrReservedName      ErrorCode = "RESERVED_NAME"
	ErrInvalidCharacters ErrorCode = "INVALID_CHARACTERS"

	// System Errors.

	ErrFileNotFound          ErrorCode = "FILE_NOT_FOUND"
	ErrPermissionDenied      ErrorCode = "PERMISSION_DENIED"
	ErrFileWriteFailed       ErrorCode = "FILE_WRITE_FAILED"
	ErrFileReadFailed        ErrorCode = "FILE_READ_FAILED"
	ErrDirectoryCreateFailed ErrorCode = "DIRECTORY_CREATE_FAILED"
	ErrDependencyNotFound    ErrorCode = "DEPENDENCY_NOT_FOUND"

	// Template Errors.

	ErrTemplateNotFound        ErrorCode = "TEMPLATE_NOT_FOUND"
	ErrTemplateExecutionFailed ErrorCode = "TEMPLATE_EXECUTION_FAILED"
	ErrTemplateSyntaxError     ErrorCode = "TEMPLATE_SYNTAX_ERROR"

	// Configuration Errors.

	ErrInvalidConfiguration    ErrorCode = "INVALID_CONFIGURATION"
	ErrInvalidFileFormat       ErrorCode = "INVALID_FILE_FORMAT"
	ErrInvalidProjectStructure ErrorCode = "INVALID_PROJECT_STRUCTURE"
	ErrMissingDependency       ErrorCode = "MISSING_DEPENDENCY"
	ErrExternalToolNotFound    ErrorCode = "EXTERNAL_TOOL_NOT_FOUND"

	// Application Errors.

	ErrInvalidProject    ErrorCode = "INVALID_PROJECT"
	ErrInvalidBinary     ErrorCode = "INVALID_BINARY"
	ErrInvalidConfig     ErrorCode = "INVALID_CONFIG"
	ErrInvalidTemplate   ErrorCode = "INVALID_TEMPLATE"
	ErrInvalidField      ErrorCode = "INVALID_FIELD"
	ErrInvalidOperation  ErrorCode = "INVALID_OPERATION"
	ErrConfigGeneration  ErrorCode = "CONFIG_GENERATION"
	ErrDependencyMissing ErrorCode = "DEPENDENCY_MISSING"

	// Template Errors.

	ErrTemplateParsing   ErrorCode = "TEMPLATE_PARSING"
	ErrTemplateRendering ErrorCode = "TEMPLATE_RENDERING"

	// File Operation Errors.

	ErrFileOperation  ErrorCode = "FILE_OPERATION"
	ErrFilePermission ErrorCode = "FILE_PERMISSION"

	// Git Errors.

	ErrGitOperation ErrorCode = "GIT_OPERATION"
	ErrGitNotFound  ErrorCode = "GIT_NOT_FOUND"
	ErrGitCommand   ErrorCode = "GIT_COMMAND"

	// Extended Validation Errors.

	ErrInvalidDockerImage ErrorCode = "INVALID_DOCKER_IMAGE"
	ErrInvalidVersion     ErrorCode = "INVALID_VERSION"
	ErrInvalidGitBranch   ErrorCode = "INVALID_GIT_BRANCH"
	ErrInvalidGitTag      ErrorCode = "INVALID_GIT_TAG"
	ErrInvalidPort        ErrorCode = "INVALID_PORT"

	// External Service Errors.

	ErrGitOperationFailed   ErrorCode = "GIT_OPERATION_FAILED"
	ErrRegistryAccessDenied ErrorCode = "REGISTRY_ACCESS_DENIED"
	ErrGitHubAPIError       ErrorCode = "GITHUB_API_ERROR"
)

// Error implements the error interface.
func (de *DomainError) Error() string {
	if de.Context != "" {
		return fmt.Sprintf("[%s] %s (context: %s)", de.Code, de.Message, de.Context)
	}

	return fmt.Sprintf("[%s] %s", de.Code, de.Message)
}

// Unwrap returns the underlying cause.
func (de *DomainError) Unwrap() error {
	return de.Cause
}

// WithContext adds context to the error.
func (de *DomainError) WithContext(context string) *DomainError {
	return &DomainError{
		Code:       de.Code,
		Message:    de.Message,
		Details:    de.Details,
		Context:    context,
		Field:      de.Field,
		Suggestion: de.Suggestion,
		Cause:      de.Cause,
	}
}

// WithField adds field information to the error.
func (de *DomainError) WithField(field string) *DomainError {
	de.Field = field

	return de
}

// WithSuggestion adds a suggestion to the error.
func (de *DomainError) WithSuggestion(suggestion string) *DomainError {
	de.Suggestion = suggestion

	return de
}

// WithCause adds an underlying cause to the error.
func (de *DomainError) WithCause(cause error) *DomainError {
	de.Cause = cause

	return de
}

// IsErrorCode checks if an error matches a specific error code.
func IsErrorCode(err error, code ErrorCode) bool {
	if domainErr, ok := errors.AsType[*DomainError](err); ok {
		return domainErr.Code == code
	}

	return false
}

// NewValidationError creates a validation error of the given code, message, and details.
func NewValidationError(code ErrorCode, message, details string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

func NewConfigurationError(code ErrorCode, message, details string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

func NewBusinessRuleError(code ErrorCode, message, details string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

func NewSystemError(code ErrorCode, message, details string, cause error) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Details: details,
		Cause:   cause,
	}
}

func NewTemplateError(code ErrorCode, message, details string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

func NewExternalServiceError(code ErrorCode, message, details string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

func NewFileError(code ErrorCode, message, details string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

func NewConfigError(code ErrorCode, message, details string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

func NewGitError(code ErrorCode, message, details string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// InvalidProjectNameError returns a validation error for an invalid project name.
func InvalidProjectNameError(value string) *DomainError {
	return NewValidationError(
		ErrInvalidProjectName,
		"Invalid project name",
		fmt.Sprintf("'%s' does not meet validation requirements", value),
	)
}

func InvalidBinaryNameError(value string) *DomainError {
	return NewValidationError(
		ErrInvalidBinaryName,
		"Invalid binary name",
		fmt.Sprintf("'%s' does not meet validation requirements", value),
	)
}

func InvalidMainPathError(value string) *DomainError {
	return NewValidationError(
		ErrInvalidMainPath,
		"Invalid main path",
		fmt.Sprintf("'%s' does not meet validation requirements", value),
	)
}

func InvalidPlatformError(value string) *DomainError {
	return NewValidationError(
		ErrInvalidPlatform,
		"Invalid platform",
		fmt.Sprintf("'%s' is not a supported platform", value),
	)
}

func InvalidArchitectureError(value string) *DomainError {
	return NewValidationError(
		ErrInvalidArchitecture,
		"Invalid architecture",
		fmt.Sprintf("'%s' is not a supported architecture", value),
	)
}

func DockerNotSupportedError(projectType ProjectType) *DomainError {
	return NewConfigurationError(
		ErrDockerNotSupported,
		"Docker not supported",
		fmt.Sprintf("Project type %s does not support Docker", projectType),
	)
}

func PlatformArchMismatchError(platform Platform, arch Architecture) *DomainError {
	return NewConfigurationError(
		ErrPlatformArchMismatch,
		"Platform-architecture mismatch",
		fmt.Sprintf("Architecture %s is not supported on platform %s", arch, platform),
	)
}

// FileNotFoundError returns a system error for a missing file at the given path.
func FileNotFoundError(path string, cause error) *DomainError {
	return NewSystemError(
		ErrFileNotFound,
		"File not found",
		fmt.Sprintf("File '%s' does not exist", path),
		cause,
	)
}

func PermissionDeniedError(path string, cause error) *DomainError {
	return NewSystemError(
		ErrPermissionDenied,
		"Permission denied",
		fmt.Sprintf("Permission denied for '%s'", path),
		cause,
	)
}

func FileWriteFailedError(path string, cause error) *DomainError {
	return NewSystemError(
		ErrFileWriteFailed,
		"File write failed",
		fmt.Sprintf("Failed to write file '%s'", path),
		cause,
	)
}

// TemplateNotFoundError returns a template error when the given template cannot be located.
func TemplateNotFoundError(template string) *DomainError {
	return NewTemplateError(
		ErrTemplateNotFound,
		"Template not found",
		fmt.Sprintf("Template '%s' not found", template),
	)
}

func TemplateExecutionFailedError(template string, cause error) *DomainError {
	return NewTemplateError(
		ErrTemplateExecutionFailed,
		"Template execution failed",
		fmt.Sprintf("Failed to execute template '%s'", template),
	).WithCause(cause)
}

// GitOperationFailedError returns an external service error when a git operation fails.
func GitOperationFailedError(operation string, cause error) *DomainError {
	return NewExternalServiceError(
		ErrGitOperationFailed,
		"Git operation failed",
		"Failed to perform "+operation,
	).WithCause(cause)
}

func RegistryAccessDeniedError(registry string) *DomainError {
	return NewExternalServiceError(
		ErrRegistryAccessDenied,
		"Registry access denied",
		fmt.Sprintf("Access denied to registry '%s'", registry),
	)
}

// GetRecoverySuggestion returns a user-friendly hint explaining how to fix the error.
func (de *DomainError) GetRecoverySuggestion() string {
	switch de.Code {
	case ErrInvalidProjectName:
		return "Use only letters, numbers, hyphens, and underscores. Must start with a letter and be 1-63 characters."
	case ErrInvalidBinaryName:
		return "Use only letters, numbers, hyphens, and underscores. Must start with a letter and be 1-63 characters. Avoid reserved Windows names."
	case ErrInvalidMainPath:
		return "Use relative path with only valid characters. Avoid parent directory references (..)."
	case ErrDockerNotSupported:
		return "Disable Docker support or choose a project type that supports containers."
	case ErrPlatformArchMismatch:
		return "Select architectures that are compatible with your target platforms."
	case ErrPermissionDenied:
		return "Check file permissions and ensure you have write access to the directory."
	case ErrFileNotFound:
		return "Verify the file exists and the path is correct."
	case ErrTemplateNotFound:
		return "Ensure the template exists and is accessible."
	default:
		return "Check the error details and try again with corrected input."
	}
}

// ErrorSeverity classifies errors for proper handling (e.g. logging, UI presentation).
type ErrorSeverity int

const (
	ErrorSeverityInfo ErrorSeverity = iota
	ErrorSeverityWarning
	ErrorSeverityError
	ErrorSeverityCritical
)

func (de *DomainError) GetSeverity() ErrorSeverity {
	switch de.Code {
	// Validation errors are warnings (user can fix)
	case ErrInvalidProjectName, ErrInvalidBinaryName, ErrInvalidMainPath,
		ErrInvalidPlatform, ErrInvalidArchitecture, ErrInvalidGitProvider:
		return ErrorSeverityWarning

	// Configuration errors are errors (more serious)
	case ErrDockerNotSupported, ErrPlatformArchMismatch, ErrInvalidStateTransition:
		return ErrorSeverityError

	// System errors are critical
	case ErrPermissionDenied, ErrFileNotFound:
		return ErrorSeverityCritical

	// Default to error
	default:
		return ErrorSeverityError
	}
}

// IsRecoverable returns true if the error can be recovered from by the user.
func (de *DomainError) IsRecoverable() bool {
	severity := de.GetSeverity()

	return severity <= ErrorSeverityError
}

// IsRetryable returns true if the operation can be retried.
func (de *DomainError) IsRetryable() bool {
	switch de.Code {
	case ErrFileWriteFailed, ErrGitOperationFailed, ErrRegistryAccessDenied:
		return true
	default:
		return false
	}
}
