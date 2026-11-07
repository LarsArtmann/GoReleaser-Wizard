package domain

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

// Security validation patterns.
var (
	// Project name pattern: alphanumeric, hyphens, underscores, dots, starts and ends with alphanumeric.
	projectNamePattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9\-._]*[a-zA-Z0-9]$`)

	// Binary name pattern: more restrictive, no special characters that could cause issues.
	binaryNamePattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9\-_]*[a-zA-Z0-9]$`)

	// Path traversal detection.
	pathTraversalPattern = regexp.MustCompile(`\.\.[/\\]`)

	// Shell metacharacters that could be dangerous.
	shellMetacharPattern = regexp.MustCompile(`[;&|<>"'$` + "`" + `\\]`)

	// Reserved names (OS-specific and Go-specific).
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

	// Dangerous file extensions.
	dangerousExtensions = map[string]bool{
		".exe": true, ".bat": true, ".cmd": true, ".com": true, ".pif": true,
		".scr": true, ".vbs": true, ".js": true, ".jar": true, ".sh": true,
		".ps1": true, ".py": true, ".rb": true, ".pl": true, ".php": true,
	}
)

// ValidateProjectName validates project name according to security rules.
func ValidateProjectName(name string) error {
	if name == "" {
		return errors.New("project name cannot be empty")
	}

	if len(name) < 1 || len(name) > 63 {
		return errors.New("project name must be 1-63 characters long")
	}

	// Check for consecutive special characters
	if strings.Contains(name, "--") || strings.Contains(name, "__") ||
		strings.Contains(name, "..") || strings.Contains(name, "__") {
		return errors.New("project name cannot contain consecutive special characters")
	}

	// Validate pattern
	if !projectNamePattern.MatchString(name) {
		return errors.New(
			"project name contains invalid characters. Use letters, numbers, hyphens, underscores, and dots",
		)
	}

	if err := checkReservedName(name, "project"); err != nil {
		return err
	}

	// Check for dangerous patterns
	if strings.HasPrefix(name, ".") || strings.HasPrefix(name, "-") {
		return errors.New("project name cannot start with special characters")
	}

	if strings.HasSuffix(name, ".") || strings.HasSuffix(name, "-") {
		return errors.New("project name cannot end with special characters")
	}

	return nil
}

// ValidateBinaryName validates binary name with strict security rules.
func ValidateBinaryName(name string) error {
	if name == "" {
		return errors.New("binary name cannot be empty")
	}

	if len(name) < 1 || len(name) > 255 {
		return errors.New("binary name must be 1-255 characters long")
	}

	// No spaces or shell metacharacters
	if shellMetacharPattern.MatchString(name) {
		return errors.New("binary name contains dangerous characters")
	}

	// Must be valid filename on all platforms
	if strings.ContainsAny(name, `<>:"/\|?*`) {
		return errors.New("binary name contains invalid filename characters")
	}

	// Validate pattern
	if !binaryNamePattern.MatchString(name) {
		return errors.New(
			"binary name must start with a letter and contain only letters, numbers, hyphens, and underscores",
		)
	}

	if err := checkReservedName(name, "binary"); err != nil {
		return err
	}

	// Check for dangerous extensions
	ext := strings.ToLower(filepath.Ext(name))
	if dangerousExtensions[ext] {
		return fmt.Errorf("binary name has dangerous extension: %s", ext)
	}

	return nil
}

// checkReservedName reports whether name is in the reservedNames set, doing the
// case-insensitive lookup and formatting the error message consistently across
// validators. label is the human-readable noun used in the error (e.g. "project").
func checkReservedName(name, label string) error {
	if reservedNames[strings.ToLower(name)] {
		return fmt.Errorf("'%s' is a reserved %s name and cannot be used", name, label)
	}

	return nil
}

// ValidateMainPath validates the main package path with security checks.
func ValidateMainPath(path string) error {
	if path == "" {
		return errors.New("main path cannot be empty")
	}

	// No path traversal attacks
	if pathTraversalPattern.MatchString(path) {
		return errors.New("path traversal not allowed")
	}

	// No absolute paths
	if filepath.IsAbs(path) {
		return errors.New("absolute paths not allowed")
	}

	// Clean the path
	cleanPath := filepath.Clean(path)

	// Check for dangerous patterns
	if strings.Contains(cleanPath, "..") {
		return errors.New("path traversal not allowed")
	}

	// Check for shell metacharacters
	if shellMetacharPattern.MatchString(cleanPath) {
		return errors.New("path contains dangerous characters")
	}

	// Check for URL scheme injection (javascript:, data:, vbscript:)
	if containsURLInjection(cleanPath) {
		return errors.New("URL scheme injection not allowed")
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

// ValidateProjectDescription validates project description.
func ValidateProjectDescription(desc string) error {
	if len(desc) > maxDescriptionLength {
		return fmt.Errorf(
			"project description too long: %d characters (max %d): %w",
			len(desc),
			maxDescriptionLength,
			errors.New("validation failed"),
		)
	}

	// Check for HTML/markdown injection attempts
	if strings.Contains(desc, "<script") || strings.Contains(desc, "javascript:") {
		return fmt.Errorf(
			"description contains dangerous HTML/script tags for desc=%q: %w",
			desc,
			errors.New("validation failed"),
		)
	}

	// Check for excessive length that might indicate injection
	if len(strings.TrimSpace(desc)) == 0 && len(desc) > 100 {
		return fmt.Errorf(
			"description has suspicious whitespace pattern len=%d: %w",
			len(desc),
			errors.New("validation failed"),
		)
	}

	return nil
}
