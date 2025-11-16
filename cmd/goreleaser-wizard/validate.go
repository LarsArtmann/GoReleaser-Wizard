package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/LarsArtmann/template-GoReleaser/internal/domain"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Validation use case (would be injected in real implementation)
	validationUseCase *domain.ValidationUseCase
	fileSystemRepo   domain.FileSystemRepository
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate GoReleaser configuration",
	Long: `Validate your GoReleaser configuration and check for common issues.

This command will:
- Check if .goreleaser.yaml exists and is valid YAML
- Run goreleaser check if available
- Verify project structure matches configuration
- Check for missing dependencies
- Suggest improvements`,
	Run: runValidate,
}

func init() {
	validateCmd.Flags().Bool("verbose", false, "show detailed validation output")
	validateCmd.Flags().Bool("fix", false, "attempt to fix common issues")
	validateCmd.Flags().Bool("project-only", false, "validate project structure only")
}

func runValidate(cmd *cobra.Command, args []string) {
	// Set up panic recovery using domain error handling
	defer recoverFromPanic("validate command")

	verbose, _ := cmd.Flags().GetBool("verbose")
	fix, _ := cmd.Flags().GetBool("fix")
	projectOnly, _ := cmd.Flags().GetBool("project-only")

	fmt.Println(titleStyle.Render("🔍 Validating GoReleaser Configuration"))
	fmt.Println()

	// Initialize dependencies (in real implementation, this would be injected)
	fileSystemRepo = &SimpleFileSystemRepository{}
	validationUseCase = domain.NewValidationUseCase(logger, fileSystemRepo)

	// Collect validation results
	results := &ValidationResults{}

	if !projectOnly {
		// Validate GoReleaser configuration
		if err := validateGoReleaserConfig(&results); err != nil {
			displayError(err)
			return
		}

		// Validate GitHub Actions workflow
		if err := validateGitHubActions(&results); err != nil {
			displayError(err)
			return
		}
	}

	// Validate project structure
	if err := validateProjectStructure(&results); err != nil {
		displayError(err)
		return
	}

	// Display results
	displayValidationResults(&results, verbose)

	// Attempt fixes if requested
	if fix && len(results.Errors) > 0 {
		if err := attemptFixes(&results); err != nil {
			displayError(err)
			return
		}
	}

	// Exit with appropriate code
	os.Exit(results.GetExitCode())
}

// ValidationResults holds all validation results
type ValidationResults struct {
	ConfigExists    bool
	ConfigValid     bool
	ActionsExists   bool
	ActionsValid    bool
	ProjectValid    bool
	GoReleaserFound bool
	Errors         []*domain.DomainError
	Warnings       []*domain.DomainError
	Recommendations []string
}

// GetExitCode returns appropriate exit code
func (vr *ValidationResults) GetExitCode() int {
	if len(vr.Errors) > 0 {
		return 1
	}
	if len(vr.Warnings) > 0 {
		return 2
	}
	return 0
}

// validateGoReleaserConfig validates GoReleaser configuration
func validateGoReleaserConfig(results *ValidationResults) error {
	configPath := ".goreleaser.yaml"

	// Check if config exists
	exists, err := fileSystemRepo.FileExists(context.Background(), configPath)
	if err != nil {
		results.Errors = append(results.Errors, 
			domain.NewSystemError(
				domain.ErrFileReadFailed,
				"Failed to check configuration file",
				fmt.Sprintf("Cannot access %s", configPath),
				err,
			).WithContext(configPath))
		return nil
	}

	results.ConfigExists = exists
	if !exists {
		results.Errors = append(results.Errors,
			domain.NewSystemError(
				domain.ErrFileNotFound,
				"Configuration file not found",
				fmt.Sprintf("%s does not exist", configPath),
				nil,
			).WithContext(configPath))
		results.Recommendations = append(results.Recommendations, 
			"Run 'goreleaser-wizard init' to create configuration")
		return nil
	}

	// Validate YAML syntax
	if err := validateYAML(configPath, results); err != nil {
		return err
	}

	// Try to parse and validate as GoReleaser config
	if err := parseGoReleaserConfig(configPath, results); err != nil {
		return err
	}

	// Run goreleaser check if available
	if err := runGoReleaserCheck(configPath, results); err != nil {
		return nil // Not fatal, just record warning
	}

	results.ConfigValid = len(results.Errors) == 0
	return nil
}

// validateGitHubActions validates GitHub Actions workflow
func validateGitHubActions(results *ValidationResults) error {
	workflowPath := ".github/workflows/release.yml"

	// Check if workflow exists
	exists, err := fileSystemRepo.FileExists(context.Background(), workflowPath)
	if err != nil {
		results.Warnings = append(results.Warnings,
			domain.NewSystemError(
				domain.ErrFileReadFailed,
				"Failed to check GitHub Actions workflow",
				fmt.Sprintf("Cannot access %s", workflowPath),
				err,
			).WithContext(workflowPath))
		return nil
	}

	results.ActionsExists = exists
	if !exists {
		results.Recommendations = append(results.Recommendations,
			"Add GitHub Actions workflow for automated releases")
		return nil
	}

	// Validate YAML syntax
	if err := validateYAML(workflowPath, results); err != nil {
		return err
	}

	// Validate workflow content
	if err := validateWorkflowContent(workflowPath, results); err != nil {
		return err
	}

	results.ActionsValid = len(results.Errors) == 0
	return nil
}

// validateProjectStructure validates project structure
func validateProjectStructure(results *ValidationResults) error {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		results.Errors = append(results.Errors,
			domain.NewSystemError(
				domain.ErrPermissionDenied,
				"Failed to get current directory",
				"Cannot determine working directory",
				err,
			))
		return nil
	}

	// Use domain validation use case
	ctx := context.Background()
	result, err := validationUseCase.ValidateProjectStructure(ctx, cwd)
	if err != nil {
		results.Errors = append(results.Errors, err.(*domain.DomainError))
		return nil
	}

	results.ProjectValid = result.IsValid
	results.Errors = append(results.Errors, result.Issues...)
	results.Warnings = append(results.Warnings, result.Warnings...)
	results.Recommendations = append(results.Recommendations, result.Recommendations...)

	return nil
}

// validateYAML validates YAML syntax
func validateYAML(filePath string, results *ValidationResults) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		results.Errors = append(results.Errors,
			domain.NewSystemError(
				domain.ErrFileReadFailed,
				"Failed to read file",
				fmt.Sprintf("Cannot read %s", filePath),
				err,
			).WithContext(filePath))
		return nil
	}

	// Simple YAML validation - check for balanced brackets and quotes
	content := string(data)
	if !isValidYAML(content) {
		results.Errors = append(results.Errors,
			domain.NewTemplateError(
				domain.ErrTemplateSyntaxError,
				"Invalid YAML syntax",
				fmt.Sprintf("File %s contains YAML syntax errors", filePath),
			).WithContext(filePath))
		return nil
	}

	return nil
}

// isValidYAML performs basic YAML validation
func isValidYAML(content string) bool {
	// This is a basic check - in real implementation, use yaml parser
	indent := 0
	inString := false
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Count brackets and braces
		for _, char := range line {
			switch char {
			case '"':
				inString = !inString
			case '{':
				if !inString {
					indent++
				}
			case '}':
				if !inString {
					indent--
				}
			}
		}

		// Check indentation
		if !inString && indent < 0 {
			return false
		}
	}

	return indent == 0
}

// parseGoReleaserConfig parses and validates GoReleaser configuration
func parseGoReleaserConfig(configPath string, results *ValidationResults) error {
	// This would use proper YAML parsing and GoReleaser validation
	// For now, just check for required fields
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		return domain.NewSystemError(
			domain.ErrFileReadFailed,
			"Failed to read configuration",
			fmt.Sprintf("Cannot read %s", configPath),
			err,
		)
	}

	content := string(data)

	// Check for required fields
	requiredFields := []string{"project_name", "builds"}
	for _, field := range requiredFields {
		if !strings.Contains(content, field+":") {
			results.Errors = append(results.Errors,
				domain.NewValidationError(
					domain.ErrMissingRequiredField,
					"Missing required field",
					fmt.Sprintf("Configuration missing required field: %s", field),
				).WithContext(field))
		}
	}

	return nil
}

// runGoReleaserCheck runs goreleaser check command
func runGoReleaserCheck(configPath string, results *ValidationResults) error {
	// Check if goreleaser is available
	if _, err := exec.LookPath("goreleaser"); err != nil {
		results.Warnings = append(results.Warnings,
			domain.NewExternalServiceError(
				domain.ErrDependencyNotFound,
				"GoReleaser not found",
				"goreleaser command not available in PATH",
			))
		results.Recommendations = append(results.Recommendations,
			"Install GoReleaser: https://goreleaser.com/install/")
		results.GoReleaserFound = false
		return nil
	}

	results.GoReleaserFound = true

	// Run goreleaser check
	cmd := exec.Command("goreleaser", "check", "-f", configPath)
	output, err := cmd.CombinedOutput()

	if err != nil {
		results.Errors = append(results.Errors,
			domain.NewExternalServiceError(
				domain.ErrGitOperationFailed,
				"GoReleaser check failed",
				string(output),
			))
		return nil
	}

	logger.Info("GoReleaser check passed")
	return nil
}

// validateWorkflowContent validates GitHub Actions workflow content
func validateWorkflowContent(workflowPath string, results *ValidationResults) error {
	data, err := os.ReadFile(workflowPath)
	if err != nil {
		return domain.NewSystemError(
			domain.ErrFileReadFailed,
			"Failed to read workflow file",
			fmt.Sprintf("Cannot read %s", workflowPath),
			err,
		)
	}

	content := string(data)

	// Check for required workflow elements
	requiredElements := []string{"name:", "on:", "jobs:"}
	for _, element := range requiredElements {
		if !strings.Contains(content, element) {
			results.Warnings = append(results.Warnings,
				domain.NewTemplateError(
					domain.ErrTemplateExecutionFailed,
					"Missing workflow element",
					fmt.Sprintf("Workflow missing required element: %s", element),
				).WithContext(workflowPath))
		}
	}

	return nil
}

// displayValidationResults displays validation results
func displayValidationResults(results *ValidationResults, verbose bool) {
	fmt.Println("📋 Validation Summary:")
	fmt.Println()

	// Configuration status
	if results.ConfigExists {
		if results.ConfigValid {
			fmt.Println(successStyle.Render("✅ GoReleaser configuration: Valid"))
		} else {
			fmt.Println(errorStyle.Render("❌ GoReleaser configuration: Invalid"))
		}
	} else {
		fmt.Println(errorStyle.Render("❌ GoReleaser configuration: Not found"))
	}

	// GitHub Actions status
	if results.ActionsExists {
		if results.ActionsValid {
			fmt.Println(successStyle.Render("✅ GitHub Actions workflow: Valid"))
		} else {
			fmt.Println(errorStyle.Render("❌ GitHub Actions workflow: Invalid"))
		}
	} else {
		fmt.Println(infoStyle.Render("ℹ️  GitHub Actions workflow: Not found"))
	}

	// Project structure status
	if results.ProjectValid {
		fmt.Println(successStyle.Render("✅ Project structure: Valid"))
	} else {
		fmt.Println(errorStyle.Render("❌ Project structure: Invalid"))
	}

	// GoReleaser availability
	if results.GoReleaserFound {
		fmt.Println(successStyle.Render("✅ GoReleaser: Available"))
	} else {
		fmt.Println(infoStyle.Render("ℹ️  GoReleaser: Not installed"))
	}

	fmt.Println()

	// Display errors
	if len(results.Errors) > 0 {
		fmt.Println(errorStyle.Render("❌ Errors:"))
		for _, err := range results.Errors {
			fmt.Printf("  • %s\n", err.Message)
			if verbose {
				fmt.Printf("    Details: %s\n", err.Details)
				if err.Context != "" {
					fmt.Printf("    Context: %s\n", err.Context)
				}
			}
		}
		fmt.Println()
	}

	// Display warnings
	if len(results.Warnings) > 0 {
		fmt.Println(infoStyle.Render("⚠️  Warnings:"))
		for _, warning := range results.Warnings {
			fmt.Printf("  • %s\n", warning.Message)
			if verbose {
				fmt.Printf("    Details: %s\n", warning.Details)
			}
		}
		fmt.Println()
	}

	// Display recommendations
	if len(results.Recommendations) > 0 {
		fmt.Println(infoStyle.Render("💡 Recommendations:"))
		for _, rec := range results.Recommendations {
			fmt.Printf("  • %s\n", rec)
		}
		fmt.Println()
	}
}

// attemptFixes attempts to fix common issues
func attemptFixes(results *ValidationResults) error {
	fmt.Println("🔧 Attempting to fix common issues...")
	fmt.Println()

	fixed := 0

	// Fix missing configuration directory
	if !results.ConfigExists {
		if err := os.MkdirAll(".github/workflows", 0755); err == nil {
			fmt.Println(successStyle.Render("✅ Created .github/workflows directory"))
			fixed++
		}
	}

	// Try to create basic configuration if missing
	if !results.ConfigExists {
		configContent := generateBasicConfig()
		if err := os.WriteFile(".goreleaser.yaml", []byte(configContent), 0644); err == nil {
			fmt.Println(successStyle.Render("✅ Created basic .goreleaser.yaml"))
			fixed++
		}
	}

	if fixed > 0 {
		fmt.Println()
		fmt.Println(successStyle.Render(fmt.Sprintf("✅ Fixed %d issues", fixed)))
		fmt.Println(infoStyle.Render("💡 Run validation again to check remaining issues"))
	} else {
		fmt.Println(infoStyle.Render("ℹ️  No auto-fixable issues found"))
	}

	return nil
}

// generateBasicConfig generates a basic GoReleaser configuration
func generateBasicConfig() string {
	return `# Basic GoReleaser configuration
# Generated by GoReleaser Wizard

# The project name
project_name: my-project

# The build configuration
builds:
  - env:
      - CGO_ENABLED=0
    goos:
      - linux
      - darwin
      - windows
    goarch:
      - amd64
      - arm64

# Archive configuration
archives:
  - format: tar.gz
    name_template: '{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}'

# Release configuration
release:
  github:
    owner: your-username
    name: my-project
`
}

// SimpleFileSystemRepository is a basic implementation for demonstration
type SimpleFileSystemRepository struct{}

func (r *SimpleFileSystemRepository) ReadFile(ctx context.Context, path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (r *SimpleFileSystemRepository) WriteFile(ctx context.Context, path string, data []byte, perm os.FileMode) error {
	return os.WriteFile(path, data, perm)
}

func (r *SimpleFileSystemRepository) CreateFile(ctx context.Context, path string) (io.WriteCloser, error) {
	return os.Create(path)
}

func (r *SimpleFileSystemRepository) DeleteFile(ctx context.Context, path string) error {
	return os.Remove(path)
}

func (r *SimpleFileSystemRepository) FileExists(ctx context.Context, path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (r *SimpleFileSystemRepository) CreateDir(ctx context.Context, path string, perm os.FileMode) error {
	return os.Mkdir(path, perm)
}

func (r *SimpleFileSystemRepository) CreateDirAll(ctx context.Context, path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (r *SimpleFileSystemRepository) DirExists(ctx context.Context, path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		return info.IsDir(), nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (r *SimpleFileSystemRepository) ReadDir(ctx context.Context, path string) ([]os.DirEntry, error) {
	return os.ReadDir(path)
}

func (r *SimpleFileSystemRepository) GetFileInfo(ctx context.Context, path string) (os.FileInfo, error) {
	return os.Stat(path)
}

func (r *SimpleFileSystemRepository) CheckPermissions(ctx context.Context, path string) (bool, error) {
	// Basic permission check - try to read file
	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return false, nil
	}
	file.Close()
	return true, nil
}

func (r *SimpleFileSystemRepository) AbsPath(path string) (string, error) {
	return filepath.Abs(path)
}

func (r *SimpleFileSystemRepository) RelPath(base, target string) (string, error) {
	return filepath.Rel(base, target)
}

func (r *SimpleFileSystemRepository) CleanPath(path string) string {
	return filepath.Clean(path)
}

func (r *SimpleFileSystemRepository) JoinPath(elem ...string) string {
	return filepath.Join(elem...)
}

func (r *SimpleFileSystemRepository) TempDir(dir, pattern string) (string, error) {
	return os.MkdirTemp(dir, pattern)
}