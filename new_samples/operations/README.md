# Operations Samples

This folder contains samples demonstrating operational patterns in Cadence.

## Samples Included

### Cancel Workflow
Demonstrates how to handle workflow cancellation gracefully, including cleanup operations.

## Prerequisites

1. Run the Cadence server (see [main README](../../README.md))
2. Register the `cadence-samples` domain:

```bash
cadence --env development --domain cadence-samples domain register
```

## Steps to Run

### 1. Start the Worker

```bash
go run .
```

### 2. Start the Workflow

```bash
cadence --env development \
  --domain cadence-samples \
  workflow start \
  --workflow_type cadence_samples.CancelWorkflow \
  --tl cadence-samples-worker \
  --et 60 \
  --input '{}'
```

Copy the workflow ID from the output.

### 3. Cancel the Workflow

```bash
cadence --env development \
  --domain cadence-samples \
  workflow cancel \
  --workflow_id <YOUR_WORKFLOW_ID>
```

### 4. Observe the Behavior

Watch the worker logs to see:
1. Activity sending heartbeats
2. Cancellation received
3. Cleanup activity running

## Key Concept

When a workflow is cancelled, use `NewDisconnectedContext` for cleanup:

```go
if cadence.IsCanceledError(retError) {
    newCtx, _ := workflow.NewDisconnectedContext(ctx)
    workflow.ExecuteActivity(newCtx, CleanupActivity).Get(ctx, nil)
}
```

## View Your Workflow

Open [localhost:8088](http://localhost:8088) and click on `cadence-samples` domain.

## References

- [Cadence Documentation](https://cadenceworkflow.io)
- [Workflow Cancellation](https://cadenceworkflow.io/docs/go-client/cancellation/)

