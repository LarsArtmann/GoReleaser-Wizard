package domain

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
)

// Validation use case implementation
type ValidationUseCase struct {
	logger Logger
	repo  FileSystemRepository
}

// NewValidationUseCase creates a new validation use case
func NewValidationUseCase(logger Logger, repo FileSystemRepository) *ValidationUseCase {
	return &ValidationUseCase{
		logger: logger,
		repo:  repo,
	}
}

// ValidateConfiguration performs comprehensive validation of project configuration
func (vu *ValidationUseCase) ValidateConfiguration(ctx context.Context, config *SafeProjectConfig) (*ValidationResult, error) {
	vu.logger.DebugContext(ctx, "Starting comprehensive configuration validation")
	
	result := &ValidationResult{
		IsValid:  true,
		Errors:   []*DomainError{},
		Warnings: []*DomainError{},
	}
	
	// Step 1: Basic field validation
	if err := vu.validateBasicFields(ctx, config); err != nil {
		result.Errors = append(result.Errors, err)
		result.IsValid = false
	}
	
	// Step 2: Type validation
	if err := vu.validateTypes(ctx, config); err != nil {
		result.Errors = append(result.Errors, err)
		result.IsValid = false
	}
	
	// Step 3: Platform-architecture compatibility
	if err := vu.validatePlatformArchCompatibility(ctx, config); err != nil {
		result.Errors = append(result.Errors, err)
		result.IsValid = false
	}
	
	// Step 4: Business rule validation
	if err := vu.validateBusinessRules(ctx, config); err != nil {
		result.Errors = append(result.Errors, err)
		result.IsValid = false
	}
	
	// Step 5: Security validation
	if err := vu.validateSecurity(ctx, config); err != nil {
		result.Errors = append(result.Errors, err)
		result.IsValid = false
	}
	
	// Step 6: Generate warnings
	vu.generateWarnings(ctx, config, result)
	
	vu.logger.DebugContext(ctx, "Validation completed", "valid", result.IsValid, "errors", len(result.Errors), "warnings", len(result.Warnings))
	
	return result, nil
}

// validateBasicFields validates basic required fields
func (vu *ValidationUseCase) validateBasicFields(ctx context.Context, config *SafeProjectConfig) *DomainError {
	// Project name validation
	if err := ValidateProjectName(config.ProjectName); err != nil {
		return NewValidationError(ErrInvalidProjectName, "Project name validation failed", err.Error()).WithContext("project_name")
	}
	
	// Binary name validation
	if err := ValidateBinaryName(config.BinaryName); err != nil {
		return NewValidationError(ErrInvalidBinaryName, "Binary name validation failed", err.Error()).WithContext("binary_name")
	}
	
	// Main path validation
	if err := ValidateMainPath(config.MainPath); err != nil {
		return NewValidationError(ErrInvalidMainPath, "Main path validation failed", err.Error()).WithContext("main_path")
	}
	
	// Project description validation (optional)
	if config.ProjectDescription != "" {
		if err := ValidateProjectDescription(config.ProjectDescription); err != nil {
			return NewValidationError(ErrInvalidProjectDescription, "Project description validation failed", err.Error()).WithContext("project_description")
		}
	}
	
	return nil
}

// validateTypes validates enum types
func (vu *ValidationUseCase) validateTypes(ctx context.Context, config *SafeProjectConfig) *DomainError {
	// Project type validation
	if !config.ProjectType.IsValid() {
		return NewValidationError(ErrInvalidProjectName, "Invalid project type", fmt.Sprintf("'%s' is not a supported project type", config.ProjectType)).WithContext("project_type")
	}
	
	// Platform validation
	if err := ValidatePlatforms(config.Platforms); err != nil {
		return NewValidationError(ErrInvalidPlatform, "Platform validation failed", err.Error()).WithContext("platforms")
	}
	
	// Architecture validation
	if err := ValidateArchitectures(config.Architectures); err != nil {
		return NewValidationError(ErrInvalidArchitecture, "Architecture validation failed", err.Error()).WithContext("architectures")
	}
	
	// Git provider validation
	if err := ValidateGitProvider(config.GitProvider); err != nil {
		return NewValidationError(ErrInvalidGitProvider, "Git provider validation failed", err.Error()).WithContext("git_provider")
	}
	
	// Docker registry validation
	if err := ValidateDockerRegistry(config.DockerRegistry); err != nil {
		return NewValidationError(ErrInvalidDockerRegistry, "Docker registry validation failed", err.Error()).WithContext("docker_registry")
	}
	
	// Action triggers validation
	if err := ValidateActionTriggers(config.ActionsOn); err != nil {
		return NewValidationError(ErrInvalidActionTrigger, "Action trigger validation failed", err.Error()).WithContext("actions_on")
	}
	
	// Build tags validation
	if len(config.BuildTags) > 0 {
		if err := ValidateBuildTags(config.BuildTags); err != nil {
			return NewValidationError(ErrInvalidBuildTag, "Build tags validation failed", err.Error()).WithContext("build_tags")
		}
	}
	
	// Configuration state validation
	if err := ValidateConfigState(config.State); err != nil {
		return NewValidationError(ErrInvalidConfigState, "Configuration state validation failed", err.Error()).WithContext("state")
	}
	
	return nil
}

// validatePlatformArchCompatibility validates platform-architecture compatibility
func (vu *ValidationUseCase) validatePlatformArchCompatibility(ctx context.Context, config *SafeProjectConfig) *DomainError {
	if err := ValidatePlatformArchCompatibility(config.Platforms, config.Architectures); err != nil {
		return NewValidationError(ErrPlatformArchMismatch, "Platform architecture compatibility failed", err.Error()).WithContext("platforms_architectures")
	}
	return nil
}

// validateBusinessRules validates domain business rules
func (vu *ValidationUseCase) validateBusinessRules(ctx context.Context, config *SafeProjectConfig) *DomainError {
	// Docker support rule
	if config.GetDockerEnabled() && !config.ProjectType.DockerSupported() {
		return NewConfigurationError(ErrDockerNotSupported, "Docker not supported for project type", fmt.Sprintf("Project type %s does not support Docker", config.ProjectType)).WithContext("docker_enabled")
	}
	
	// Main path requirement rule
	if config.ProjectType.RequiresMainPath() && config.MainPath == "" {
		return NewConfigurationError(ErrMainPathRequired, "Main path required for project type", fmt.Sprintf("Project type %s requires a main path", config.ProjectType)).WithContext("main_path")
	}
	
	// State transition rule
	if !config.State.AllowsGeneration() && config.GetGenerateActions() {
		return NewConfigurationError(ErrInvalidStateTransition, "State transition invalid", fmt.Sprintf("Configuration in state '%s' cannot generate actions", config.State)).WithContext("generate_actions")
	}
	
	// Docker registry URL validation
	if config.GetDockerEnabled() {
		if err := ValidateDockerRegistryURL(config.DockerRegistry, config.DockerImage); err != nil {
			return NewValidationError(ErrInvalidURLPattern, "Docker registry URL validation failed", err.Error()).WithContext("docker_image")
		}
	}
	
	return nil
}

// validateSecurity performs security-focused validation
func (vu *ValidationUseCase) validateSecurity(ctx context.Context, config *SafeProjectConfig) *DomainError {
	// Check for potential path traversal in main path
	if containsPathTraversal(config.MainPath) {
		return NewBusinessRuleError(ErrInvalidCharacters, "Path traversal detected", "Main path contains potentially dangerous path traversal sequences").WithContext("main_path")
	}
	
	// Check for shell metacharacters in binary name
	if containsShellMetacharacters(config.BinaryName) {
		return NewBusinessRuleError(ErrInvalidCharacters, "Shell metacharacters detected", "Binary name contains potentially dangerous shell metacharacters").WithContext("binary_name")
	}
	
	// Check Docker image name for security issues
	if config.DockerImage != "" {
		if containsURLInjection(config.DockerImage) {
			return NewBusinessRuleError(ErrInvalidCharacters, "URL injection detected", "Docker image name contains potentially dangerous URL injection sequences").WithContext("docker_image")
		}
	}
	
	return nil
}

// generateWarnings generates validation warnings
func (vu *ValidationUseCase) generateWarnings(ctx context.Context, config *SafeProjectConfig, result *ValidationResult) {
	// Warning for single platform
	if len(config.Platforms) == 1 {
		warning := NewBusinessRuleError(ErrMissingRequiredField, "Single platform configuration", "Consider targeting multiple platforms for broader compatibility").WithContext("platforms")
		result.Warnings = append(result.Warnings, warning)
	}
	
	// Warning for missing Docker image name
	if config.GetDockerEnabled() && config.DockerImage == "" {
	}
	
	// Warning for mismatched CGO setting
	if config.CGOStatus.ToBool() != config.ProjectType.DefaultCGOEnabled() {
		warning := NewConfigurationError(ErrInvalidStateTransition, "CGO setting mismatched", "CGO setting differs from project type default").WithContext("cgo_enabled")
		result.Warnings = append(result.Warnings, warning)
	}
	
	// Warning for missing version information
	if !config.LDFlags {
		warning := NewConfigurationError(ErrMissingRequiredField, "LD flags disabled", "Version information injection is disabled").WithContext("ldflags")
		result.Warnings = append(result.Warnings, warning)
	}
}

// ValidateProjectStructure validates project directory structure
func (vu *ValidationUseCase) ValidateProjectStructure(ctx context.Context, projectPath string) (*ProjectValidationResult, error) {
	vu.logger.DebugContext(ctx, "Validating project structure", "path", projectPath)
	
	result := &ProjectValidationResult{
		IsValid:  true,
		Issues:   []*DomainError{},
		Warnings: []*DomainError{},
	}
	
	// Check if project directory exists
	exists, err := vu.repo.DirExists(ctx, projectPath)
	if err != nil {
		return nil, NewSystemError(ErrFileNotFound, "Failed to check project directory", projectPath, err)
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
	if err := vu.validateProjectRequirements(ctx, info); err != nil {
		result.Issues = append(result.Issues, err)
		result.IsValid = false
	}
	
	// Generate recommendations
	vu.generateProjectRecommendations(ctx, info, result)
	
	vu.logger.DebugContext(ctx, "Project structure validation completed", "valid", result.IsValid, "issues", len(result.Issues))
	
	return result, nil
}

// analyzeProjectStructure analyzes the project directory structure
func (vu *ValidationUseCase) analyzeProjectStructure(ctx context.Context, projectPath string) (*ProjectInfo, error) {
	info := &ProjectInfo{
		Path: projectPath,
	}
	
	// Check for go.mod
	exists, err := vu.repo.FileExists(ctx, vu.repo.JoinPath(projectPath, "go.mod"))
	if err != nil {
		return nil, NewSystemError(ErrFileReadFailed, "Failed to check for go.mod", projectPath, err)
	}
	info.HasGoMod = exists
	
	if exists {
		// Parse go.mod for module name and dependencies
		if err := vu.parseGoMod(ctx, vu.repo.JoinPath(projectPath, "go.mod"), info); err != nil {
			vu.logger.WarnContext(ctx, "Failed to parse go.mod", "error", err)
		}
	}
	
	// Find main.go file
	mainPath, err := vu.findMainFile(ctx, projectPath)
	if err != nil {
		vu.logger.WarnContext(ctx, "Failed to find main.go", "error", err)
	}
	if mainPath != "" {
		info.HasMainFile = true
		info.MainFilePath = mainPath
	}
	
	// Determine project type and binary name
	vu.inferProjectType(ctx, info)
	
	return info, nil
}

// parseGoMod parses go.mod file for module information
func (vu *ValidationUseCase) parseGoMod(ctx context.Context, goModPath string, info *ProjectInfo) error {
	data, err := vu.repo.ReadFile(ctx, goModPath)
	if err != nil {
		return NewSystemError(ErrFileReadFailed, "Failed to read go.mod", goModPath, err)
	}
	
	// Simple parsing for module name
	content := string(data)
	lines := strings.Split(content, "\n")
	for _, line := range lines {
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

// findMainFile searches for main.go in common locations
func (vu *ValidationUseCase) findMainFile(ctx context.Context, projectPath string) (string, error) {
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
			return path, nil
		}
	}
	
	return "", nil
}

// inferProjectType infers project type and binary name from structure
func (vu *ValidationUseCase) inferProjectType(ctx context.Context, info *ProjectInfo) {
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
		info.ProjectType = ProjectTypeAPI
	} else if strings.Contains(info.MainFilePath, "web/") {
		info.ProjectType = ProjectTypeWeb
	} else {
		info.ProjectType = ProjectTypeLibrary
	}
	
	info.Buildable = info.HasMainFile && info.HasGoMod
}

// validateProjectRequirements validates project against requirements
func (vu *ValidationUseCase) validateProjectRequirements(ctx context.Context, info *ProjectInfo) *DomainError {
	if !info.HasGoMod {
		return NewSystemError(ErrDependencyNotFound, "Go module not found", "Project must have a go.mod file", nil)
	}
	
	if !info.HasMainFile {
		return NewSystemError(ErrFileNotFound, "Main file not found", "Project must have a main.go file for compilation", nil)
	}
	
	return nil
}

// generateProjectRecommendations generates recommendations for project improvement
func (vu *ValidationUseCase) generateProjectRecommendations(ctx context.Context, info *ProjectInfo, result *ProjectValidationResult) {
	// Recommendation for missing GitHub Actions
	if _, err := vu.repo.DirExists(ctx, vu.repo.JoinPath(info.Path, ".github", "workflows")); err != nil {
		result.Recommendations = append(result.Recommendations, "Add GitHub Actions workflow for automated builds and releases")
	}
	
	// Recommendation for missing README
	if _, err := vu.repo.FileExists(ctx, vu.repo.JoinPath(info.Path, "README.md")); err != nil {
		result.Recommendations = append(result.Recommendations, "Add README.md with project documentation")
	}
	
	// Recommendation for missing GoReleaser config
	if _, err := vu.repo.FileExists(ctx, vu.repo.JoinPath(info.Path, ".goreleaser.yaml")); err != nil {
		result.Recommendations = append(result.Recommendations, "Add .goreleaser.yaml configuration for automated releases")
	}
	
	// Recommendation for missing Dockerfile
	if _, err := vu.repo.FileExists(ctx, vu.repo.JoinPath(info.Path, "Dockerfile")); err != nil {
		result.Recommendations = append(result.Recommendations, "Add Dockerfile for containerized builds")
	}
}

// Utility functions for security validation
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
	return strings.Contains(url, "javascript:") || strings.Contains(url, "data:") || strings.Contains(url, "vbscript:")
}