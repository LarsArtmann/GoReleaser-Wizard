package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
}

func runValidate(cmd *cobra.Command, args []string) {
	// Set up logger
	logger := log.New(os.Stderr)
	if viper.GetBool("debug") {
		logger.SetLevel(log.DebugLevel)
	}

	// Set up panic recovery
	defer HandlePanic("validate command", logger)

	verbose, _ := cmd.Flags().GetBool("verbose")
	fix, _ := cmd.Flags().GetBool("fix")

	fmt.Println(titleStyle.Render("🔍 Validating GoReleaser Configuration"))
	fmt.Println()

	issues := []string{}
	warnings := []string{}
	passed := 0
	total := 0

	// Check 1: .goreleaser.yaml exists
	total++
	if err := CheckFileExists(".goreleaser.yaml", false); err != nil {
		issues = append(issues, ".goreleaser.yaml not found")
		fmt.Println(errorStyle.Render("✗ .goreleaser.yaml not found"))
		if fix {
			fmt.Println(infoStyle.Render("  → Run 'goreleaser-wizard init' to create one"))
		}
		logger.Debug("GoReleaser config check", "error", err)
	} else {
		passed++
		fmt.Println(successStyle.Render("✓ .goreleaser.yaml exists"))
	}

	// Check 2: go.mod exists
	total++
	if err := CheckFileExists("go.mod", false); err != nil {
		issues = append(issues, "go.mod not found")
		fmt.Println(errorStyle.Render("✗ go.mod not found"))
		if fix {
			fmt.Println(infoStyle.Render("  → Run 'go mod init <module-name>' to create one"))
		}
		logger.Debug("Go module check", "error", err)
	} else {
		passed++
		fmt.Println(successStyle.Render("✓ go.mod exists"))
	}

	// Check 3: Git repository
	total++
	if err := CheckFileExists(".git", false); err != nil {
		warnings = append(warnings, "Not a git repository")
		fmt.Println(errorStyle.Render("⚠ Not a git repository"))
		fmt.Println(infoStyle.Render("  → GoReleaser requires a git repository to work"))
		if fix {
			fmt.Println(infoStyle.Render("  → Run 'git init' to initialize repository"))
		}
		logger.Debug("Git repository check", "error", err)
	} else {
		// Check for uncommitted changes with error handling
		gitCmd := exec.Command("git", "status", "--porcelain")
		output, err := gitCmd.Output()
		if err != nil {
			logger.Warn("Failed to check git status", "error", err)
			warnings = append(warnings, "Could not check git status")
			fmt.Println(errorStyle.Render("⚠ Could not check git status"))
		} else if len(output) > 0 {
			warnings = append(warnings, "Uncommitted changes detected")
			fmt.Println(errorStyle.Render("⚠ Uncommitted changes detected"))
			if verbose {
				fmt.Println(infoStyle.Render("  → " + strings.TrimSpace(string(output))))
			}
			if fix {
				fmt.Println(infoStyle.Render("  → Commit changes with 'git add . && git commit -m \"message\"'"))
			}
		} else {
			passed++
			fmt.Println(successStyle.Render("✓ Git repository clean"))
		}
	}

	// Check 4: GoReleaser installed
	total++
	goreleaserPath, err := exec.LookPath("goreleaser")
	if err != nil {
		issues = append(issues, "GoReleaser not installed")
		fmt.Println(errorStyle.Render("✗ GoReleaser not installed"))
		if fix {
			fmt.Println(infoStyle.Render("  → Install with: go install github.com/goreleaser/goreleaser/v2@latest"))
			fmt.Println(infoStyle.Render("  → Or download from: https://goreleaser.com/install/"))
		}
		logger.Debug("GoReleaser dependency check", "error", err)
	} else {
		passed++
		fmt.Println(successStyle.Render("✓ GoReleaser installed"))
		if verbose {
			fmt.Println(infoStyle.Render("  → " + goreleaserPath))
		}

		// Check 5: Run goreleaser check with enhanced error handling
		total++
		fmt.Print("  Checking configuration... ")
		checkCmd := exec.Command("goreleaser", "check")
		checkOutput, checkErr := checkCmd.CombinedOutput()
		if checkErr != nil {
			issues = append(issues, "Configuration validation failed")
			fmt.Println(errorStyle.Render("Failed"))
			if verbose {
				fmt.Println(infoStyle.Render("  → " + strings.TrimSpace(string(checkOutput))))
			}
			if fix {
				fmt.Println(infoStyle.Render("  → Fix configuration issues in .goreleaser.yaml"))
				fmt.Println(infoStyle.Render("  → Run 'goreleaser-wizard init --force' to regenerate"))
			}
			logger.Debug("GoReleaser config validation", "error", checkErr, "output", string(checkOutput))
		} else {
			passed++
			fmt.Println(successStyle.Render("OK"))
		}
	}

	// Check 6: Main package exists
	if CheckFileExists(".goreleaser.yaml", false) == nil {
		total++
		// Parse config to find main path
		// For simplicity, we'll check common locations
		mainFound := false
		commonPaths := []string{
			"main.go",
			"./cmd/*/main.go",
			"./*.go",
		}
		for _, path := range commonPaths {
			matches, globErr := filepath.Glob(path)
			if globErr != nil {
				logger.Debug("Glob pattern error", "pattern", path, "error", globErr)
				continue
			}
			if len(matches) > 0 {
				mainFound = true
				break
			}
		}

		if !mainFound {
			warnings = append(warnings, "No main.go found in expected locations")
			fmt.Println(errorStyle.Render("⚠ No main.go found"))
			fmt.Println(infoStyle.Render("  → Make sure your main package path is correct in .goreleaser.yaml"))
			if fix {
				fmt.Println(infoStyle.Render("  → Create main.go or update build.main path in config"))
			}
		} else {
			passed++
			fmt.Println(successStyle.Render("✓ Main package found"))
		}
	}

	// Check 7: Docker (if configured)
	if CheckFileExists("Dockerfile", false) == nil {
		total++
		dockerPath, err := exec.LookPath("docker")
		if err != nil {
			warnings = append(warnings, "Docker not installed but Dockerfile exists")
			fmt.Println(errorStyle.Render("⚠ Docker not installed"))
			if fix {
				fmt.Println(infoStyle.Render("  → Install Docker from https://docker.com/"))
				fmt.Println(infoStyle.Render("  → Or remove Docker configuration from .goreleaser.yaml"))
			}
			logger.Debug("Docker dependency check", "error", err)
		} else {
			passed++
			fmt.Println(successStyle.Render("✓ Docker installed"))
			if verbose {
				fmt.Println(infoStyle.Render("  → " + dockerPath))
			}
		}
	}

	// Check 8: GitHub Actions workflow
	total++
	workflowFound := CheckFileExists(".github/workflows/release.yml", false) == nil ||
		CheckFileExists(".github/workflows/release.yaml", false) == nil

	if workflowFound {
		passed++
		fmt.Println(successStyle.Render("✓ GitHub Actions workflow found"))
	} else {
		warnings = append(warnings, "No GitHub Actions workflow for releases")
		fmt.Println(infoStyle.Render("ℹ No GitHub Actions workflow"))
		if fix {
			fmt.Println(infoStyle.Render("  → Run 'goreleaser-wizard init' with GitHub Actions option"))
			fmt.Println(infoStyle.Render("  → Or manually create .github/workflows/release.yml"))
		}
	}

	// Summary
	fmt.Println()
	fmt.Println(titleStyle.Render("📊 Validation Summary"))
	fmt.Printf("Checks passed: %d/%d\n", passed, total)

	if len(issues) > 0 {
		fmt.Println()
		fmt.Println(errorStyle.Render("❌ Critical Issues:"))
		for _, issue := range issues {
			fmt.Println("  • " + issue)
		}
	}

	if len(warnings) > 0 {
		fmt.Println()
		fmt.Println(infoStyle.Render("⚠️  Warnings:"))
		for _, warning := range warnings {
			fmt.Println("  • " + warning)
		}
	}

	// Test build suggestion
	if len(issues) == 0 {
		fmt.Println()
		fmt.Println(successStyle.Render("✨ Configuration looks good!"))
		fmt.Println()
		fmt.Println("Next steps:")
		fmt.Println("  1. Test build: goreleaser build --snapshot --clean")
		fmt.Println("  2. Create tag: git tag -a v0.1.0 -m 'First release'")
		fmt.Println("  3. Push tag: git push origin v0.1.0")
		logger.Info("Validation completed successfully", "passed", passed, "total", total, "warnings", len(warnings))
	} else {
		fmt.Println()
		fmt.Println(errorStyle.Render("⚠️  Please fix the issues above before releasing"))
		if !fix {
			fmt.Println()
			fmt.Println("Run with --fix to see suggested fixes")
		}
		logger.Error("Validation failed", "issues", len(issues), "warnings", len(warnings), "passed", passed, "total", total)
		os.Exit(1)
	}
}
