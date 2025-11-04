package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/charmbracelet/log"
)

// TestJobManager tests the job manager functionality
func TestJobManager(t *testing.T) {
	logger := log.New(os.Stderr)
	jm := NewJobManager(logger)

	// Test basic job manager
	tests := []struct {
		name     string
		setup    func(*JobManager)
		parallel bool
		wantErr  bool
	}{
		{
			name: "sequential_success",
			setup: func(jm *JobManager) {
				jm.SetParallel(false)
				jm.AddJob(NewProjectValidationJob(".", logger))
			},
			parallel: false,
			wantErr:  false,
		},
		{
			name: "parallel_success",
			setup: func(jm *JobManager) {
				jm.SetParallel(true)
				jm.SetMaxJobs(2)
				jm.AddJob(NewProjectValidationJob(".", logger))
			},
			parallel: true,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(jm)

			// Execute jobs
			ctx := context.Background()
			err := jm.ExecuteJobs(ctx)

			if (err != nil) != tt.wantErr {
				t.Errorf("ExecuteJobs() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Check results
			results := jm.GetResults()
			if len(results) == 0 {
				t.Error("Expected at least one result")
			}

			// Clear for next test
			jm.Clear()
		})
	}
}

// TestJobExecution tests individual job execution
func TestJobExecution(t *testing.T) {
	logger := log.New(os.Stderr)

	// Create temporary directory for testing
	tmpDir, _ := os.MkdirTemp("", "job-test")
	defer os.RemoveAll(tmpDir)

	// Create basic project
	goMod := `module github.com/user/job-test
go 1.21
`
	os.WriteFile(tmpDir+"/go.mod", []byte(goMod), 0644)
	os.WriteFile(tmpDir+"/main.go", []byte("package main\n\nfunc main() {}"), 0644)

	// Change to test directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tmpDir)

	tests := []struct {
		name    string
		job     Job
		wantErr bool
	}{
		{
			name:    "project_validation_success",
			job:     NewProjectValidationJob(".", logger),
			wantErr: false,
		},
		{
			name: "config_generation_success",
			job: NewConfigGenerationJob(&ProjectConfig{
				ProjectName: "job-test",
				BinaryName:  "job-test",
				MainPath:    ".",
				GitProvider: "GitHub",
			}, false, logger),
			wantErr: false,
		},
		{
			name:    "dependency_check_failure",
			job:     NewDependencyCheckJob([]string{"nonexistent-binary-xyz"}, logger),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := tt.job.Execute(ctx)

			if (err != nil) != tt.wantErr {
				t.Errorf("%s.Execute() error = %v, wantErr %v", tt.job.Name(), err, tt.wantErr)
			}
		})
	}
}

// TestJobRollback tests job rollback functionality
func TestJobRollback(t *testing.T) {
	logger := log.New(os.Stderr)

	// Create temporary directory for testing
	tmpDir, _ := os.MkdirTemp("", "rollback-test")
	defer os.RemoveAll(tmpDir)

	// Change to test directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tmpDir)

	// Create basic project
	goMod := `module github.com/user/rollback-test
go 1.21
`
	os.WriteFile(tmpDir+"/go.mod", []byte(goMod), 0644)
	os.WriteFile(tmpDir+"/main.go", []byte("package main\n\nfunc main() {}"), 0644)

	tests := []struct {
		name            string
		job             Job
		executeRollback bool
	}{
		{
			name: "config_generation_rollback",
			job: NewConfigGenerationJob(&ProjectConfig{
				ProjectName: "rollback-test",
				BinaryName:  "rollback-test",
				MainPath:    ".",
				GitProvider: "GitHub",
			}, false, logger),
			executeRollback: true,
		},
		{
			name:            "project_validation_rollback",
			job:             NewProjectValidationJob(".", logger),
			executeRollback: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// Execute job first
			err := tt.job.Execute(ctx)
			if err != nil && tt.job.Name() != "Dependency Check" {
				t.Errorf("%s.Execute() failed: %v", tt.job.Name(), err)
				return
			}

			// Test rollback
			if tt.executeRollback {
				err = tt.job.Rollback(ctx)
				if err != nil {
					t.Errorf("%s.Rollback() failed: %v", tt.job.Name(), err)
				}
			}
		})
	}
}

// TestWorkflow tests workflow functionality
func TestWorkflow(t *testing.T) {
	logger := log.New(os.Stderr)

	// Create temporary directory for testing
	tmpDir, _ := os.MkdirTemp("", "workflow-test")
	defer os.RemoveAll(tmpDir)

	// Create basic project
	goMod := `module github.com/user/workflow-test
go 1.21
`
	os.WriteFile(tmpDir+"/go.mod", []byte(goMod), 0644)
	os.WriteFile(tmpDir+"/main.go", []byte("package main\n\nfunc main() {}"), 0644)

	// Change to test directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tmpDir)

	tests := []struct {
		name     string
		workflow *Workflow
		wantErr  bool
	}{
		{
			name: "validation_only_workflow",
			workflow: func() *Workflow {
				wf := NewWorkflow("Validation Test", "Test validation workflow", logger)
				wf.JobManager.AddJob(NewProjectValidationJob(".", logger))
				wf.SetTimeout(5 * time.Minute)
				return wf
			}(),
			wantErr: false,
		},
		{
			name: "config_only_workflow",
			workflow: func() *Workflow {
				wf := NewWorkflow("Config Test", "Test config generation workflow", logger)
				wf.JobManager.AddJob(NewProjectValidationJob(".", logger))
				wf.JobManager.AddJob(NewConfigGenerationJob(&ProjectConfig{
					ProjectName: "workflow-test",
					BinaryName:  "workflow-test",
					MainPath:    ".",
					GitProvider: "GitHub",
				}, false, logger))
				wf.SetTimeout(5 * time.Minute)
				return wf
			}(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := tt.workflow.Execute(ctx)

			if (err != nil) != tt.wantErr {
				t.Errorf("Workflow.Execute() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Check results
			results := tt.workflow.GetResults()
			if len(results) == 0 {
				t.Error("Expected at least one result")
			}

			// Check statistics
			stats := tt.workflow.GetStatistics()
			if stats["total_results"] == nil {
				t.Error("Expected statistics to contain total_results")
			}
		})
	}
}

// TestWorkflowBuilder tests workflow builder functionality
func TestWorkflowBuilder(t *testing.T) {
	logger := log.New(os.Stderr)
	wb := NewWorkflowBuilder(logger)

	// Create temporary directory for testing
	tmpDir, _ := os.MkdirTemp("", "builder-test")
	defer os.RemoveAll(tmpDir)

	// Create basic project
	goMod := `module github.com/user/builder-test
go 1.21
`
	os.WriteFile(tmpDir+"/go.mod", []byte(goMod), 0644)
	os.WriteFile(tmpDir+"/main.go", []byte("package main\n\nfunc main() {}"), 0644)

	// Change to test directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tmpDir)

	config := &ProjectConfig{
		ProjectName:        "builder-test",
		ProjectDescription: "A test project for workflow builder",
		ProjectType:        "CLI Application",
		BinaryName:         "builder-test",
		MainPath:           ".",
		Platforms:          []string{"linux", "darwin"},
		Architectures:      []string{"amd64"},
		CGOEnabled:         false,
		GitProvider:        "GitHub",
		GenerateActions:    true,
		ActionsOn:          []string{"On version tags (v*)"},
	}

	tests := []struct {
		name         string
		workflowType WorkflowType
		force        bool
		wantErr      bool
	}{
		{
			name:         "full_wizard_workflow",
			workflowType: WorkflowTypeFullWizard,
			force:        false,
			wantErr:      false,
		},
		{
			name:         "config_only_workflow",
			workflowType: WorkflowTypeConfigOnly,
			force:        false,
			wantErr:      false,
		},
		{
			name:         "validation_only_workflow",
			workflowType: WorkflowTypeValidationOnly,
			force:        false,
			wantErr:      false,
		},
		{
			name:         "unsupported_workflow",
			workflowType: "unsupported",
			force:        false,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workflow, err := wb.BuildWorkflow(tt.workflowType, config, tt.force)

			if (err != nil) != tt.wantErr {
				t.Errorf("BuildWorkflow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Execute workflow
				ctx := context.Background()
				err = workflow.Execute(ctx)
				if err != nil {
					t.Errorf("Workflow.Execute() error = %v", err)
				}
			}
		})
	}
}

// TestConcurrentJobExecution tests concurrent job execution
func TestConcurrentJobExecution(t *testing.T) {
	logger := log.New(os.Stderr)
	jm := NewJobManager(logger)

	// Create temporary directory for testing
	tmpDir, _ := os.MkdirTemp("", "concurrent-test")
	defer os.RemoveAll(tmpDir)

	// Create multiple test projects
	for i := range 3 {
		projectDir := filepath.Join(tmpDir, fmt.Sprintf("project-%d", i))
		os.MkdirAll(projectDir, 0755)

		goMod := fmt.Sprintf("module github.com/user/concurrent-test-%d\ngo 1.21\n", i)
		os.WriteFile(filepath.Join(projectDir, "go.mod"), []byte(goMod), 0644)
		os.WriteFile(filepath.Join(projectDir, "main.go"), []byte("package main\n\nfunc main() {}"), 0644)
	}

	// Change to test directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tmpDir)

	// Add validation jobs for all projects
	for i := range 3 {
		projectDir := fmt.Sprintf("project-%d", i)
		jm.AddJob(NewProjectValidationJob(projectDir, logger))
	}

	// Set parallel execution
	jm.SetParallel(true)
	jm.SetMaxJobs(2)

	// Execute jobs
	ctx := context.Background()
	err := jm.ExecuteJobs(ctx)

	if err != nil {
		t.Errorf("Concurrent job execution failed: %v", err)
	}

	// Check results
	results := jm.GetResults()
	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}

	// Check that all jobs completed
	completed := jm.GetCompletedResults()
	if len(completed) != 3 {
		t.Errorf("Expected 3 completed jobs, got %d", len(completed))
	}
}
