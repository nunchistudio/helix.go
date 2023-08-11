package event

import (
	"strings"

	"go.opentelemetry.io/otel/baggage"
)

/*
Page holds the details about the webpage from which the event is triggered from.
*/
type Page struct {
	Path     string `json:"path,omitempty"`
	Referrer string `json:"referrer,omitempty"`
	Search   string `json:"search,omitempty"`
	Title    string `json:"title,omitempty"`
	URL      string `json:"url,omitempty"`
}

/*
injectEventPageToFlatMap injects values found in a Page object to a flat map
representation of an Event.
*/
func injectEventPageToFlatMap(page Page, flatten map[string]string) {
	if flatten == nil {
		flatten = make(map[string]string)
	}

	if page.Path != "" {
		flatten["event.page.path"] = page.Path
	}

	if page.Referrer != "" {
		flatten["event.page.referrer"] = page.Referrer
	}

	if page.Search != "" {
		flatten["event.page.search"] = page.Search
	}

	if page.Title != "" {
		flatten["event.page.title"] = page.Title
	}

	if page.URL != "" {
		flatten["event.page.url"] = page.URL
	}
}

/*
applyEventPageFromBaggageMember extracts the value of a Baggage member given its
key and applies it to an Event's Page. This assumes the Baggage member's key starts
with "event.page.".
*/
func applyEventPageFromBaggageMember(m baggage.Member, e *Event) {
	split := strings.Split(m.Key(), ".")

	switch split[2] {
	case "path":
		e.Page.Path = m.Value()
	case "referrer":
		e.Page.Referrer = m.Value()
	case "search":
		e.Page.Search = m.Value()
	case "title":
		e.Page.Title = m.Value()
	case "url":
		e.Page.URL = m.Value()
	}
}
