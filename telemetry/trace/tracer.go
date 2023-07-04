package trace

import (
	"context"

	"go.nunchi.studio/helix/internal/tracer"

	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/trace"
)

/*
Start creates a Span and a context containing the newly-created Span.

If the context provided contains a Span then the newly-created Span will be a
child of that Span, otherwise it will be a root Span.

Any Span that is created must also be ended. This is the responsibility of the
user. Implementations of this API may leak memory or other resources if Spans are
not ended.
*/
func Start(ctx context.Context, kind SpanKind, name string) (context.Context, *Span) {

	// Create a new Baggage populated with members retrieved from the context.
	b, err := baggage.New(tracer.FromContextToBaggageMembers(ctx)...)
	if err != nil {
		return ctx, nil
	}

	// Create a new context including the Baggage previously created.
	ctx = baggage.ContextWithBaggage(ctx, b)

	// Populate the Span attributes retrieved from the context.
	ctx, span := tracer.Tracer().Start(ctx, name, trace.WithSpanKind(trace.SpanKind(kind)))
	for _, attr := range tracer.FromContextToSpanAttributes(ctx) {
		span.SetAttributes(attr)
	}

	s := &Span{
		client: span,
	}

	return ctx, s
}
