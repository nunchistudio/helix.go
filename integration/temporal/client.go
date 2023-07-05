package temporal

import (
	"context"
	"fmt"

	"go.nunchi.studio/helix/telemetry/trace"

	"go.temporal.io/api/enums/v1"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
)

/*
iclient is the internal client used as Temporal client. It implements the Client
interface and allows to wrap the Temporal's client functions for automatic tracing
and error recording.
*/
type iclient struct {
	config *Config
	client client.Client
}

/*
Client exposes an opinionated way to interact with Temporal's client capabilities.
*/
type Client interface {
	ExecuteWorkflow(ctx context.Context, opts client.StartWorkflowOptions, workflowType string, payload ...any) (client.WorkflowRun, error)
	GetWorkflow(ctx context.Context, workflowID string, runID string) client.WorkflowRun
	SignalWorkflow(ctx context.Context, workflowID string, runID string, signalName string, arg any) error
	SignalWithStartWorkflow(ctx context.Context, workflowID string, signalName string, signalArg any, opts client.StartWorkflowOptions, workflowType string, payload any) (client.WorkflowRun, error)
	CancelWorkflow(ctx context.Context, workflowID string, runID string) error
	TerminateWorkflow(ctx context.Context, workflowID string, runID string, reason string) error
	GetWorkflowHistory(ctx context.Context, workflowID string, runID string, isLongPoll bool, filterType enums.HistoryEventFilterType) client.HistoryEventIterator
	CompleteActivity(ctx context.Context, namespace string, workflowID string, runID string, activityID string, result any, err error) error
	RecordActivityHeartbeat(ctx context.Context, namespace string, workflowID string, runID string, activityID string) error

	ListClosedWorkflow(ctx context.Context, request *workflowservice.ListClosedWorkflowExecutionsRequest) (*workflowservice.ListClosedWorkflowExecutionsResponse, error)
	ListOpenWorkflow(ctx context.Context, request *workflowservice.ListOpenWorkflowExecutionsRequest) (*workflowservice.ListOpenWorkflowExecutionsResponse, error)
	ListWorkflow(ctx context.Context, request *workflowservice.ListWorkflowExecutionsRequest) (*workflowservice.ListWorkflowExecutionsResponse, error)
	ListArchivedWorkflow(ctx context.Context, request *workflowservice.ListArchivedWorkflowExecutionsRequest) (*workflowservice.ListArchivedWorkflowExecutionsResponse, error)
	ScanWorkflow(ctx context.Context, request *workflowservice.ScanWorkflowExecutionsRequest) (*workflowservice.ScanWorkflowExecutionsResponse, error)
	CountWorkflow(ctx context.Context, request *workflowservice.CountWorkflowExecutionsRequest) (*workflowservice.CountWorkflowExecutionsResponse, error)
	GetSearchAttributes(ctx context.Context) (*workflowservice.GetSearchAttributesResponse, error)
	QueryWorkflow(ctx context.Context, request *client.QueryWorkflowWithOptionsRequest) (*client.QueryWorkflowWithOptionsResponse, error)
	DescribeWorkflowExecution(ctx context.Context, workflowID string, runID string) (*workflowservice.DescribeWorkflowExecutionResponse, error)
	DescribeTaskQueue(ctx context.Context, taskqueue string, taskqueueType enums.TaskQueueType) (*workflowservice.DescribeTaskQueueResponse, error)
	ResetWorkflowExecution(ctx context.Context, request *workflowservice.ResetWorkflowExecutionRequest) (*workflowservice.ResetWorkflowExecutionResponse, error)

	CheckHealth(ctx context.Context, request *client.CheckHealthRequest) (*client.CheckHealthResponse, error)

	ScheduleClient() ScheduleClient
}

/*
ExecuteWorkflow starts a workflow execution and return a WorkflowRun instance and
error.

It automatically handles tracing and error recording via interceptor.
*/
func (c *iclient) ExecuteWorkflow(ctx context.Context, opts client.StartWorkflowOptions, workflowType string, payload ...any) (client.WorkflowRun, error) {
	return c.client.ExecuteWorkflow(ctx, opts, workflowType, payload...)
}

/*
GetWorkflow retrieves a workflow execution and return a WorkflowRun instance.

It automatically handles tracing.
*/
func (c *iclient) GetWorkflow(ctx context.Context, workflowID string, runID string) client.WorkflowRun {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: GetWorkflow", humanized))
	defer span.End()

	run := c.client.GetWorkflow(ctx, workflowID, runID)
	setDefaultAttributes(span, c.config)
	setClientRequestAttributes(span, clientRequest{
		workflowID:    workflowID,
		workflowRunID: runID,
	})

	return run
}

/*
SignalWorkflow sends a signals to a workflow in execution.

It automatically handles tracing and error recording via interceptor.
*/
func (c *iclient) SignalWorkflow(ctx context.Context, workflowID string, runID string, signalName string, arg any) error {
	return c.client.SignalWorkflow(ctx, workflowID, runID, signalName, arg)
}

/*
SignalWithStartWorkflow sends a signal to a running workflow. If the workflow is
not running or not found, it starts the workflow and then sends the signal in
transaction.

It automatically handles tracing and error recording via interceptor.
*/
func (c *iclient) SignalWithStartWorkflow(ctx context.Context, workflowID string, signalName string, arg any, opts client.StartWorkflowOptions, workflowType string, payload any) (client.WorkflowRun, error) {
	return c.client.SignalWithStartWorkflow(ctx, workflowID, signalName, arg, opts, workflowType, payload)
}

/*
CancelWorkflow requests cancellation of a workflow in execution. Cancellation
request closes the channel returned by the workflow.Context.Done() of the workflow
that is target of the request.

It automatically handles tracing and error recording.
*/
func (c *iclient) CancelWorkflow(ctx context.Context, workflowID string, runID string) error {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: CancelWorkflow", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to cancel workflow", err)
		}
	}()

	err = c.client.CancelWorkflow(ctx, workflowID, runID)
	setDefaultAttributes(span, c.config)
	setClientRequestAttributes(span, clientRequest{
		workflowID:    workflowID,
		workflowRunID: runID,
	})

	return err
}

/*
TerminateWorkflow terminates a workflow execution. Terminate stops a workflow
execution immediately without letting the workflow to perform any cleanup.

It automatically handles tracing and error recording.
*/
func (c *iclient) TerminateWorkflow(ctx context.Context, workflowID string, runID string, reason string) error {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: TerminateWorkflow", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to terminate workflow", err)
		}
	}()

	err = c.client.TerminateWorkflow(ctx, workflowID, runID, reason)
	setDefaultAttributes(span, c.config)
	setClientRequestAttributes(span, clientRequest{
		workflowID:    workflowID,
		workflowRunID: runID,
	})

	return err
}

/*
GetWorkflowHistory retrieves history events of a particular workflow.

It automatically handles tracing.
*/
func (c *iclient) GetWorkflowHistory(ctx context.Context, workflowID string, runID string, isLongPoll bool, filterType enums.HistoryEventFilterType) client.HistoryEventIterator {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: GetWorkflowHistory", humanized))
	defer span.End()

	iter := c.client.GetWorkflowHistory(ctx, workflowID, runID, isLongPoll, filterType)
	setDefaultAttributes(span, c.config)
	setClientRequestAttributes(span, clientRequest{
		workflowID:    workflowID,
		workflowRunID: runID,
	})

	return iter
}

/*
CompleteActivity reports activity completed.

It automatically handles tracing and error recording.
*/
func (c *iclient) CompleteActivity(ctx context.Context, namespace string, workflowID string, runID string, activityID string, result any, encountered error) error {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: CompleteActivity", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to complete activity", err)
		}
	}()

	err = c.client.CompleteActivityByID(ctx, namespace, workflowID, runID, activityID, result, encountered)
	setDefaultAttributes(span, c.config)
	setClientRequestAttributes(span, clientRequest{
		namespace:     namespace,
		workflowID:    workflowID,
		workflowRunID: runID,
		activityID:    activityID,
	})

	return err
}

/*
RecordActivityHeartbeat records heartbeat for an activity.

It automatically handles tracing and error recording.
*/
func (c *iclient) RecordActivityHeartbeat(ctx context.Context, namespace string, workflowID string, runID string, activityID string) error {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: RecordActivityHeartbeat", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to record activity heartbeat", err)
		}
	}()

	err = c.client.RecordActivityHeartbeatByID(ctx, namespace, workflowID, runID, activityID)
	setDefaultAttributes(span, c.config)
	setClientRequestAttributes(span, clientRequest{
		namespace:     namespace,
		workflowID:    workflowID,
		workflowRunID: runID,
		activityID:    activityID,
	})

	return err
}

/*
ListClosedWorkflow retrieves closed workflow executions based on request filters.
Retrieved workflow executions are sorted by close time in descending order.

It automatically handles tracing and error recording.
*/
func (c *iclient) ListClosedWorkflow(ctx context.Context, request *workflowservice.ListClosedWorkflowExecutionsRequest) (*workflowservice.ListClosedWorkflowExecutionsResponse, error) {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: ListClosedWorkflow", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to list closed workflows", err)
		}
	}()

	res, err := c.client.ListClosedWorkflow(ctx, request)
	setDefaultAttributes(span, c.config)
	if request != nil {
		setClientRequestAttributes(span, clientRequest{
			namespace: request.Namespace,
		})
	}

	return res, err
}

/*
ListOpenWorkflow retrieves open workflow executions based on request filters.
Retrieved workflow executions are sorted by start time in descending order.

It automatically handles tracing and error recording.
*/
func (c *iclient) ListOpenWorkflow(ctx context.Context, request *workflowservice.ListOpenWorkflowExecutionsRequest) (*workflowservice.ListOpenWorkflowExecutionsResponse, error) {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: ListOpenWorkflow", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to list open workflows", err)
		}
	}()

	res, err := c.client.ListOpenWorkflow(ctx, request)
	setDefaultAttributes(span, c.config)
	if request != nil {
		setClientRequestAttributes(span, clientRequest{
			namespace: request.Namespace,
		})
	}

	return res, err
}

/*
ListWorkflow retrieves workflow executions based on query.

It automatically handles tracing and error recording.
*/
func (c *iclient) ListWorkflow(ctx context.Context, request *workflowservice.ListWorkflowExecutionsRequest) (*workflowservice.ListWorkflowExecutionsResponse, error) {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: ListWorkflow", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to list workflows", err)
		}
	}()

	res, err := c.client.ListWorkflow(ctx, request)
	setDefaultAttributes(span, c.config)
	if request != nil {
		setClientRequestAttributes(span, clientRequest{
			namespace: request.Namespace,
			query:     request.Query,
		})
	}

	return res, err
}

/*
ListArchivedWorkflow retrieves archived workflow executions based on query. This
API will return a bad request if Temporal cluster or target namespace is not
configured for visibility archival or read is not enabled. However, differen
visibility archivers have different limitations on the query. Please check the
documentation of the visibility archiver used by your namespace to see what kind
of queries are accept and whether retrieved workflow executions are ordered or not.

It automatically handles tracing and error recording.
*/
func (c *iclient) ListArchivedWorkflow(ctx context.Context, request *workflowservice.ListArchivedWorkflowExecutionsRequest) (*workflowservice.ListArchivedWorkflowExecutionsResponse, error) {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: ListArchivedWorkflow", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to list archived workflows", err)
		}
	}()

	res, err := c.client.ListArchivedWorkflow(ctx, request)
	setDefaultAttributes(span, c.config)
	if request != nil {
		setClientRequestAttributes(span, clientRequest{
			namespace: request.Namespace,
			query:     request.Query,
		})
	}

	return res, err
}

/*
ScanWorkflow retrieves workflow executions based on query. ScanWorkflow should be
used when retrieving large amount of workflows and order is not needed. It will
use more ElasticSearch resources than ListWorkflow, but will be several times
faster when retrieving millions of workflows.

It automatically handles tracing and error recording.
*/
func (c *iclient) ScanWorkflow(ctx context.Context, request *workflowservice.ScanWorkflowExecutionsRequest) (*workflowservice.ScanWorkflowExecutionsResponse, error) {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: ScanWorkflow", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to scan workflows", err)
		}
	}()

	res, err := c.client.ScanWorkflow(ctx, request)
	setDefaultAttributes(span, c.config)
	if request != nil {
		setClientRequestAttributes(span, clientRequest{
			namespace: request.Namespace,
			query:     request.Query,
		})
	}

	return res, err
}

/*
CountWorkflow retrieves the number of workflow executions based on query.

It automatically handles tracing and error recording.
*/
func (c *iclient) CountWorkflow(ctx context.Context, request *workflowservice.CountWorkflowExecutionsRequest) (*workflowservice.CountWorkflowExecutionsResponse, error) {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: CountWorkflow", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to count workflows", err)
		}
	}()

	res, err := c.client.CountWorkflow(ctx, request)
	setDefaultAttributes(span, c.config)
	if request != nil {
		setClientRequestAttributes(span, clientRequest{
			namespace: request.Namespace,
			query:     request.Query,
		})
	}

	return res, err
}

/*
GetSearchAttributes returns valid search attributes keys and value types. The
search attributes can be used in query of List/Scan/Count APIs. Adding new search
attributes requires Temporal server to update dynamic config.

It automatically handles tracing and error recording.
*/
func (c *iclient) GetSearchAttributes(ctx context.Context) (*workflowservice.GetSearchAttributesResponse, error) {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: GetSearchAttributes", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to get search attributes", err)
		}
	}()

	res, err := c.client.GetSearchAttributes(ctx)
	setDefaultAttributes(span, c.config)

	return res, err
}

/*
QueryWorkflow queries a given workflow's last execution and returns the query
result synchronously.

It automatically handles tracing and error recording.
*/
func (c *iclient) QueryWorkflow(ctx context.Context, request *client.QueryWorkflowWithOptionsRequest) (*client.QueryWorkflowWithOptionsResponse, error) {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: QueryWorkflow", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to query workflow", err)
		}
	}()

	res, err := c.client.QueryWorkflowWithOptions(ctx, request)
	setDefaultAttributes(span, c.config)
	if request != nil {
		setClientRequestAttributes(span, clientRequest{
			workflowID:    request.WorkflowID,
			workflowRunID: request.RunID,
		})
	}

	return res, err
}

/*
DescribeWorkflowExecution returns information about the specified workflow
execution.

It automatically handles tracing and error recording.
*/
func (c *iclient) DescribeWorkflowExecution(ctx context.Context, workflowID, runID string) (*workflowservice.DescribeWorkflowExecutionResponse, error) {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: DescribeWorkflowExecution", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to describe workflow execution", err)
		}
	}()

	res, err := c.client.DescribeWorkflowExecution(ctx, workflowID, runID)
	setDefaultAttributes(span, c.config)
	setClientRequestAttributes(span, clientRequest{
		workflowID:    workflowID,
		workflowRunID: runID,
	})

	return res, err
}

/*
DescribeTaskQueue returns information about the target taskqueue, right now this
API returns the pollers which polled this taskqueue in last few minutes.

It automatically handles tracing and error recording.
*/
func (c *iclient) DescribeTaskQueue(ctx context.Context, taskqueue string, taskqueueType enums.TaskQueueType) (*workflowservice.DescribeTaskQueueResponse, error) {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: DescribeTaskQueue", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to describe task queue", err)
		}
	}()

	res, err := c.client.DescribeTaskQueue(ctx, taskqueue, taskqueueType)
	setDefaultAttributes(span, c.config)
	setClientRequestAttributes(span, clientRequest{
		taskqueue: taskqueue,
	})

	return res, err
}

/*
ResetWorkflowExecution resets an existing workflow execution. It will immediately
terminate the current execution instance.

It automatically handles tracing and error recording.
*/
func (c *iclient) ResetWorkflowExecution(ctx context.Context, request *workflowservice.ResetWorkflowExecutionRequest) (*workflowservice.ResetWorkflowExecutionResponse, error) {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: ResetWorkflowExecution", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to reset workflow execution", err)
		}
	}()

	res, err := c.client.ResetWorkflowExecution(ctx, request)
	setDefaultAttributes(span, c.config)
	if request != nil {
		setClientRequestAttributes(span, clientRequest{
			namespace:     request.Namespace,
			workflowID:    request.WorkflowExecution.WorkflowId,
			workflowRunID: request.WorkflowExecution.RunId,
		})
	}

	return res, err
}

/*
CheckHealth performs a server health check using the gRPC health check API. If the
check fails, an error is returned.

It automatically handles tracing and error recording.
*/
func (c *iclient) CheckHealth(ctx context.Context, request *client.CheckHealthRequest) (*client.CheckHealthResponse, error) {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: CheckHealth", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to check health", err)
		}
	}()

	res, err := c.client.CheckHealth(ctx, request)
	setDefaultAttributes(span, c.config)

	return res, err
}

/*
ScheduleClient creates a new shedule client with the same gRPC connection as this
client.
*/
func (c *iclient) ScheduleClient() ScheduleClient {
	sc := &scheduleclient{
		config: c.config,
		client: c.client.ScheduleClient(),
	}

	return sc
}
