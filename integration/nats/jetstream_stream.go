package nats

import (
	"context"
	"fmt"

	"go.nunchi.studio/helix/telemetry/trace"

	"github.com/nats-io/nats.go/jetstream"
)

/*
stream implements the Stream interface and allows to wrap the NATS JetStream stream
management functions for automatic tracing and error recording.
*/
type stream struct {
	config jetstream.StreamConfig
	client jetstream.Stream
}

/*
Stream exposes an opinionated way to interact with NATS JetStream stream management
capabilities.
*/
type Stream interface {
	CreateOrUpdateConsumer(ctx context.Context, config jetstream.ConsumerConfig) (Consumer, error)
	OrderedConsumer(ctx context.Context, config jetstream.OrderedConsumerConfig) (Consumer, error)
	Consumer(ctx context.Context, consumername string) (Consumer, error)
	DeleteConsumer(ctx context.Context, consumername string) error
	Purge(ctx context.Context, opts ...jetstream.StreamPurgeOpt) error
	GetMsg(ctx context.Context, seq uint64, opts ...jetstream.GetMsgOpt) (*jetstream.RawStreamMsg, error)
	GetLastMsgForSubject(ctx context.Context, subject string) (*jetstream.RawStreamMsg, error)
	DeleteMsg(ctx context.Context, seq uint64) error
	SecureDeleteMsg(ctx context.Context, seq uint64) error
}

/*
CreateOrUpdateConsumer creates a consumer on a given stream with given config.
If consumer already exists, it will be updated (if possible). Consumer interface
is returned, serving as a hook to operate on a consumer (e.g. fetch messages).

It automatically handles tracing and error recording.
*/
func (s *stream) CreateOrUpdateConsumer(ctx context.Context, config jetstream.ConsumerConfig) (Consumer, error) {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Stream / CreateOrUpdateConsumer", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to create or update consumer for stream", err)
		}
	}()

	created, err := s.client.CreateOrUpdateConsumer(ctx, config)
	setStreamAttributes(span, s.config)
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
func (s *stream) OrderedConsumer(ctx context.Context, config jetstream.OrderedConsumerConfig) (Consumer, error) {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Stream / OrderedConsumer", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to get ordered consumer for stream", err)
		}
	}()

	created, err := s.client.OrderedConsumer(ctx, config)
	setStreamAttributes(span, s.config)
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
Consumer returns a Consumer interface for an existing consumer.

It automatically handles tracing and error recording.
*/
func (s *stream) Consumer(ctx context.Context, consumername string) (Consumer, error) {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Stream / Consumer", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to get consumer for stream", err)
		}
	}()

	found, err := s.client.Consumer(ctx, consumername)
	setStreamAttributes(span, s.config)
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
DeleteConsumer removes a consumer of a stream.

It automatically handles tracing and error recording.
*/
func (s *stream) DeleteConsumer(ctx context.Context, consumername string) error {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Stream / DeleteConsumer", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to delete consumer for stream", err)
		}
	}()

	err = s.client.DeleteConsumer(ctx, consumername)
	setStreamAttributes(span, s.config)
	setConsumerAttributes(span, jetstream.ConsumerConfig{Name: consumername})

	return err
}

/*
Purge removes messages from a stream.

It automatically handles tracing and error recording.
*/
func (s *stream) Purge(ctx context.Context, opts ...jetstream.StreamPurgeOpt) error {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Stream / Purge", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to purge stream", err)
		}
	}()

	err = s.client.Purge(ctx, opts...)
	setStreamAttributes(span, s.config)

	return err
}

/*
GetMsg retrieves a raw stream message stored in JetStream by sequence number.

It automatically handles tracing and error recording.
*/
func (s *stream) GetMsg(ctx context.Context, seq uint64, opts ...jetstream.GetMsgOpt) (*jetstream.RawStreamMsg, error) {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Stream / GetMsg", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to get raw stream message", err)
		}
	}()

	msg, err := s.client.GetMsg(ctx, seq, opts...)
	setStreamAttributes(span, s.config)

	return msg, err
}

/*
GetLastMsgForSubject retrieves the last raw stream message stored in JetStream by
subject.

It automatically handles tracing and error recording.
*/
func (s *stream) GetLastMsgForSubject(ctx context.Context, subject string) (*jetstream.RawStreamMsg, error) {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Stream / GetLastMsgForSubject", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to get latest raw stream message", err)
		}
	}()

	msg, err := s.client.GetLastMsgForSubject(ctx, subject)
	setStreamAttributes(span, s.config)

	return msg, err
}

/*
DeleteMsg deletes a message from a stream. The message is marked as erased, but
not overwritten.

It automatically handles tracing and error recording.
*/
func (s *stream) DeleteMsg(ctx context.Context, seq uint64) error {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Stream / DeleteMsg", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to delete message in stream", err)
		}
	}()

	err = s.client.DeleteMsg(ctx, seq)
	setStreamAttributes(span, s.config)

	return err
}

/*
SecureDeleteMsg deletes a message from a stream. The deleted message is overwritten
with random data. As a result, this operation is slower than DeleteMsg().

It automatically handles tracing and error recording.
*/
func (s *stream) SecureDeleteMsg(ctx context.Context, seq uint64) error {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Stream / SecureDeleteMsg", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to securely delete message in stream", err)
		}
	}()

	err = s.client.SecureDeleteMsg(ctx, seq)
	setStreamAttributes(span, s.config)

	return err
}
