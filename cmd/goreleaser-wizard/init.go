package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"charm.land/log/v2"
	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
	"github.com/LarsArtmann/GoReleaser-Wizard/internal/git"
	"github.com/spf13/cobra"
)

// runInitWizard runs the init command.
func runInitWizard(cmd *cobra.Command, _ []string) {
	defer recoverFromPanic("init wizard")

	applyGitHubOverrides(cmd)

	force := getBoolFlag(cmd, "force")
	interactive := getBoolFlag(cmd, "interactive")

	printCommandHeader("🧙 Initializing GoReleaser Configuration")

	// If interactive mode requested but not in a terminal, error with helpful message
	if interactive && !IsTerminal() {
		err := domain.NewValidationError(
			domain.ErrInvalidConfiguration,
			"Interactive mode requires a terminal",
			NonInteractiveHelp,
		).WithContext("init wizard")
		displayError(err)

		return
	}

	// Detect project information
	config := &domain.SafeProjectConfig{}

	err := detectProjectInfo(config)
	if err != nil {
		displayError(err)

		return
	}

	// Apply defaults to ensure complete configuration
	config.ApplyDefaults()

	// If interactive mode, run the TUI wizard
	if interactive {
		err := RunTUIWizard(config)
		if err != nil {
			displayError(err)

			return
		}
	}

	// Create and execute workflow
	logger := log.New(os.Stderr)
	if !ExecuteWorkflow(WorkflowTypeFullWizard, config, force, logger) {
		return
	}

	fmt.Println()
	fmt.Println(successStyle.Render("🎉 GoReleaser configuration initialized successfully!"))
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  • Review the generated .goreleaser.yaml")
	fmt.Println("  • Run 'goreleaser-wizard validate' to check the configuration")

	for _, tip := range providerWorkflowTips(config.CIProvider, false) {
		fmt.Println(tip)
	}

	fmt.Println("  • Commit the configuration to your repository")
}

// fileReadError creates a standardized file read error.
func fileReadError(context, message, details string, originalErr error) error {
	return domain.NewSystemError(
		domain.ErrFileReadFailed,
		message,
		details,
		originalErr,
	).WithContext(context)
}

// detectProjectInfo detects project information from the current directory.
func detectProjectInfo(config *domain.SafeProjectConfig) error {
	// Get current working directory
	wd, err := os.Getwd()
	if err != nil {
		return fileReadError("detect_project", "Failed to get working directory",
			"Could not determine current directory", err)
	}

	// Check for go.mod file
	goModPath := filepath.Join(wd, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		return domain.NewValidationError(
			domain.ErrFileNotFound,
			"Not a Go project",
			"go.mod file not found. Please run this command in a Go project directory.",
		).WithContext("detect_project")
	}

	// Read go.mod to get module name
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return fileReadError(goModPath, "Failed to read go.mod",
			"Could not read go.mod file", err)
	}

	// Parse module name from go.mod
	scanner := bufio.NewScanner(strings.NewReader(string(content)))

	var moduleName string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "module ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				moduleName = parts[1]
			}

			break
		}
	}

	if moduleName == "" {
		return domain.NewValidationError(
			domain.ErrInvalidFileFormat,
			fmt.Sprintf("Invalid go.mod format (moduleName=%q)", moduleName),
			"Could not find module declaration in go.mod",
		).WithContext(goModPath)
	}

	// Extract project name from module path
	parts := strings.Split(moduleName, "/")
	projectName := parts[len(parts)-1]

	// Detect main path
	mainPath, binaryName, projectType, err := detectMainStructure(wd)
	if err != nil {
		return err
	}

	// Set configuration values
	config.ProjectName = projectName
	config.MainPath = mainPath
	// Use project name as binary name when main.go is in root (not temp dir name)
	if mainPath == "." {
		config.BinaryName = projectName
	} else {
		config.BinaryName = binaryName
	}

	config.ProjectType = domain.ProjectType(projectType)
	config.ProjectDescription = git.GetGitHubRepoDescription()

	return nil
}

// detectMainStructure detects the main.go structure and determines project type.
func detectMainStructure(wd string) (string, string, string, error) {
	const cli = "cli"

	// Check for main.go in root
	if _, err := os.Stat(filepath.Join(wd, "main.go")); err == nil {
		return ".", filepath.Base(wd), cli, nil
	}

	// Check for cmd/*/main.go structure
	if mainPath, binaryName, ok := detectCmdStructure(wd); ok {
		return mainPath, binaryName, cli, nil
	}

	// Check for other common structures
	return detectAlternativeStructure(wd)
}

// detectCmdStructure checks for cmd/*/main.go pattern and returns the best match.
func detectCmdStructure(wd string) (string, string, bool) {
	cmdDir := filepath.Join(wd, "cmd")

	entries, err := os.ReadDir(cmdDir)
	if err != nil {
		return "", "", false
	}

	// Collect all directories with main.go
	var cmdDirs []string

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		mainFile := filepath.Join(cmdDir, entry.Name(), "main.go")
		if _, err := os.Stat(mainFile); err == nil {
			cmdDirs = append(cmdDirs, entry.Name())
		}
	}

	if len(cmdDirs) == 0 {
		return "", "", false
	}

	// Single binary: use its directory name
	if len(cmdDirs) == 1 {
		return "./cmd/" + cmdDirs[0], cmdDirs[0], true
	}

	// Multiple binaries: prefer the one matching project name or first alphabetically
	selectedBinary := selectBestBinary(cmdDirs, filepath.Base(wd))

	return "./cmd/" + selectedBinary, selectedBinary, true
}

// selectBestBinary chooses the best binary from multiple cmd directories.
func selectBestBinary(cmdDirs []string, projectName string) string {
	binaryName := cmdDirs[0]

	for _, dir := range cmdDirs {
		if projectName == dir {
			return dir
		}

		if dir < binaryName {
			binaryName = dir
		}
	}

	return binaryName
}

// globProjects resolves a glob pattern under wd. It is a variable so tests can
// inject malformed-pattern failures without crafting filesystem states.
//
//nolint:gochecknoglobals // test seam; swapped only inside tests with defer-restore
var globProjects = func(wd, pattern string) ([]string, error) {
	return filepath.Glob(filepath.Join(wd, pattern))
}

// detectAlternativeStructure checks for alternative project structures.
func detectAlternativeStructure(wd string) (string, string, string, error) {
	const cli = "cli"

	patterns := []struct {
		path     string
		binName  string
		projType string
	}{
		{"src/*/main.go", "", "library"},
		{"pkg/*/main.go", "", "library"},
		{"*.go", filepath.Base(wd), cli},
	}

	for _, pattern := range patterns {
		matches, err := globProjects(wd, pattern.path)
		if err != nil {
			return "", "", "", domain.NewSystemError(
				domain.ErrInvalidProjectStructure,
				fmt.Sprintf("Failed to scan project structure (pattern=%q)", pattern.path),
				"The file-matching pattern used to detect the project layout is malformed",
				err,
			).WithContext(wd)
		}

		if len(matches) == 0 {
			continue
		}

		if pattern.binName == "" {
			// Extract binary name from path
			selectedBinary := extractBinaryFromPath(matches[0])

			return pattern.path, selectedBinary, pattern.projType, nil
		}

		return pattern.path, pattern.binName, pattern.projType, nil
	}

	// Default fallback
	return ".", filepath.Base(wd), cli, nil
}

// extractBinaryFromPath extracts the binary name from a matched file path.
// For src/<bin>/main.go layouts the binary is the directory holding main.go,
// matching detectCmdStructure's cmd/<bin>/main.go convention.
func extractBinaryFromPath(matchPath string) string {
	parts := strings.Split(matchPath, string(filepath.Separator))
	for i, part := range parts {
		if (part == "src" || part == "pkg") && i+1 < len(parts) {
			return parts[i+1]
		}
	}

	return strings.TrimSuffix(filepath.Base(matchPath), ".go")
}

func init() {
	configureInitCommand(initCmd)
}
