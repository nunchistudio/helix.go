package openfeature

import (
	"context"

	"go.nunchi.studio/helix/integration"

	"github.com/open-feature/go-sdk/pkg/openfeature"
)

/*
Ensure *connection complies to the integration.Integration type.
*/
var _ integration.Integration = (*connection)(nil)

/*
String returns the string representation of the OpenFeature integration.
*/
func (conn *connection) String() string {
	return identifier
}

/*
Start does nothing since the OpenFeature integration only exposes a client to
communicate with an OpenFeature provider.
*/
func (conn *connection) Start(ctx context.Context) error {
	return nil
}

/*
Close tries to gracefully shutdown the OpenFeature client provider.
*/
func (conn *connection) Close(ctx context.Context) error {
	openfeature.Shutdown()
	return nil
}

/*
Status always returns a `200` status.
*/
func (conn *connection) Status(ctx context.Context) (int, error) {
	return 200, nil
}
