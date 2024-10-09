package event

import (
	"fmt"
	"strconv"
	"strings"

	"go.opentelemetry.io/otel/baggage"
)

/*
Subscription holds the details about the account/customer from which the event
has been triggered. It's useful for tracking customer usages.
*/
type Subscription struct {
	ID          string            `json:"id,omitempty"`
	TenantID    string            `json:"tenant_id,omitempty"`
	CustomerID  string            `json:"customer_id,omitempty"`
	PlanID      string            `json:"plan_id,omitempty"`
	Usage       string            `json:"usage,omitempty"`
	IncrementBy float64           `json:"increment_by,omitempty"`
	Flags       map[string]string `json:"flags,omitempty"`
}

/*
injectEventSubscriptionsToFlatMap injects values found in a slice of Subscription
objects to a flat map representation of an Event.
*/
func injectEventSubscriptionsToFlatMap(subs []Subscription, flatten map[string]string) {
	if flatten == nil {
		flatten = make(map[string]string)
	}

	for i, sub := range subs {
		if sub.ID != "" {
			flatten[fmt.Sprintf("event.subscriptions[%d].id", i)] = sub.ID
		}

		if sub.TenantID != "" {
			flatten[fmt.Sprintf("event.subscriptions[%d].tenant_id", i)] = sub.TenantID
		}

		if sub.CustomerID != "" {
			flatten[fmt.Sprintf("event.subscriptions[%d].customer_id", i)] = sub.CustomerID
		}

		if sub.PlanID != "" {
			flatten[fmt.Sprintf("event.subscriptions[%d].plan_id", i)] = sub.PlanID
		}

		if sub.Usage != "" {
			flatten[fmt.Sprintf("event.subscriptions[%d].usage", i)] = sub.Usage
		}

		if sub.IncrementBy != 0 {
			flatten[fmt.Sprintf("event.subscriptions[%d].increment_by", i)] = fmt.Sprintf("%f", sub.IncrementBy)
		}

		if sub.Flags != nil {
			for k, v := range sub.Flags {
				flatten[fmt.Sprintf("event.subscriptions[%d].flags.%s", i, k)] = v
			}
		}
	}
}

/*
applyEventSubscriptionsFromBaggageMember extracts the value of a Baggage member
given its key and applies it to an Event's Subscriptions. This assumes the Baggage
member's key starts with "event.subscriptions.".
*/
func applyEventSubscriptionsFromBaggageMember(m baggage.Member, e *Event) {
	split := strings.Split(m.Key(), ".")

	// Make sure to append a new subscription if the index found is greater than
	// the current length of the Subscriptions slice. Since the Baggage members
	// are not ordered, a key with index 1 may be called before one with index 0,
	// such as "event.subscriptions[1].id" called before "event.subscriptions[0].id".
	i, _ := strconv.Atoi(split[2])
	for i > len(e.Subscriptions)-1 {
		e.Subscriptions = append(e.Subscriptions, Subscription{})
	}

	switch split[3] {
	case "id":
		e.Subscriptions[i].ID = m.Value()
	case "tenant_id":
		e.Subscriptions[i].TenantID = m.Value()
	case "customer_id":
		e.Subscriptions[i].CustomerID = m.Value()
	case "plan_id":
		e.Subscriptions[i].PlanID = m.Value()
	case "usage":
		e.Subscriptions[i].Usage = m.Value()
	case "increment_by":
		e.Subscriptions[i].IncrementBy, _ = strconv.ParseFloat(m.Value(), 64)
	case "flags":
		if e.Subscriptions[i].Flags == nil {
			e.Subscriptions[i].Flags = make(map[string]string)
		}

		e.Subscriptions[i].Flags[split[4]] = m.Value()
	}
}
