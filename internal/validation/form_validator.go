package validation

import (
	"fmt"
	"strings"
)

// FormValidator provides validation functions for huh forms
type FormValidator struct {
	errors map[string]string
}

// NewFormValidator creates a new form validator
func NewFormValidator() *FormValidator {
	return &FormValidator{
		errors: make(map[string]string),
	}
}

// ValidateProjectName creates a huh-compatible validator for project names
func (fv *FormValidator) ValidateProjectName() func(string) error {
	return func(value string) error {
		if err := ValidateProjectName(value); err != nil {
			fv.errors["project_name"] = err.Error()
			return err
		}
		delete(fv.errors, "project_name")
		return nil
	}
}

// ValidateBinaryName creates a huh-compatible validator for binary names
func (fv *FormValidator) ValidateBinaryName() func(string) error {
	return func(value string) error {
		if err := ValidateBinaryName(value); err != nil {
			fv.errors["binary_name"] = err.Error()
			return err
		}
		delete(fv.errors, "binary_name")
		return nil
	}
}

// ValidateMainPath creates a huh-compatible validator for main path
func (fv *FormValidator) ValidateMainPath() func(string) error {
	return func(value string) error {
		if err := ValidateMainPath(value); err != nil {
			fv.errors["main_path"] = err.Error()
			return err
		}
		delete(fv.errors, "main_path")
		return nil
	}
}

// ValidateProjectDescription creates a huh-compatible validator for project description
func (fv *FormValidator) ValidateProjectDescription() func(string) error {
	return func(value string) error {
		if err := ValidateProjectDescription(value); err != nil {
			fv.errors["project_description"] = err.Error()
			return err
		}
		delete(fv.errors, "project_description")
		return nil
	}
}

// ValidateDockerRegistry creates a huh-compatible validator for Docker registry
func (fv *FormValidator) ValidateDockerRegistry() func(string) error {
	return func(value string) error {
		if err := ValidateDockerRegistry(value); err != nil {
			fv.errors["docker_registry"] = err.Error()
			return err
		}
		delete(fv.errors, "docker_registry")
		return nil
	}
}

// ValidateBuildTags creates a huh-compatible validator for build tags
func (fv *FormValidator) ValidateBuildTags(tags []string) error {
	if err := ValidateBuildTags(tags); err != nil {
		fv.errors["build_tags"] = err.Error()
		return err
	}
	delete(fv.errors, "build_tags")
	return nil
}

// GetErrors returns all current validation errors
func (fv *FormValidator) GetErrors() map[string]string {
	return fv.errors
}

// HasErrors returns true if there are validation errors
func (fv *FormValidator) HasErrors() bool {
	return len(fv.errors) > 0
}

// GetErrorSummary returns a formatted summary of all errors
func (fv *FormValidator) GetErrorSummary() string {
	if !fv.HasErrors() {
		return ""
	}

	var errors []string
	for field, message := range fv.errors {
		errors = append(errors, fmt.Sprintf("• %s: %s", field, message))
	}

	return "Validation errors:\n" + strings.Join(errors, "\n")
}

// ClearErrors clears all current validation errors
func (fv *FormValidator) ClearErrors() {
	fv.errors = make(map[string]string)
}

// SanitizeAndValidate sanitizes input and validates it
func (fv *FormValidator) SanitizeAndValidate(input string, validator func(string) error) (string, error) {
	// First sanitize the input
	sanitized := SanitizeInput(input)

	// Then validate it
	if err := validator(sanitized); err != nil {
		return "", err
	}

	return sanitized, nil
}

// GetFieldError gets a specific field error
func (fv *FormValidator) GetFieldError(field string) string {
	return fv.errors[field]
}

// SetFieldError manually sets a field error
func (fv *FormValidator) SetFieldError(field, message string) {
	fv.errors[field] = message
}

// ValidateRequired validates that a field is not empty
func (fv *FormValidator) ValidateRequired(fieldName string) func(string) error {
	return func(value string) error {
		if strings.TrimSpace(value) == "" {
			err := fmt.Errorf("%s is required", fieldName)
			fv.errors[fieldName] = err.Error()
			return err
		}
		delete(fv.errors, fieldName)
		return nil
	}
}

// ValidateLength validates string length within bounds
func (fv *FormValidator) ValidateLength(min, max int, fieldName string) func(string) error {
	return func(value string) error {
		length := len(strings.TrimSpace(value))
		if length < min || length > max {
			err := fmt.Errorf("%s must be between %d and %d characters", fieldName, min, max)
			fv.errors[fieldName] = err.Error()
			return err
		}
		delete(fv.errors, fieldName)
		return nil
	}
}

// ValidateNoShellMetacharacters validates that input doesn't contain shell metacharacters
func (fv *FormValidator) ValidateNoShellMetacharacters(fieldName string) func(string) error {
	return func(value string) error {
		if shellMetacharPattern.MatchString(value) {
			err := fmt.Errorf("%s contains dangerous characters", fieldName)
			fv.errors[fieldName] = err.Error()
			return err
		}
		delete(fv.errors, fieldName)
		return nil
	}
}

// ValidateNoPathTraversal validates that input doesn't contain path traversal
func (fv *FormValidator) ValidateNoPathTraversal(fieldName string) func(string) error {
	return func(value string) error {
		if pathTraversalPattern.MatchString(value) {
			err := fmt.Errorf("%s contains path traversal attempts", fieldName)
			fv.errors[fieldName] = err.Error()
			return err
		}
		delete(fv.errors, fieldName)
		return nil
	}
}
