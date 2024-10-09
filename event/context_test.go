package event

import (
	"context"
	"net/url"
	"testing"

	"go.nunchi.studio/helix/internal/contextkey"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/baggage"
)

func TestEventFromContext(t *testing.T) {
	testcases := []struct {
		ctx      context.Context
		baggage  func() baggage.Baggage
		expected Event
		success  bool
	}{
		{
			ctx:      context.WithValue(context.Background(), contextkey.Event, "not an Event"),
			expected: Event{},
			success:  false,
		},
		{
			ctx:      context.WithValue(context.Background(), contextkey.Event, Event{}),
			expected: Event{},
			success:  true,
		},
		{
			ctx: context.WithValue(context.Background(), contextkey.Event, Event{
				Name: "testing",
			}),
			expected: Event{
				Name: "testing",
			},
			success: true,
		},
		{
			ctx: context.Background(),
			baggage: func() baggage.Baggage {
				b, _ := baggage.New()
				memberName, _ := baggage.NewMember("event.name", "testing")
				memberParamsFilters0, _ := baggage.NewMember("event.params.filters.0", "a")
				memberParamsFilters1, _ := baggage.NewMember("event.params.filters.1", "b")
				memberParamsQuery0, _ := baggage.NewMember("event.params.query.0", "search_query")

				b, _ = b.SetMember(memberName)
				b, _ = b.SetMember(memberParamsFilters0)
				b, _ = b.SetMember(memberParamsFilters1)
				b, _ = b.SetMember(memberParamsQuery0)
				return b
			},
			expected: Event{
				Name: "testing",
				Params: url.Values{
					"filters": []string{"a", "b"},
					"query":   []string{"search_query"},
				},
			},
			success: true,
		},
	}

	for _, tc := range testcases {
		if tc.baggage != nil {
			tc.ctx = baggage.ContextWithBaggage(tc.ctx, tc.baggage())
		}

		actual, ok := EventFromContext(tc.ctx)

		assert.Equal(t, tc.expected, actual)
		assert.Equal(t, tc.success, ok)
	}
}

func TestContextWithEvent(t *testing.T) {
	testcases := []struct {
		ctx      context.Context
		input    Event
		expected context.Context
		success  bool
	}{
		{
			ctx: context.Background(),
			input: Event{
				Name: "testing",
			},
			expected: context.WithValue(context.Background(), contextkey.Event, Event{
				Name: "testing",
			}),
			success: true,
		},
	}

	for _, tc := range testcases {
		actual := ContextWithEvent(tc.ctx, tc.input)

		assert.Equal(t, tc.expected, actual)
	}
}
