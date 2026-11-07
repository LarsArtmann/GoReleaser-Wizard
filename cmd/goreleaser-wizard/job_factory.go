package main

import (
	"charm.land/log/v2"
)

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

	// Add GitHub Actions generation job (skipped for forgejo/none providers)
	if config.GetGenerateActions() && config.CIProvider.GeneratesGitHubWorkflows() {
		jobs = append(jobs, NewGitHubActionsGenerationJob(config, jf.logger))
	}

	// Add Forgejo Actions generation job
	if config.ShouldGenerateForgejoWorkflows() {
		jobs = append(jobs, NewForgejoWorkflowsGenerationJob(config, jf.logger))
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
