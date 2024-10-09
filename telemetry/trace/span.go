package trace

import (
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

/*
Span is the individual component of a Trace. It represents a single named and
timed operation of a workflow that is traced. A tracer is used to create a Span
and it is then up to the operation the Span represents to properly end the Span
when the operation itself ends.
*/
type Span struct {

	// client is the underlying OpenTelemetry client Span.
	client trace.Span

	// hasError is used to keep track if an error has been recorded in the Span,
	// allowing to set appropriate status when needed.
	hasError bool
}

/*
SpanKind is the role a Span plays in a Trace.
*/
type SpanKind int

/*
SpanKindInternal is a SpanKind for a Span that represents an internal operation
within an application.
*/
const SpanKindInternal SpanKind = 1

/*
SpanKindServer is a SpanKind for a Span that represents the operation of handling
a request from a client.
*/
const SpanKindServer SpanKind = 2

/*
SpanKindClient is a SpanKind for a Span that represents the operation of client
making a request to a server.
*/
const SpanKindClient SpanKind = 3

/*
SpanKindProducer is a SpanKind for a Span that represents the operation of a
producer sending a message to a message broker. Unlike SpanKindClient and
SpanKindServer, there is often no direct relationship between this kind of Span
and a SpanKindConsumer kind. A SpanKindProducer Span will end once the message
is accepted by the message broker which might not overlap with the processing of
that message.
*/
const SpanKindProducer SpanKind = 4

/*
SpanKindConsumer is a SpanKind for a Span that represents the operation of a
consumer receiving a message from a message broker. Like SpanKindProducer Spans,
there is often no direct relationship between this Span and the Span that produced
the message.
*/
const SpanKindConsumer SpanKind = 5

/*
SetStringAttribute sets a string attribute to the span.
*/
func (s *Span) SetStringAttribute(key string, value string) {
	s.client.SetAttributes(attribute.String(key, value))
}

/*
SetSliceStringAttribute sets a slice of string attributes to the span.
*/
func (s *Span) SetSliceStringAttribute(key string, values []string) {
	for i, value := range values {
		s.client.SetAttributes(attribute.String(fmt.Sprintf("%s[%d]", key, i), value))
	}
}

/*
SetBoolAttribute sets a boolean attribute to the span.
*/
func (s *Span) SetBoolAttribute(key string, value bool) {
	s.client.SetAttributes(attribute.Bool(key, value))
}

/*
SetIntAttribute sets a integer attribute to the span.
*/
func (s *Span) SetIntAttribute(key string, value int64) {
	s.client.SetAttributes(attribute.Int64(key, value))
}

/*
SetFloatAttribute sets a float attribute to the span.
*/
func (s *Span) SetFloatAttribute(key string, value float64) {
	s.client.SetAttributes(attribute.Float64(key, value))
}

/*
RecordError will record the error as an exception span event for this Span.
*/
func (s *Span) RecordError(msg string, err error) {
	s.hasError = true

	s.client.RecordError(err)
	s.client.SetStatus(codes.Error, msg)
}

/*
AddEvent adds an event to the Span with the provided name.
*/
func (s *Span) AddEvent(name string) {
	s.client.AddEvent(name)
}

/*
Context returns the original OpenTelemetry span's context.
*/
func (s *Span) Context() trace.SpanContext {
	return s.client.SpanContext()
}

/*
End sets the appropriate status and completes the Span. The Span is considered
complete and ready to be delivered through the rest of the telemetry pipeline
after this method is called. Therefore, updates to the Span are not allowed after
this method has been called.
*/
func (s *Span) End() {
	if !s.hasError {
		s.client.SetStatus(codes.Ok, "")
	}

	s.client.End()
}
