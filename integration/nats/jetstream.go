package nats

import (
	"context"
	"fmt"

	"go.nunchi.studio/helix/telemetry/trace"

	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

/*
MsgHandler is like nats.MsgHandler but allows to pass a context for leveraging
automatic distributed tracing with OpenTelemetry.
*/
type MsgHandler func(ctx context.Context, msg *nats.Msg)

/*
JetStream exposes an opinionated way to interact with NATS JetStream. All functions
are wrapped with a context because some of them automatically do distributed
tracing (by using the said context) as well as error recording within traces.

Interfaces wrapped:
  - nats.JetStream
  - nats.JetStreamManager
  - nats.KeyValue
  - nats.KeyValueManager
*/
type JetStream interface {
	Publish(ctx context.Context, msg *nats.Msg, opts ...nats.PubOpt) (*nats.PubAck, error)
	PublishAsync(ctx context.Context, msg *nats.Msg, opts ...nats.PubOpt) (nats.PubAckFuture, error)
	PublishAsyncPending(ctx context.Context) int
	PublishAsyncComplete(ctx context.Context) <-chan struct{}
	Subscribe(ctx context.Context, subject string, cb MsgHandler, opts ...nats.SubOpt) (*nats.Subscription, error)
	SubscribeSync(ctx context.Context, subject string, opts ...nats.SubOpt) (*nats.Subscription, error)
	ChanSubscribe(ctx context.Context, subject string, ch chan *nats.Msg, opts ...nats.SubOpt) (*nats.Subscription, error)
	ChanQueueSubscribe(ctx context.Context, subject string, queue string, ch chan *nats.Msg, opts ...nats.SubOpt) (*nats.Subscription, error)
	QueueSubscribe(ctx context.Context, subject string, queue string, cb MsgHandler, opts ...nats.SubOpt) (*nats.Subscription, error)
	QueueSubscribeSync(ctx context.Context, subject string, queue string, opts ...nats.SubOpt) (*nats.Subscription, error)
	PullSubscribe(ctx context.Context, subject string, durable string, opts ...nats.SubOpt) (*nats.Subscription, error)

	AddStream(ctx context.Context, cfg *nats.StreamConfig, opts ...nats.JSOpt) (*nats.StreamInfo, error)
	UpdateStream(ctx context.Context, cfg *nats.StreamConfig, opts ...nats.JSOpt) (*nats.StreamInfo, error)
	DeleteStream(ctx context.Context, stream string, opts ...nats.JSOpt) error
	StreamInfo(ctx context.Context, stream string, opts ...nats.JSOpt) (*nats.StreamInfo, error)
	PurgeStream(ctx context.Context, stream string, opts ...nats.JSOpt) error
	Streams(ctx context.Context, opts ...nats.JSOpt) <-chan *nats.StreamInfo
	StreamNames(ctx context.Context, opts ...nats.JSOpt) <-chan string
	GetMsg(ctx context.Context, stream string, seq uint64, opts ...nats.JSOpt) (*nats.RawStreamMsg, error)
	GetLastMsg(ctx context.Context, stream string, subject string, opts ...nats.JSOpt) (*nats.RawStreamMsg, error)
	DeleteMsg(ctx context.Context, stream string, seq uint64, opts ...nats.JSOpt) error
	SecureDeleteMsg(ctx context.Context, stream string, seq uint64, opts ...nats.JSOpt) error
	AddConsumer(ctx context.Context, stream string, cfg *nats.ConsumerConfig, opts ...nats.JSOpt) (*nats.ConsumerInfo, error)
	UpdateConsumer(ctx context.Context, stream string, cfg *nats.ConsumerConfig, opts ...nats.JSOpt) (*nats.ConsumerInfo, error)
	DeleteConsumer(ctx context.Context, stream string, consumer string, opts ...nats.JSOpt) error
	ConsumerInfo(ctx context.Context, stream string, name string, opts ...nats.JSOpt) (*nats.ConsumerInfo, error)
	Consumers(ctx context.Context, stream string, opts ...nats.JSOpt) <-chan *nats.ConsumerInfo
	ConsumerNames(ctx context.Context, stream string, opts ...nats.JSOpt) <-chan string
	AccountInfo(ctx context.Context, opts ...nats.JSOpt) (*nats.AccountInfo, error)
	StreamNameBySubject(ctx context.Context, subject string, opts ...nats.JSOpt) (string, error)

	KeyValue(ctx context.Context, bucket string) (KeyValue, error)
	CreateKeyValue(ctx context.Context, cfg *nats.KeyValueConfig) (KeyValue, error)
	DeleteKeyValue(ctx context.Context, bucket string) error
	KeyValueStoreNames(ctx context.Context) <-chan string
	KeyValueStores(ctx context.Context) <-chan nats.KeyValueStatus
}

/*
Publish publishes a message to NATS JetStream.

It automatically handles tracing and error recording.
*/
func (conn *connection) Publish(ctx context.Context, msg *nats.Msg, opts ...nats.PubOpt) (*nats.PubAck, error) {
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

	ack, err := conn.jetstream.PublishMsg(msg, opts...)
	setMsgAttributes(span, msg)

	return ack, err
}

/*
PublishAsync publishes a message to NATS JetStream and returns a nats.PubAckFuture.
The message should not be changed until the PubAckFuture has been processed.

It automatically handles tracing and error recording.
*/
func (conn *connection) PublishAsync(ctx context.Context, msg *nats.Msg, opts ...nats.PubOpt) (nats.PubAckFuture, error) {
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
Subscribe creates an async Subscription for JetStream. The stream and consumer
names can be provided with the nats.Bind() option. For creating an ephemeral
(where the consumer name is picked by the server), you can provide the stream
name with nats.BindStream(). If no stream name is specified, the library will
attempt to figure out which stream the subscription is for.

The callback function passed automatically handles tracing and error recording.
*/
func (conn *connection) Subscribe(ctx context.Context, subject string, cb MsgHandler, opts ...nats.SubOpt) (*nats.Subscription, error) {
	wrapped := func(msg *nats.Msg) {
		if msg.Header == nil {
			msg.Header = make(nats.Header)
		}

		ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(msg.Header))
		ctx, span := trace.Start(ctx, trace.SpanKindConsumer, fmt.Sprintf("%s: Subscribe", humanized))
		defer span.End()

		setMsgAttributes(span, msg)
		cb(ctx, msg)
	}

	return conn.jetstream.Subscribe(subject, wrapped, opts...)
}

/*
SubscribeSync creates a nats.Subscription that can be used to process messages
synchronously.
*/
func (conn *connection) SubscribeSync(ctx context.Context, subject string, opts ...nats.SubOpt) (*nats.Subscription, error) {
	return conn.jetstream.SubscribeSync(subject, opts...)
}

/*
ChanSubscribe creates channel based nats.Subscription.
*/
func (conn *connection) ChanSubscribe(ctx context.Context, subject string, ch chan *nats.Msg, opts ...nats.SubOpt) (*nats.Subscription, error) {
	return conn.jetstream.ChanSubscribe(subject, ch, opts...)
}

/*
ChanQueueSubscribe creates channel based nats.Subscription with a queue group.
*/
func (conn *connection) ChanQueueSubscribe(ctx context.Context, subject string, queue string, ch chan *nats.Msg, opts ...nats.SubOpt) (*nats.Subscription, error) {
	return conn.jetstream.ChanQueueSubscribe(subject, queue, ch, opts...)
}

/*
QueueSubscribe creates a nats.Subscription with a queue group. If no optional
durable name nor binding options are specified, the queue name will be used as a
durable name.

The callback function passed automatically handles tracing and error recording.
*/
func (conn *connection) QueueSubscribe(ctx context.Context, subject string, queue string, cb MsgHandler, opts ...nats.SubOpt) (*nats.Subscription, error) {
	wrapped := func(msg *nats.Msg) {
		if msg.Header == nil {
			msg.Header = make(nats.Header)
		}

		ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(msg.Header))
		ctx, span := trace.Start(ctx, trace.SpanKindConsumer, fmt.Sprintf("%s: QueueSubscribe", humanized))
		defer span.End()

		setMsgAttributes(span, msg)
		cb(ctx, msg)
	}

	return conn.jetstream.QueueSubscribe(subject, queue, wrapped, opts...)
}

/*
QueueSubscribeSync creates a nats.Subscription with a queue group that can be
used to process messages synchronously.
*/
func (conn *connection) QueueSubscribeSync(ctx context.Context, subject string, queue string, opts ...nats.SubOpt) (*nats.Subscription, error) {
	return conn.jetstream.QueueSubscribeSync(subject, queue, opts...)
}

/*
PullSubscribe creates a nats.Subscription that can fetch messages.
*/
func (conn *connection) PullSubscribe(ctx context.Context, subject string, durable string, opts ...nats.SubOpt) (*nats.Subscription, error) {
	return conn.jetstream.PullSubscribe(subject, durable, opts...)
}

/*
AddStream creates a stream.

It automatically handles tracing and error recording.
*/
func (conn *connection) AddStream(ctx context.Context, cfg *nats.StreamConfig, opts ...nats.JSOpt) (*nats.StreamInfo, error) {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: AddStream", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to add stream", err)
		}
	}()

	info, err := conn.jetstream.AddStream(cfg, opts...)
	setStreamAttributes(span, cfg)

	return info, err
}

/*
UpdateStream updates a stream.

It automatically handles tracing and error recording.
*/
func (conn *connection) UpdateStream(ctx context.Context, cfg *nats.StreamConfig, opts ...nats.JSOpt) (*nats.StreamInfo, error) {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: UpdateStream", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to update stream", err)
		}
	}()

	info, err := conn.jetstream.UpdateStream(cfg, opts...)
	setStreamAttributes(span, cfg)

	return info, err
}

/*
DeleteStream deletes a stream.

It automatically handles tracing and error recording.
*/
func (conn *connection) DeleteStream(ctx context.Context, stream string, opts ...nats.JSOpt) error {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: DeleteStream", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to delete stream", err)
		}
	}()

	err = conn.jetstream.DeleteStream(stream, opts...)
	setStreamAttributes(span, &nats.StreamConfig{
		Name: stream,
	})

	return err
}

/*
StreamInfo retrieves information from a stream.

It automatically handles tracing and error recording.
*/
func (conn *connection) StreamInfo(ctx context.Context, stream string, opts ...nats.JSOpt) (*nats.StreamInfo, error) {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: StreamInfo", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to get stream", err)
		}
	}()

	info, err := conn.jetstream.StreamInfo(stream, opts...)
	setStreamAttributes(span, &nats.StreamConfig{
		Name: stream,
	})

	return info, err
}

/*
PurgeStream purges a stream messages.

It automatically handles tracing and error recording.
*/
func (conn *connection) PurgeStream(ctx context.Context, stream string, opts ...nats.JSOpt) error {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: PurgeStream", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to purge stream", err)
		}
	}()

	err = conn.jetstream.PurgeStream(stream, opts...)
	setStreamAttributes(span, &nats.StreamConfig{
		Name: stream,
	})

	return err
}

/*
Streams can be used to retrieve a list of nats.StreamInfo objects.
*/
func (conn *connection) Streams(ctx context.Context, opts ...nats.JSOpt) <-chan *nats.StreamInfo {
	return conn.jetstream.Streams(opts...)
}

/*
StreamNames is used to retrieve a list of stream names.
*/
func (conn *connection) StreamNames(ctx context.Context, opts ...nats.JSOpt) <-chan string {
	return conn.jetstream.StreamNames(opts...)
}

/*
GetMsg retrieves a raw stream message stored in JetStream by sequence number.
Use options nats.DirectGet() or nats.DirectGetNext() to trigger retrieval directly
from a distributed group of servers (leader and replicas). The stream must have
been created/updated with the AllowDirect boolean.

It automatically handles tracing and error recording.
*/
func (conn *connection) GetMsg(ctx context.Context, stream string, seq uint64, opts ...nats.JSOpt) (*nats.RawStreamMsg, error) {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: GetMsg", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to get message", err)
		}
	}()

	raw, err := conn.jetstream.GetMsg(stream, seq, opts...)
	setStreamAttributes(span, &nats.StreamConfig{
		Name: stream,
	})

	return raw, err
}

/*
GetLastMsg retrieves the last raw stream message stored in JetStream by subject.
Use option nats.DirectGet() to trigger retrieval directly from a distributed group
of servers (leader and replicas). The stream must have been created/updated with
the AllowDirect boolean.

It automatically handles tracing and error recording.
*/
func (conn *connection) GetLastMsg(ctx context.Context, stream string, subject string, opts ...nats.JSOpt) (*nats.RawStreamMsg, error) {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: GetLastMsg", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to get last message", err)
		}
	}()

	raw, err := conn.jetstream.GetLastMsg(stream, subject, opts...)
	setStreamAttributes(span, &nats.StreamConfig{
		Name:     stream,
		Subjects: []string{subject},
	})

	return raw, err
}

/*
DeleteMsg deletes a message from a stream. The message is marked as erased, but
its value is not overwritten.

It automatically handles tracing and error recording.
*/
func (conn *connection) DeleteMsg(ctx context.Context, stream string, seq uint64, opts ...nats.JSOpt) error {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: DeleteMsg", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to delete message", err)
		}
	}()

	err = conn.jetstream.DeleteMsg(stream, seq, opts...)
	setStreamAttributes(span, &nats.StreamConfig{
		Name: stream,
	})

	return err
}

/*
SecureDeleteMsg deletes a message from a stream. The deleted message is overwritten
with random data As a result, this operation is slower than DeleteMsg().

It automatically handles tracing and error recording.
*/
func (conn *connection) SecureDeleteMsg(ctx context.Context, stream string, seq uint64, opts ...nats.JSOpt) error {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: SecureDeleteMsg", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to securely delete message", err)
		}
	}()

	err = conn.jetstream.SecureDeleteMsg(stream, seq, opts...)
	setStreamAttributes(span, &nats.StreamConfig{
		Name: stream,
	})

	return err
}

/*
AddConsumer adds a consumer to a stream.

It automatically handles tracing and error recording.
*/
func (conn *connection) AddConsumer(ctx context.Context, stream string, cfg *nats.ConsumerConfig, opts ...nats.JSOpt) (*nats.ConsumerInfo, error) {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: AddConsumer", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to add consumer", err)
		}
	}()

	info, err := conn.jetstream.AddConsumer(stream, cfg, opts...)
	setConsumerAttributes(span, stream, cfg)

	return info, err
}

/*
UpdateConsumer updates an existing consumer.

It automatically handles tracing and error recording.
*/
func (conn *connection) UpdateConsumer(ctx context.Context, stream string, cfg *nats.ConsumerConfig, opts ...nats.JSOpt) (*nats.ConsumerInfo, error) {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: UpdateConsumer", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to update consumer", err)
		}
	}()

	info, err := conn.jetstream.UpdateConsumer(stream, cfg, opts...)
	setConsumerAttributes(span, stream, cfg)

	return info, err
}

/*
DeleteConsumer deletes a consumer.

It automatically handles tracing and error recording.
*/
func (conn *connection) DeleteConsumer(ctx context.Context, stream string, name string, opts ...nats.JSOpt) error {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: DeleteConsumer", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to delete consumer", err)
		}
	}()

	err = conn.jetstream.DeleteConsumer(stream, name, opts...)
	setConsumerAttributes(span, stream, &nats.ConsumerConfig{
		Name: name,
	})

	return err
}

/*
ConsumerInfo retrieves information of a consumer from a stream.

It automatically handles tracing and error recording.
*/
func (conn *connection) ConsumerInfo(ctx context.Context, stream string, name string, opts ...nats.JSOpt) (*nats.ConsumerInfo, error) {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: ConsumerInfo", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to get consumer", err)
		}
	}()

	info, err := conn.jetstream.ConsumerInfo(stream, name, opts...)
	setConsumerAttributes(span, stream, &nats.ConsumerConfig{
		Name: name,
	})

	return info, err
}

/*
Consumers is used to retrieve a list of nats.ConsumerInfo objects.
*/
func (conn *connection) Consumers(ctx context.Context, stream string, opts ...nats.JSOpt) <-chan *nats.ConsumerInfo {
	return conn.jetstream.Consumers(stream, opts...)
}

/*
ConsumerNames is used to retrieve a list of consumer names.
*/
func (conn *connection) ConsumerNames(ctx context.Context, stream string, opts ...nats.JSOpt) <-chan string {
	return conn.jetstream.ConsumerNames(stream, opts...)
}

/*
AccountInfo retrieves info about the NATS JetStream usage from an account.

It automatically handles tracing and error recording.
*/
func (conn *connection) AccountInfo(ctx context.Context, opts ...nats.JSOpt) (*nats.AccountInfo, error) {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: AccountInfo", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to get account", err)
		}
	}()

	info, err := conn.jetstream.AccountInfo(opts...)
	return info, err
}

/*
StreamNameBySubject returns a stream matching given subject.

It automatically handles tracing and error recording.
*/
func (conn *connection) StreamNameBySubject(ctx context.Context, subject string, opts ...nats.JSOpt) (string, error) {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: StreamNameBySubject", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to get stream", err)
		}
	}()

	name, err := conn.jetstream.StreamNameBySubject(subject, opts...)
	setStreamAttributes(span, &nats.StreamConfig{
		Name:     name,
		Subjects: []string{subject},
	})

	return name, err
}

/*
KeyValue will lookup and bind to an existing key-value store.

It automatically handles tracing and error recording.
*/
func (conn *connection) KeyValue(ctx context.Context, bucket string) (KeyValue, error) {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: KeyValue", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to get key-value store", err)
		}
	}()

	store, err := conn.jetstream.KeyValue(bucket)
	if err != nil {
		return nil, err
	}

	kv := &keyvalue{
		bucket: bucket,
		store:  store,
	}

	setKeyValueAttributes(span, "", &nats.KeyValueConfig{
		Bucket: bucket,
	})

	return kv, nil
}

/*
CreateKeyValue creates a key-value store.

It automatically handles tracing and error recording.
*/
func (conn *connection) CreateKeyValue(ctx context.Context, cfg *nats.KeyValueConfig) (KeyValue, error) {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: CreateKeyValue", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to create key-value store", err)
		}
	}()

	store, err := conn.jetstream.CreateKeyValue(cfg)
	if err != nil {
		return nil, err
	}

	kv := &keyvalue{
		store: store,
	}

	setKeyValueAttributes(span, "", cfg)
	return kv, nil
}

/*
DeleteKeyValue deletes a key-value store.

It automatically handles tracing and error recording.
*/
func (conn *connection) DeleteKeyValue(ctx context.Context, bucket string) error {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: DeleteKeyValue", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to delete key-value store", err)
		}
	}()

	err = conn.jetstream.DeleteKeyValue(bucket)
	setKeyValueAttributes(span, "", &nats.KeyValueConfig{
		Bucket: bucket,
	})

	return err
}

/*
KeyValueStoreNames retrieves a list of key-value store names.
*/
func (conn *connection) KeyValueStoreNames(ctx context.Context) <-chan string {
	return conn.jetstream.KeyValueStoreNames()
}

/*
KeyValueStores retrieves a list of key-value store statuses.
*/
func (conn *connection) KeyValueStores(ctx context.Context) <-chan nats.KeyValueStatus {
	return conn.jetstream.KeyValueStores()
}
