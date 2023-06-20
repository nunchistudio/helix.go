package trace

import (
	"context"

	"go.nunchi.studio/helix/event"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
)

/*
fromContextToBaggageMembers tries to extract an event from the context to add
each field as a baggage member to the trace.
*/
func fromContextToBaggageMembers(ctx context.Context) []baggage.Member {
	var members []baggage.Member

	// Try to extract the event from the context.
	ectx, ok := event.EventFromContext(ctx)
	if !ok {
		return members
	}

	// Transform the nasted object to a flatten map of string. It is a required
	// step to pass them from service to service.
	mapped := event.ToFlatMap(ectx)
	for k, v := range mapped {
		m, err := baggage.NewMember(k, v)
		if err == nil {
			members = append(members, m)
		}
	}

	return members
}

/*
fromContextToSpanAttributes tries to extract an event from the context to add
each field as an attribute to the current span.
*/
func fromContextToSpanAttributes(ctx context.Context) []attribute.KeyValue {
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
