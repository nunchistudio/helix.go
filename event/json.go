package event

import (
	"encoding/json"
)

/*
Key is the key that shall be present in a JSON-encoded value representing an Event.

Example:

	{
	  "key": "value",
	  "event": {
	    "name": "subscribed"
	  }
	}
*/
const Key string = "event"

/*
EventFromJSON returns the Event found at the "event" key in the JSON-encoded data
passed, if any. Returns true if an Event has been found, false otherwise.
*/
func EventFromJSON(input json.RawMessage) (Event, bool) {
	var e Event
	var mapped map[string]any
	err := json.Unmarshal(input, &mapped)
	if err != nil {
		return e, false
	}

	v, exists := mapped[Key]
	if !exists {
		return e, false
	}

	b, err := json.Marshal(v)
	if err != nil {
		return e, false
	}

	err = json.Unmarshal(b, &e)
	if err != nil {
		return e, false
	}

	return e, true
}

/*
ToFlatMap returns a flatten map for a given Event. Keys are prefixed with "event.",
and struct level are separated by a ".". All values are stringified.

This is primarily designed for the telemetry packages, allowing to pass contextual
information about an event using Go's context or HTTP headers, but can be useful
in some other use cases.

Example:

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

Will produce:

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
*/
func ToFlatMap(e Event) map[string]string {
	var flatten = make(map[string]string)

	injectEventToFlatMap(e, flatten)
	return flatten
}
