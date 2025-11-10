package validation

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

// Security validation patterns
var (
	// Project name pattern: alphanumeric, hyphens, underscores, dots, starts and ends with alphanumeric
	projectNamePattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9\-._]*[a-zA-Z0-9]$`)

	// Binary name pattern: more restrictive, no special characters that could cause issues
	binaryNamePattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9\-_]*[a-zA-Z0-9]$`)

	// Build tag pattern (Go build tag syntax)
	buildTagPattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`)

	// Docker registry pattern - more restrictive to prevent double dots
	dockerRegistryPattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9\-.]*[a-zA-Z0-9](\.[a-zA-Z][a-zA-Z0-9\-.]*[a-zA-Z0-9])*(:[0-9]+)?(/[a-zA-Z0-9][a-zA-Z0-9\-.]*[a-zA-Z0-9])*$`)

	// Path traversal detection
	pathTraversalPattern = regexp.MustCompile(`\.\.[/\\]`)

	// Shell metacharacters that could be dangerous
	shellMetacharPattern = regexp.MustCompile(`[;&|<>"'$` + "`" + `\\]`)

	// Reserved names (OS-specific and Go-specific)
	reservedNames = map[string]bool{
		// Windows reserved
		"con": true, "prn": true, "aux": true, "nul": true,
		"com1": true, "com2": true, "com3": true, "com4": true,
		"com5": true, "com6": true, "com7": true, "com8": true,
		"com9": true, "lpt1": true, "lpt2": true, "lpt3": true,
		"lpt4": true, "lpt5": true, "lpt6": true, "lpt7": true,
		"lpt8": true, "lpt9": true,

		// Go/Build system reserved
		"go": true, "test": true, "vendor": true, "internal": true,
		"main": true, "init": true, "close": true, "copy": true,

		// Unix special files
		"etc": true, "usr": true, "var": true, "bin": true, "sbin": true,
		"lib": true, "lib64": true, "dev": true, "proc": true, "sys": true,
		"root": true, "home": true, "tmp": true, "opt": true, "srv": true,
		"mnt": true, "media": true, "run": true,
	}

	// Dangerous file extensions
	dangerousExtensions = map[string]bool{
		".exe": true, ".bat": true, ".cmd": true, ".com": true, ".pif": true,
		".scr": true, ".vbs": true, ".js": true, ".jar": true, ".sh": true,
		".ps1": true, ".py": true, ".rb": true, ".pl": true, ".php": true,
	}
)

// ValidationError represents a validation error with context
type ValidationError struct {
	Field   string
	Value   string
	Message string
	Code    string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation failed for field '%s': %s", e.Field, e.Message)
}

// ValidationResult contains validation results
type ValidationResult struct {
	Valid  bool
	Errors []ValidationError
}

// AddError adds a validation error to the result
func (vr *ValidationResult) AddError(field, value, message, code string) {
	vr.Valid = false
	vr.Errors = append(vr.Errors, ValidationError{
		Field:   field,
		Value:   value,
		Message: message,
		Code:    code,
	})
}

// ValidateProjectName validates project name according to security rules
func ValidateProjectName(name string) error {
	if name == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	if len(name) < 1 || len(name) > 63 {
		return fmt.Errorf("project name must be 1-63 characters long")
	}

	// Check for consecutive special characters
	if strings.Contains(name, "--") || strings.Contains(name, "__") ||
		strings.Contains(name, "..") || strings.Contains(name, "__") {
		return fmt.Errorf("project name cannot contain consecutive special characters")
	}

	// Validate pattern
	if !projectNamePattern.MatchString(name) {
		return fmt.Errorf("project name contains invalid characters. Use letters, numbers, hyphens, underscores, and dots")
	}

	// Check for reserved names (case-insensitive)
	lowerName := strings.ToLower(name)
	if reservedNames[lowerName] {
		return fmt.Errorf("'%s' is a reserved name and cannot be used", name)
	}

	// Check for dangerous patterns
	if strings.HasPrefix(name, ".") || strings.HasPrefix(name, "-") {
		return fmt.Errorf("project name cannot start with special characters")
	}

	if strings.HasSuffix(name, ".") || strings.HasSuffix(name, "-") {
		return fmt.Errorf("project name cannot end with special characters")
	}

	return nil
}

// ValidateBinaryName validates binary name with strict security rules
func ValidateBinaryName(name string) error {
	if name == "" {
		return fmt.Errorf("binary name cannot be empty")
	}

	if len(name) < 1 || len(name) > 255 {
		return fmt.Errorf("binary name must be 1-255 characters long")
	}

	// No spaces or shell metacharacters
	if shellMetacharPattern.MatchString(name) {
		return fmt.Errorf("binary name contains dangerous characters")
	}

	// Must be valid filename on all platforms
	if strings.ContainsAny(name, `<>:"/\|?*`) {
		return fmt.Errorf("binary name contains invalid filename characters")
	}

	// Validate pattern
	if !binaryNamePattern.MatchString(name) {
		return fmt.Errorf("binary name must start with a letter and contain only letters, numbers, hyphens, and underscores")
	}

	// Check for reserved names (case-insensitive)
	lowerName := strings.ToLower(name)
	if reservedNames[lowerName] {
		return fmt.Errorf("'%s' is a reserved name and cannot be used", name)
	}

	// Check for dangerous extensions
	ext := strings.ToLower(filepath.Ext(name))
	if dangerousExtensions[ext] {
		return fmt.Errorf("binary name has dangerous extension: %s", ext)
	}

	return nil
}

// ValidateMainPath validates the main package path with security checks
func ValidateMainPath(path string) error {
	if path == "" {
		return fmt.Errorf("main path cannot be empty")
	}

	// No path traversal attacks
	if pathTraversalPattern.MatchString(path) {
		return fmt.Errorf("path traversal not allowed")
	}

	// No absolute paths
	if filepath.IsAbs(path) {
		return fmt.Errorf("absolute paths not allowed")
	}

	// Clean the path
	cleanPath := filepath.Clean(path)

	// Check for dangerous patterns
	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("path traversal not allowed")
	}

	// Check for shell metacharacters
	if shellMetacharPattern.MatchString(cleanPath) {
		return fmt.Errorf("path contains dangerous characters")
	}

	// Validate path doesn't point to system directories
	parts := strings.SplitSeq(cleanPath, string(filepath.Separator))
	for part := range parts {
		if part != "" {
			lowerPart := strings.ToLower(part)
			if reservedNames[lowerPart] {
				return fmt.Errorf("path contains reserved directory name: %s", part)
			}
		}
	}

	return nil
}

// ValidateProjectDescription validates project description
func ValidateProjectDescription(desc string) error {
	if len(desc) > 255 {
		return fmt.Errorf("project description must be 255 characters or less")
	}

	// Check for HTML/markdown injection attempts
	if strings.Contains(desc, "<script") || strings.Contains(desc, "javascript:") {
		return fmt.Errorf("description contains potentially dangerous content")
	}

	// Check for excessive length that might indicate injection
	if len(strings.TrimSpace(desc)) == 0 && len(desc) > 100 {
		return fmt.Errorf("description contains suspicious content")
	}

	return nil
}

// ValidateBuildTags validates Go build tags
func ValidateBuildTags(tags []string) error {
	for _, tag := range tags {
		if tag == "" {
			return fmt.Errorf("build tag cannot be empty")
		}

		if len(tag) > 50 {
			return fmt.Errorf("build tag '%s' is too long (max 50 characters)", tag)
		}

		if !buildTagPattern.MatchString(tag) {
			return fmt.Errorf("build tag '%s' contains invalid characters", tag)
		}

		// Check for dangerous patterns
		if strings.Contains(tag, "..") || strings.Contains(tag, "/") {
			return fmt.Errorf("build tag '%s' contains dangerous patterns", tag)
		}
	}

	return nil
}

// ValidateDockerRegistry validates Docker registry URL
func ValidateDockerRegistry(registry string) error {
	if registry == "" {
		return fmt.Errorf("docker registry cannot be empty")
	}

	registry = strings.TrimSpace(registry)

	// Check for invalid patterns like double dots
	if strings.Contains(registry, "..") {
		return fmt.Errorf("invalid docker registry format")
	}

	// Basic URL pattern check
	if !dockerRegistryPattern.MatchString(registry) {
		return fmt.Errorf("invalid docker registry format")
	}

	// No credentials in URL
	if strings.Contains(registry, "@") || strings.Contains(registry, ":") {
		// Allow port but not credentials
		parts := strings.Split(registry, ":")
		if len(parts) > 2 {
			return fmt.Errorf("docker registry should not contain credentials")
		}
	}

	// No dangerous protocols
	if strings.HasPrefix(registry, "http://") {
		// Allow HTTP only for localhost
		if !strings.Contains(registry, "localhost") && !strings.Contains(registry, "127.0.0.1") {
			return fmt.Errorf("insecure HTTP registry only allowed for localhost")
		}
	}

	return nil
}

// ValidateGitProvider validates git provider selection
func ValidateGitProvider(provider string) error {
	validProviders := []string{
		"github", "gitlab", "bitbucket", "gitea", "self-hosted", "local",
	}

	for _, valid := range validProviders {
		if strings.EqualFold(provider, valid) {
			return nil
		}
	}

	return fmt.Errorf("invalid git provider: %s", provider)
}

// SanitizeInput sanitizes user input to prevent injection
func SanitizeInput(input string) string {
	// Remove null bytes
	input = strings.ReplaceAll(input, "\x00", "")

	// Remove control characters except newlines, tabs, and carriage returns
	var result strings.Builder
	for _, r := range input {
		if unicode.IsControl(r) && r != '\n' && r != '\t' && r != '\r' {
			continue
		}
		result.WriteRune(r)
	}

	// Custom trim: remove leading/trailing spaces and tabs but preserve newlines
	resultStr := result.String()
	// Trim leading spaces and tabs
	for len(resultStr) > 0 && (resultStr[0] == ' ' || resultStr[0] == '\t') {
		resultStr = resultStr[1:]
	}
	// Trim trailing spaces and tabs
	for len(resultStr) > 0 && (resultStr[len(resultStr)-1] == ' ' || resultStr[len(resultStr)-1] == '\t') {
		resultStr = resultStr[:len(resultStr)-1]
	}

	return resultStr
}

// ValidateConfiguration performs comprehensive validation of project configuration
func ValidateConfiguration(config any) *ValidationResult {
	result := &ValidationResult{Valid: true}

	// This will be implemented when we integrate with the actual config struct
	// For now, it's a placeholder that will be expanded

	return result
}
