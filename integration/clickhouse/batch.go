package clickhouse

import (
	"context"
	"fmt"

	"go.nunchi.studio/helix/telemetry/trace"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

/*
batch implements the Batch interface and allows to wrap the ClickHouse batch
functions for automatic tracing and error recording.
*/
type batch struct {
	config *Config
	client driver.Batch
}

/*
Batch exposes an opinionated way to interact with a ClickHouse batch capabilities.
All functions are wrapped with a context because some of them automatically
do distributed tracing (by using the said context) as well as error recording
within traces.
*/
type Batch interface {
	Append(ctx context.Context, v any) error
	Abort(ctx context.Context) error
	Flush(ctx context.Context) error
	Send(ctx context.Context) error
	IsSent(ctx context.Context) bool
}

/*
Append appends a new value to the batch.

Example:

	b.Append(ctx, &MyStruct{})
*/
func (b *batch) Append(ctx context.Context, value any) error {
	return b.client.AppendStruct(value)
}

/*
Abort tries to abort the batch.

It automatically handles tracing and error recording.
*/
func (b *batch) Abort(ctx context.Context) error {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Batch / Abort", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to abort batch", err)
		}
	}()

	err = b.client.Abort()
	setDefaultAttributes(span, b.config)

	return err
}

/*
Flush tries to flush the batch.

It automatically handles tracing and error recording.
*/
func (b *batch) Flush(ctx context.Context) error {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Batch / Flush", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to flush batch", err)
		}
	}()

	err = b.client.Flush()
	setDefaultAttributes(span, b.config)

	return err
}

/*
Send tries to send the batch.

It automatically handles tracing and error recording.
*/
func (b *batch) Send(ctx context.Context) error {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Batch / Send", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to send batch", err)
		}
	}()

	err = b.client.Send()
	setDefaultAttributes(span, b.config)

	return err
}

/*
IsSent informs if the batch has already been sent.

It automatically handles tracing and error recording.
*/
func (b *batch) IsSent(ctx context.Context) bool {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Batch / IsSent", humanized))
	defer span.End()

	sent := b.client.IsSent()
	setDefaultAttributes(span, b.config)

	return sent
}
