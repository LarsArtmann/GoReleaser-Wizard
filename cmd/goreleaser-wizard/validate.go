package main

import (
	"fmt"
	"os"

	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
	"github.com/spf13/cobra"
)

// printCommandHeader prints a styled command heading.
func printCommandHeader(title string) {
	fmt.Println(titleStyle.Render(title))
	fmt.Println()
}

var fileSystemRepo domain.FileSystemRepository

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
	configureValidateCommand(validateCmd)
}

// configureValidateCommand adds the validate command's flags. Both the real
// command and tests wire flags through it, so flag lookups can never hit an
// unregistered name.
func configureValidateCommand(cmd *cobra.Command) {
	cmd.Flags().Bool("verbose", false, "show detailed validation output")
	cmd.Flags().Bool("fix", false, "attempt to fix common issues")
	cmd.Flags().Bool("project-only", false, "validate project structure only")
}

func runValidate(cmd *cobra.Command, _ []string) {
	exitCode := performValidation(cmd)
	os.Exit(exitCode)
}

// performValidation executes the validation logic and returns the exit code.
// This function is separated from runValidate to allow testing without os.Exit.
func performValidation(cmd *cobra.Command) int {
	// Set up panic recovery using domain error handling
	defer recoverFromPanic("validate command")

	verbose := getBoolFlag(cmd, "verbose")
	fix := getBoolFlag(cmd, "fix")
	projectOnly := getBoolFlag(cmd, "project-only")

	printCommandHeader("🔍 Validating GoReleaser Configuration")

	// Initialize dependencies (in real implementation, this would be injected)
	fileSystemRepo = &SimpleFileSystemRepository{}

	// Collect validation results
	results := &ValidationResults{}

	if !projectOnly {
		// Validate GoReleaser configuration
		err := validateGoReleaserConfig(results)
		if err != nil {
			displayError(err)

			return 1
		}

		// Validate GitHub Actions workflow
		validateGitHubActions(results)
	}

	// Validate project structure
	validateProjectStructure(results)

	// Display results
	displayValidationResults(results, verbose)

	// Attempt fixes if requested
	if fix && len(results.Errors) > 0 {
		attemptFixes(results)
	}

	return results.GetExitCode()
}
