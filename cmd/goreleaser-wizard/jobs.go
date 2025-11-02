package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/charmbracelet/log"
)

// ConfigGenerationJob generates GoReleaser configuration
type ConfigGenerationJob struct {
	id     string
	config *ProjectConfig
	force  bool
	logger *log.Logger
}

// NewConfigGenerationJob creates a new config generation job
func NewConfigGenerationJob(config *ProjectConfig, force bool, logger *log.Logger) *ConfigGenerationJob {
	return &ConfigGenerationJob{
		id:     "config-generation",
		config: config,
		force:  force,
		logger: logger,
	}
}

func (j *ConfigGenerationJob) ID() string {
	return j.id
}

func (j *ConfigGenerationJob) Name() string {
	return "Generate GoReleaser Configuration"
}

func (j *ConfigGenerationJob) Execute(ctx context.Context) error {
	j.logger.Info("Generating GoReleaser configuration")
	
	// Check if context is cancelled
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Validate config
	if j.config.ProjectName == "" {
		return fmt.Errorf("project name is required")
	}

	// Check existing files
	if !j.force {
		if _, err := os.Stat(".goreleaser.yaml"); err == nil {
			return fmt.Errorf(".goreleaser.yaml already exists (use --force to overwrite)")
		}
	}

	// Generate configuration
	err := generateGoReleaserConfig(j.config)
	if err != nil {
		return fmt.Errorf("failed to generate GoReleaser config: %w", err)
	}

	j.logger.Info("GoReleaser configuration generated successfully")
	return nil
}

func (j *ConfigGenerationJob) Rollback(ctx context.Context) error {
	j.logger.Info("Rolling back GoReleaser configuration generation")
	
	// Check if context is cancelled
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Remove generated config
	if _, err := os.Stat(".goreleaser.yaml"); err == nil {
		// Check if backup exists
		if _, err := os.Stat(".goreleaser.yaml.backup"); err == nil {
			// Restore backup
			err := os.Rename(".goreleaser.yaml.backup", ".goreleaser.yaml")
			if err != nil {
				j.logger.Errorf("Failed to restore backup: %v", err)
				return err
			}
			j.logger.Info("Restored backup configuration")
		} else {
			// Remove generated file
			err := os.Remove(".goreleaser.yaml")
			if err != nil {
				j.logger.Errorf("Failed to remove generated config: %v", err)
				return err
			}
			j.logger.Info("Removed generated configuration")
		}
	}

	return nil
}

// GitHubActionsGenerationJob generates GitHub Actions workflow
type GitHubActionsGenerationJob struct {
	id     string
	config *ProjectConfig
	logger *log.Logger
}

// NewGitHubActionsGenerationJob creates a new GitHub Actions generation job
func NewGitHubActionsGenerationJob(config *ProjectConfig, logger *log.Logger) *GitHubActionsGenerationJob {
	return &GitHubActionsGenerationJob{
		id:     "github-actions-generation",
		config: config,
		logger: logger,
	}
}

func (j *GitHubActionsGenerationJob) ID() string {
	return j.id
}

func (j *GitHubActionsGenerationJob) Name() string {
	return "Generate GitHub Actions Workflow"
}

func (j *GitHubActionsGenerationJob) Execute(ctx context.Context) error {
	j.logger.Info("Generating GitHub Actions workflow")
	
	// Check if context is cancelled
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Check if GitHub Actions is enabled
	if !j.config.GenerateActions {
		j.logger.Info("GitHub Actions generation is disabled, skipping")
		return nil
	}

	// Create .github/workflows directory
	workflowDir := ".github/workflows"
	err := os.MkdirAll(workflowDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create workflow directory: %w", err)
	}

	// Generate workflow
	err = generateGitHubActions(j.config)
	if err != nil {
		return fmt.Errorf("failed to generate GitHub Actions workflow: %w", err)
	}

	j.logger.Info("GitHub Actions workflow generated successfully")
	return nil
}

func (j *GitHubActionsGenerationJob) Rollback(ctx context.Context) error {
	j.logger.Info("Rolling back GitHub Actions workflow generation")
	
	// Check if context is cancelled
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Remove generated workflow
	workflowPath := filepath.Join(".github", "workflows", "release.yml")
	if _, err := os.Stat(workflowPath); err == nil {
		err := os.Remove(workflowPath)
		if err != nil {
			j.logger.Errorf("Failed to remove generated workflow: %v", err)
			return err
		}
		j.logger.Info("Removed generated workflow")
	}

	// Try to remove .github directory if empty
	workflowDir := filepath.Join(".github", "workflows")
	workflowFiles, err := filepath.Glob(filepath.Join(workflowDir, "*.yml"))
	if err == nil && len(workflowFiles) == 0 {
		os.Remove(workflowDir)
		os.Remove(".github")
		j.logger.Info("Removed empty .github directory")
	}

	return nil
}

// ProjectValidationJob validates project structure
type ProjectValidationJob struct {
	id       string
	projectDir string
	logger   *log.Logger
}

// NewProjectValidationJob creates a new project validation job
func NewProjectValidationJob(projectDir string, logger *log.Logger) *ProjectValidationJob {
	return &ProjectValidationJob{
		id:         "project-validation",
		projectDir: projectDir,
		logger:     logger,
	}
}

func (j *ProjectValidationJob) ID() string {
	return j.id
}

func (j *ProjectValidationJob) Name() string {
	return "Validate Project Structure"
}

func (j *ProjectValidationJob) Execute(ctx context.Context) error {
	j.logger.Info("Validating project structure")
	
	// Check if context is cancelled
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Check if project directory exists
	if _, err := os.Stat(j.projectDir); os.IsNotExist(err) {
		return fmt.Errorf("project directory does not exist: %s", j.projectDir)
	}

	// Check for go.mod
	goModPath := filepath.Join(j.projectDir, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		return fmt.Errorf("go.mod not found in project directory")
	}

	// Check for main package
	mainPaths := []string{
		filepath.Join(j.projectDir, "main.go"),
		filepath.Join(j.projectDir, "cmd", "*", "main.go"),
	}
	
	var mainFound bool
	for _, mainPath := range mainPaths {
		matches, err := filepath.Glob(mainPath)
		if err == nil && len(matches) > 0 {
			mainFound = true
			break
		}
	}
	
	if !mainFound {
		return fmt.Errorf("no main.go found in project (expected at main.go or cmd/*/main.go)")
	}

	j.logger.Info("Project structure validation passed")
	return nil
}

func (j *ProjectValidationJob) Rollback(ctx context.Context) error {
	// Validation job doesn't create any files, so rollback is a no-op
	j.logger.Info("Project validation rollback is a no-op")
	return nil
}

// DependencyCheckJob checks for required dependencies
type DependencyCheckJob struct {
	id        string
	dependencies []string
	logger    *log.Logger
}

// NewDependencyCheckJob creates a new dependency check job
func NewDependencyCheckJob(dependencies []string, logger *log.Logger) *DependencyCheckJob {
	return &DependencyCheckJob{
		id:          "dependency-check",
		dependencies: dependencies,
		logger:      logger,
	}
}

func (j *DependencyCheckJob) ID() string {
	return j.id
}

func (j *DependencyCheckJob) Name() string {
	return "Check Dependencies"
}

func (j *DependencyCheckJob) Execute(ctx context.Context) error {
	j.logger.Info("Checking dependencies")
	
	// Check if context is cancelled
	if ctx.Err() != nil {
		return ctx.Err()
	}

	var missingDeps []string
	
	for _, dep := range j.dependencies {
		path, err := exec.LookPath(dep)
		if err != nil {
			missingDeps = append(missingDeps, dep)
			j.logger.Warnf("Dependency not found: %s", dep)
		} else {
			j.logger.Debugf("Found dependency: %s at %s", dep, path)
		}
	}
	
	if len(missingDeps) > 0 {
		return fmt.Errorf("missing dependencies: %v", missingDeps)
	}

	j.logger.Info("All dependencies are available")
	return nil
}

func (j *DependencyCheckJob) Rollback(ctx context.Context) error {
	// Dependency check doesn't modify state, so rollback is a no-op
	j.logger.Info("Dependency check rollback is a no-op")
	return nil
}

// JobFactory creates jobs for common wizard operations
type JobFactory struct {
	logger *log.Logger
}

// NewJobFactory creates a new job factory
func NewJobFactory(logger *log.Logger) *JobFactory {
	return &JobFactory{
		logger: logger,
	}
}

// CreateFullWizardJobs creates all jobs for a complete wizard operation
func (jf *JobFactory) CreateFullWizardJobs(config *ProjectConfig, force bool) []Job {
	var jobs []Job
	
	// Add project validation job
	jobs = append(jobs, NewProjectValidationJob(".", jf.logger))
	
	// Add dependency check job
	dependencies := []string{"go"}
	if config.DockerEnabled {
		dependencies = append(dependencies, "docker")
	}
	if config.Signing {
		dependencies = append(dependencies, "cosign")
	}
	jobs = append(jobs, NewDependencyCheckJob(dependencies, jf.logger))
	
	// Add config generation job
	jobs = append(jobs, NewConfigGenerationJob(config, force, jf.logger))
	
	// Add GitHub Actions generation job
	if config.GenerateActions {
		jobs = append(jobs, NewGitHubActionsGenerationJob(config, jf.logger))
	}
	
	return jobs
}

// CreateConfigOnlyJobs creates jobs for config generation only
func (jf *JobFactory) CreateConfigOnlyJobs(config *ProjectConfig, force bool) []Job {
	return []Job{
		NewProjectValidationJob(".", jf.logger),
		NewConfigGenerationJob(config, force, jf.logger),
	}
}

// CreateValidationOnlyJob creates a validation-only job
func (jf *JobFactory) CreateValidationOnlyJob(projectDir string) Job {
	return NewProjectValidationJob(projectDir, jf.logger)
}