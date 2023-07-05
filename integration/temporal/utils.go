package temporal

import (
	"fmt"
	"unicode"

	"go.nunchi.studio/helix/telemetry/trace"

	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/workflow"
)

/*
setDefaultAttributes sets integration attributes to a trace span.
*/
func setDefaultAttributes(span *trace.Span, cfg *Config) {
	if cfg != nil {
		span.SetStringAttribute(fmt.Sprintf("%s.server.address", identifier), cfg.Address)
		span.SetStringAttribute(fmt.Sprintf("%s.namespace", identifier), cfg.Namespace)
	}
}

/*
setWorkflowAttributes sets workflow attributes to a trace span. It uses the trace
type from the OpenTelemetry package since this happens in the interceptor and
we only have access via type assertion.
*/
func setWorkflowAttributes(span oteltrace.Span, info *workflow.Info) {
	if info != nil {
		span.SetAttributes(attribute.String(fmt.Sprintf("%s.worker.taskqueue", identifier), info.TaskQueueName))
		span.SetAttributes(attribute.String(fmt.Sprintf("%s.workflow.namespace", identifier), info.Namespace))
		span.SetAttributes(attribute.String(fmt.Sprintf("%s.workflow.type", identifier), info.WorkflowType.Name))
		span.SetAttributes(attribute.Int(fmt.Sprintf("%s.workflow.attempt", identifier), int(info.Attempt)))
	}
}

/*
setActivityAttributes sets activity attributes to a trace span. It uses the trace
type from the OpenTelemetry package since this happens in the interceptor and
we only have access via type assertion.
*/
func setActivityAttributes(span oteltrace.Span, info activity.Info) {
	span.SetAttributes(attribute.String(fmt.Sprintf("%s.worker.taskqueue", identifier), info.TaskQueue))
	span.SetAttributes(attribute.String(fmt.Sprintf("%s.workflow.namespace", identifier), info.WorkflowNamespace))
	span.SetAttributes(attribute.String(fmt.Sprintf("%s.workflow.type", identifier), info.WorkflowType.Name))
	span.SetAttributes(attribute.String(fmt.Sprintf("%s.activity.type", identifier), info.ActivityType.Name))
	span.SetAttributes(attribute.Int(fmt.Sprintf("%s.activity.attempt", identifier), int(info.Attempt)))
}

/*
setScheduleAttributes sets workflow schedule attributes to a trace span.
*/
func setScheduleAttributes(span *trace.Span, id string) {
	span.SetStringAttribute(fmt.Sprintf("%s.schedule.id", identifier), id)
}

/*
clientRequest holds information about a request made by the Temporal client.
*/
type clientRequest struct {
	namespace     string
	workflowID    string
	workflowRunID string
	activityID    string
	signalName    string
	query         string
	taskqueue     string
}

/*
setClientRequestAttributes sets a client's request attributes to a trace span.
*/
func setClientRequestAttributes(span *trace.Span, req clientRequest) {
	if req.namespace != "" {
		span.SetStringAttribute(fmt.Sprintf("%s.namespace", identifier), req.namespace)
	}

	if req.workflowID != "" {
		span.SetStringAttribute(fmt.Sprintf("%s.workflow.id", identifier), req.workflowID)
	}

	if req.workflowRunID != "" {
		span.SetStringAttribute(fmt.Sprintf("%s.workflow.run_id", identifier), req.workflowRunID)
	}

	if req.activityID != "" {
		span.SetStringAttribute(fmt.Sprintf("%s.activity.id", identifier), req.activityID)
	}

	if req.signalName != "" {
		span.SetStringAttribute(fmt.Sprintf("%s.signal.name", identifier), req.signalName)
	}

	if req.query != "" {
		span.SetStringAttribute(fmt.Sprintf("%s.search.query", identifier), req.query)
	}

	if req.taskqueue != "" {
		span.SetStringAttribute(fmt.Sprintf("%s.worker.taskqueue", identifier), req.taskqueue)
	}
}

/*
normalizeErrorMessage normalizes an error returned by the Temporal client to match
the format of helix.go. This is only used inside Start and Close for a better
readability in the terminal. Otherwise, functions return native Temporal errors.

Example:

	"dial tcp 127.0.0.1:7233: connect: connection refused"

Becomes:

	"Dial tcp 127.0.0.1:7233: connect: connection refused"
*/
func normalizeErrorMessage(err error) string {
	var msg string = err.Error()
	runes := []rune(msg)
	runes[0] = unicode.ToUpper(runes[0])

	return string(runes)
}
