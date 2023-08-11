package event

import (
	"strings"

	"go.opentelemetry.io/otel/baggage"
)

/*
Referrer holds the details about the marketing referrer from which a client is
executing the event from.
*/
type Referrer struct {
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
	Link string `json:"link,omitempty"`
}

/*
injectEventReferrerToFlatMap injects values found in a Referrer object to a flat
map representation of an Event.
*/
func injectEventReferrerToFlatMap(ref Referrer, flatten map[string]string) {
	if flatten == nil {
		flatten = make(map[string]string)
	}

	if ref.Type != "" {
		flatten["event.referrer.type"] = ref.Type
	}

	if ref.Name != "" {
		flatten["event.referrer.name"] = ref.Name
	}

	if ref.URL != "" {
		flatten["event.referrer.url"] = ref.URL
	}

	if ref.Link != "" {
		flatten["event.referrer.link"] = ref.Link
	}
}

/*
applyEventReferrerFromBaggageMember extracts the value of a Baggage member given
its key and applies it to an Event's Referrer. This assumes the Baggage member's
key starts with "event.referrer.".
*/
func applyEventReferrerFromBaggageMember(m baggage.Member, e *Event) {
	split := strings.Split(m.Key(), ".")

	switch split[2] {
	case "type":
		e.Referrer.Type = m.Value()
	case "name":
		e.Referrer.Name = m.Value()
	case "url":
		e.Referrer.URL = m.Value()
	case "link":
		e.Referrer.Link = m.Value()
	}
}
