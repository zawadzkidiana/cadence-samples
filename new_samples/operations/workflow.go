package main

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/cadence"
	"go.uber.org/cadence/activity"
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

// CancelWorkflow demonstrates how to handle workflow cancellation gracefully.
func CancelWorkflow(ctx workflow.Context) (retError error) {
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute * 30,
		HeartbeatTimeout:       time.Second * 5,
		WaitForCancellation:    true,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)
	logger.Info("cancel workflow started")

	// Defer cleanup: runs when workflow is cancelled
	defer func() {
		if cadence.IsCanceledError(retError) {
			// When workflow is canceled, get a new disconnected context for cleanup
			newCtx, _ := workflow.NewDisconnectedContext(ctx)
			err := workflow.ExecuteActivity(newCtx, CleanupActivity).Get(ctx, nil)
			if err != nil {
				logger.Error("Cleanup activity failed", zap.Error(err))
				retError = err
				return
			}
			retError = nil
			logger.Info("Workflow completed with cleanup.")
		}
	}()

	// This activity runs until cancelled
	var result string
	err := workflow.ExecuteActivity(ctx, ActivityToBeCanceled).Get(ctx, &result)
	if err != nil && !cadence.IsCanceledError(err) {
		logger.Error("Error from ActivityToBeCanceled", zap.Error(err))
		return err
	}
	logger.Info(fmt.Sprintf("ActivityToBeCanceled returns %v, %v", result, err))

	// This activity will be skipped because context is already cancelled
	err = workflow.ExecuteActivity(ctx, ActivityToBeSkipped).Get(ctx, nil)
	if err != nil && !cadence.IsCanceledError(err) {
		logger.Error("Error from ActivityToBeSkipped", zap.Error(err))
	}

	return err
}

// ActivityToBeCanceled runs indefinitely, sending heartbeats until cancelled.
func ActivityToBeCanceled(ctx context.Context) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("activity started, to cancel workflow, use CLI: 'cadence workflow cancel -w <WorkflowID>'")
	for {
		select {
		case <-time.After(1 * time.Second):
			logger.Info("heart beating...")
			activity.RecordHeartbeat(ctx, "")
		case <-ctx.Done():
			logger.Info("context is cancelled")
			return "I am canceled by Done", ctx.Err()
		}
	}
}

// CleanupActivity performs cleanup operations after cancellation.
func CleanupActivity(ctx context.Context) error {
	logger := activity.GetLogger(ctx)
	logger.Info("cleanupActivity started")
	return nil
}

// ActivityToBeSkipped will not run if the workflow is already cancelled.
func ActivityToBeSkipped(ctx context.Context) error {
	logger := activity.GetLogger(ctx)
	logger.Info("this activity will be skipped due to cancellation")
	return nil
}

