package temporal

import (
	"context"
	"fmt"
	"strings"

	"go.nunchi.studio/helix/internal/contextkey"
	"go.nunchi.studio/helix/internal/tracer"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	oteltrace "go.opentelemetry.io/otel/trace"
	"go.temporal.io/sdk/contrib/opentelemetry"
	"go.temporal.io/sdk/interceptor"
)

/*
buildTracer tries to build the Temporal custom tracer from Config.
*/
func buildTracer(cfg Config) (interceptor.Tracer, error) {
	return opentelemetry.NewTracer(opentelemetry.TracerOptions{
		Tracer:            tracer.Tracer(),
		TextMapPropagator: otel.GetTextMapPropagator(),
		SpanContextKey:    contextkey.Span,
		SpanStarter: func(ctx context.Context, t oteltrace.Tracer, spanName string, opts ...oteltrace.SpanStartOption) oteltrace.Span {

			// Create a new Baggage populated with members retrieved from the context.
			b, _ := baggage.New(tracer.FromContextToBaggageMembers(ctx)...)

			// Create a new context including the Baggage previously created.
			ctx = baggage.ContextWithBaggage(ctx, b)

			// By default, the Temporal Go client includes the name of the workflow or
			// activity in the traces, such as "RunWorkflow:myworkflow". Only keep
			// the action name.
			split := strings.Split(spanName, ":")
			name := split[0]

			// Populate the Span attributes retrieved from the context.
			ctx, span := tracer.Tracer().Start(ctx, fmt.Sprintf("%s: %s", humanized, name))
			for _, attr := range tracer.FromContextToSpanAttributes(ctx) {
				span.SetAttributes(attr)
			}

			span.SetAttributes(attribute.String(fmt.Sprintf("%s.server.address", identifier), cfg.Address))
			span.SetAttributes(attribute.String(fmt.Sprintf("%s.namespace", identifier), cfg.Namespace))
			return span
		},
	})
}

/*
customtracer embeds interceptor.Tracer and override the StartSpan method, allowing
to override/rename default trace attributes set by the Temporal client.
*/
type customtracer struct {
	interceptor.Tracer
}

/*
StartSpan starts and returns a span with the given options, but with default
trace attributes renamed.
*/
func (m customtracer) StartSpan(opts *interceptor.TracerStartSpanOptions) (interceptor.TracerSpan, error) {
	if v := opts.Tags["temporalWorkflowID"]; v != "" {
		opts.Tags["temporal.workflow.id"] = v
		delete(opts.Tags, "temporalWorkflowID")
	}

	if v := opts.Tags["temporalRunID"]; v != "" {
		opts.Tags["temporal.workflow.run_id"] = v
		delete(opts.Tags, "temporalRunID")
	}

	if v := opts.Tags["temporalActivityID"]; v != "" {
		opts.Tags["temporal.activity.id"] = v
		delete(opts.Tags, "temporalActivityID")
	}

	return m.Tracer.StartSpan(opts)
}
