package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
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

// ProjectConfig holds all the configuration options
type ProjectConfig struct {
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
	Signing        bool
	Homebrew       bool
	Snap           bool
	SBOM           bool

	// GitHub Actions
	GenerateActions bool
	ActionsOn       []string

	// Advanced
	ProVersion     bool
	Compression    string
	Hooks          []string
	SkipValidation bool
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Interactive wizard to create GoReleaser configuration",
	Long: `The init command launches an interactive wizard that guides you through
creating a perfect GoReleaser configuration for your project.

It will:
- Detect your project structure
- Ask relevant questions based on your project type
- Generate optimized .goreleaser.yaml
- Optionally create GitHub Actions workflow
- Apply best practices automatically`,
	Run: runInitWizard,
}

func init() {
	initCmd.Flags().Bool("force", false, "overwrite existing configuration")
	initCmd.Flags().Bool("minimal", false, "create minimal configuration")
	initCmd.Flags().Bool("pro", false, "use GoReleaser Pro features")
}

func runInitWizard(cmd *cobra.Command, args []string) {
	// Set up logger
	logger := log.New(os.Stderr)
	if viper.GetBool("debug") {
		logger.SetLevel(log.DebugLevel)
	}

	// Set up panic recovery
	defer HandlePanic("init wizard", logger)

	fmt.Println(titleStyle.Render("🚀 GoReleaser Configuration Wizard"))
	fmt.Println(infoStyle.Render("Let's create the perfect GoReleaser config for your project!\n"))

	// Check if config already exists
	force, _ := cmd.Flags().GetBool("force")
	if !force {
		if err := CheckFileExists(".goreleaser.yaml", false); err == nil {
			// File exists and is accessible
			logger.Warn("Configuration already exists", "file", ".goreleaser.yaml")
			err := NewWizardError(
				ErrConfiguration,
				".goreleaser.yaml already exists",
				"Configuration file found in current directory",
				"Use --force to overwrite existing configuration",
				nil,
			)
			LogAndDisplayError(err, logger)
			return
		}
	}

	config := &ProjectConfig{}

	// Detect project info
	detectProjectInfo(config)

	// Run interactive forms with enhanced error handling
	if err := askBasicInfo(config); err != nil {
		LogAndDisplayError(UserInputError("basic information", err), logger)
		return
	}

	if err := askBuildOptions(config); err != nil {
		LogAndDisplayError(UserInputError("build options", err), logger)
		return
	}

	if err := askReleaseOptions(config); err != nil {
		LogAndDisplayError(UserInputError("release options", err), logger)
		return
	}

	if err := askAdvancedOptions(config); err != nil {
		LogAndDisplayError(UserInputError("advanced options", err), logger)
		return
	}

	// Generate configuration
	fmt.Println("\n" + infoStyle.Render("Generating configuration..."))

	if err := generateGoReleaserConfig(config); err != nil {
		LogAndDisplayError(TemplateError("goreleaser.yaml", err), logger)
		return
	}
	fmt.Println(successStyle.Render("✓ Created .goreleaser.yaml"))

	if config.GenerateActions {
		if err := generateGitHubActions(config); err != nil {
			LogAndDisplayError(TemplateError("github actions workflow", err), logger)
			return
		}
		fmt.Println(successStyle.Render("✓ Created .github/workflows/release.yml"))
	}

	// Show next steps
	fmt.Println("\n" + titleStyle.Render("✨ Setup Complete!"))
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Review the generated .goreleaser.yaml")
	fmt.Println("  2. Run 'goreleaser-wizard validate' to check configuration")
	fmt.Println("  3. Test with 'goreleaser build --snapshot --clean'")
	fmt.Println("  4. Create a git tag and push to trigger release")
	fmt.Println("\nFor more info: https://goreleaser.com")
}

func detectProjectInfo(config *ProjectConfig) {
	// Try to detect project name from go.mod with error handling
	data, err := SafeReadFile("go.mod")
	if err == nil {
		lines := strings.SplitSeq(string(data), "\n")
		for line := range lines {
			if after, ok := strings.CutPrefix(line, "module "); ok {
				module := after
				parts := strings.Split(module, "/")
				if len(parts) > 0 {
					config.ProjectName = parts[len(parts)-1]
				}
				break
			}
		}
	}

	// Detect main.go location with error handling
	if CheckFileExists("main.go", false) == nil {
		config.MainPath = "."
		config.ProjectType = "CLI Application"
	} else if config.ProjectName != "" && CheckFileExists("cmd/"+config.ProjectName+"/main.go", false) == nil {
		config.MainPath = "./cmd/" + config.ProjectName
		config.ProjectType = "CLI Application"
	} else {
		// Look for any main.go in cmd/ with error handling
		if entries, err := os.ReadDir("cmd"); err == nil {
			for _, entry := range entries {
				if entry.IsDir() {
					mainPath := filepath.Join("cmd", entry.Name(), "main.go")
					if CheckFileExists(mainPath, false) == nil {
						config.MainPath = "./cmd/" + entry.Name()
						config.BinaryName = entry.Name()
						config.ProjectType = "CLI Application"
						break
					}
				}
			}
		}
	}

	// Default binary name
	if config.BinaryName == "" && config.ProjectName != "" {
		config.BinaryName = config.ProjectName
	}
}

func askBasicInfo(config *ProjectConfig) error {
	var projectTypes = []string{
		"CLI Application",
		"Web Service",
		"Library with CLI",
		"Multiple Binaries",
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Project Name").
				Description("Name of your project").
				Value(&config.ProjectName).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("project name is required")
					}
					return nil
				}),

			huh.NewInput().
				Title("Project Description").
				Description("Brief description of your project").
				Value(&config.ProjectDescription).
				Placeholder("A fantastic Go application"),

			huh.NewSelect[string]().
				Title("Project Type").
				Description("What kind of project is this?").
				Options(huh.NewOptions(projectTypes...)...).
				Value(&config.ProjectType),

			huh.NewInput().
				Title("Binary Name").
				Description("Name of the compiled binary").
				Value(&config.BinaryName).
				Placeholder(config.ProjectName),

			huh.NewInput().
				Title("Main Package Path").
				Description("Path to main.go (e.g., . or ./cmd/app)").
				Value(&config.MainPath).
				Placeholder("./cmd/"+config.BinaryName),
		).Title("Basic Information"),
	)

	return form.Run()
}

func askBuildOptions(config *ProjectConfig) error {
	var (
		platformOptions = []string{
			"linux",
			"darwin",
			"windows",
			"freebsd",
			"openbsd",
			"netbsd",
			"dragonfly",
			"android",
			"ios",
		}
		archOptions = []string{
			"amd64",
			"arm64",
			"arm",
			"386",
			"ppc64le",
			"s390x",
			"mips",
			"mipsle",
			"mips64",
			"mips64le",
			"riscv64",
			"wasm",
		}
	)

	// Set defaults
	config.Platforms = []string{"linux", "darwin", "windows"}
	config.Architectures = []string{"amd64", "arm64"}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Target Platforms").
				Description("Which operating systems to build for?").
				Options(huh.NewOptions(platformOptions...)...).
				Value(&config.Platforms),

			huh.NewMultiSelect[string]().
				Title("Target Architectures").
				Description("Which CPU architectures to build for?").
				Options(huh.NewOptions(archOptions...)...).
				Value(&config.Architectures),

			huh.NewConfirm().
				Title("Enable CGO?").
				Description("Required for SQLite and some C libraries").
				Value(&config.CGOEnabled).
				Affirmative("Yes").
				Negative("No (recommended)"),

			huh.NewConfirm().
				Title("Embed Version Info?").
				Description("Add version, commit, and date to binary").
				Value(&config.LDFlags).
				Affirmative("Yes (recommended)").
				Negative("No"),
		).Title("Build Options"),
	)

	return form.Run()
}

func askReleaseOptions(config *ProjectConfig) error {
	gitProviders := []string{
		"GitHub",
		"GitLab",
		"Gitea",
		"Local Only",
	}

	config.GitProvider = "GitHub" // default

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Git Provider").
				Description("Where is your repository hosted?").
				Options(huh.NewOptions(gitProviders...)...).
				Value(&config.GitProvider),

			huh.NewConfirm().
				Title("Docker Images?").
				Description("Build and push Docker images").
				Value(&config.DockerEnabled).
				Affirmative("Yes").
				Negative("No"),

			huh.NewConfirm().
				Title("Code Signing?").
				Description("Sign releases with cosign (keyless)").
				Value(&config.Signing).
				Affirmative("Yes").
				Negative("No"),

			huh.NewConfirm().
				Title("Generate SBOM?").
				Description("Software Bill of Materials for security").
				Value(&config.SBOM).
				Affirmative("Yes").
				Negative("No"),
		).Title("Release Options"),
	)

	if err := form.Run(); err != nil {
		return err
	}

	// Ask about Docker registry if Docker is enabled
	if config.DockerEnabled {
		registryForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Docker Registry").
					Description("Registry to push images (e.g., ghcr.io/username)").
					Value(&config.DockerRegistry).
					Placeholder("ghcr.io/" + config.ProjectName),
			),
		)
		if err := registryForm.Run(); err != nil {
			return err
		}
	}

	// Ask about package managers
	if config.GitProvider == "GitHub" {
		pmForm := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Homebrew Tap?").
					Description("Create Homebrew formula for macOS/Linux").
					Value(&config.Homebrew).
					Affirmative("Yes").
					Negative("No"),

				huh.NewConfirm().
					Title("Snap Package?").
					Description("Create Snap package for Linux").
					Value(&config.Snap).
					Affirmative("Yes").
					Negative("No"),
			).Title("Package Managers"),
		)
		if err := pmForm.Run(); err != nil {
			return err
		}
	}

	return nil
}

func askAdvancedOptions(config *ProjectConfig) error {
	compressionOptions := []string{
		"none",
		"gzip",
		"upx (smaller but slower)",
	}

	config.Compression = "gzip" // default

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Generate GitHub Actions?").
				Description("Create workflow for automated releases").
				Value(&config.GenerateActions).
				Affirmative("Yes (recommended)").
				Negative("No"),

			huh.NewSelect[string]().
				Title("Archive Compression").
				Description("Compression for release archives").
				Options(huh.NewOptions(compressionOptions...)...).
				Value(&config.Compression),

			huh.NewConfirm().
				Title("GoReleaser Pro?").
				Description("Use Pro features (requires license)").
				Value(&config.ProVersion).
				Affirmative("Yes").
				Negative("No (free version)"),
		).Title("Advanced Options"),
	)

	if err := form.Run(); err != nil {
		return err
	}

	// Ask about GitHub Actions triggers if enabled
	if config.GenerateActions {
		triggerOptions := []string{
			"On version tags (v*)",
			"On all tags",
			"Manual trigger only",
			"On push to main",
		}

		config.ActionsOn = []string{"On version tags (v*)"}

		triggerForm := huh.NewForm(
			huh.NewGroup(
				huh.NewMultiSelect[string]().
					Title("GitHub Actions Triggers").
					Description("When should releases be created?").
					Options(huh.NewOptions(triggerOptions...)...).
					Value(&config.ActionsOn),
			),
		)
		if err := triggerForm.Run(); err != nil {
			return err
		}
	}

	return nil
}
