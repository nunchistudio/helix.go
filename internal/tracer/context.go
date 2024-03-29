package tracer

import (
	"context"
	"strings"

	"go.nunchi.studio/helix/event"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
)

/*
FromContextToBaggageMembers tries to extract an event from the context to add
each field as a baggage member to the trace.
*/
func FromContextToBaggageMembers(ctx context.Context) []baggage.Member {
	var members []baggage.Member

	// Try to extract the event from the context.
	ectx, ok := event.EventFromContext(ctx)
	if !ok {
		return members
	}

	// Transform the nasted object to a flatten map of string. It is a required
	// step to pass them from service to service. This also replaces keys part of
	// an array (to be compatible with Baggage specificiation), such as transforming
	// "event.subscriptions[0].id" to "event.subscriptions.0.id".
	mapped := event.ToFlatMap(ectx)
	for k, v := range mapped {
		k = strings.ReplaceAll(k, "[", ".")
		k = strings.ReplaceAll(k, "].", ".")
		k = strings.ReplaceAll(k, "]", "")

		m, err := baggage.NewMember(k, v)
		if err == nil {
			members = append(members, m)
		}
	}

	return members
}

/*
FromContextToSpanAttributes tries to extract an event from the context to add
each field as an attribute to the current span.
*/
func FromContextToSpanAttributes(ctx context.Context) []attribute.KeyValue {
	var attributes []attribute.KeyValue

	// Try to extract the event from the context.
	ectx, ok := event.EventFromContext(ctx)
	if !ok {
		return attributes
	}

	// Transform the nasted object to a flatten map of string. It is a required
	// step to pass them from service to service.
	mapped := event.ToFlatMap(ectx)
	for k, v := range mapped {
		attributes = append(attributes, attribute.String(k, v))
	}

	return attributes
}
