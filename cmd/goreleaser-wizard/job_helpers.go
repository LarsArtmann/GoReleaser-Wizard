package main

import (
	"context"

	"charm.land/log/v2"
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

func checkJobContext(ctx context.Context, logger *log.Logger, message string) error {
	logger.Info(message)

	return checkContextCancellation(ctx, "context cancelled")
}
