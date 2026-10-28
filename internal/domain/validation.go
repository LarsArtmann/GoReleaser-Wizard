package domain

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

// ValidationUseCase is the implementation of the validation use case.
type ValidationUseCase struct {
	logger Logger
	repo   FileSystemRepository
}

// NewValidationUseCase creates a new validation use case.
func NewValidationUseCase(logger Logger, repo FileSystemRepository) *ValidationUseCase {
	return &ValidationUseCase{
		logger: logger,
		repo:   repo,
	}
}

// addValidationError adds an error to the result if err is not nil.
func addValidationError(result *ValidationResult, err error) {
	if err != nil {
		if domainErr, ok := errors.AsType[*DomainError](err); ok {
			result.Errors = append(result.Errors, domainErr)
		}

		result.IsValid = false
	}
}

// ValidateConfiguration performs comprehensive validation of project configuration.
func (vu *ValidationUseCase) ValidateConfiguration(
	ctx context.Context,
	config *SafeProjectConfig,
) (*ValidationResult, error) {
	vu.logger.DebugContext(ctx, "Starting comprehensive configuration validation")

	result := &ValidationResult{
		IsValid:  true,
		Errors:   []*DomainError{},
		Warnings: []*DomainError{},
	}

	addValidationError(result, vu.validateBasicFields(config))
	addValidationError(result, vu.validateTypes(config))
	addValidationError(result, vu.validatePlatformArchCompatibility(config))
	addValidationError(result, vu.validateBusinessRules(config))
	addValidationError(result, vu.validateSecurity(config))

	// Step 6: Generate warnings
	vu.generateWarnings(config, result)

	vu.logger.DebugContext(
		ctx,
		"Validation completed",
		"valid",
		result.IsValid,
		"errors",
		len(result.Errors),
		"warnings",
		len(result.Warnings),
	)

	return result, nil
}

// validateBasicFields validates basic required fields.
func (vu *ValidationUseCase) validateBasicFields(
	config *SafeProjectConfig,
) *DomainError {
	// Project name validation
	err := ValidateProjectName(config.ProjectName)
	if err != nil {
		return NewValidationError(
			ErrInvalidProjectName,
			"Project name validation failed",
			err.Error(),
		).WithContext("project_name")
	}

	// Binary name validation
	err = ValidateBinaryName(config.BinaryName)
	if err != nil {
		return NewValidationError(
			ErrInvalidBinaryName,
			"Binary name validation failed",
			err.Error(),
		).WithContext("binary_name")
	}

	// Main path validation
	err = ValidateMainPath(config.MainPath)
	if err != nil {
		return NewValidationError(
			ErrInvalidMainPath,
			"Main path validation failed",
			err.Error(),
		).WithContext("main_path")
	}

	// Project description validation (optional)
	if config.ProjectDescription != "" {
		err = ValidateProjectDescription(config.ProjectDescription)
		if err != nil {
			return NewValidationError(
				ErrInvalidProjectDescription,
				"Project description validation failed",
				err.Error(),
			).WithContext("project_description")
		}
	}

	return nil
}

// validateTypes validates enum types.
func (vu *ValidationUseCase) validateTypes(
	config *SafeProjectConfig,
) *DomainError {
	// Project type validation
	if !config.ProjectType.IsValid() {
		return NewValidationError(
			ErrInvalidProjectName,
			"Invalid project type",
			fmt.Sprintf("'%s' is not a supported project type", config.ProjectType),
		).WithContext("project_type")
	}

	// Platform validation
	err := ValidatePlatforms(config.Platforms)
	if err != nil {
		return NewValidationError(
			ErrInvalidPlatform,
			"Platform validation failed",
			err.Error(),
		).WithContext("platforms")
	}

	// Architecture validation
	err = ValidateArchitectures(config.Architectures)
	if err != nil {
		return NewValidationError(
			ErrInvalidArchitecture,
			"Architecture validation failed",
			err.Error(),
		).WithContext("architectures")
	}

	// Git provider validation
	err = ValidateGitProvider(config.GitProvider)
	if err != nil {
		return NewValidationError(
			ErrInvalidGitProvider,
			"Git provider validation failed",
			err.Error(),
		).WithContext("git_provider")
	}

	// Docker registry validation
	err = ValidateDockerRegistry(config.DockerRegistry)
	if err != nil {
		return NewValidationError(
			ErrInvalidDockerRegistry,
			"Docker registry validation failed",
			err.Error(),
		).WithContext("docker_registry")
	}

	// Action triggers validation
	err = ValidateActionTriggers(config.ActionsOn)
	if err != nil {
		return NewValidationError(
			ErrInvalidActionTrigger,
			"Action trigger validation failed",
			err.Error(),
		).WithContext("actions_on")
	}

	// Build tags validation
	if len(config.BuildTags) > 0 {
		err = ValidateBuildTags(config.BuildTags)
		if err != nil {
			return NewValidationError(
				ErrInvalidBuildTag,
				"Build tags validation failed",
				err.Error(),
			).WithContext("build_tags")
		}
	}

	// Configuration state validation
	err = ValidateConfigState(config.State)
	if err != nil {
		return NewValidationError(
			ErrInvalidConfigState,
			"Configuration state validation failed",
			err.Error(),
		).WithContext("state")
	}

	return nil
}

// validatePlatformArchCompatibility validates platform-architecture compatibility.
func (vu *ValidationUseCase) validatePlatformArchCompatibility(
	config *SafeProjectConfig,
) *DomainError {
	err := ValidatePlatformArchCompatibility(config.Platforms, config.Architectures)
	if err != nil {
		return NewValidationError(
			ErrPlatformArchMismatch,
			"Platform architecture compatibility failed",
			err.Error(),
		).WithContext("platforms_architectures")
	}

	return nil
}

// validateBusinessRules validates domain business rules.
func (vu *ValidationUseCase) validateBusinessRules(
	config *SafeProjectConfig,
) *DomainError {
	// Docker support rule
	if config.GetDockerEnabled() && !config.ProjectType.DockerSupported() {
		return NewConfigurationError(
			ErrDockerNotSupported,
			"Docker not supported for project type",
			fmt.Sprintf("Project type %s does not support Docker", config.ProjectType),
		).WithContext("docker_enabled")
	}

	// Main path requirement rule
	if config.ProjectType.RequiresMainPath() && config.MainPath == "" {
		return NewConfigurationError(
			ErrMainPathRequired,
			"Main path required for project type",
			fmt.Sprintf("Project type %s requires a main path", config.ProjectType),
		).WithContext("main_path")
	}

	// State transition rule
	if !config.State.AllowsGeneration() && config.GetGenerateActions() {
		return NewConfigurationError(
			ErrInvalidStateTransition,
			"State transition invalid",
			fmt.Sprintf("Configuration in state '%s' cannot generate actions", config.State),
		).WithContext("generate_actions")
	}

	// Docker registry URL validation
	if config.GetDockerEnabled() {
		err := ValidateDockerRegistryURL(config.DockerRegistry, config.DockerImage)
		if err != nil {
			return NewValidationError(
				ErrInvalidURLPattern,
				"Docker registry URL validation failed",
				err.Error(),
			).WithContext("docker_image")
		}
	}

	return nil
}

// validateSecurity performs security-focused validation.
func (vu *ValidationUseCase) validateSecurity(
	config *SafeProjectConfig,
) *DomainError {
	// Check for potential path traversal in main path
	if containsPathTraversal(config.MainPath) {
		return NewBusinessRuleError(
			ErrInvalidCharacters,
			"Path traversal detected",
			"Main path contains potentially dangerous path traversal sequences",
		).WithContext("main_path")
	}

	// Check for shell metacharacters in binary name
	if containsShellMetacharacters(config.BinaryName) {
		return NewBusinessRuleError(
			ErrInvalidCharacters,
			"Shell metacharacters detected",
			"Binary name contains potentially dangerous shell metacharacters",
		).WithContext("binary_name")
	}

	// Check Docker image name for security issues
	if config.DockerImage != "" {
		if containsURLInjection(config.DockerImage) {
			return NewBusinessRuleError(
				ErrInvalidCharacters,
				"URL injection detected",
				"Docker image name contains potentially dangerous URL injection sequences",
			).WithContext("docker_image")
		}
	}

	return nil
}

// generateWarnings generates validation warnings.
func (vu *ValidationUseCase) generateWarnings(
	config *SafeProjectConfig,
	result *ValidationResult,
) {
	// Warning for single platform
	if len(config.Platforms) == 1 {
		vu.addWarning(result, NewBusinessRuleError(
			ErrMissingRequiredField,
			"Single platform configuration",
			"Consider targeting multiple platforms for broader compatibility",
		).WithContext("platforms"))
	}

	// Warning for missing Docker image name
	if config.GetDockerEnabled() && config.DockerImage == "" {
		vu.addWarning(result, NewConfigurationError(
			ErrMissingRequiredField,
			"Missing Docker image name",
			"Docker is enabled but no image name is specified",
		).WithContext("docker_image"))
	}

	// Warning for mismatched CGO setting
	if config.CGOStatus.ToBool() != config.ProjectType.DefaultCGOEnabled() {
		vu.addWarning(result, NewConfigurationError(
			ErrInvalidStateTransition,
			"CGO setting mismatched",
			"CGO setting differs from project type default",
		).WithContext("cgo_enabled"))
	}

	// Warning for missing version information
	if !config.LDFlags {
		vu.addWarning(result, NewConfigurationError(
			ErrMissingRequiredField,
			"LD flags disabled",
			"Version information injection is disabled",
		).WithContext("ldflags"))
	}
}

// addWarning adds a warning to the validation result.
func (vu *ValidationUseCase) addWarning(result *ValidationResult, warning *DomainError) {
	result.Warnings = append(result.Warnings, warning)
}

// ValidateProjectStructure validates project directory structure.
func (vu *ValidationUseCase) ValidateProjectStructure(
	ctx context.Context,
	projectPath string,
) (*ProjectValidationResult, error) {
	vu.logger.DebugContext(ctx, "Validating project structure", "path", projectPath)

	result := &ProjectValidationResult{
		IsValid:  true,
		Issues:   []*DomainError{},
		Warnings: []*DomainError{},
	}

	// Check if project directory exists
	exists, err := vu.repo.DirExists(ctx, projectPath)
	if err != nil {
		return nil, NewSystemError(
			ErrFileNotFound,
			"Failed to check project directory",
			projectPath,
			err,
		)
	}

	if !exists {
		return nil, NewSystemError(ErrFileNotFound, "Project directory not found", projectPath, nil)
	}

	// Analyze project structure
	info, err := vu.analyzeProjectStructure(ctx, projectPath)
	if err != nil {
		return nil, err
	}

	result.Info = info

	// Validate project structure
	validationErr := vu.validateProjectRequirements(info)
	if validationErr != nil {
		result.Issues = append(result.Issues, validationErr)
		result.IsValid = false
	}

	// Generate recommendations
	vu.generateProjectRecommendations(ctx, info, result)

	vu.logger.DebugContext(
		ctx,
		"Project structure validation completed",
		"valid",
		result.IsValid,
		"issues",
		len(result.Issues),
	)

	return result, nil
}

// analyzeProjectStructure analyzes the project directory structure.
func (vu *ValidationUseCase) analyzeProjectStructure(
	ctx context.Context,
	projectPath string,
) (*ProjectInfo, error) {
	info := &ProjectInfo{
		Path: projectPath,
	}

	// Check for go.mod
	exists, err := vu.repo.FileExists(ctx, vu.repo.JoinPath(projectPath, "go.mod"))
	if err != nil {
		return nil, NewSystemError(
			ErrFileReadFailed,
			"Failed to check for go.mod",
			projectPath,
			err,
		)
	}

	info.HasGoMod = exists

	if exists {
		// Parse go.mod for module name and dependencies
		err := vu.parseGoMod(ctx, vu.repo.JoinPath(projectPath, "go.mod"), info)
		if err != nil {
			vu.logger.WarnContext(ctx, "Failed to parse go.mod", "error", err)
		}
	}

	// Find main.go file
	mainPath := vu.findMainFile(ctx, projectPath)

	if mainPath != "" {
		info.HasMainFile = true
		info.MainFilePath = mainPath
	}

	// Determine project type and binary name
	vu.inferProjectType(info)

	return info, nil
}

// parseGoMod parses go.mod file for module information.
func (vu *ValidationUseCase) parseGoMod(
	ctx context.Context,
	goModPath string,
	info *ProjectInfo,
) error {
	data, err := vu.repo.ReadFile(ctx, goModPath)
	if err != nil {
		return NewSystemError(ErrFileReadFailed, "Failed to read go.mod", goModPath, err)
	}

	// Simple parsing for module name
	content := string(data)

	lines := strings.SplitSeq(content, "\n")
	for line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "module ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				info.Name = strings.Trim(parts[1], `"`)

				break
			}
		}
	}

	return nil
}

// findMainFile searches for main.go in common locations.
func (vu *ValidationUseCase) findMainFile(ctx context.Context, projectPath string) string {
	commonPaths := []string{
		"main.go",
		"cmd/main.go",
		"cmd/" + filepath.Base(projectPath) + "/main.go",
		"src/main.go",
	}

	for _, path := range commonPaths {
		fullPath := vu.repo.JoinPath(projectPath, path)

		exists, err := vu.repo.FileExists(ctx, fullPath)
		if err != nil {
			continue
		}

		if exists {
			return path
		}
	}

	return ""
}

// inferProjectType infers project type and binary name from structure.
func (vu *ValidationUseCase) inferProjectType(info *ProjectInfo) {
	// Infer binary name from module name
	if info.Name != "" && info.MainFilePath != "" {
		// Use directory name from main path as binary name
		dir := filepath.Dir(info.MainFilePath)
		if dir != "." {
			info.BinaryName = filepath.Base(dir)
		} else {
			// Use module name last part
			parts := strings.Split(info.Name, "/")
			info.BinaryName = parts[len(parts)-1]
		}
	}

	// Infer project type from structure
	if strings.Contains(info.MainFilePath, "cmd/") {
		info.ProjectType = ProjectTypeCLI
	} else if strings.Contains(info.MainFilePath, "server/") || strings.Contains(info.MainFilePath, "api/") {
		info.ProjectType = ProjectTypeWebAPI
	} else if strings.Contains(info.MainFilePath, "web/") {
		info.ProjectType = ProjectTypeWebAPI
	} else {
		info.ProjectType = ProjectTypeLibrary
	}

	info.Buildable = info.HasMainFile && info.HasGoMod
}

// validateProjectRequirements validates project against requirements.
func (vu *ValidationUseCase) validateProjectRequirements(
	info *ProjectInfo,
) *DomainError {
	if !info.HasGoMod {
		return NewSystemError(
			ErrDependencyNotFound,
			"Go module not found",
			"Project must have a go.mod file",
			nil,
		)
	}

	if !info.HasMainFile {
		return NewSystemError(
			ErrFileNotFound,
			"Main file not found",
			"Project must have a main.go file for compilation",
			nil,
		)
	}

	return nil
}

// generateProjectRecommendations generates recommendations for project improvement.
func (vu *ValidationUseCase) generateProjectRecommendations(
	ctx context.Context,
	info *ProjectInfo,
	result *ProjectValidationResult,
) {
	// Recommendation for missing GitHub Actions
	if _, err := vu.repo.DirExists(
		ctx,
		vu.repo.JoinPath(info.Path, ".github", "workflows"),
	); err != nil {
		result.Recommendations = append(
			result.Recommendations,
			"Add GitHub Actions workflow for automated builds and releases",
		)
	}

	// Check for missing files and add recommendations
	checkMissingFileAndRecommend(
		ctx,
		vu.repo,
		info.Path,
		"README.md",
		"Add README.md with project documentation",
		&result.Recommendations,
	)
	checkMissingFileAndRecommend(
		ctx,
		vu.repo,
		info.Path,
		".goreleaser.yaml",
		"Add .goreleaser.yaml configuration for automated releases",
		&result.Recommendations,
	)
	checkMissingFileAndRecommend(
		ctx,
		vu.repo,
		info.Path,
		"Dockerfile",
		"Add Dockerfile for containerized builds",
		&result.Recommendations,
	)
}

// Utility functions for security validation.
func containsPathTraversal(path string) bool {
	return strings.Contains(path, "..") || strings.Contains(path, `\`)
}

func containsShellMetacharacters(value string) bool {
	shellMetachars := []string{"|", "&", ";", "<", ">", "`", "$", "(", ")", "{", "}"}
	for _, char := range shellMetachars {
		if strings.Contains(value, char) {
			return true
		}
	}

	return false
}

func containsURLInjection(url string) bool {
	return strings.Contains(url, "javascript:") || strings.Contains(url, "data:") ||
		strings.Contains(url, "vbscript:")
}

// checkMissingFileAndRecommend checks if a file exists and adds a recommendation if it's missing.
func checkMissingFileAndRecommend(
	ctx context.Context,
	repo FileSystemRepository,
	basePath, filename, recommendation string,
	recommendations *[]string,
) {
	if _, err := repo.FileExists(ctx, repo.JoinPath(basePath, filename)); err != nil {
		*recommendations = append(*recommendations, recommendation)
	}
}
