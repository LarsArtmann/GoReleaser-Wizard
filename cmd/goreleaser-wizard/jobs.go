package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"charm.land/log/v2"
	"github.com/LarsArtmann/GoReleaser-Wizard/cmd/goreleaser-wizard/generators"
	"github.com/LarsArtmann/GoReleaser-Wizard/cmd/goreleaser-wizard/types"
	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
)

// noOpRollback returns a no-op rollback function for jobs that don't modify state.
func noOpRollback(logger *log.Logger, jobName string) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		logger.Info(jobName + " rollback is a no-op")

		return nil
	}
}

// noOpRollbackHelper provides a no-op rollback implementation for validation jobs.
type noOpRollbackHelper struct {
	logger *log.Logger
	name   string
}

func (h *noOpRollbackHelper) Rollback(ctx context.Context) error {
	return noOpRollback(h.logger, h.name)(ctx)
}

const workflowDirPermission = 0o755

// generateGoReleaserConfig generates the GoReleaser configuration file in the
// current directory via the typed generator (single template source of truth).
func generateGoReleaserConfig(config *domain.SafeProjectConfig, logger *log.Logger) error {
	if err := config.Validate(); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	generator := generators.NewGoReleaserGenerator(config, &LoggerAdapter{logger: logger})

	if err := generator.Generate(context.Background()); err != nil {
		return fmt.Errorf("failed to generate GoReleaser config: %w", err)
	}

	return nil
}

// generateGitHubActions generates the GitHub Actions release workflow in the
// current directory via the typed generator (single template source of truth).
func generateGitHubActions(config *domain.SafeProjectConfig, logger *log.Logger) error {
	if !config.ShouldGenerateActionsFiles() {
		return errors.New("GitHub Actions generation is not enabled")
	}

	generator := generators.NewGitHubActionsGenerator(config, &LoggerAdapter{logger: logger})

	if err := generator.Generate(context.Background()); err != nil {
		return fmt.Errorf("failed to generate GitHub Actions workflow: %w", err)
	}

	return nil
}

// ConfigGenerationJob generates GoReleaser configuration.
type ConfigGenerationJob struct {
	id     string
	config *domain.SafeProjectConfig
	force  bool
	logger *log.Logger
}

// NewConfigGenerationJob creates a new config generation job.
func NewConfigGenerationJob(
	config *domain.SafeProjectConfig,
	force bool,
	logger *log.Logger,
) *ConfigGenerationJob {
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

	if err := checkContextCancellation(ctx, "context cancelled"); err != nil {
		return err
	}

	// Transition to Processing state
	j.config.State = domain.ConfigStateProcessing

	// Validate config
	if j.config.ProjectName == "" {
		j.config.State = domain.ConfigStateInvalid

		return errors.New("project name is required")
	}

	// Check existing files
	if !j.force {
		if _, err := os.Stat(goreleaserConfigFilename); err == nil {
			j.config.State = domain.ConfigStateDraft

			return errors.New(".goreleaser.yaml already exists (use --force to overwrite)")
		}
	}

	// Apply defaults to ensure config is complete
	j.config.ApplyDefaults()

	// Generate configuration
	err := generateGoReleaserConfig(j.config, j.logger)
	if err != nil {
		j.config.State = domain.ConfigStateInvalid

		return fmt.Errorf("failed to generate GoReleaser config: %w", err)
	}

	// Transition to Generated state
	j.config.State = domain.ConfigStateGenerated

	j.logger.Info("GoReleaser configuration generated successfully")

	return nil
}

func (j *ConfigGenerationJob) Rollback(ctx context.Context) error {
	j.logger.Info("Rollback: removing generated .goreleaser.yaml")

	err := os.Remove(goreleaserConfigFilename)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove .goreleaser.yaml: %w", err)
	}

	j.logger.Info("Rollback: .goreleaser.yaml removed")

	return nil
}

// GitHubActionsGenerationJob generates GitHub Actions workflow.
type GitHubActionsGenerationJob struct {
	id     string
	config *domain.SafeProjectConfig
	logger *log.Logger
}

// NewGitHubActionsGenerationJob creates a new GitHub Actions generation job.
func NewGitHubActionsGenerationJob(
	config *domain.SafeProjectConfig,
	logger *log.Logger,
) *GitHubActionsGenerationJob {
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

	if err := checkContextCancellation(ctx, "context cancelled"); err != nil {
		return err
	}

	// Check if GitHub Actions is enabled
	if !j.config.GetGenerateActions() {
		j.logger.Info("GitHub Actions generation is disabled, skipping")

		return nil
	}

	// Create .github/workflows directory
	workflowDir := ".github/workflows"

	err := os.MkdirAll(workflowDir, 0o755)
	if err != nil {
		return fmt.Errorf("failed to create workflow directory: %w", err)
	}

	// Generate workflow
	err = generateGitHubActions(j.config, j.logger)
	if err != nil {
		return fmt.Errorf("failed to generate GitHub Actions workflow: %w", err)
	}

	j.logger.Info("GitHub Actions workflow generated successfully")

	return nil
}

func (j *GitHubActionsGenerationJob) Rollback(ctx context.Context) error {
	j.logger.Info("Rolling back GitHub Actions workflow generation")

	if err := checkContextCancellation(ctx, "context cancelled"); err != nil {
		return err
	}

	// Remove generated workflow
	workflowPath := releaseWorkflowTargetPath
	if _, err := os.Stat(workflowPath); err == nil {
		err := os.Remove(workflowPath)
		if err != nil {
			j.logger.Errorf("Failed to remove generated workflow: %v", err)

			return fmt.Errorf("failed to remove generated workflow: %w", err)
		}

		j.logger.Info("Removed generated workflow")
	}

	// Try to remove .github directory if empty
	workflowDir := filepath.Join(".github", "workflows")

	workflowFiles, err := filepath.Glob(filepath.Join(workflowDir, "*.yml"))
	if err == nil && len(workflowFiles) == 0 {
		_ = os.Remove(workflowDir)
		_ = os.Remove(".github")

		j.logger.Info("Removed empty .github directory")
	}

	return nil
}

// DockerfileGenerationJob generates a Dockerfile when Docker support is enabled.
type DockerfileGenerationJob struct {
	id     string
	config *domain.SafeProjectConfig
	force  bool
	logger *log.Logger
}

// NewDockerfileGenerationJob creates a new Dockerfile generation job.
func NewDockerfileGenerationJob(
	config *domain.SafeProjectConfig,
	force bool,
	logger *log.Logger,
) *DockerfileGenerationJob {
	return &DockerfileGenerationJob{
		id:     "dockerfile-generation",
		config: config,
		force:  force,
		logger: logger,
	}
}

func (j *DockerfileGenerationJob) ID() string {
	return j.id
}

func (j *DockerfileGenerationJob) Name() string {
	return "Generate Dockerfile"
}

func (j *DockerfileGenerationJob) Execute(ctx context.Context) error {
	j.logger.Info("Generating Dockerfile")

	if err := checkContextCancellation(ctx, "context cancelled"); err != nil {
		return err
	}

	if !j.config.DockerSupport.IsEnabled() {
		j.logger.Info("Docker support is disabled, skipping Dockerfile generation")

		return nil
	}

	if !j.force {
		if _, err := os.Stat(dockerfileFilename); err == nil {
			return errors.New("dockerfile already exists (use --force to overwrite)")
		}
	}

	generator := generators.NewDockerfileGenerator(j.config, &LoggerAdapter{logger: j.logger})

	if err := generator.Generate(ctx); err != nil {
		return fmt.Errorf("failed to generate Dockerfile: %w", err)
	}

	j.logger.Info("Dockerfile generated successfully")

	return nil
}

func (j *DockerfileGenerationJob) Rollback(ctx context.Context) error {
	j.logger.Info("Rollback: removing generated Dockerfile")

	if err := os.Remove(dockerfileFilename); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove Dockerfile: %w", err)
	}

	return nil
}

// ProjectValidationJob validates project structure.
type ProjectValidationJob struct {
	noOpRollbackHelper

	id         string
	projectDir string
}

// NewProjectValidationJob creates a new project validation job.
func NewProjectValidationJob(projectDir string, logger *log.Logger) *ProjectValidationJob {
	return &ProjectValidationJob{
		logger:     logger,
		name:       "Project validation",
		id:         "project-validation",
		projectDir: projectDir,
	}
}

func (j *ProjectValidationJob) ID() string {
	return j.id
}

func (j *ProjectValidationJob) Name() string {
	return "Validate Project Structure"
}

func checkJobContext(ctx context.Context, logger *log.Logger, message string) error {
	logger.Info(message)

	return checkContextCancellation(ctx, "context cancelled")
}

func (j *ProjectValidationJob) Execute(ctx context.Context) error {
	if err := checkJobContext(ctx, j.logger, "Validating project structure"); err != nil {
		return err
	}

	// Check if project directory exists
	if _, err := os.Stat(j.projectDir); os.IsNotExist(err) {
		return fmt.Errorf("project directory does not exist: %s", j.projectDir)
	}

	// Check for go.mod
	goModPath := filepath.Join(j.projectDir, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		return errors.New("go.mod not found in project directory")
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
		return errors.New("no main.go found in project (expected at main.go or cmd/*/main.go)")
	}

	j.logger.Info("Project structure validation passed")

	return nil
}

// GenerationPreflightJob fails the workflow before anything is written when a
// target artifact already exists, so generation is atomic: it never leaves a
// project half-generated because a later target collided with an existing file.
type GenerationPreflightJob struct {
	noOpRollbackHelper

	id      string
	force   bool
	targets []string
}

// NewGenerationPreflightJob creates a job that verifies none of the given
// target paths exist unless force is set. It also warns when GitHub owner or
// repository resolution fell back to placeholders.
func NewGenerationPreflightJob(force bool, targets []string, logger *log.Logger) *GenerationPreflightJob {
	return &GenerationPreflightJob{
		logger:  logger,
		name:    "Generation preflight",
		id:      "generation-preflight",
		force:   force,
		targets: targets,
	}
}

func (j *GenerationPreflightJob) ID() string {
	return j.id
}

func (j *GenerationPreflightJob) Name() string {
	return "Check Generation Targets"
}

func (j *GenerationPreflightJob) Execute(ctx context.Context) error {
	if err := checkJobContext(ctx, j.logger, "Checking generation targets"); err != nil {
		return err
	}

	if !j.force {
		var existing []string

		for _, target := range j.targets {
			if _, err := os.Stat(target); err == nil {
				existing = append(existing, target)
			}
		}

		if len(existing) > 0 {
			return fmt.Errorf(
				"refusing to overwrite existing file(s): %s (use --force to overwrite)",
				strings.Join(existing, ", "),
			)
		}
	}

	if types.HasPlaceholderGitHubTarget() {
		j.logger.Warnf(
			"GitHub owner/repo not detected: no git remote found and no --github-owner/--github-repo given. "+
				"The generated release section targets the placeholder repository %q/%q; "+
				"add a git remote or re-run with --github-owner and --github-repo before releasing.",
			types.PlaceholderGitHubOwner, types.PlaceholderGitHubRepo,
		)
	}

	return nil
}

// generationTargets lists every file the wizard may write for the given
// configuration. includeActions controls whether the GitHub Actions workflow
// is among the targets (the config-only workflow never writes it).
func generationTargets(config *domain.SafeProjectConfig, includeActions bool) []string {
	targets := []string{goreleaserConfigFilename}

	if includeActions && config.GetGenerateActions() {
		targets = append(targets, releaseWorkflowTargetPath)
	}

	if config.DockerSupport.IsEnabled() {
		targets = append(targets, dockerfileFilename)
	}

	return targets
}

// DependencyCheckJob checks for required dependencies.
type DependencyCheckJob struct {
	noOpRollbackHelper

	id           string
	dependencies []string
}

// NewDependencyCheckJob creates a new dependency check job.
func NewDependencyCheckJob(dependencies []string, logger *log.Logger) *DependencyCheckJob {
	return &DependencyCheckJob{
		logger:       logger,
		name:         "Dependency check",
		id:           "dependency-check",
		dependencies: dependencies,
	}
}

func (j *DependencyCheckJob) ID() string {
	return j.id
}

func (j *DependencyCheckJob) Name() string {
	return "Check Dependencies"
}

func (j *DependencyCheckJob) Execute(ctx context.Context) error {
	if err := checkJobContext(ctx, j.logger, "Checking dependencies"); err != nil {
		return err
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
	return noOpRollback(j.logger, "Dependency check")(ctx)
}

// JobFactory creates jobs for common wizard operations.
type JobFactory struct {
	logger *log.Logger
}

// NewJobFactory creates a new job factory.
func NewJobFactory(logger *log.Logger) *JobFactory {
	return &JobFactory{
		logger: logger,
	}
}

// CreateFullWizardJobs creates all jobs for a complete wizard operation.
func (jf *JobFactory) CreateFullWizardJobs(config *ProjectConfig, force bool) []Job {
	var jobs []Job

	// Add project validation job
	jobs = append(jobs, NewProjectValidationJob(".", jf.logger))

	// Fail before writing anything if a target artifact already exists
	jobs = append(jobs,
		NewGenerationPreflightJob(force, generationTargets(config, true), jf.logger))

	// Add dependency check job
	dependencies := []string{"go"}
	if config.GetDockerEnabled() {
		dependencies = append(dependencies, "docker")
	}

	if config.SigningLevel.RequiresCosign() {
		dependencies = append(dependencies, "cosign")
	}

	jobs = append(jobs, NewDependencyCheckJob(dependencies, jf.logger))

	// Add config generation job
	jobs = append(jobs, NewConfigGenerationJob(config, force, jf.logger))

	// Add GitHub Actions generation job
	if config.GetGenerateActions() {
		jobs = append(jobs, NewGitHubActionsGenerationJob(config, jf.logger))
	}

	// Add Dockerfile generation job when Docker support is enabled
	if config.DockerSupport.IsEnabled() {
		jobs = append(jobs, NewDockerfileGenerationJob(config, force, jf.logger))
	}

	return jobs
}

// CreateConfigOnlyJobs creates jobs for config generation only.
func (jf *JobFactory) CreateConfigOnlyJobs(config *ProjectConfig, force bool) []Job {
	jobs := []Job{
		NewProjectValidationJob(".", jf.logger),
		NewGenerationPreflightJob(force, generationTargets(config, false), jf.logger),
		NewConfigGenerationJob(config, force, jf.logger),
	}

	// A Docker-enabled config references a Dockerfile; generate it in the
	// config-only path too so the generated config is actually releasable.
	if config.DockerSupport.IsEnabled() {
		jobs = append(jobs, NewDockerfileGenerationJob(config, force, jf.logger))
	}

	return jobs
}

// CreateValidationOnlyJob creates a validation-only job.
func (jf *JobFactory) CreateValidationOnlyJob(projectDir string) *ProjectValidationJob {
	return NewProjectValidationJob(projectDir, jf.logger)
}
