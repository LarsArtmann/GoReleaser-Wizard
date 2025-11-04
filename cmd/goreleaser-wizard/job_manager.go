package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/charmbracelet/log"
)

// Job represents a wizard operation job
type Job interface {
	ID() string
	Name() string
	Execute(ctx context.Context) error
	Rollback(ctx context.Context) error
}

// JobStatus represents the status of a job
type JobStatus int

const (
	JobStatusPending JobStatus = iota
	JobStatusRunning
	JobStatusCompleted
	JobStatusFailed
	JobStatusRolledBack
)

func (js JobStatus) String() string {
	switch js {
	case JobStatusPending:
		return "pending"
	case JobStatusRunning:
		return "running"
	case JobStatusCompleted:
		return "completed"
	case JobStatusFailed:
		return "failed"
	case JobStatusRolledBack:
		return "rolled_back"
	default:
		return "unknown"
	}
}

// JobResult represents the result of a job execution
type JobResult struct {
	Job      Job
	Status   JobStatus
	Error    error
	Duration time.Duration
	Started  time.Time
	Finished time.Time
	Output   string
}

// JobManager manages the execution of wizard jobs
type JobManager struct {
	jobs        []Job
	results     []JobResult
	mu          sync.Mutex
	logger      *log.Logger
	parallel    bool
	maxJobs     int
	currentJobs int
}

// NewJobManager creates a new job manager
func NewJobManager(logger *log.Logger) *JobManager {
	return &JobManager{
		jobs:        make([]Job, 0),
		results:     make([]JobResult, 0),
		logger:      logger,
		parallel:    false,
		maxJobs:     3, // Default max parallel jobs
		currentJobs: 0,
	}
}

// SetParallel sets whether jobs should run in parallel
func (jm *JobManager) SetParallel(parallel bool) {
	jm.mu.Lock()
	defer jm.mu.Unlock()
	jm.parallel = parallel
}

// SetMaxJobs sets the maximum number of parallel jobs
func (jm *JobManager) SetMaxJobs(maxJobs int) {
	jm.mu.Lock()
	defer jm.mu.Unlock()
	jm.maxJobs = maxJobs
}

// AddJob adds a job to the manager
func (jm *JobManager) AddJob(job Job) {
	jm.mu.Lock()
	defer jm.mu.Unlock()
	jm.jobs = append(jm.jobs, job)
}

// ExecuteJobs executes all jobs according to the manager settings
func (jm *JobManager) ExecuteJobs(ctx context.Context) error {
	if jm.parallel {
		return jm.executeParallel(ctx)
	}
	return jm.executeSequential(ctx)
}

// executeSequential executes jobs one by one
func (jm *JobManager) executeSequential(ctx context.Context) error {
	for _, job := range jm.jobs {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		result := jm.executeJob(ctx, job)
		jm.addResult(result)

		if result.Status == JobStatusFailed {
			return fmt.Errorf("job %s failed: %w", job.Name(), result.Error)
		}
	}

	return nil
}

// executeParallel executes jobs in parallel with concurrency limits
func (jm *JobManager) executeParallel(ctx context.Context) error {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, jm.maxJobs)
	errChan := make(chan error, len(jm.jobs))

	for _, job := range jm.jobs {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		wg.Add(1)
		go func(j Job) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			result := jm.executeJob(ctx, j)
			jm.addResult(result)

			if result.Status == JobStatusFailed {
				errChan <- fmt.Errorf("job %s failed: %w", j.Name(), result.Error)
			}
		}(job)
	}

	wg.Wait()
	close(errChan)

	// Check for errors
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

// executeJob executes a single job and records the result
func (jm *JobManager) executeJob(ctx context.Context, job Job) JobResult {
	start := time.Now()

	// Update job status
	jm.updateJobStatus(job.ID(), JobStatusRunning)
	jm.logger.Infof("Executing job: %s", job.Name())

	// Execute the job
	err := job.Execute(ctx)
	duration := time.Since(start)
	finished := time.Now()

	status := JobStatusCompleted
	if err != nil {
		status = JobStatusFailed
		jm.logger.Errorf("Job %s failed: %v", job.Name(), err)
	} else {
		jm.logger.Infof("Job %s completed successfully", job.Name())
	}

	result := JobResult{
		Job:      job,
		Status:   status,
		Error:    err,
		Duration: duration,
		Started:  start,
		Finished: finished,
		Output:   fmt.Sprintf("Job %s %s", job.Name(), status),
	}

	// Update job status
	jm.updateJobStatus(job.ID(), status)

	return result
}

// updateJobStatus updates the status of a job (for UI updates)
func (jm *JobManager) updateJobStatus(jobID string, status JobStatus) {
	// This could be extended to update a UI or event system
	jm.logger.Debugf("Job %s status: %s", jobID, status)
}

// addResult adds a result to the results list
func (jm *JobManager) addResult(result JobResult) {
	jm.mu.Lock()
	defer jm.mu.Unlock()
	jm.results = append(jm.results, result)
}

// GetResults returns all job results
func (jm *JobManager) GetResults() []JobResult {
	jm.mu.Lock()
	defer jm.mu.Unlock()

	results := make([]JobResult, len(jm.results))
	copy(results, jm.results)
	return results
}

// GetCompletedResults returns only completed job results
func (jm *JobManager) GetCompletedResults() []JobResult {
	jm.mu.Lock()
	defer jm.mu.Unlock()

	completed := make([]JobResult, 0)
	for _, result := range jm.results {
		if result.Status == JobStatusCompleted {
			completed = append(completed, result)
		}
	}
	return completed
}

// GetFailedResults returns only failed job results
func (jm *JobManager) GetFailedResults() []JobResult {
	jm.mu.Lock()
	defer jm.mu.Unlock()

	failed := make([]JobResult, 0)
	for _, result := range jm.results {
		if result.Status == JobStatusFailed {
			failed = append(failed, result)
		}
	}
	return failed
}

// RollbackFailedJobs rolls back all failed jobs
func (jm *JobManager) RollbackFailedJobs(ctx context.Context) error {
	failed := jm.GetFailedResults()

	jm.logger.Infof("Rolling back %d failed jobs", len(failed))

	for i := len(failed) - 1; i >= 0; i-- { // Rollback in reverse order
		result := failed[i]

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		jm.logger.Infof("Rolling back job: %s", result.Job.Name())

		err := result.Job.Rollback(ctx)
		if err != nil {
			jm.logger.Errorf("Failed to rollback job %s: %v", result.Job.Name(), err)
			// Continue with other rollbacks
		} else {
			jm.updateJobStatus(result.Job.ID(), JobStatusRolledBack)
			jm.logger.Infof("Successfully rolled back job: %s", result.Job.Name())
		}
	}

	return nil
}

// Clear clears all jobs and results
func (jm *JobManager) Clear() {
	jm.mu.Lock()
	defer jm.mu.Unlock()
	jm.jobs = make([]Job, 0)
	jm.results = make([]JobResult, 0)
}

// GetStatistics returns job execution statistics
func (jm *JobManager) GetStatistics() map[string]interface{} {
	jm.mu.Lock()
	defer jm.mu.Unlock()

	stats := map[string]interface{}{
		"total_jobs":     len(jm.jobs),
		"total_results":  len(jm.results),
		"completed":      0,
		"failed":         0,
		"total_duration": time.Duration(0),
	}

	var totalDuration time.Duration
	for _, result := range jm.results {
		totalDuration += result.Duration
		switch result.Status {
		case JobStatusCompleted:
			stats["completed"] = stats["completed"].(int) + 1
		case JobStatusFailed:
			stats["failed"] = stats["failed"].(int) + 1
		}
	}

	stats["total_duration"] = totalDuration
	stats["average_duration"] = time.Duration(0)
	if len(jm.results) > 0 {
		stats["average_duration"] = totalDuration / time.Duration(len(jm.results))
	}

	return stats
}
