package integration

import (
	"context"
)

/*
Integration describes the lifecycle of an integration.
*/
type Integration interface {

	// String returns the string representation of the integration.
	//
	// Examples:
	//
	//   "nats"
	//   "vault"
	String() string

	// Start starts/opens a connection with the integration, if applicable. This
	// function can be blocking, for starting a server or a worker for example.
	// The service package executes each Start function of attached integrations
	// in their own goroutine, and returns an error as soon as a goroutine returns
	// a non-nil error.
	Start(ctx context.Context) error

	// Close closes the connection with the integration, if applicable.
	Close(ctx context.Context) error
}
