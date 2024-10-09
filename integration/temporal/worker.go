package temporal

import (
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

/*
iworker is the internal worker used as Temporal worker. It implements the Worker
interface and allows to wrap the Temporal's worker functions for following best
practices.
*/
type iworker struct {
	worker worker.Worker
}

/*
Worker exposes an opinionated way to interact with Temporal's worker capabilities.
*/
type Worker interface {
	RegisterWorkflow(w any, opts workflow.RegisterOptions)
	RegisterActivity(a any, opts activity.RegisterOptions)
}

/*
RegisterWorkflow registers a workflow function with the worker.
*/
func (iw *iworker) RegisterWorkflow(w any, opts workflow.RegisterOptions) {
	iw.worker.RegisterWorkflowWithOptions(w, opts)
}

/*
RegisterActivity registers an activity function with the worker.
*/
func (iw *iworker) RegisterActivity(a any, opts activity.RegisterOptions) {
	iw.worker.RegisterActivityWithOptions(a, opts)
}
