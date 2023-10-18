package nats

import (
	"context"
	"fmt"
	"strings"
	"unicode"

	"go.nunchi.studio/helix/telemetry/log"
	"go.nunchi.studio/helix/telemetry/trace"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
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
setJetStreamMsgAttributes sets NATS JetStream message attributes to a trace span.
*/
func setJetStreamMsgAttributes(span *trace.Span, msg jetstream.Msg) {
	span.SetStringAttribute(fmt.Sprintf("%s.message.subject", identifier), msg.Subject())
}

/*
setConsumerAttributes sets NATS consumer attributes to a trace span.
*/
func setConsumerAttributes(span *trace.Span, cfg jetstream.ConsumerConfig) {
	span.SetStringAttribute(fmt.Sprintf("%s.jetstream.consumer.name", identifier), cfg.Name)
	span.SetBoolAttribute(fmt.Sprintf("%s.jetstream.consumer.ordered", identifier), false)
	span.SetSliceStringAttribute(fmt.Sprintf("%s.jetstream.consumer.subjects", identifier), cfg.FilterSubjects)
}

/*
setOrderedConsumerAttributes sets NATS ordered consumer attributes to a trace span.
*/
func setOrderedConsumerAttributes(span *trace.Span, cfg jetstream.OrderedConsumerConfig) {
	span.SetBoolAttribute(fmt.Sprintf("%s.jetstream.consumer.ordered", identifier), true)
	span.SetSliceStringAttribute(fmt.Sprintf("%s.jetstream.consumer.subjects", identifier), cfg.FilterSubjects)
}

/*
setStreamAttributes sets NATS stream attributes to a trace span.
*/
func setStreamAttributes(span *trace.Span, cfg jetstream.StreamConfig) {
	span.SetStringAttribute(fmt.Sprintf("%s.jetstream.stream.name", identifier), cfg.Name)
	span.SetSliceStringAttribute(fmt.Sprintf("%s.jetstream.stream.subjects", identifier), cfg.Subjects)
}

/*
setKeyValueAttributes sets NATS Key-Value attributes to a trace span.
*/
func setKeyValueAttributes(span *trace.Span, key string, cfg jetstream.KeyValueConfig) {
	if key != "" {
		span.SetStringAttribute(fmt.Sprintf("%s.jetstream.kv.key", identifier), key)
	}

	span.SetStringAttribute(fmt.Sprintf("%s.jetstream.kv.bucket.name", identifier), cfg.Bucket)
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
