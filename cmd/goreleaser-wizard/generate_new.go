package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/LarsArtmann/template-GoReleaser/internal/domain"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate GoReleaser configuration with flags (non-interactive)",
	Long: `Generate GoReleaser configuration using command-line flags instead of
the interactive wizard. Useful for CI/CD pipelines and automation.`,
	Run: runGenerate,
}

func init() {
	generateCmd.Flags().String("name", "", "project name")
	generateCmd.Flags().String("description", "", "project description")
	generateCmd.Flags().String("binary", "", "binary name")
	generateCmd.Flags().String("main", ".", "path to main.go")
	generateCmd.Flags().StringSlice("platforms", []string{"linux", "darwin", "windows"}, "target platforms")
	generateCmd.Flags().StringSlice("architectures", []string{"amd64", "arm64"}, "target architectures")
	generateCmd.Flags().Bool("docker", false, "enable Docker builds")
	generateCmd.Flags().Bool("signing", false, "enable code signing")
	generateCmd.Flags().Bool("github-action", false, "generate GitHub Actions workflow")
	generateCmd.Flags().Bool("force", false, "overwrite existing files")
	generateCmd.Flags().String("project-type", "cli", "project type")
	generateCmd.Flags().String("git-provider", "github", "git provider")
}

func runGenerate(cmd *cobra.Command, args []string) {
	// Set up panic recovery using domain error handling
	defer recoverFromPanic("generate command")

	// Parse flags and create domain-safe configuration
	safeConfig, err := parseGenerateFlags(cmd)
	if err != nil {
		displayError(err)
		return
	}

	// Validate configuration using domain types
	if err := validateConfiguration(safeConfig); err != nil {
		displayError(err)
		return
	}

	// Check for existing files
	force, _ := cmd.Flags().GetBool("force")
	if err := checkExistingFiles(safeConfig, force); err != nil {
		displayError(err)
		return
	}

	// Generate configuration files
	if err := generateConfigurationFiles(safeConfig); err != nil {
		displayError(err)
		return
	}

	fmt.Println()
	fmt.Println(successStyle.Render("✅ Configuration generated successfully!"))
	fmt.Println()
	fmt.Println("Generated files:")
	fmt.Println("  📄 .goreleaser.yaml")
	if safeConfig.GenerateActions {
		fmt.Println("  🔄 .github/workflows/release.yml")
	}
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  1. Review generated files")
	fmt.Println("  2. Test with: goreleaser check")
	fmt.Println("  3. Create a release")
}

// parseGenerateFlags parses command-line flags into domain-safe configuration
func parseGenerateFlags(cmd *cobra.Command) (*domain.SafeProjectConfig, error) {
	safeConfig := domain.NewSafeProjectConfig()

	// Parse basic information
	name, _ := cmd.Flags().GetString("name")
	description, _ := cmd.Flags().GetString("description")
	binary, _ := cmd.Flags().GetString("binary")
	main, _ := cmd.Flags().GetString("main")
	projectTypeStr, _ := cmd.Flags().GetString("project-type")

	if name == "" {
		return nil, domain.NewValidationError(
			domain.ErrMissingRequiredField,
			"Project name required",
			"Use --name to specify project name",
		)
	}

	safeConfig.ProjectName = name
	safeConfig.ProjectDescription = description
	safeConfig.BinaryName = binary
	safeConfig.MainPath = main

	// Convert and validate project type
	projectType := convertToProjectType(projectTypeStr)
	if !projectType.IsValid() {
		return nil, domain.NewValidationError(
			domain.ErrInvalidProjectName,
			"Invalid project type",
			fmt.Sprintf("'%s' is not a valid project type", projectTypeStr),
		)
	}
	safeConfig.ProjectType = projectType

	// Parse build configuration
	platforms, _ := cmd.Flags().GetStringSlice("platforms")
	architectures, _ := cmd.Flags().GetStringSlice("architectures")

	// Convert platforms
	configPlatforms := make([]domain.Platform, len(platforms))
	for i, p := range platforms {
		platform := domain.Platform(strings.ToLower(p))
		if !platform.IsValid() {
			return nil, domain.NewValidationError(
				domain.ErrInvalidPlatform,
				"Invalid platform",
				fmt.Sprintf("'%s' is not a valid platform", p),
			)
		}
		configPlatforms[i] = platform
	}
	safeConfig.Platforms = configPlatforms

	// Convert architectures
	configArchitectures := make([]domain.Architecture, len(architectures))
	for i, a := range architectures {
		arch := domain.Architecture(strings.ToLower(a))
		if !arch.IsValid() {
			return nil, domain.NewValidationError(
				domain.ErrInvalidArchitecture,
				"Invalid architecture",
				fmt.Sprintf("'%s' is not a valid architecture", a),
			)
		}
		configArchitectures[i] = arch
	}
	safeConfig.Architectures = configArchitectures

	// Parse release configuration
	docker, _ := cmd.Flags().GetBool("docker")
	signing, _ := cmd.Flags().GetBool("signing")
	gitProviderStr, _ := cmd.Flags().GetString("git-provider")

	// Convert git provider
	gitProvider := convertToGitProvider(gitProviderStr)
	if !gitProvider.IsValid() {
		return nil, domain.NewValidationError(
			domain.ErrInvalidGitProvider,
			"Invalid git provider",
			fmt.Sprintf("'%s' is not a valid git provider", gitProviderStr),
		)
	}

	safeConfig.DockerEnabled = docker
	safeConfig.DockerRegistry = gitProvider.DefaultRegistry()
	safeConfig.Signing = signing
	safeConfig.GitProvider = gitProvider

	// Parse CI/CD configuration
	githubAction, _ := cmd.Flags().GetBool("github-action")
	safeConfig.GenerateActions = githubAction
	if githubAction {
		safeConfig.ActionsOn = []domain.ActionTrigger{domain.ActionTriggerVersionTags}
	}

	// Apply defaults
	safeConfig.ApplyDefaults()

	// Set binary name if not provided
	if safeConfig.BinaryName == "" {
		safeConfig.BinaryName = safeConfig.ProjectName
	}

	return safeConfig, nil
}

// validateConfiguration validates configuration using domain types
func validateConfiguration(config *domain.SafeProjectConfig) error {
	// Update state to processing
	config.State = domain.ConfigStateProcessing

	// Validate invariants
	if err := config.ValidateInvariants(); err != nil {
		config.State = domain.ConfigStateInvalid
		return err
	}

	// Mark as valid
	config.State = domain.ConfigStateValid
	return nil
}

// checkExistingFiles checks for existing configuration files
func checkExistingFiles(config *domain.SafeProjectConfig, force bool) *domain.DomainError {
	configFile := ".goreleaser.yaml"
	actionsFile := ".github/workflows/release.yml"

	// Check GoReleaser config
	if exists, err := fileExists(configFile); err != nil {
		return domain.NewSystemError(
			domain.ErrFileReadFailed,
			"Failed to check configuration file",
			fmt.Sprintf("Cannot access %s", configFile),
			err,
		).WithContext(configFile)
	} else if exists && !force {
		return domain.NewConfigurationError(
			domain.ErrFileWriteFailed,
			"Configuration already exists",
			fmt.Sprintf("File %s already exists", configFile),
			"Use --force to overwrite",
		).WithContext(configFile)
	}

	// Check GitHub Actions
	if config.GenerateActions {
		if exists, err := fileExists(actionsFile); err != nil {
			return domain.NewSystemError(
				domain.ErrFileReadFailed,
				"Failed to check GitHub Actions workflow",
				fmt.Sprintf("Cannot access %s", actionsFile),
				err,
			).WithContext(actionsFile)
		} else if exists && !force {
			return domain.NewConfigurationError(
				domain.ErrFileWriteFailed,
				"GitHub Actions workflow already exists",
				fmt.Sprintf("File %s already exists", actionsFile),
				"Use --force to overwrite",
			).WithContext(actionsFile)
		}
	}

	return nil
}

// generateConfigurationFiles generates configuration files
func generateConfigurationFiles(config *domain.SafeProjectConfig) error {
	// Update state
	config.State = domain.ConfigStateProcessing

	// Generate GoReleaser configuration
	goreleaserYAML := generateGoReleaserConfigFromFlags(config)
	if err := writeFile(".goreleaser.yaml", goreleaserYAML, 0644); err != nil {
		config.State = domain.ConfigStateInvalid
		return err
	}

	// Generate GitHub Actions workflow
	if config.GenerateActions {
		actionsYAML := generateGitHubActionsFromFlags(config)
		if err := writeFile(".github/workflows/release.yml", actionsYAML, 0644); err != nil {
			config.State = domain.ConfigStateInvalid
			return err
		}
	}

	// Mark as generated
	config.State = domain.ConfigStateGenerated
	return nil
}

// generateGoReleaserConfigFromFlags generates GoReleaser configuration from flags
func generateGoReleaserConfigFromFlags(config *domain.SafeProjectConfig) string {
	var builder strings.Builder

	// Header
	builder.WriteString("# GoReleaser Configuration\n")
	builder.WriteString("# Generated by GoReleaser Wizard\n")
	builder.WriteString(fmt.Sprintf("# Project: %s\n", config.ProjectName))
	builder.WriteString("\n")

	// Project information
	builder.WriteString("project_name: ")
	builder.WriteString(config.ProjectName)
	builder.WriteString("\n")

	if config.ProjectDescription != "" {
		builder.WriteString("project_description: ")
		builder.WriteString(config.ProjectDescription)
		builder.WriteString("\n")
	}

	// Build configuration
	builder.WriteString("\nbuilds:\n")
	for _, platform := range config.Platforms {
		for _, arch := range config.Architectures {
			builder.WriteString("  - env:\n")
			builder.WriteString(fmt.Sprintf("      - CGO_ENABLED=%t\n", config.CGOEnabled))
			builder.WriteString("    goos:\n")
			builder.WriteString(fmt.Sprintf("      - %s\n", platform))
			builder.WriteString("    goarch:\n")
			builder.WriteString(fmt.Sprintf("      - %s\n", arch))
		}
	}

	// Archive configuration
	builder.WriteString("\narchives:\n")
	builder.WriteString("  - format: tar.gz\n")
	builder.WriteString("    name_template: '{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}'\n")

	// Release configuration
	builder.WriteString("\nrelease:\n")
	builder.WriteString("  github:\n")
	builder.WriteString("    owner: YOUR_USERNAME\n")
	builder.WriteString("    name: ")
	builder.WriteString(config.ProjectName)
	builder.WriteString("\n")

	// Docker configuration
	if config.DockerEnabled {
		builder.WriteString("\ndockers:\n")
		for _, platform := range config.Platforms {
			builder.WriteString("  - image_templates:\n")
			if config.DockerRegistry == domain.DockerRegistryGitHub {
				builder.WriteString(fmt.Sprintf("      - 'ghcr.io/your-username/%s:{{ .Tag }}'\n", config.ProjectName))
			} else {
				builder.WriteString(fmt.Sprintf("      - 'your-username/%s:{{ .Tag }}'\n", config.ProjectName))
			}
		}
	}

	// SBOM configuration
	if config.SBOM {
		builder.WriteString("\nsbom:\n")
		builder.WriteString("  artifacts: archive\n")
	}

	return builder.String()
}

// generateGitHubActionsFromFlags generates GitHub Actions workflow
func generateGitHubActionsFromFlags(config *domain.SafeProjectConfig) string {
	var builder strings.Builder

	// Header
	builder.WriteString("name: Release\n")
	builder.WriteString("\n")

	// Triggers
	builder.WriteString("on:\n")
	triggers := domain.GenerateGitHubActionsTriggersYAML(config.ActionsOn)
	for _, trigger := range strings.Split(triggers, "\n") {
		if strings.TrimSpace(trigger) != "" {
			builder.WriteString("  ")
			builder.WriteString(strings.TrimSpace(trigger))
			builder.WriteString("\n")
		}
	}

	// Jobs
	builder.WriteString("\njobs:\n")
	builder.WriteString("  goreleaser:\n")
	builder.WriteString("    runs-on: ubuntu-latest\n")
	builder.WriteString("    steps:\n")

	// Checkout
	builder.WriteString("      - name: Checkout\n")
	builder.WriteString("        uses: actions/checkout@v4\n")
	builder.WriteString("        with:\n")
	builder.WriteString("          fetch-depth: 0\n")

	// Setup Go
	builder.WriteString("\n      - name: Setup Go\n")
	builder.WriteString("        uses: actions/setup-go@v4\n")
	builder.WriteString("        with:\n")
	builder.WriteString("          go-version: 'stable'\n")

	// GoReleaser
	builder.WriteString("\n      - name: Run GoReleaser\n")
	builder.WriteString("        uses: goreleaser/goreleaser-action@v5\n")
	builder.WriteString("        with:\n")
	builder.WriteString("          version: latest\n")
	builder.WriteString("          args: release --clean\n")

	if config.DockerEnabled {
		builder.WriteString("        env:\n")
		if config.DockerRegistry == domain.DockerRegistryGitHub {
			builder.WriteString("          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}\n")
		} else {
			builder.WriteString("          DOCKER_USERNAME: ${{ secrets.DOCKER_USERNAME }}\n")
			builder.WriteString("          DOCKER_PASSWORD: ${{ secrets.DOCKER_PASSWORD }}\n")
		}
	}

	return builder.String()
}

// Display flag help and usage information
func displayGenerateHelp() {
	fmt.Println("Example usage:")
	fmt.Println("  goreleaser-wizard generate --name my-project --type cli")
	fmt.Println("  goreleaser-wizard generate --name my-api --type api --docker")
	fmt.Println("  goreleaser-wizard generate --name my-lib --type library")
	fmt.Println()
	fmt.Println("Available flags:")
	fmt.Println("  --name            Project name (required)")
	fmt.Println("  --description     Project description")
	fmt.Println("  --binary          Binary name (defaults to project name)")
	fmt.Println("  --main            Path to main.go (default: ./)")
	fmt.Println("  --project-type    Project type (cli, web, library, api, desktop)")
	fmt.Println("  --platforms       Target platforms (linux, darwin, windows, freebsd, openbsd, netbsd)")
	fmt.Println("  --architectures   Target architectures (amd64, arm64, 386, arm)")
	fmt.Println("  --docker          Enable Docker builds")
	fmt.Println("  --signing         Enable code signing")
	fmt.Println("  --github-action   Generate GitHub Actions workflow")
	fmt.Println("  --git-provider     Git provider (github, gitlab, bitbucket, gitea, self-hosted)")
	fmt.Println("  --force           Overwrite existing files")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  # CLI application")
	fmt.Println("  goreleaser-wizard generate --name my-cli --type cli --platforms linux,darwin")
	fmt.Println()
	fmt.Println("  # Web service with Docker")
	fmt.Println("  goreleaser-wizard generate --name my-api --type api --docker --signing")
	fmt.Println()
	fmt.Println("  # Library with minimal configuration")
	fmt.Println("  goreleaser-wizard generate --name my-lib --type library")
}