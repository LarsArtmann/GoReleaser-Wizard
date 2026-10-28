package main

import (
	"fmt"
	"os"

	"charm.land/log/v2"
	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
	"github.com/spf13/cobra"
)

// runGenerate runs the generate command.
func runGenerate(cmd *cobra.Command, args []string) {
	defer recoverFromPanic("generate command")

	applyGitHubOverrides(cmd)

	force := getBoolFlag(cmd, "force")
	configOnly := getBoolFlag(cmd, "config-only")

	printCommandHeader("⚙️ Generating GoReleaser Configuration")

	// Detect project information
	projectConfig := &domain.SafeProjectConfig{}

	err := detectProjectInfo(projectConfig)
	if err != nil {
		displayError(err)

		return
	}

	// Override with flags if provided
	if name := getStringFlag(cmd, "project-name"); name != "" {
		projectConfig.ProjectName = name
	}

	if path := getStringFlag(cmd, "main-path"); path != "" {
		projectConfig.MainPath = path
	}

	if binName := getStringFlag(cmd, "binary-name"); binName != "" {
		projectConfig.BinaryName = binName
	}

	if projType := getStringFlag(cmd, "project-type"); projType != "" {
		projectConfig.ProjectType = domain.ProjectType(projType)
	}

	if provider := getStringFlag(cmd, "ci-provider"); provider != "" {
		projectConfig.CIProvider = domain.ConvertToCIProvider(provider)
	}

	// Create workflow based on options
	var workflowType WorkflowType
	if configOnly {
		workflowType = WorkflowTypeConfigOnly
	} else {
		workflowType = WorkflowTypeFullWizard
	}

	logger := log.New(os.Stderr)
	if !ExecuteWorkflow(workflowType, projectConfig, force, logger) {
		return
	}

	fmt.Println()

	if configOnly {
		fmt.Println(successStyle.Render("📄 GoReleaser configuration generated successfully!"))
	} else {
		fmt.Println(
			successStyle.Render("🎉 GoReleaser configuration and workflow generated successfully!"),
		)
	}

	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  • Review generated files")
	fmt.Println("  • Run 'goreleaser-wizard validate' to check configuration")

	if !configOnly {
		fmt.Println("  • Commit .github/workflows/release.yml to enable CI/CD")
	}

	fmt.Println("  • Test configuration with 'goreleaser check'")
}

func init() {
	configureGenerateCommand(generateCmd)
}
