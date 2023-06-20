package event

import (
	"net"
	"net/url"
	"time"
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
	Name         string            `json:"name,omitempty"`
	Meta         map[string]string `json:"meta,omitempty"`
	Params       url.Values        `json:"params,omitempty"`
	IsAnonymous  bool              `json:"is_anonymous"`
	UserID       string            `json:"user_id,omitempty"`
	GroupID      string            `json:"group_id,omitempty"`
	Subscription Subscription      `json:"subscription,omitempty"`
	App          App               `json:"app,omitempty"`
	Library      Library           `json:"library,omitempty"`
	Campaign     Campaign          `json:"campaign,omitempty"`
	Referrer     Referrer          `json:"referrer,omitempty"`
	Cloud        Cloud             `json:"cloud,omitempty"`
	Device       Device            `json:"device,omitempty"`
	OS           OS                `json:"os,omitempty"`
	Location     Location          `json:"location,omitempty"`
	Network      Network           `json:"network,omitempty"`
	Page         Page              `json:"page,omitempty"`
	Screen       Screen            `json:"screen,omitempty"`
	IP           net.IP            `json:"ip,omitempty"`
	Locale       string            `json:"locale,omitempty"`
	Timezone     string            `json:"timezone,omitempty"`
	UserAgent    string            `json:"user_agent,omitempty"`
	Timestamp    time.Time         `json:"timestamp,omitempty"`
}
