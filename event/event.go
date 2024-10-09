package event

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/otel/baggage"
)

/*
Event is a dictionary of information that provides useful context about an event.
An Event shall be present as much as possible when passing data across services,
allowing to better understand the origin of an event.

Event should be used for data that you’re okay with potentially exposing to anyone
who inspects your network traffic. This is because it’s stored in HTTP headers
for distributed tracing. If your relevant network traffic is entirely within your
own network, then this caveat may not apply.

This is heavily inspired by the following references, and was adapted to better
fit this ecosystem:

  - The Segment's Context described at:
    https://segment.com/docs/connections/spec/common/#context
  - The Elastic Common Schema described at:
    https://www.elastic.co/guide/en/ecs/current/ecs-field-reference.html
*/
type Event struct {
	ID            string            `json:"id,omitempty"`
	Name          string            `json:"name,omitempty"`
	Meta          map[string]string `json:"meta,omitempty"`
	Params        url.Values        `json:"params,omitempty"`
	IsAnonymous   bool              `json:"is_anonymous"`
	UserID        string            `json:"user_id,omitempty"`
	GroupID       string            `json:"group_id,omitempty"`
	TenantID      string            `json:"tenant_id,omitempty"`
	IP            net.IP            `json:"ip,omitempty"`
	UserAgent     string            `json:"user_agent,omitempty"`
	Locale        string            `json:"locale,omitempty"`
	Timezone      string            `json:"timezone,omitempty"`
	Timestamp     time.Time         `json:"timestamp,omitempty"`
	App           App               `json:"app,omitempty"`
	Campaign      Campaign          `json:"campaign,omitempty"`
	Cloud         Cloud             `json:"cloud,omitempty"`
	Device        Device            `json:"device,omitempty"`
	Library       Library           `json:"library,omitempty"`
	Location      Location          `json:"location,omitempty"`
	Network       Network           `json:"network,omitempty"`
	OS            OS                `json:"os,omitempty"`
	Page          Page              `json:"page,omitempty"`
	Referrer      Referrer          `json:"referrer,omitempty"`
	Screen        Screen            `json:"screen,omitempty"`
	Subscriptions []Subscription    `json:"subscriptions,omitempty"`
}

/*
injectEventToFlatMap injects values found in an Event object to a flat map
representation of an Event. Top-level keys are handled here, while objects
are handled in their own functions for better clarity and maintainability.
*/
func injectEventToFlatMap(e Event, flatten map[string]string) {
	if flatten == nil {
		flatten = make(map[string]string)
	}

	flatten["event.id"] = e.ID
	flatten["event.name"] = e.Name

	if e.Meta != nil {
		for k, v := range e.Meta {
			flatten[fmt.Sprintf("event.meta.%s", k)] = v
		}
	}

	if e.Params != nil {
		for k, v := range e.Params {
			split := strings.Split(k, ".")
			for i, s := range v {
				flatten[fmt.Sprintf("event.params.%s[%d]", split[0], i)] = s
			}
		}
	}

	flatten["event.is_anonymous"] = strconv.FormatBool(e.IsAnonymous)
	flatten["event.user_id"] = e.UserID
	flatten["event.group_id"] = e.GroupID
	flatten["event.tenant_id"] = e.TenantID
	flatten["event.ip"] = e.IP.String()
	flatten["event.user_agent"] = e.UserAgent
	flatten["event.locale"] = e.Locale
	flatten["event.timezone"] = e.Timezone
	if !e.Timestamp.IsZero() {
		flatten["event.timestamp"] = e.Timestamp.Format(time.RFC3339Nano)
	}

	injectEventAppToFlatMap(e.App, flatten)
	injectEventCampaignToFlatMap(e.Campaign, flatten)
	injectEventCloudToFlatMap(e.Cloud, flatten)
	injectEventDeviceToFlatMap(e.Device, flatten)
	injectEventLibraryToFlatMap(e.Library, flatten)
	injectEventLocationToFlatMap(e.Location, flatten)
	injectEventNetworkToFlatMap(e.Network, flatten)
	injectEventOSToFlatMap(e.OS, flatten)
	injectEventPageToFlatMap(e.Page, flatten)
	injectEventReferrerToFlatMap(e.Referrer, flatten)
	injectEventScreenToFlatMap(e.Screen, flatten)
	injectEventSubscriptionsToFlatMap(e.Subscriptions, flatten)

	for k, v := range flatten {
		if v == "" || v == "false" || v == "0" || v == "0E+00" || v == "0.000000" || v == "<nil>" {
			delete(flatten, k)
		}
	}
}

/*
extractEventFromBaggage extracts the value of a Baggage and returns the Event
found. This assumes the Baggage members' key starts with "event.". Top-level keys
are handled here, while objects are handled in their own functions for better
clarity and maintainability.
*/
func extractEventFromBaggage(b baggage.Baggage) Event {
	var e Event

	for _, m := range b.Members() {
		if !strings.HasPrefix(m.Key(), "event.") {
			continue
		}

		if strings.HasPrefix(m.Key(), "event.meta.") {
			if e.Meta == nil {
				e.Meta = make(map[string]string)
			}

			e.Meta[strings.TrimPrefix(m.Key(), "event.meta.")] = m.Value()
			continue
		}

		if strings.HasPrefix(m.Key(), "event.params.") {
			if e.Params == nil {
				e.Params = make(url.Values)
			}

			key := strings.TrimPrefix(m.Key(), "event.params.")
			split := strings.Split(key, ".")

			e.Params.Add(split[0], m.Value())
			continue
		}

		split := strings.Split(m.Key(), ".")
		switch split[1] {
		case "id":
			e.ID = b.Member("event.id").Value()
		case "name":
			e.Name = b.Member("event.name").Value()
		case "is_anonymous":
			e.IsAnonymous, _ = strconv.ParseBool(b.Member("event.is_anonymous").Value())
		case "user_id":
			e.UserID = b.Member("event.user_id").Value()
		case "group_id":
			e.GroupID = b.Member("event.group_id").Value()
		case "tenant_id":
			e.TenantID = b.Member("event.tenant_id").Value()
		case "ip":
			e.IP = net.ParseIP(b.Member("event.ip").Value())
		case "user_agent":
			e.UserAgent = b.Member("event.user_agent").Value()
		case "locale":
			e.Locale = b.Member("event.locale").Value()
		case "timezone":
			e.Timezone = b.Member("event.timezone").Value()
		case "timestamp":
			e.Timestamp, _ = time.Parse(time.RFC3339Nano, b.Member("event.timestamp").Value())

		case "app":
			applyEventAppFromBaggageMember(m, &e)
		case "campaign":
			applyEventCampaignFromBaggageMember(m, &e)
		case "cloud":
			applyEventCloudFromBaggageMember(m, &e)
		case "device":
			applyEventDeviceFromBaggageMember(m, &e)
		case "library":
			applyEventLibraryFromBaggageMember(m, &e)
		case "location":
			applyEventLocationFromBaggageMember(m, &e)
		case "network":
			applyEventNetworkFromBaggageMember(m, &e)
		case "os":
			applyEventOSFromBaggageMember(m, &e)
		case "page":
			applyEventPageFromBaggageMember(m, &e)
		case "referrer":
			applyEventReferrerFromBaggageMember(m, &e)
		case "screen":
			applyEventScreenFromBaggageMember(m, &e)
		case "subscriptions":
			applyEventSubscriptionsFromBaggageMember(m, &e)
		}
	}

	return e
}
