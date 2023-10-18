package nats

import (
	"context"
	"fmt"

	"go.nunchi.studio/helix/telemetry/trace"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

/*
JetStream exposes an opinionated way to interact with NATS JetStream.
*/
type JetStream interface {
	Publish(ctx context.Context, msg *nats.Msg, opts ...jetstream.PublishOpt) (*jetstream.PubAck, error)
	PublishAsync(ctx context.Context, msg *nats.Msg, opts ...jetstream.PublishOpt) (jetstream.PubAckFuture, error)
	PublishAsyncPending(ctx context.Context) int
	PublishAsyncComplete(ctx context.Context) <-chan struct{}

	Stream(ctx context.Context, streamname string) (Stream, error)
	CreateStream(ctx context.Context, config jetstream.StreamConfig) (Stream, error)
	UpdateStream(ctx context.Context, config jetstream.StreamConfig) (Stream, error)
	DeleteStream(ctx context.Context, streamname string) error

	Consumer(ctx context.Context, streamname string, consumername string) (Consumer, error)
	CreateOrUpdateConsumer(ctx context.Context, streamname string, config jetstream.ConsumerConfig) (Consumer, error)
	OrderedConsumer(ctx context.Context, streamname string, config jetstream.OrderedConsumerConfig) (Consumer, error)
	DeleteConsumer(ctx context.Context, streamname string, consumername string) error

	KeyValue(ctx context.Context, bucket string) (KeyValue, error)
	CreateKeyValue(ctx context.Context, config jetstream.KeyValueConfig) (KeyValue, error)
	DeleteKeyValue(ctx context.Context, bucket string) error
}

/*
Publish publishes a message to NATS JetStream.

It automatically handles tracing and error recording.
*/
func (conn *connection) Publish(ctx context.Context, msg *nats.Msg, opts ...jetstream.PublishOpt) (*jetstream.PubAck, error) {
	ctx, span := trace.Start(ctx, trace.SpanKindProducer, fmt.Sprintf("%s: Publish", humanized))
	if msg.Header == nil {
		msg.Header = make(nats.Header)
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(msg.Header))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to publish", err)
		}
	}()

	ack, err := conn.jetstream.PublishMsg(ctx, msg, opts...)
	setMsgAttributes(span, msg)

	return ack, err
}

/*
PublishAsync publishes a message to NATS JetStream and returns a nats.PubAckFuture.
The message should not be changed until the PubAckFuture has been processed.

It automatically handles tracing and error recording.
*/
func (conn *connection) PublishAsync(ctx context.Context, msg *nats.Msg, opts ...jetstream.PublishOpt) (jetstream.PubAckFuture, error) {
	ctx, span := trace.Start(ctx, trace.SpanKindProducer, fmt.Sprintf("%s: PublishAsync", humanized))
	if msg.Header == nil {
		msg.Header = make(nats.Header)
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(msg.Header))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to publish", err)
		}
	}()

	ackf, err := conn.jetstream.PublishMsgAsync(msg, opts...)
	setMsgAttributes(span, msg)

	return ackf, err
}

/*
PublishAsyncPending returns the number of async publishes outstanding for the
current NATS JetStream context.
*/
func (conn *connection) PublishAsyncPending(ctx context.Context) int {
	return conn.jetstream.PublishAsyncPending()
}

/*
PublishAsyncComplete returns a channel that will be closed when all outstanding
messages are ack'd.
*/
func (conn *connection) PublishAsyncComplete(ctx context.Context) <-chan struct{} {
	return conn.jetstream.PublishAsyncComplete()
}

/*
Stream returns a stream hook for a given stream name.

It automatically handles tracing and error recording.
*/
func (conn *connection) Stream(ctx context.Context, streamname string) (Stream, error) {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Stream", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to get stream", err)
		}
	}()

	found, err := conn.jetstream.Stream(ctx, streamname)
	setStreamAttributes(span, jetstream.StreamConfig{Name: streamname})
	if err != nil {
		return nil, err
	}

	stream := &stream{
		config: jetstream.StreamConfig{
			Name: streamname,
		},
		client: found,
	}

	return stream, nil
}

/*
CreateStream creates a new stream with given config and returns a hook to operate
on it.

It automatically handles tracing and error recording.
*/
func (conn *connection) CreateStream(ctx context.Context, config jetstream.StreamConfig) (Stream, error) {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: CreateStream", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to create stream", err)
		}
	}()

	created, err := conn.jetstream.CreateStream(ctx, config)
	setStreamAttributes(span, config)
	if err != nil {
		return nil, err
	}

	stream := &stream{
		config: config,
		client: created,
	}

	return stream, nil
}

/*
UpdateStream updates an existing stream with given config and returns a hook to
operate on it.

It automatically handles tracing and error recording.
*/
func (conn *connection) UpdateStream(ctx context.Context, config jetstream.StreamConfig) (Stream, error) {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: UpdateStream", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to update stream", err)
		}
	}()

	updated, err := conn.jetstream.UpdateStream(ctx, config)
	setStreamAttributes(span, config)
	if err != nil {
		return nil, err
	}

	stream := &stream{
		config: config,
		client: updated,
	}

	return stream, nil
}

/*
DeleteStream deletes a stream given its name.

It automatically handles tracing and error recording.
*/
func (conn *connection) DeleteStream(ctx context.Context, streamname string) error {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: DeleteStream", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to delete stream", err)
		}
	}()

	err = conn.jetstream.DeleteStream(ctx, streamname)
	setStreamAttributes(span, jetstream.StreamConfig{Name: streamname})

	return err
}

/*
Consumer returns a hook to an existing consumer, allowing processing of messages.

It automatically handles tracing and error recording.
*/
func (conn *connection) Consumer(ctx context.Context, streamname string, consumername string) (Consumer, error) {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Consumer", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to get consumer", err)
		}
	}()

	found, err := conn.jetstream.Consumer(ctx, streamname, consumername)
	setStreamAttributes(span, jetstream.StreamConfig{Name: streamname})
	setConsumerAttributes(span, jetstream.ConsumerConfig{Name: consumername})
	if err != nil {
		return nil, err
	}

	consumer := &consumer{
		config: jetstream.ConsumerConfig{
			Name: consumername,
		},
		client: found,
	}

	return consumer, nil
}

/*
CreateOrUpdateConsumer  creates a consumer on a given stream with given config.
If consumer already exists, jetstream.ErrConsumerExists is returned. Consumer
interface is returned, serving as a hook to operate on a consumer.

It automatically handles tracing and error recording.
*/
func (conn *connection) CreateOrUpdateConsumer(ctx context.Context, streamname string, config jetstream.ConsumerConfig) (Consumer, error) {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: CreateOrUpdateConsumer", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to create or update consumer", err)
		}
	}()

	created, err := conn.jetstream.CreateOrUpdateConsumer(ctx, streamname, config)
	setStreamAttributes(span, jetstream.StreamConfig{Name: streamname})
	setConsumerAttributes(span, config)
	if err != nil {
		return nil, err
	}

	consumer := &consumer{
		config: config,
		client: created,
	}

	return consumer, nil
}

/*
OrderedConsumer returns an OrderedConsumer instance. OrderedConsumer allows fetching
messages from a stream (just like standard consumer), for in order delivery of
messages. Underlying consumer is re-created when necessary, without additional
client code.

It automatically handles tracing and error recording.
*/
func (conn *connection) OrderedConsumer(ctx context.Context, streamname string, config jetstream.OrderedConsumerConfig) (Consumer, error) {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: OrderedConsumer", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to get ordered consumer", err)
		}
	}()

	created, err := conn.jetstream.OrderedConsumer(ctx, streamname, config)
	setStreamAttributes(span, jetstream.StreamConfig{Name: streamname})
	setOrderedConsumerAttributes(span, config)
	if err != nil {
		return nil, err
	}

	consumer := &consumer{
		client: created,
	}

	return consumer, nil
}

/*
DeleteConsumer removes a consumer with given name from a stream.

It automatically handles tracing and error recording.
*/
func (conn *connection) DeleteConsumer(ctx context.Context, streamname string, consumername string) error {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: DeleteConsumer", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to delete consumer", err)
		}
	}()

	err = conn.jetstream.DeleteConsumer(ctx, streamname, consumername)
	setStreamAttributes(span, jetstream.StreamConfig{Name: streamname})
	setConsumerAttributes(span, jetstream.ConsumerConfig{Name: consumername})

	return err
}

/*
KeyValue will lookup and bind to an existing key-value store.

It automatically handles tracing and error recording.
*/
func (conn *connection) KeyValue(ctx context.Context, bucket string) (KeyValue, error) {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: KeyValue", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to get key-value store", err)
		}
	}()

	store, err := conn.jetstream.KeyValue(ctx, bucket)
	setKeyValueAttributes(span, "", jetstream.KeyValueConfig{Bucket: bucket})
	if err != nil {
		return nil, err
	}

	kv := &keyvalue{
		bucket: bucket,
		store:  store,
	}

	return kv, nil
}

/*
CreateKeyValue creates a key-value store.

It automatically handles tracing and error recording.
*/
func (conn *connection) CreateKeyValue(ctx context.Context, config jetstream.KeyValueConfig) (KeyValue, error) {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: CreateKeyValue", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to create key-value store", err)
		}
	}()

	store, err := conn.jetstream.CreateKeyValue(ctx, config)
	setKeyValueAttributes(span, "", config)
	if err != nil {
		return nil, err
	}

	kv := &keyvalue{
		store: store,
	}

	return kv, nil
}

/*
DeleteKeyValue deletes a key-value store.

It automatically handles tracing and error recording.
*/
func (conn *connection) DeleteKeyValue(ctx context.Context, bucket string) error {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: DeleteKeyValue", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to delete key-value store", err)
		}
	}()

	err = conn.jetstream.DeleteKeyValue(ctx, bucket)
	setKeyValueAttributes(span, "", jetstream.KeyValueConfig{Bucket: bucket})

	return err
}
