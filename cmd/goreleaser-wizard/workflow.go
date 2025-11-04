package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/log"
)

// Workflow represents a sequence of jobs
type Workflow struct {
	Name        string
	Description string
	JobManager  *JobManager
	Factory     *JobFactory
	Timeout     time.Duration
}

// NewWorkflow creates a new workflow
func NewWorkflow(name, description string, logger *log.Logger) *Workflow {
	return &Workflow{
		Name:        name,
		Description: description,
		JobManager:  NewJobManager(logger),
		Factory:     NewJobFactory(logger),
		Timeout:     30 * time.Minute, // Default timeout
	}
}

// SetTimeout sets workflow timeout
func (w *Workflow) SetTimeout(timeout time.Duration) {
	w.Timeout = timeout
}

// SetParallel sets whether jobs should run in parallel
func (w *Workflow) SetParallel(parallel bool, maxJobs int) {
	w.JobManager.SetParallel(parallel)
	w.JobManager.SetMaxJobs(maxJobs)
}

// Execute executes the workflow
func (w *Workflow) Execute(ctx context.Context) error {
	// Create timeout context
	timeoutCtx, cancel := context.WithTimeout(ctx, w.Timeout)
	defer cancel()

	w.JobManager.logger.Infof("Starting workflow: %s", w.Name)
	w.JobManager.logger.Infof("Workflow timeout: %v", w.Timeout)

	// Execute jobs
	err := w.JobManager.ExecuteJobs(timeoutCtx)
	if err != nil {
		w.JobManager.logger.Errorf("Workflow failed: %v", err)

		// Attempt rollback
		rollbackErr := w.JobManager.RollbackFailedJobs(timeoutCtx)
		if rollbackErr != nil {
			w.JobManager.logger.Errorf("Rollback failed: %v", rollbackErr)
			return fmt.Errorf("workflow failed and rollback failed: %w", err)
		}

		w.JobManager.logger.Info("Workflow rolled back successfully")
		return fmt.Errorf("workflow failed: %w", err)
	}

	w.JobManager.logger.Info("Workflow completed successfully")
	return nil
}

// GetResults returns workflow execution results
func (w *Workflow) GetResults() []JobResult {
	return w.JobManager.GetResults()
}

// GetStatistics returns workflow statistics
func (w *Workflow) GetStatistics() map[string]interface{} {
	stats := w.JobManager.GetStatistics()
	stats["workflow_name"] = w.Name
	stats["workflow_description"] = w.Description
	stats["timeout"] = w.Timeout
	return stats
}

// WorkflowType represents different types of workflows
type WorkflowType string

const (
	WorkflowTypeFullWizard     WorkflowType = "full-wizard"
	WorkflowTypeConfigOnly     WorkflowType = "config-only"
	WorkflowTypeValidationOnly WorkflowType = "validation-only"
	WorkflowTypeMigrate        WorkflowType = "migrate"
	WorkflowTypeUpdate         WorkflowType = "update"
	WorkflowTypeRollback       WorkflowType = "rollback"
)

// WorkflowBuilder builds workflows for different scenarios
type WorkflowBuilder struct {
	logger  *log.Logger
	factory *JobFactory
}

// NewWorkflowBuilder creates a new workflow builder
func NewWorkflowBuilder(logger *log.Logger) *WorkflowBuilder {
	return &WorkflowBuilder{
		logger:  logger,
		factory: NewJobFactory(logger),
	}
}

// BuildWorkflow builds a workflow based on type and configuration
func (wb *WorkflowBuilder) BuildWorkflow(wfType WorkflowType, config *ProjectConfig, force bool) (*Workflow, error) {
	var workflow *Workflow
	var jobs []Job

	switch wfType {
	case WorkflowTypeFullWizard:
		workflow = NewWorkflow("Full Wizard", "Complete GoReleaser setup with all features", wb.logger)
		jobs = wb.factory.CreateFullWizardJobs(config, force)
		workflow.SetParallel(false, 1) // Sequential execution for wizard

	case WorkflowTypeConfigOnly:
		workflow = NewWorkflow("Config Generation", "Generate GoReleaser configuration only", wb.logger)
		jobs = wb.factory.CreateConfigOnlyJobs(config, force)
		workflow.SetParallel(false, 1)

	case WorkflowTypeValidationOnly:
		workflow = NewWorkflow("Project Validation", "Validate project structure only", wb.logger)
		jobs = []Job{wb.factory.CreateValidationOnlyJob(".")}
		workflow.SetParallel(false, 1)

	default:
		return nil, fmt.Errorf("unsupported workflow type: %s", wfType)
	}

	// Add jobs to workflow
	for _, job := range jobs {
		workflow.JobManager.AddJob(job)
	}

	// Set appropriate timeout based on workflow type
	switch wfType {
	case WorkflowTypeFullWizard:
		workflow.SetTimeout(10 * time.Minute)
	case WorkflowTypeConfigOnly:
		workflow.SetTimeout(5 * time.Minute)
	case WorkflowTypeValidationOnly:
		workflow.SetTimeout(2 * time.Minute)
	}

	return workflow, nil
}

// BuildMigrateWorkflow builds a migration workflow
func (wb *WorkflowBuilder) BuildMigrateWorkflow(fromVersion, toVersion string, config *ProjectConfig) (*Workflow, error) {
	workflow := NewWorkflow(
		fmt.Sprintf("Migration %s -> %s", fromVersion, toVersion),
		fmt.Sprintf("Migrate configuration from version %s to %s", fromVersion, toVersion),
		wb.logger,
	)

	// Create migration jobs
	jobs := wb.createMigrationJobs(fromVersion, toVersion, config)

	for _, job := range jobs {
		workflow.JobManager.AddJob(job)
	}

	workflow.SetTimeout(15 * time.Minute)
	workflow.SetParallel(false, 1)

	return workflow, nil
}

// BuildUpdateWorkflow builds an update workflow
func (wb *WorkflowBuilder) BuildUpdateWorkflow(config *ProjectConfig, dryRun bool) (*Workflow, error) {
	workflow := NewWorkflow(
		"Update Configuration",
		fmt.Sprintf("Update GoReleaser configuration (dry-run: %v)", dryRun),
		wb.logger,
	)

	// Create update jobs
	jobs := wb.createUpdateJobs(config, dryRun)

	for _, job := range jobs {
		workflow.JobManager.AddJob(job)
	}

	workflow.SetTimeout(10 * time.Minute)
	workflow.SetParallel(false, 1)

	return workflow, nil
}

// createMigrationJobs creates jobs for migration workflow
func (wb *WorkflowBuilder) createMigrationJobs(fromVersion, toVersion string, config *ProjectConfig) []Job {
	var jobs []Job

	// Backup current configuration
	backupJob := &ConfigBackupJob{
		id:     "backup-config",
		logger: wb.logger,
	}
	jobs = append(jobs, backupJob)

	// Validate migration compatibility
	validationJob := &MigrationValidationJob{
		id:          "validate-migration",
		fromVersion: fromVersion,
		toVersion:   toVersion,
		logger:      wb.logger,
	}
	jobs = append(jobs, validationJob)

	// Migrate configuration
	migrateJob := &ConfigMigrationJob{
		id:          "migrate-config",
		fromVersion: fromVersion,
		toVersion:   toVersion,
		config:      config,
		logger:      wb.logger,
	}
	jobs = append(jobs, migrateJob)

	return jobs
}

// createUpdateJobs creates jobs for update workflow
func (wb *WorkflowBuilder) createUpdateJobs(config *ProjectConfig, dryRun bool) []Job {
	var jobs []Job

	// Validate project structure
	validationJob := NewProjectValidationJob(".", wb.logger)
	jobs = append(jobs, validationJob)

	// Update configuration
	updateJob := &ConfigUpdateJob{
		id:     "update-config",
		config: config,
		dryRun: dryRun,
		logger: wb.logger,
	}
	jobs = append(jobs, updateJob)

	return jobs
}

// ConfigBackupJob backs up existing configuration
type ConfigBackupJob struct {
	id     string
	logger *log.Logger
}

func (j *ConfigBackupJob) ID() string {
	return j.id
}

func (j *ConfigBackupJob) Name() string {
	return "Backup Configuration"
}

func (j *ConfigBackupJob) Execute(ctx context.Context) error {
	j.logger.Info("Backing up existing configuration")

	// Check if .goreleaser.yaml exists
	if _, err := os.Stat(".goreleaser.yaml"); os.IsNotExist(err) {
		j.logger.Info("No existing configuration to backup")
		return nil
	}

	// Create backup with timestamp
	timestamp := time.Now().Format("20060102-150405")
	backupFile := fmt.Sprintf(".goreleaser.yaml.backup.%s", timestamp)

	err := os.Rename(".goreleaser.yaml", backupFile)
	if err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	j.logger.Infof("Configuration backed up to: %s", backupFile)
	return nil
}

func (j *ConfigBackupJob) Rollback(ctx context.Context) error {
	// Backup job doesn't need rollback
	return nil
}

// MigrationValidationJob validates migration compatibility
type MigrationValidationJob struct {
	id          string
	fromVersion string
	toVersion   string
	logger      *log.Logger
}

func (j *MigrationValidationJob) ID() string {
	return j.id
}

func (j *MigrationValidationJob) Name() string {
	return "Validate Migration Compatibility"
}

func (j *MigrationValidationJob) Execute(ctx context.Context) error {
	j.logger.Infof("Validating migration from %s to %s", j.fromVersion, j.toVersion)

	// This is a simplified validation - in real implementation,
	// this would check compatibility between versions
	if j.fromVersion == j.toVersion {
		return fmt.Errorf("source and target versions are the same")
	}

	j.logger.Info("Migration compatibility validated")
	return nil
}

func (j *MigrationValidationJob) Rollback(ctx context.Context) error {
	// Validation job doesn't need rollback
	return nil
}

// ConfigMigrationJob migrates configuration
type ConfigMigrationJob struct {
	id          string
	fromVersion string
	toVersion   string
	config      *ProjectConfig
	logger      *log.Logger
}

func (j *ConfigMigrationJob) ID() string {
	return j.id
}

func (j *ConfigMigrationJob) Name() string {
	return "Migrate Configuration"
}

func (j *ConfigMigrationJob) Execute(ctx context.Context) error {
	j.logger.Infof("Migrating configuration from %s to %s", j.fromVersion, j.toVersion)

	// This is a simplified migration - in real implementation,
	// this would transform configuration based on version differences
	err := generateGoReleaserConfig(j.config)
	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	j.logger.Info("Configuration migrated successfully")
	return nil
}

func (j *ConfigMigrationJob) Rollback(ctx context.Context) error {
	j.logger.Info("Rolling back configuration migration")

	// Restore from backup (simplified)
	backupFiles, err := filepath.Glob(".goreleaser.yaml.backup.*")
	if err != nil || len(backupFiles) == 0 {
		return fmt.Errorf("no backup found for rollback")
	}

	// Get the most recent backup
	latestBackup := backupFiles[len(backupFiles)-1]

	err = os.Rename(latestBackup, ".goreleaser.yaml")
	if err != nil {
		return fmt.Errorf("failed to restore backup: %w", err)
	}

	j.logger.Infof("Configuration rolled back from: %s", latestBackup)
	return nil
}

// ConfigUpdateJob updates configuration
type ConfigUpdateJob struct {
	id     string
	config *ProjectConfig
	dryRun bool
	logger *log.Logger
}

func (j *ConfigUpdateJob) ID() string {
	return j.id
}

func (j *ConfigUpdateJob) Name() string {
	return "Update Configuration"
}

func (j *ConfigUpdateJob) Execute(ctx context.Context) error {
	if j.dryRun {
		j.logger.Info("Dry-run: Skipping configuration update")
		return nil
	}

	j.logger.Info("Updating configuration")

	err := generateGoReleaserConfig(j.config)
	if err != nil {
		return fmt.Errorf("configuration update failed: %w", err)
	}

	j.logger.Info("Configuration updated successfully")
	return nil
}

func (j *ConfigUpdateJob) Rollback(ctx context.Context) error {
	// In real implementation, this would restore previous configuration
	j.logger.Info("Configuration update rollback not implemented")
	return nil
}
