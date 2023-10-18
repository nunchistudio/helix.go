package nats

import (
	"context"
	"fmt"

	"go.nunchi.studio/helix/telemetry/trace"

	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

/*
messagescontext implements the MessagesContext interface and allows to wrap the
NATS JetStream messages' iterator for automatic tracing and error recording.
*/
type messagescontext struct {
	client jetstream.MessagesContext
}

/*
MessagesContext exposes an opinionated way to interact with a NATS JetStream
messages' iterator.
*/
type MessagesContext interface {
	Next(ctx context.Context) (context.Context, jetstream.Msg, error)
	Stop(ctx context.Context)
}

/*
Next retrieves next message on a stream. It will block until the next message is
available. The context returned contains trace details set when producing the
message received, allowing to chain spans within the same trace.

It automatically handles tracing and error recording.
*/
func (mc *messagescontext) Next(ctx context.Context) (context.Context, jetstream.Msg, error) {
	msg, err := mc.client.Next()

	ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(msg.Headers()))
	ctx, span := trace.Start(ctx, trace.SpanKindConsumer, fmt.Sprintf("%s: Consumer Iterator / Message", humanized))
	defer span.End()

	return ctx, msg, err
}

/*
Stop closes the iterator and cancels subscription.

It automatically handles tracing.
*/
func (mc *messagescontext) Stop(ctx context.Context) {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Consumer Iterator / Stop", humanized))
	defer span.End()

	mc.client.Stop()
}
