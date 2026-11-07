package main

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// validateGoReleaserConfig validates GoReleaser configuration.
func validateGoReleaserConfig(results *ValidationResults) error {
	configPath := goreleaserConfigFilename

	// Check if config exists
	exists, err := fileSystemRepo.FileExists(context.Background(), configPath)
	if err != nil {
		results.Errors = append(results.Errors,
			NewSystemError(
				"FILE_READ_FAILED",
				"Failed to check configuration file",
				"Cannot access "+configPath,
				err,
			).WithContext(configPath))

		return nil
	}

	results.ConfigExists = exists
	if !exists {
		results.Errors = append(results.Errors,
			NewSystemError(
				"FILE_NOT_FOUND",
				"Configuration file not found",
				configPath+" does not exist",
				nil,
			).WithContext(configPath))
		results.Recommendations = append(results.Recommendations,
			"Run 'goreleaser-wizard init' to create configuration")

		return nil
	}

	// Check if config is valid YAML
	if err := validateYAML(configPath, results); err != nil {
		results.ConfigValid = false

		return nil
	}

	// Detect deprecated keys with migration hints (works without the
	// goreleaser binary, unlike the check below).
	validateDeprecations(results, configPath)

	// Check if goreleaser is available
	if _, err := exec.LookPath("goreleaser"); err == nil {
		results.GoReleaserFound = true
		// Run goreleaser check
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		cmd := exec.CommandContext(ctx, "goreleaser", "check", "--config", configPath)
		if output, err := cmd.CombinedOutput(); err != nil {
			results.ConfigValid = false
			results.Errors = append(results.Errors,
				NewValidationError(
					"INVALID_CONFIGURATION",
					"GoReleaser configuration is invalid",
					string(output),
				).WithContext(configPath))
		} else {
			results.ConfigValid = true
		}
	} else {
		results.Warnings = append(results.Warnings,
			NewValidationError(
				"EXTERNAL_TOOL_NOT_FOUND",
				"GoReleaser not found",
				"Install GoReleaser for configuration validation",
			))
		// Assume valid if no goreleaser available
		results.ConfigValid = true
	}

	return nil
}

// validateGitHubActions validates GitHub Actions workflow.
func validateGitHubActions(results *ValidationResults) {
	workflowPath := releaseWorkflowTargetPath

	// Check if workflow exists
	exists, err := fileSystemRepo.FileExists(context.Background(), workflowPath)
	if err != nil {
		results.Errors = append(results.Errors,
			NewSystemError(
				"FILE_READ_FAILED",
				"Failed to check workflow file",
				"Cannot access "+workflowPath,
				err,
			).WithContext(workflowPath))

		return
	}

	results.ActionsExists = exists
	if !exists {
		results.Recommendations = append(results.Recommendations,
			"Consider adding GitHub Actions workflow for automated releases")

		return
	}

	// Validate workflow content
	if err := validateWorkflowContent(workflowPath, results); err != nil {
		results.ActionsValid = false

		if domainErr, ok := errors.AsType[*DomainError](err); ok {
			results.Errors = append(results.Errors, domainErr)
		}

		return
	}

	results.ActionsValid = true
}

// validateProjectStructure validates project structure.
func validateProjectStructure(results *ValidationResults) {
	// Check for go.mod
	if exists, err := fileSystemRepo.FileExists(context.Background(), "go.mod"); err != nil {
		results.Errors = append(results.Errors,
			NewSystemError(
				"FILE_READ_FAILED",
				"Failed to check go.mod",
				"Cannot access go.mod file",
				err,
			).WithContext("go.mod"))

		return
	} else if !exists {
		results.Errors = append(results.Errors,
			NewValidationError(
				"MISSING_DEPENDENCY",
				"go.mod not found",
				"Run 'go mod init' to initialize Go module",
			).WithContext("go.mod"))
		results.ProjectValid = false

		return
	}

	// Check for main.go directory
	mainPath := "cmd"
	if exists, err := fileSystemRepo.DirExists(context.Background(), mainPath); err != nil {
		results.Errors = append(results.Errors,
			NewSystemError(
				"FILE_READ_FAILED",
				"Failed to check main directory",
				"Cannot access "+mainPath,
				err,
			).WithContext(mainPath))

		return
	} else if !exists {
		// Try common alternatives
		alternatives := []string{"main", "src"}
		for _, alt := range alternatives {
			altExists, err := fileSystemRepo.DirExists(context.Background(), alt)
			if err != nil {
				// A failed stat is not "missing": report it and keep checking
				// the remaining alternatives instead of silently skipping.
				results.Warnings = append(results.Warnings,
					NewSystemError(
						"FILE_READ_FAILED",
						"Failed to check alternative main directory",
						"Cannot access "+alt,
						err,
					).WithContext(alt))

				continue
			}

			if altExists {
				results.Warnings = append(results.Warnings,
					NewValidationError(
						"INVALID_PROJECT_STRUCTURE",
						"Non-standard main directory",
						fmt.Sprintf("Found %s instead of %s", alt, mainPath),
					).WithContext(alt))
				results.ProjectValid = true

				return
			}
		}

		results.Errors = append(results.Errors,
			NewValidationError(
				"INVALID_PROJECT_STRUCTURE",
				"Main directory not found",
				"Create cmd directory with main package",
			).WithContext(mainPath))
		results.ProjectValid = false

		return
	}

	results.ProjectValid = true
}

// validateYAML validates YAML file.
func validateYAML(filePath string, results *ValidationResults) error {
	content, err := fileSystemRepo.ReadFile(context.Background(), filePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	if len(lines) == 0 {
		return addValidationError(results, filePath, "Empty YAML file", filePath+" is empty")
	}

	if !hasValidYAMLStructure(lines) {
		return addValidationError(
			results,
			filePath,
			"Invalid YAML structure",
			filePath+" does not appear to be valid YAML",
		)
	}

	return nil
}

func hasValidYAMLStructure(lines []string) bool {
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, ":") && !strings.HasPrefix(trimmed, "#") {
			return true
		}
	}

	return false
}

func addValidationError(results *ValidationResults, filePath, title, details string) error {
	results.Errors = append(results.Errors,
		NewValidationError(
			"INVALID_FILE_FORMAT",
			title,
			details,
		).WithContext(filePath))

	return nil
}
