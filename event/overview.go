/*
Package event exposes some objects and utilities to get/set contextual information
about an event. In a distributed architecture where events are flowing across
internal services and third-party integrations, it is highly encouraged to pass
an Event through a context.Context in order to trace them from end-to-end.
helix.go core and integrations rely on Go contexts to manage logs and traces
across services.

This package must not import any other package of this ecosystem.
*/
package event
