package validation

import (
	"fmt"
	"regexp"
	"strings"
)

// TemplateEscaper provides secure escaping for different output formats
type TemplateEscaper struct{}

// NewTemplateEscaper creates a new template escaper
func NewTemplateEscaper() *TemplateEscaper {
	return &TemplateEscaper{}
}

// EscapeYAML escapes values for safe YAML output
func (te *TemplateEscaper) EscapeYAML(value string) string {
	if value == "" {
		return ""
	}

	// First sanitize the input
	value = SanitizeInput(value)

	// YAML escaping rules - remove the premature escaping

	// For multi-line strings, use YAML literal block style if needed
	if strings.Contains(value, "\n") {
		// Check if it needs literal block style
		if strings.ContainsAny(value, ":{}[],&*#?|-<>'\"%@`") {
			return fmt.Sprintf("|-\n%s", te.indentYAMLLines(value))
		}
		return fmt.Sprintf("|-\n%s", te.indentYAMLLines(value))
	}

	// For single line, check if it needs quoting
	if strings.ContainsAny(value, ":{}[],&*#?|-<>'\"%@`") ||
		strings.HasPrefix(value, " ") || strings.HasSuffix(value, " ") ||
		strings.HasPrefix(value, "!") || strings.HasPrefix(value, "&") ||
		looksLikeNumber(value) {
		return fmt.Sprintf("'%s'", strings.ReplaceAll(value, "'", "''"))
	}

	return value
}

// EscapeShell escapes values for safe shell command usage
func (te *TemplateEscaper) EscapeShell(value string) string {
	if value == "" {
		return ""
	}

	// Sanitize first
	value = SanitizeInput(value)

	// Check for dangerous shell patterns
	if containsShellInjection(value) {
		return "" // Don't escape potentially dangerous values
	}

	// Basic shell escaping for POSIX shells
	value = strings.ReplaceAll(value, "'", "''")

	// Wrap in single quotes for maximum safety
	return fmt.Sprintf("'%s'", value)
}

// EscapeGitHubActions escapes values for GitHub Actions workflow files
func (te *TemplateEscaper) EscapeGitHubActions(value string) string {
	if value == "" {
		return ""
	}

	// Sanitize first
	value = SanitizeInput(value)

	// GitHub Actions YAML escaping
	value = te.EscapeYAML(value)

	// Additional GitHub Actions specific escaping
	if strings.Contains(value, "${{") {
		// Escape GitHub Actions expression syntax
		value = strings.ReplaceAll(value, "${{", "${{ '' }}${{")
	}

	return value
}

// EscapeJSON escapes values for JSON output
func (te *TemplateEscaper) EscapeJSON(value string) string {
	if value == "" {
		return `""`
	}

	// Sanitize first
	value = SanitizeInput(value)

	// JSON escaping
	value = strings.ReplaceAll(value, `\`, `\\`)
	value = strings.ReplaceAll(value, `"`, `\"`)
	value = strings.ReplaceAll(value, "\n", `\n`)
	value = strings.ReplaceAll(value, "\r", `\r`)
	value = strings.ReplaceAll(value, "\t", `\t`)

	return fmt.Sprintf(`"%s"`, value)
}

// EscapeDockerLabel escapes values for Docker labels
func (te *TemplateEscaper) EscapeDockerLabel(value string) string {
	if value == "" {
		return ""
	}

	// Sanitize first
	value = SanitizeInput(value)

	// Docker labels have specific restrictions
	// Must match pattern: [a-zA-Z0-9._-]+
	if !isValidDockerLabel(value) {
		// Sanitize to valid characters only
		var result strings.Builder
		for _, r := range value {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
				(r >= '0' && r <= '9') || r == '.' || r == '_' || r == '-' {
				result.WriteRune(r)
			} else {
				result.WriteRune('-') // Replace invalid chars with dash
			}
		}
		value = result.String()
	}

	// Ensure it doesn't start with digit, dot, or dash
	if len(value) > 0 && (value[0] >= '0' && value[0] <= '9' ||
		value[0] == '.' || value[0] == '-') {
		value = "label-" + value
	}

	return value
}

// ValidateTemplateContent validates generated content for security
func (te *TemplateEscaper) ValidateTemplateContent(content, templateType string) error {
	// Check for injection patterns
	dangerousPatterns := []string{
		"<script",
		"javascript:",
		"vbscript:",
		"onload=",
		"onerror=",
		"onclick=",
		"${{",
		"${SCRIPT",
		"<%",
		"<%=",
		"`", // Backticks for command substitution
		"$(",
		";rm",
		"|rm",
		"&&rm",
		"||rm",
	}

	lowerContent := strings.ToLower(content)
	for _, pattern := range dangerousPatterns {
		if strings.Contains(lowerContent, pattern) {
			return fmt.Errorf("template contains potentially dangerous content: %s", pattern)
		}
	}

	// Type-specific validations
	switch templateType {
	case "yaml":
		if te.containsYAMLInjection(content) {
			return fmt.Errorf("YAML template contains potential injection")
		}
	case "shell":
		if containsShellInjection(content) {
			return fmt.Errorf("shell script contains potential injection")
		}
	case "github-actions":
		if te.containsGitHubActionsInjection(content) {
			return fmt.Errorf("GitHub Actions workflow contains potential injection")
		}
	}

	return nil
}

// Helper functions

func (te *TemplateEscaper) indentYAMLLines(value string) string {
	lines := strings.Split(value, "\n")
	var result strings.Builder

	for i, line := range lines {
		if i > 0 {
			result.WriteString("  ") // 2 spaces for YAML continuation
		}
		result.WriteString(line)
		if i < len(lines)-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}

func looksLikeNumber(value string) bool {
	if value == "" {
		return false
	}

	// Check if it's a number (which would need quoting in YAML)
	hasDigits := false
	for _, r := range value {
		if r >= '0' && r <= '9' {
			hasDigits = true
		} else if r != '.' && r != '-' && r != '+' && r != 'e' && r != 'E' {
			return false
		}
	}

	return hasDigits
}

func containsShellInjection(value string) bool {
	shellPatterns := []string{
		";", "|", "&", "<", ">", "`", "$(", "${",
		"rm ", "del ", "format ", "shutdown", "reboot",
		">/dev/", "</dev/", "2>&1",
	}

	lowerValue := strings.ToLower(value)
	for _, pattern := range shellPatterns {
		if strings.Contains(lowerValue, pattern) {
			return true
		}
	}

	return false
}

func isValidDockerLabel(value string) bool {
	if value == "" {
		return false
	}

	// Docker label validation regex pattern
	dockerLabelPattern := regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
	return dockerLabelPattern.MatchString(value)
}

func (te *TemplateEscaper) containsYAMLInjection(content string) bool {
	yamlPatterns := []string{
		"!!", "!!map", "!!seq", "!!str", "!!int",
		"anchor:", "alias:", "<<:", "&", "*",
	}

	lowerContent := strings.ToLower(content)
	for _, pattern := range yamlPatterns {
		if strings.Contains(lowerContent, pattern) {
			return true
		}
	}

	return false
}

func (te *TemplateEscaper) containsGitHubActionsInjection(content string) bool {
	githubPatterns := []string{
		"${{", "::set-output", "::add-path", "::error", "::warning",
		"$GITHUB_", "github.token", "secrets.",
	}

	lowerContent := strings.ToLower(content)
	for _, pattern := range githubPatterns {
		if strings.Contains(lowerContent, pattern) {
			return true
		}
	}

	return false
}
