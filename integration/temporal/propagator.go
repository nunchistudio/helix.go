package temporal

import (
	"context"

	"go.nunchi.studio/helix/event"
	"go.nunchi.studio/helix/internal/contextkey"
	"go.nunchi.studio/helix/internal/tracer"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/converter"
	"go.temporal.io/sdk/workflow"
)

/*
custompropagator implements the workflow.ContextPropagator interface, allowing
to set custom context propagation logic across Temporal workflows and activities.
*/
type custompropagator struct{}

/*
Inject injects information from a Go context into headers.
*/
func (p *custompropagator) Inject(ctx context.Context, hw workflow.HeaderWriter) error {
	e, ok := event.EventFromContext(ctx)
	if !ok {
		return nil
	}

	// Retrieve the current span, and set Event's attributes.
	span := ctx.Value(contextkey.Span).(trace.Span)
	for k, v := range event.ToFlatMap(e) {
		span.SetAttributes(attribute.String(k, v))
	}

	// Transform the Event found to a Temporal payload so we can set it in header
	// right after.
	payload, err := converter.GetDefaultDataConverter().ToPayload(e)
	if err != nil {
		return err
	}

	hw.Set(event.Key, payload)
	return nil
}

/*
InjectFromWorkflow injects information from a workflow's context into headers.
*/
func (p *custompropagator) InjectFromWorkflow(ctx workflow.Context, hw workflow.HeaderWriter) error {
	e, ok := ctx.Value(contextkey.Event).(event.Event)
	if !ok {
		return nil
	}

	// Retrieve the current span, and set Event's attributes.
	span := ctx.Value(contextkey.Span).(trace.Span)
	for k, v := range event.ToFlatMap(e) {
		span.SetAttributes(attribute.String(k, v))
	}

	// Also set the workflow's attributes from its info.
	setWorkflowAttributes(span, workflow.GetInfo(ctx))

	// Transform the Event found to a Temporal payload so we can set it in header
	// right after.
	payload, err := converter.GetDefaultDataConverter().ToPayload(e)
	if err != nil {
		return err
	}

	hw.Set(event.Key, payload)
	return nil
}

/*
Extract extracts context information from headers and returns a context object.
*/
func (p *custompropagator) Extract(ctx context.Context, hr workflow.HeaderReader) (context.Context, error) {
	if value, ok := hr.Get(event.Key); ok {
		var e event.Event
		if err := converter.GetDefaultDataConverter().FromPayload(value, &e); err != nil {
			return ctx, nil
		}

		// Retrieve the current span, and set Event's attributes. Make sure a span
		// is set.
		span, ok := ctx.Value(contextkey.Span).(trace.Span)
		if span == nil || !ok {
			_, span = tracer.Tracer().Start(ctx, "")
			ctx = context.WithValue(ctx, contextkey.Span, span)
		}

		for k, v := range event.ToFlatMap(e) {
			span.SetAttributes(attribute.String(k, v))
		}

		// Also set the activity's attributes from its info.
		setActivityAttributes(span, activity.GetInfo(ctx))

		// Add the Event to the context returned.
		ctx = context.WithValue(ctx, contextkey.Event, e)
	}

	return ctx, nil
}

/*
ExtractToWorkflow extracts context information from headers and returns a workflow
context.
*/
func (p *custompropagator) ExtractToWorkflow(ctx workflow.Context, hr workflow.HeaderReader) (workflow.Context, error) {
	if value, ok := hr.Get(event.Key); ok {
		var e event.Event
		if err := converter.GetDefaultDataConverter().FromPayload(value, &e); err != nil {
			return ctx, nil
		}

		// Retrieve the current span, and set Event's attributes. Make sure a span
		// is set.
		span, ok := ctx.Value(contextkey.Span).(trace.Span)
		if span == nil || !ok {
			_, span = tracer.Tracer().Start(context.Background(), "")
			ctx = workflow.WithValue(ctx, contextkey.Span, span)
		}

		for k, v := range event.ToFlatMap(e) {
			span.SetAttributes(attribute.String(k, v))
		}

		// Also set the workflow's attributes from its info.
		setWorkflowAttributes(span, workflow.GetInfo(ctx))

		// Add the Event to the context returned.
		ctx = workflow.WithValue(ctx, contextkey.Event, e)
	}

	return ctx, nil
}
