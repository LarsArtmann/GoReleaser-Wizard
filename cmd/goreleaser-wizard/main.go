package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Build-time variables set by GoReleaser
	version        = "dev"
	commit         = "none"
	date           = "unknown"
	builtBy        = "unknown"
	gitDescription = ""
	gitState       = ""

	cfgFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goreleaser-wizard",
	Short: "Interactive setup wizard for GoReleaser",
	Long: `GoReleaser Wizard is an interactive CLI tool that helps you create
perfect GoReleaser configurations for your Go projects.

It guides you through the configuration process with smart defaults
and best practices, generating both .goreleaser.yaml and GitHub Actions
workflows tailored to your project's needs.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	// Set up logger for error handling
	logger := log.New(os.Stderr)
	if viper.GetBool("debug") {
		logger.SetLevel(log.DebugLevel)
	}

	// Set up panic recovery
	defer HandlePanic("command execution", logger)

	if err := rootCmd.Execute(); err != nil {
		LogAndDisplayError(err, logger)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.goreleaser-wizard.yaml)")
	rootCmd.PersistentFlags().Bool("no-color", false, "disable color output")
	rootCmd.PersistentFlags().Bool("debug", false, "enable debug output")

	// Bind flags to viper
	viper.BindPFlag("no-color", rootCmd.PersistentFlags().Lookup("no-color"))
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	// Add commands
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(generateCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	logger := log.New(os.Stderr)
	if viper.GetBool("debug") {
		logger.SetLevel(log.DebugLevel)
	}

	// Set up panic recovery for config initialization
	defer HandlePanic("config initialization", logger)

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
		// Validate the config file exists and is readable
		if err := CheckFileExists(cfgFile, true); err != nil {
			LogAndDisplayError(
				NewWizardError(
					ErrConfiguration,
					"Custom config file not found",
					fmt.Sprintf("Could not access config file: %s", cfgFile),
					"Check the file path and permissions",
					err,
				),
				logger,
			)
			return
		}
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			LogAndDisplayError(
				NewWizardError(
					ErrPermission,
					"Unable to determine user home directory",
					"System could not determine the user's home directory",
					"Check your system's home directory configuration",
					err,
				),
				logger,
			)
			return
		}

		// Search config in home directory with name ".goreleaser-wizard" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".goreleaser-wizard")
	}

	viper.SetEnvPrefix("GORELEASER_WIZARD")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		// Only log if it's not a "file not found" error for optional config
		if cfgFile != "" || !os.IsNotExist(err) {
			logger.Warn("Config file error", "error", err, "file", viper.ConfigFileUsed())
		}
	} else if viper.GetBool("debug") {
		logger.Info("Using config file", "file", viper.ConfigFileUsed())
	}
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("GoReleaser Wizard %s\n", version)
		fmt.Printf("  Build Date: %s\n", date)
		fmt.Printf("  Git Commit: %s\n", commit)
		fmt.Printf("  Built By: %s\n", builtBy)
		if gitState != "" {
			fmt.Printf("  Git State: %s\n", gitState)
		}
		if gitDescription != "" {
			fmt.Printf("  Git Summary: %s\n", gitDescription)
		}
	},
}

func main() {
	// Set up global panic recovery
	logger := log.New(os.Stderr)
	defer HandlePanic("main", logger)

	Execute()
}
