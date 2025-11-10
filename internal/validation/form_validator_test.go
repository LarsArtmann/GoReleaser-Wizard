package validation

import (
	"testing"
)

func TestFormValidator(t *testing.T) {
	fv := NewFormValidator()

	// Test initial state
	if fv.HasErrors() {
		t.Errorf("NewFormValidator() should not have errors initially")
	}

	if len(fv.GetErrors()) != 0 {
		t.Errorf("NewFormValidator() should return empty errors map")
	}
}

func TestFormValidatorValidateProjectName(t *testing.T) {
	fv := NewFormValidator()

	// Test valid name
	err := fv.ValidateProjectName()("myproject")
	if err != nil {
		t.Errorf("ValidateProjectName() valid input error = %v", err)
	}

	if fv.HasErrors() {
		t.Errorf("ValidateProjectName() should not have errors for valid input")
	}

	// Test invalid name
	err = fv.ValidateProjectName()("con")
	if err == nil {
		t.Errorf("ValidateProjectName() should error for reserved name")
	}

	if !fv.HasErrors() {
		t.Errorf("ValidateProjectName() should have errors for invalid input")
	}

	// Check specific error
	if fv.GetFieldError("project_name") == "" {
		t.Errorf("ValidateProjectName() should set project_name error")
	}

	// Clear errors and test
	fv.ClearErrors()
	if fv.HasErrors() {
		t.Errorf("ClearErrors() should clear all errors")
	}
}

func TestFormValidatorValidateBinaryName(t *testing.T) {
	fv := NewFormValidator()

	// Test valid binary name
	err := fv.ValidateBinaryName()("myapp")
	if err != nil {
		t.Errorf("ValidateBinaryName() valid input error = %v", err)
	}

	// Test invalid binary name with dangerous chars
	err = fv.ValidateBinaryName()("my;app")
	if err == nil {
		t.Errorf("ValidateBinaryName() should error for dangerous characters")
	}

	if fv.GetFieldError("binary_name") == "" {
		t.Errorf("ValidateBinaryName() should set binary_name error")
	}
}

func TestFormValidatorValidateMainPath(t *testing.T) {
	fv := NewFormValidator()

	// Test valid path
	err := fv.ValidateMainPath()("./cmd/app")
	if err != nil {
		t.Errorf("ValidateMainPath() valid input error = %v", err)
	}

	// Test path traversal
	err = fv.ValidateMainPath()("../../../etc/passwd")
	if err == nil {
		t.Errorf("ValidateMainPath() should error for path traversal")
	}

	if fv.GetFieldError("main_path") == "" {
		t.Errorf("ValidateMainPath() should set main_path error")
	}
}

func TestFormValidatorValidateProjectDescription(t *testing.T) {
	fv := NewFormValidator()

	// Test valid description
	err := fv.ValidateProjectDescription()("A great app")
	if err != nil {
		t.Errorf("ValidateProjectDescription() valid input error = %v", err)
	}

	// Test script injection
	err = fv.ValidateProjectDescription()("<script>alert('xss')</script>")
	if err == nil {
		t.Errorf("ValidateProjectDescription() should error for script injection")
	}

	if fv.GetFieldError("project_description") == "" {
		t.Errorf("ValidateProjectDescription() should set project_description error")
	}
}

func TestFormValidatorValidateDockerRegistry(t *testing.T) {
	fv := NewFormValidator()

	// Test valid registry
	err := fv.ValidateDockerRegistry()("ghcr.io/username/app")
	if err != nil {
		t.Errorf("ValidateDockerRegistry() valid input error = %v", err)
	}

	// Test insecure HTTP (non-localhost)
	err = fv.ValidateDockerRegistry()("http://registry.example.com/app")
	if err == nil {
		t.Errorf("ValidateDockerRegistry() should error for insecure HTTP")
	}

	if fv.GetFieldError("docker_registry") == "" {
		t.Errorf("ValidateDockerRegistry() should set docker_registry error")
	}
}

func TestFormValidatorValidateBuildTags(t *testing.T) {
	fv := NewFormValidator()

	// Test valid tags
	err := fv.ValidateBuildTags([]string{"prod", "linux"})
	if err != nil {
		t.Errorf("ValidateBuildTags() valid input error = %v", err)
	}

	// Test invalid tags
	err = fv.ValidateBuildTags([]string{"prod;rm"})
	if err == nil {
		t.Errorf("ValidateBuildTags() should error for dangerous tags")
	}

	if fv.GetFieldError("build_tags") == "" {
		t.Errorf("ValidateBuildTags() should set build_tags error")
	}
}

func TestFormValidatorErrorSummary(t *testing.T) {
	fv := NewFormValidator()

	// Empty summary
	summary := fv.GetErrorSummary()
	if summary != "" {
		t.Errorf("GetErrorSummary() should return empty string when no errors")
	}

	// Add some errors
	_ = fv.ValidateProjectName()("invalid..name")
	_ = fv.ValidateBinaryName()("con")

	summary = fv.GetErrorSummary()
	if summary == "" {
		t.Errorf("GetErrorSummary() should return non-empty string when errors exist")
	}

	if !contains(summary, "project_name") {
		t.Errorf("GetErrorSummary() should contain project_name error")
	}

	if !contains(summary, "binary_name") {
		t.Errorf("GetErrorSummary() should contain binary_name error")
	}
}

func TestFormValidatorSanitizeAndValidate(t *testing.T) {
	fv := NewFormValidator()

	// Test valid input
	sanitized, err := fv.SanitizeAndValidate("  myproject  ", fv.ValidateProjectName())
	if err != nil {
		t.Errorf("SanitizeAndValidate() error = %v", err)
	}

	if sanitized != "myproject" {
		t.Errorf("SanitizeAndValidate() = %v, want %v", sanitized, "myproject")
	}

	// Test invalid input
	_, err = fv.SanitizeAndValidate("con", fv.ValidateProjectName())
	if err == nil {
		t.Errorf("SanitizeAndValidate() should error for invalid input")
	}
}

func TestFormValidatorValidateRequired(t *testing.T) {
	fv := NewFormValidator()

	validator := fv.ValidateRequired("Test Field")

	// Test empty input
	err := validator("")
	if err == nil {
		t.Errorf("ValidateRequired() should error for empty input")
	}

	// Test valid input
	err = validator("test")
	if err != nil {
		t.Errorf("ValidateRequired() should not error for valid input")
	}
}

func TestFormValidatorValidateLength(t *testing.T) {
	fv := NewFormValidator()

	validator := fv.ValidateLength(3, 10, "Test Field")

	// Test too short
	err := validator("ab")
	if err == nil {
		t.Errorf("ValidateLength() should error for too short input")
	}

	// Test too long
	err = validator("abcdefghijk")
	if err == nil {
		t.Errorf("ValidateLength() should error for too long input")
	}

	// Test valid length
	err = validator("abcde")
	if err != nil {
		t.Errorf("ValidateLength() should not error for valid length input")
	}
}

func TestFormValidatorValidateNoShellMetacharacters(t *testing.T) {
	fv := NewFormValidator()

	validator := fv.ValidateNoShellMetacharacters("Test Field")

	// Test dangerous characters
	err := validator("test;rm")
	if err == nil {
		t.Errorf("ValidateNoShellMetacharacters() should error for dangerous characters")
	}

	// Test safe input
	err = validator("test-safe")
	if err != nil {
		t.Errorf("ValidateNoShellMetacharacters() should not error for safe input")
	}
}

func TestFormValidatorValidateNoPathTraversal(t *testing.T) {
	fv := NewFormValidator()

	validator := fv.ValidateNoPathTraversal("Test Field")

	// Test path traversal
	err := validator("../../../etc")
	if err == nil {
		t.Errorf("ValidateNoPathTraversal() should error for path traversal")
	}

	// Test safe path
	err = validator("./cmd/app")
	if err != nil {
		t.Errorf("ValidateNoPathTraversal() should not error for safe path")
	}
}

func TestFormValidatorSetFieldError(t *testing.T) {
	fv := NewFormValidator()

	fv.SetFieldError("test_field", "test error message")

	if !fv.HasErrors() {
		t.Errorf("SetFieldError() should result in HasErrors() returning true")
	}

	if fv.GetFieldError("test_field") != "test error message" {
		t.Errorf("GetFieldError() = %v, want %v", fv.GetFieldError("test_field"), "test error message")
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			(len(s) > len(substr) &&
				findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
