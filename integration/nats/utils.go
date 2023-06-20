package nats

import (
	"context"
	"fmt"
	"strings"
	"unicode"

	"go.nunchi.studio/helix/telemetry/log"
	"go.nunchi.studio/helix/telemetry/trace"

	"github.com/nats-io/nats.go"
)

/*
setMsgAttributes sets NATS message attributes to a trace span.
*/
func setMsgAttributes(span *trace.Span, msg *nats.Msg) {
	if msg != nil {
		span.SetStringAttribute(fmt.Sprintf("%s.message.subject", identifier), msg.Subject)
		span.SetStringAttribute(fmt.Sprintf("%s.subscription.subject", identifier), msg.Sub.Subject)
		span.SetStringAttribute(fmt.Sprintf("%s.subscription.queue", identifier), msg.Sub.Queue)
	}
}

/*
setConsumerAttributes sets NATS consumer attributes to a trace span.
*/
func setConsumerAttributes(span *trace.Span, stream string, cfg *nats.ConsumerConfig) {
	span.SetStringAttribute(fmt.Sprintf("%s.jetstream.stream.name", identifier), stream)

	if cfg != nil {
		span.SetStringAttribute(fmt.Sprintf("%s.consumer.name", identifier), cfg.Name)
		span.SetStringAttribute(fmt.Sprintf("%s.consumer.group", identifier), cfg.DeliverGroup)
		span.SetStringAttribute(fmt.Sprintf("%s.consumer.subject", identifier), cfg.DeliverSubject)
	}
}

/*
setStreamAttributes sets NATS stream attributes to a trace span.
*/
func setStreamAttributes(span *trace.Span, cfg *nats.StreamConfig) {
	if cfg != nil {
		span.SetStringAttribute(fmt.Sprintf("%s.jetstream.stream.name", identifier), cfg.Name)
		span.SetSliceStringAttribute(fmt.Sprintf("%s.jetstream.stream.subjects", identifier), cfg.Subjects)
	}
}

/*
setKeyValueAttributes sets NATS Key-Value attributes to a trace span.
*/
func setKeyValueAttributes(span *trace.Span, key string, cfg *nats.KeyValueConfig) {
	if key != "" {
		span.SetStringAttribute(fmt.Sprintf("%s.jetstream.kv.key", identifier), key)
	}

	if cfg != nil {
		span.SetStringAttribute(fmt.Sprintf("%s.jetstream.kv.bucket.name", identifier), cfg.Bucket)
	}
}

/*
normalizeErrorMessage normalizes an error returned by the NATS client to match
the format of helix.go. This is only used inside Start and Close for a better
readability in the terminal. Otherwise, functions return native NATS errors.

Example:

	"nats: no servers available for connection"

Becomes:

	"No servers available for connection"
*/
func normalizeErrorMessage(err error) string {
	var msg string = strings.TrimPrefix(err.Error(), "nats: ")
	runes := []rune(msg)
	runes[0] = unicode.ToUpper(runes[0])

	return string(runes)
}

/*
asyncErrorHandler handle asynchronous encountered errors.
*/
func asyncErrorHandler(conn *nats.Conn, sub *nats.Subscription, err error) {
	var msg string = normalizeErrorMessage(err)
	msg += fmt.Sprintf(" for subscription on %q", sub.Subject)

	if sub.Queue != "" {
		msg += fmt.Sprintf(" and queue %q", sub.Queue)
	}

	log.Error(context.TODO(), msg)
}
