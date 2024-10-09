package event

import (
	"strings"

	"go.opentelemetry.io/otel/baggage"
)

/*
Campaign holds the details about the marketing campaign from which a client is
executing the event from.
*/
type Campaign struct {
	Name    string `json:"name,omitempty"`
	Source  string `json:"source,omitempty"`
	Medium  string `json:"medium,omitempty"`
	Term    string `json:"term,omitempty"`
	Content string `json:"content,omitempty"`
}

/*
injectEventCampaignToFlatMap injects values found in a Campaign object to a flat
map representation of an Event.
*/
func injectEventCampaignToFlatMap(campaign Campaign, flatten map[string]string) {
	if flatten == nil {
		flatten = make(map[string]string)
	}

	if campaign.Name != "" {
		flatten["event.campaign.name"] = campaign.Name
	}

	if campaign.Source != "" {
		flatten["event.campaign.source"] = campaign.Source
	}

	if campaign.Medium != "" {
		flatten["event.campaign.medium"] = campaign.Medium
	}

	if campaign.Term != "" {
		flatten["event.campaign.term"] = campaign.Term
	}

	if campaign.Content != "" {
		flatten["event.campaign.content"] = campaign.Content
	}
}

/*
applyEventCampaignFromBaggageMember extracts the value of a Baggage member given
its key and applies it to an Event's Campaign. This assumes the Baggage member's
key starts with "event.campaign.".
*/
func applyEventCampaignFromBaggageMember(m baggage.Member, e *Event) {
	split := strings.Split(m.Key(), ".")

	switch split[2] {
	case "name":
		e.Campaign.Name = m.Value()
	case "source":
		e.Campaign.Source = m.Value()
	case "medium":
		e.Campaign.Medium = m.Value()
	case "term":
		e.Campaign.Term = m.Value()
	case "content":
		e.Campaign.Content = m.Value()
	}
}
