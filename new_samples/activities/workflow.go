package main

import (
	"context"
	"time"

	"go.uber.org/cadence/activity"
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

// =============================================================================
// Dynamic Workflow
// =============================================================================

// DynamicGreetingActivityName is the registered name for the activity.
// This demonstrates how to invoke activities by string name rather than function reference.
const DynamicGreetingActivityName = "cadence_samples.DynamicGreetingActivity"

type dynamicWorkflowInput struct {
	Message string `json:"message"`
}

// DynamicWorkflow demonstrates calling activities using string names for dynamic behavior.
// Instead of passing the function directly to ExecuteActivity, we pass the activity name.
// This is useful for plugin systems or configuration-driven workflows.
func DynamicWorkflow(ctx workflow.Context, input dynamicWorkflowInput) (string, error) {
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	logger.Info("DynamicWorkflow started")

	var greetingMsg string
	// Note: We pass the activity NAME (string) instead of the function reference
	err := workflow.ExecuteActivity(ctx, DynamicGreetingActivityName, input.Message).Get(ctx, &greetingMsg)
	if err != nil {
		logger.Error("DynamicGreetingActivity failed", zap.Error(err))
		return "", err
	}

	logger.Info("Workflow result", zap.String("greeting", greetingMsg))
	return greetingMsg, nil
}

// DynamicGreetingActivity is a simple activity that returns a greeting message.
func DynamicGreetingActivity(ctx context.Context, message string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("DynamicGreetingActivity started.")
	return "Hello, " + message, nil
}

// =============================================================================
// Parallel Pick First Workflow
// =============================================================================

type parallelBranchInput struct {
	Message string `json:"message"`
}

// ParallelBranchPickFirstWorkflow demonstrates running multiple activities in parallel
// and using the result from whichever completes first.
func ParallelBranchPickFirstWorkflow(ctx workflow.Context) (string, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("ParallelBranchPickFirstWorkflow started")

	selector := workflow.NewSelector(ctx)
	var firstResp string

	// Use a cancel handler to cancel all remaining activities once one completes
	childCtx, cancelHandler := workflow.WithCancel(ctx)
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
		HeartbeatTimeout:       time.Second * 20,
		WaitForCancellation:    true,
	}
	childCtx = workflow.WithActivityOptions(childCtx, ao)

	// Run two activities in parallel with different delays
	f1 := workflow.ExecuteActivity(childCtx, ParallelActivity, parallelBranchInput{Message: "first activity"}, time.Second*10)
	f2 := workflow.ExecuteActivity(childCtx, ParallelActivity, parallelBranchInput{Message: "second activity"}, time.Second*2)
	pendingFutures := []workflow.Future{f1, f2}

	selector.AddFuture(f1, func(f workflow.Future) {
		f.Get(ctx, &firstResp)
	}).AddFuture(f2, func(f workflow.Future) {
		f.Get(ctx, &firstResp)
	})

	// Wait for any of the futures to complete
	selector.Select(ctx)

	// Cancel all other pending activities
	cancelHandler()

	// Wait for pending activities to finish cancellation
	for _, f := range pendingFutures {
		err := f.Get(ctx, &firstResp)
		if err != nil {
			return "", err
		}
	}

	logger.Info("ParallelBranchPickFirstWorkflow completed")
	return firstResp, nil
}

// ParallelActivity is a simple activity that returns a greeting.
func ParallelActivity(ctx context.Context, input parallelBranchInput) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("ParallelActivity started")
	return "Hello " + input.Message, nil
}

