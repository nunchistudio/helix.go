package event

import (
	"context"

	"go.nunchi.studio/helix/internal/contextkey"

	"go.opentelemetry.io/otel/baggage"
)

/*
EventFromContext returns the Event found in the context passed, if any. If no
Event has been found, it tries to find and build one if a Baggage was found in
the context. Returns true if an Event has been found, false otherwise.
*/
func EventFromContext(ctx context.Context) (Event, bool) {
	e, ok := ctx.Value(contextkey.Event).(Event)
	if !ok {
		return eventFromBaggage(baggage.FromContext(ctx))
	}

	return e, true
}

/*
ContextWithEvent returns a copy of the context passed with the Event associated
to it.
*/
func ContextWithEvent(ctx context.Context, e Event) context.Context {
	return context.WithValue(ctx, contextkey.Event, e)
}

/*
eventFromBaggage returns the Event found in the Baggage passed, if any. Returns
true if an Event has been found, false otherwise. An event is considered found
if — and only if — the name is not empty.

Example:

	map[string]string = {
	  "event.name"                         = "subscribed",
	  "event.user_id"                      = "user_2N6YZQLcYy2SPtmHiII69yHp0WE,
	  "event.params.filters[0]"            = "a",
	  "event.params.filters[1]"            = "b",
	  "event.params.filters[2]"            = "c",
	  "event.subscriptions[0].id"          = "sub_2N6YZQXgQAv87zMmvlHxePCSsRs",
	  "event.subscriptions[0].customer_id" = "cus_2N6YZMi3sBDPQBZrZJoYBwhNQNv",
	  "event.subscriptions[0].plan_id"     = "plan_2N6YZSE1SkWT9DrlXlswLhJ5K5Q",
	}

Will produce:

		Event{
		  Name:   "subscribed",
		  UserID: "user_2N6YZQLcYy2SPtmHiII69yHp0WE,
		  Params: url.Values{
		    "filters": []string{"a", "b", "c"},
			},
		  Subscriptions: []Subscription{
	      {
	        ID:         "sub_2N6YZQXgQAv87zMmvlHxePCSsRs",
	        CustomerID: "cus_2N6YZMi3sBDPQBZrZJoYBwhNQNv",
	        PlanID:     "plan_2N6YZSE1SkWT9DrlXlswLhJ5K5Q",
	      },
		  },
		}
*/
func eventFromBaggage(b baggage.Baggage) (Event, bool) {
	e := extractEventFromBaggage(b)
	if e.Name == "" {
		return e, false
	}

	return e, true
}
