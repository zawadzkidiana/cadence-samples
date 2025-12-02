# Activity Samples

This folder contains samples demonstrating different activity patterns in Cadence.

## Samples Included

### 1. Dynamic Workflow
Demonstrates calling activities by string name instead of function reference.

### 2. Parallel Pick First Workflow
Demonstrates running multiple activities in parallel and using the first result.

## Prerequisites

1. Run the Cadence server (see [main README](../../README.md))
2. Register the `cadence-samples` domain:

```bash
cadence --env development --domain cadence-samples domain register
```

## Steps to Run

### Start the Worker

```bash
go run .
```

This worker handles both workflow types.

---

## Dynamic Workflow

Invoke activities by string name for plugin systems or configuration-driven workflows.

### Start the Workflow

```bash
cadence --env development \
  --domain cadence-samples \
  workflow start \
  --workflow_type cadence_samples.DynamicWorkflow \
  --tl cadence-samples-worker \
  --et 60 \
  --input '{"message":"Cadence"}'
```

### Key Concept

```go
// Instead of: workflow.ExecuteActivity(ctx, MyActivity, input)
// Use:        workflow.ExecuteActivity(ctx, "activity.name.string", input)
```

---

## Parallel Pick First Workflow

Run multiple activities in parallel, use the first result, cancel the rest.

### Start the Workflow

```bash
cadence --env development \
  --domain cadence-samples \
  workflow start \
  --workflow_type cadence_samples.ParallelBranchPickFirstWorkflow \
  --tl cadence-samples-worker \
  --et 60 \
  --input '{}'
```

### Key Concept

```go
selector := workflow.NewSelector(ctx)
selector.AddFuture(f1, handler1)
selector.AddFuture(f2, handler2)
selector.Select(ctx)  // Blocks until one completes
cancelHandler()       // Cancel the rest
```

---

## View Your Workflows

Open [localhost:8088](http://localhost:8088) and click on `cadence-samples` domain.

## References

- [Cadence Documentation](https://cadenceworkflow.io)
- [Activity Basics](https://cadenceworkflow.io/docs/concepts/activities/)

