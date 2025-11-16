package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/LarsArtmann/template-GoReleaser/internal/domain"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Style definitions
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("99")).
			MarginBottom(1)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86"))
)

	// Legacy ProjectConfig for backward compatibility during migration
	// TODO: Remove after migration complete
	ProjectConfig = struct {
		// Basic Info
		ProjectName        string
		ProjectDescription string
		ProjectType        string
		BinaryName         string
		MainPath           string

		// Build Options
		Platforms     []string
		Architectures []string
		CGOEnabled    bool
		BuildTags     []string
		LDFlags       bool

		// Release Options
		GitProvider    string
		DockerEnabled  bool
		DockerRegistry string
		DockerImage    string
		Signing        bool
		Homebrew       bool
		Snap           bool
		SBOM           bool

		// CI/CD Options
		GenerateActions bool
		ActionsOn       []string

		// Advanced Options
		ProVersion bool
	}{}

	// Form validator using domain types
	formValidator *FormValidator
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize GoReleaser configuration interactively",
	Long: `Initialize GoReleaser configuration through an interactive wizard.

This command will guide you through setting up your GoReleaser configuration
with smart defaults and best practices tailored to your project type.`,
	Run: runInit,
}

func init() {
	initCmd.Flags().Bool("force", false, "overwrite existing configuration")
	initCmd.Flags().Bool("migrate", false, "migrate from existing configuration")
}

func runInit(cmd *cobra.Command, args []string) {
	// Set up panic recovery using domain error handling
	defer recoverFromPanic("init command")

	force, _ := cmd.Flags().GetBool("force")
	migrate, _ := cmd.Flags().GetBool("migrate")

	fmt.Println(titleStyle.Render("🧙‍♂️ Welcome to GoReleaser Wizard"))
	fmt.Println()

	// Check for existing configuration
	if err := checkExistingConfiguration(force); err != nil {
		displayError(err)
		return
	}

	// Initialize form validator with domain types
	formValidator = NewFormValidator()

	// Start the interactive wizard
	config, err := runInteractiveWizard()
	if err != nil {
		displayError(err)
		return
	}

	// Create domain-safe configuration
	safeConfig, err := createSafeConfig(config)
	if err != nil {
		displayError(err)
		return
	}

	// Validate configuration using domain types
	if err := validateConfiguration(safeConfig); err != nil {
		displayError(err)
		return
	}

	// Generate configuration files
	if err := generateConfigurationFiles(safeConfig); err != nil {
		displayError(err)
		return
	}

	fmt.Println()
	fmt.Println(successStyle.Render("✅ Configuration created successfully!"))
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  1. Review the generated .goreleaser.yaml")
	fmt.Println("  2. Test with: goreleaser check")
	fmt.Println("  3. Create a release with: goreleaser release --snapshot")
}

// checkExistingConfiguration checks for existing configuration files
func checkExistingConfiguration(force bool) *domain.DomainError {
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
			"Use --force to overwrite or --migrate to upgrade",
		).WithContext(configFile)
	}

	// Check GitHub Actions
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

	return nil
}

// runInteractiveWizard runs the interactive configuration wizard
func runInteractiveWizard() (*ProjectConfig, error) {
	config := &ProjectConfig{}

	// Project Information
	if err := collectProjectInfo(config); err != nil {
		return nil, err
	}

	// Build Configuration
	if err := collectBuildConfiguration(config); err != nil {
		return nil, err
	}

	// Release Configuration
	if err := collectReleaseConfiguration(config); err != nil {
		return nil, err
	}

	// CI/CD Configuration
	if err := collectCIConfiguration(config); err != nil {
		return nil, err
	}

	return config, nil
}

// collectProjectInfo collects basic project information
func collectProjectInfo(config *ProjectConfig) error {
	return huh.NewForm(
		huh.NewGroup("Project Information"),
		huh.NewInput().
			Title("Project Name").
			Description("The name of your Go project").
			Placeholder("my-awesome-project").
			Validate(formValidator.ValidateProjectName()).
			Value(&config.ProjectName),
		huh.NewSelect[string]().
			Title("Project Type").
			Description("What type of project are you building?").
			Options(
				huh.NewOption("🖥️  CLI Application", "CLI Application"),
				huh.NewOption("🌐  Web Service", "Web Service"),
				huh.NewOption("📚  Library", "Library"),
				huh.NewOption("🔌  API Service", "API Service"),
				huh.NewOption("🖱️  Desktop Application", "Desktop Application"),
			).
			Validate(func(value string) error {
				// Convert to domain type for validation
				projectType := convertToProjectType(value)
				if !projectType.IsValid() {
					return domain.NewValidationError(
						domain.ErrInvalidProjectName,
						"Invalid project type",
						fmt.Sprintf("'%s' is not a valid project type", value),
					)
				}
				return nil
			}).
			Value(&config.ProjectType),
		huh.NewInput().
			Title("Binary Name").
			Description("The name of your compiled binary").
			Placeholder(config.ProjectName).
			Validate(formValidator.ValidateBinaryName()).
			Value(&config.BinaryName),
		huh.NewInput().
			Title("Main Go File Path").
			Description("Path to your main.go file (relative to project root)").
			Placeholder("./cmd/main").
			Validate(formValidator.ValidateMainPath()).
			Value(&config.MainPath),
		huh.NewText().
			Title("Project Description").
			Description("A brief description of your project (optional)").
			Placeholder("A Go project that does amazing things").
			Validate(formValidator.ValidateProjectDescription()).
			Value(&config.ProjectDescription),
	).
		WithTheme(huh.ThemeCatppuccin()).
		Run()

	return nil
}

// collectBuildConfiguration collects build configuration
func collectBuildConfiguration(config *ProjectConfig) error {
	// Get project type for defaults
	projectType := convertToProjectType(config.ProjectType)
	defaultPlatforms := getProjectTypePlatforms(projectType)
	defaultArchitectures := []string{"amd64", "arm64"}
	defaultCGO := projectType.DefaultCGOEnabled()

	return huh.NewForm(
		huh.NewGroup("Build Configuration"),
		huh.NewMultiSelect[string]().
			Title("Target Platforms").
			Description("Which platforms should you build for?").
			Options(
				huh.NewOption("🐧  Linux", "linux"),
				huh.NewOption("🍎  macOS", "darwin"),
				huh.NewOption("🪟  Windows", "windows"),
				huh.NewOption("🦌  FreeBSD", "freebsd"),
				huh.NewOption("🐬  OpenBSD", "openbsd"),
				huh.NewOption("🧊  NetBSD", "netbsd"),
			).
			Validate(func(values []string) error {
				if len(values) == 0 {
					return domain.NewValidationError(
						domain.ErrMissingRequiredField,
						"At least one platform required",
						"Please select at least one target platform",
					)
				}
				return nil
			}).
			Value(&config.Platforms).
			Height(7),
		huh.NewMultiSelect[string]().
			Title("Target Architectures").
			Description("Which CPU architectures should you build for?").
			Options(
				huh.NewOption("x64  (Intel/AMD 64-bit)", "amd64"),
				huh.NewOption("ARM64 (Apple Silicon, ARM 64-bit)", "arm64"),
				huh.NewOption("x86  (Intel/AMD 32-bit)", "386"),
				huh.NewOption("ARM  (32-bit ARM)", "arm"),
			).
			Validate(func(values []string) error {
				if len(values) == 0 {
					return domain.NewValidationError(
						domain.ErrMissingRequiredField,
						"At least one architecture required",
						"Please select at least one target architecture",
					)
				}
				return nil
			}).
			Value(&config.Architectures).
			Height(6),
		huh.NewConfirm().
			Title("Enable CGO").
			Description(fmt.Sprintf("Enable CGO compilation? (Recommended: %t)", defaultCGO)).
			Value(&config.CGOEnabled).
			Default(defaultCGO),
		huh.NewConfirm().
			Title("Enable LDFlags").
			Description("Inject version information using ldflags?").
			Value(&config.LDFlags).
			Default(true),
	).
		WithTheme(huh.ThemeCatppuccin()).
		Run()
}

// collectReleaseConfiguration collects release configuration
func collectReleaseConfiguration(config *ProjectConfig) error {
	return huh.NewForm(
		huh.NewGroup("Release Configuration"),
		huh.NewSelect[string]().
			Title("Git Provider").
			Description("Which git hosting service are you using?").
			Options(
				huh.NewOption("🐙  GitHub", "GitHub"),
				huh.NewOption("🦊  GitLab", "GitLab"),
				huh.NewOption("🪣  Bitbucket", "Bitbucket"),
				huh.NewOption("🕊️  Gitea", "Gitea"),
				huh.NewOption("🏠  Self-hosted", "Self-hosted"),
			).
			Validate(func(value string) error {
				provider := convertToGitProvider(value)
				if !provider.IsValid() {
					return domain.NewValidationError(
						domain.ErrInvalidGitProvider,
						"Invalid git provider",
						fmt.Sprintf("'%s' is not a valid git provider", value),
					)
				}
				return nil
			}).
			Value(&config.GitProvider),
		huh.NewConfirm().
			Title("Enable Docker").
			Description("Build and publish Docker images?").
			Value(&config.DockerEnabled),
		huh.NewConfirm().
			Title("Enable Code Signing").
			Description("Sign releases with Cosign?").
			Value(&config.Signing),
		huh.NewConfirm().
			Title("Enable Homebrew").
			Description("Generate Homebrew formula?").
			Value(&config.Homebrew),
		huh.NewConfirm().
			Title("Enable Snap").
			Description("Generate Snap package?").
			Value(&config.Snap),
		huh.NewConfirm().
			Title("Enable SBOM").
			Description("Generate Software Bill of Materials?").
			Value(&config.SBOM).
			Default(true),
	).
		WithTheme(huh.ThemeCatppuccin()).
		Run()
}

// collectCIConfiguration collects CI/CD configuration
func collectCIConfiguration(config *ProjectConfig) error {
	return huh.NewForm(
		huh.NewGroup("CI/CD Configuration"),
		huh.NewConfirm().
			Title("Generate GitHub Actions").
			Description("Create GitHub Actions workflow?").
			Value(&config.GenerateActions).
			Default(true),
		huh.NewMultiSelect[string]().
			Title("Trigger Actions On").
			Description("When should the workflow run?").
			Options(
				huh.NewOption("🏷️  Version tags (v1.0.0, v2.0.0)", "version tags"),
				huh.NewOption("🔖  All tags", "all tags"),
				huh.NewOption("✋  Manual trigger", "manual"),
				huh.NewOption("🌿  Push to main branch", "main"),
				huh.NewOption("📦  Published release", "release"),
			).
			Validate(func(values []string) error {
				if config.GenerateActions && len(values) == 0 {
					return domain.NewValidationError(
						domain.ErrMissingRequiredField,
						"At least one trigger required",
						"Please select at least one workflow trigger",
					)
				}
				return nil
			}).
			Value(&config.ActionsOn).
			Height(6),
		huh.NewConfirm().
			Title("Enable GoReleaser Pro").
			Description("Use GoReleaser Pro features?").
			Value(&config.ProVersion).
			Default(false),
	).
		WithTheme(huh.ThemeCatppuccin()).
		Run()
}

// createSafeConfig converts legacy config to domain-safe config
func createSafeConfig(legacy *ProjectConfig) (*domain.SafeProjectConfig, error) {
	safeConfig := domain.NewSafeProjectConfig()

	// Map basic information
	safeConfig.ProjectName = legacy.ProjectName
	safeConfig.ProjectDescription = legacy.ProjectDescription
	safeConfig.BinaryName = legacy.BinaryName
	safeConfig.MainPath = legacy.MainPath

	// Convert and validate project type
	projectType := convertToProjectType(legacy.ProjectType)
	if !projectType.IsValid() {
		return nil, domain.NewValidationError(
			domain.ErrInvalidProjectName,
			"Invalid project type",
			fmt.Sprintf("'%s' is not a valid project type", legacy.ProjectType),
		)
	}
	safeConfig.ProjectType = projectType

	// Convert platforms
	platforms := make([]domain.Platform, len(legacy.Platforms))
	for i, p := range legacy.Platforms {
		platform := domain.Platform(p)
		if !platform.IsValid() {
			return nil, domain.NewValidationError(
				domain.ErrInvalidPlatform,
				"Invalid platform",
				fmt.Sprintf("'%s' is not a valid platform", p),
			)
		}
		platforms[i] = platform
	}
	safeConfig.Platforms = platforms

	// Convert architectures
	architectures := make([]domain.Architecture, len(legacy.Architectures))
	for i, a := range legacy.Architectures {
		arch := domain.Architecture(a)
		if !arch.IsValid() {
			return nil, domain.NewValidationError(
				domain.ErrInvalidArchitecture,
				"Invalid architecture",
				fmt.Sprintf("'%s' is not a valid architecture", a),
			)
		}
		architectures[i] = arch
	}
	safeConfig.Architectures = architectures

	// Map build configuration
	safeConfig.CGOEnabled = legacy.CGOEnabled
	safeConfig.LDFlags = legacy.LDFlags
	safeConfig.BuildTags = convertBuildTags(legacy.BuildTags)

	// Convert and validate git provider
	gitProvider := convertToGitProvider(legacy.GitProvider)
	if !gitProvider.IsValid() {
		return nil, domain.NewValidationError(
			domain.ErrInvalidGitProvider,
			"Invalid git provider",
			fmt.Sprintf("'%s' is not a valid git provider", legacy.GitProvider),
		)
	}
	safeConfig.GitProvider = gitProvider

	// Map release configuration
	safeConfig.DockerEnabled = legacy.DockerEnabled
	safeConfig.DockerRegistry = domain.DockerRegistryDockerHub // Default
	safeConfig.DockerImage = legacy.DockerRegistry
	safeConfig.Signing = legacy.Signing
	safeConfig.Homebrew = legacy.Homebrew
	safeConfig.Snap = legacy.Snap
	safeConfig.SBOM = legacy.SBOM

	// Map CI/CD configuration
	safeConfig.GenerateActions = legacy.GenerateActions
	safeConfig.ActionsOn = convertActionTriggers(legacy.ActionsOn)
	safeConfig.ProVersion = legacy.ProVersion

	// Apply defaults and validate
	safeConfig.ApplyDefaults()

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

// generateConfigurationFiles generates configuration files
func generateConfigurationFiles(config *domain.SafeProjectConfig) error {
	// Update state
	config.State = domain.ConfigStateProcessing

	// Generate GoReleaser configuration (placeholder)
	goreleaserYAML := generateGoReleaserConfig(config)
	if err := writeFile(".goreleaser.yaml", goreleaserYAML, 0644); err != nil {
		config.State = domain.ConfigStateInvalid
		return err
	}

	// Generate GitHub Actions workflow (placeholder)
	if config.GenerateActions {
		actionsYAML := generateGitHubActions(config)
		if err := writeFile(".github/workflows/release.yml", actionsYAML, 0644); err != nil {
			config.State = domain.ConfigStateInvalid
			return err
		}
	}

	// Mark as generated
	config.State = domain.ConfigStateGenerated
	return nil
}

// Utility functions for conversion
func convertToProjectType(displayName string) domain.ProjectType {
	switch strings.ToLower(strings.TrimSpace(displayName)) {
	case "cli application", "cli":
		return domain.ProjectTypeCLI
	case "web service", "web":
		return domain.ProjectTypeWeb
	case "library":
		return domain.ProjectTypeLibrary
	case "api service", "api":
		return domain.ProjectTypeAPI
	case "desktop application", "desktop":
		return domain.ProjectTypeDesktop
	default:
		return domain.ProjectTypeCLI
	}
}

func convertToGitProvider(displayName string) domain.GitProvider {
	switch strings.ToLower(strings.TrimSpace(displayName)) {
	case "github":
		return domain.GitProviderGitHub
	case "gitlab":
		return domain.GitProviderGitLab
	case "bitbucket":
		return domain.GitProviderBitbucket
	case "gitea":
		return domain.GitProviderGitea
	case "self-hosted":
		return domain.GitProviderSelfHosted
	default:
		return domain.GitProviderGitHub
	}
}

func convertBuildTags(tags []string) []domain.BuildTag {
	buildTags := make([]domain.BuildTag, len(tags))
	for i, tag := range tags {
		buildTags[i] = domain.BuildTag{Name: tag}
	}
	return buildTags
}

func convertActionTriggers(triggers []string) []domain.ActionTrigger {
	actionTriggers := make([]domain.ActionTrigger, len(triggers))
	for i, trigger := range triggers {
		switch strings.ToLower(strings.TrimSpace(trigger)) {
		case "version tags":
			actionTriggers[i] = domain.ActionTriggerVersionTags
		case "all tags":
			actionTriggers[i] = domain.ActionTriggerAllTags
		case "manual":
			actionTriggers[i] = domain.ActionTriggerManual
		case "main":
			actionTriggers[i] = domain.ActionTriggerMain
		case "release":
			actionTriggers[i] = domain.ActionTriggerRelease
		default:
			actionTriggers[i] = domain.ActionTriggerVersionTags
		}
	}
	return actionTriggers
}

func getProjectTypePlatforms(projectType domain.ProjectType) []string {
	platforms := projectType.RecommendedPlatforms()
	result := make([]string, len(platforms))
	for i, p := range platforms {
		result[i] = string(p)
	}
	return result
}

// Placeholder functions - will be implemented with proper template rendering
func generateGoReleaserConfig(config *domain.SafeProjectConfig) string {
	return fmt.Sprintf(`# GoReleaser Configuration
# Generated by GoReleaser Wizard

project_name: %s
project_type: %s
binary: %s

builds:
  - env:
      - CGO_ENABLED=%t
    goos:
%s    goarch:
%s
`, config.ProjectName, config.ProjectType, config.BinaryName, config.CGOEnabled,
		formatPlatforms(config.Platforms), formatArchitectures(config.Architectures))
}

func generateGitHubActions(config *domain.SafeProjectConfig) string {
	return fmt.Sprintf(`name: Release
on:
%s
jobs:
  goreleaser:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: goreleaser/goreleaser-action@v5
        with:
          version: latest
          args: release --clean
`, domain.GenerateGitHubActionsTriggersYAML(config.ActionsOn))
}

func formatPlatforms(platforms []domain.Platform) string {
	result := make([]string, len(platforms))
	for i, p := range platforms {
		result[i] = fmt.Sprintf("      - %s", string(p))
	}
	return strings.Join(result, "\n")
}

func formatArchitectures(architectures []domain.Architecture) string {
	result := make([]string, len(architectures))
	for i, a := range architectures {
		result[i] = fmt.Sprintf("      - %s", string(a))
	}
	return strings.Join(result, "\n")
}

// File utility functions
func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func writeFile(path string, content string, perm os.FileMode) error {
	// Create directory if needed
	dir := strings.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return domain.NewSystemError(
				domain.ErrDirectoryCreateFailed,
				"Failed to create directory",
				fmt.Sprintf("Cannot create directory %s", dir),
				err,
			).WithContext(dir)
		}
	}

	// Write file
	if err := os.WriteFile(path, []byte(content), perm); err != nil {
		return domain.NewSystemError(
			domain.ErrFileWriteFailed,
			"Failed to write file",
			fmt.Sprintf("Cannot write to %s", path),
			err,
		).WithContext(path)
	}

	return nil
}

// FormValidator for Huh forms using domain types
type FormValidator struct {
	errors map[string]string
}

func NewFormValidator() *FormValidator {
	return &FormValidator{
		errors: make(map[string]string),
	}
}

func (fv *FormValidator) ValidateProjectName() func(string) error {
	return func(value string) error {
		if err := domain.ValidateProjectName(value); err != nil {
			fv.errors["project_name"] = err.Error()
			return err
		}
		delete(fv.errors, "project_name")
		return nil
	}
}

func (fv *FormValidator) ValidateBinaryName() func(string) error {
	return func(value string) error {
		if err := domain.ValidateBinaryName(value); err != nil {
			fv.errors["binary_name"] = err.Error()
			return err
		}
		delete(fv.errors, "binary_name")
		return nil
	}
}

func (fv *FormValidator) ValidateMainPath() func(string) error {
	return func(value string) error {
		if err := domain.ValidateMainPath(value); err != nil {
			fv.errors["main_path"] = err.Error()
			return err
		}
		delete(fv.errors, "main_path")
		return nil
	}
}

func (fv *FormValidator) ValidateProjectDescription() func(string) error {
	return func(value string) error {
		if err := domain.ValidateProjectDescription(value); err != nil {
			fv.errors["project_description"] = err.Error()
			return err
		}
		delete(fv.errors, "project_description")
		return nil
	}
}

func (fv *FormValidator) GetErrors() map[string]string {
	return fv.errors
}

func (fv *FormValidator) HasErrors() bool {
	return len(fv.errors) > 0
}

func (fv *FormValidator) ClearErrors() {
	fv.errors = make(map[string]string)
}