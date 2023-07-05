package temporal

import (
	"context"

	"go.nunchi.studio/helix/event"
	"go.nunchi.studio/helix/internal/contextkey"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.temporal.io/sdk/workflow"
)

/*
EventFromWorkflow tries to retrieve an Event from the workflow's context. Returns
true if an Event has been found, false otherwise.

If an Event was found, it is added to the span attributes.
*/
func EventFromWorkflow(ctx workflow.Context) (event.Event, bool) {
	var e event.Event

	if val := ctx.Value(contextkey.Event); val != nil {
		e, ok := val.(event.Event)
		if !ok {
			return e, false
		}

		span, ok := ctx.Value(contextkey.Span).(trace.Span)
		if !ok {
			return e, false
		}

		for k, v := range event.ToFlatMap(e) {
			span.SetAttributes(attribute.String(k, v))
		}
	}

	return e, true
}

/*
EventFromActivity tries to retrieve an Event from the activity's context. Returns
true if an Event has been found, false otherwise.

If an Event was found, it is added to the span attributes.
*/
func EventFromActivity(ctx context.Context) (event.Event, bool) {
	var e event.Event

	if val := ctx.Value(contextkey.Event); val != nil {
		e, ok := val.(event.Event)
		if !ok {
			return e, false
		}

		span, ok := ctx.Value(contextkey.Span).(trace.Span)
		if !ok {
			return e, false
		}

		for k, v := range event.ToFlatMap(e) {
			span.SetAttributes(attribute.String(k, v))
		}
	}

	return e, true
}
