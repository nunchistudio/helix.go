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
MsgHandler is like jetstream.MessageHandler but allows to pass a context for
leveraging automatic and distributed tracing with OpenTelemetry.
*/
type MsgHandler func(ctx context.Context, msg jetstream.Msg)

/*
consumer implements the Consumer interface and allows to wrap the NATS JetStream
consumer functions for automatic tracing and error recording.
*/
type consumer struct {
	config jetstream.ConsumerConfig
	client jetstream.Consumer
}

/*
Consumer exposes an opinionated way to interact with NATS JetStream consumer
capabilities.
*/
type Consumer interface {
	Fetch(ctx context.Context, batch int, opts ...jetstream.FetchOpt) (jetstream.MessageBatch, error)
	FetchBytes(ctx context.Context, maxBytes int, opts ...jetstream.FetchOpt) (jetstream.MessageBatch, error)
	FetchNoWait(ctx context.Context, batch int) (jetstream.MessageBatch, error)
	Consume(ctx context.Context, handler MsgHandler, opts ...jetstream.PullConsumeOpt) (jetstream.ConsumeContext, error)
	Messages(ctx context.Context, opts ...jetstream.PullMessagesOpt) (MessagesContext, error)
}

/*
Fetch is used to retrieve up to a provided number of messages from a stream. This
method will always send a single request and wait until either all messages are
retrieved or request times out.

It automatically handles tracing and error recording.
*/
func (c *consumer) Fetch(ctx context.Context, batch int, opts ...jetstream.FetchOpt) (jetstream.MessageBatch, error) {
	_, span := trace.Start(ctx, trace.SpanKindConsumer, fmt.Sprintf("%s: Consumer / Fetch", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to fetch messages from stream", err)
		}
	}()

	msg, err := c.client.Fetch(batch, opts...)
	setConsumerAttributes(span, c.config)

	return msg, err
}

/*
FetchBytes is used to retrieve up to a provided bytes from the stream. This method
will always send a single request and wait until provided number of bytes is
exceeded or request times out.

It automatically handles tracing and error recording.
*/
func (c *consumer) FetchBytes(ctx context.Context, maxBytes int, opts ...jetstream.FetchOpt) (jetstream.MessageBatch, error) {
	_, span := trace.Start(ctx, trace.SpanKindConsumer, fmt.Sprintf("%s: Consumer / FetchBytes", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to fetch bytes from stream", err)
		}
	}()

	msg, err := c.client.FetchBytes(maxBytes, opts...)
	setConsumerAttributes(span, c.config)

	return msg, err
}

/*
FetchNoWait is used to retrieve up to a provided number of messages from a stream.
This method will always send a single request and immediately return up to a
provided number of messages.

It automatically handles tracing and error recording.
*/
func (c *consumer) FetchNoWait(ctx context.Context, batch int) (jetstream.MessageBatch, error) {
	_, span := trace.Start(ctx, trace.SpanKindConsumer, fmt.Sprintf("%s: Consumer / FetchNoWait", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to fetch messages from stream", err)
		}
	}()

	msg, err := c.client.FetchNoWait(batch)
	setConsumerAttributes(span, c.config)

	return msg, err
}

/*
Consume can be used to continuously receive messages and handle them with the
provided callback function

The handler function passed is wrapped to automatically handles tracing and error
recording.
*/
func (c *consumer) Consume(ctx context.Context, handler MsgHandler, opts ...jetstream.PullConsumeOpt) (jetstream.ConsumeContext, error) {
	wrapped := func(msg jetstream.Msg) {
		ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(msg.Headers()))
		ctx, span := trace.Start(ctx, trace.SpanKindConsumer, fmt.Sprintf("%s: Consumer / Consume", humanized))
		defer span.End()

		setJetStreamMsgAttributes(span, msg)
		setConsumerAttributes(span, c.config)
		handler(ctx, msg)
	}

	return c.client.Consume(wrapped, opts...)
}

/*
Messages returns a MessagesContext, allowing continuously iterating over messages
on a stream.

The iterator automatically handles tracing and error recording on each message
received.
*/
func (c *consumer) Messages(ctx context.Context, opts ...jetstream.PullMessagesOpt) (MessagesContext, error) {
	messages, err := c.client.Messages(opts...)
	mtctx := &messagescontext{
		client: messages,
	}

	return mtctx, err
}
