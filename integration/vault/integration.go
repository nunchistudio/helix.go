package vault

import (
	"context"

	"go.nunchi.studio/helix/integration"
)

/*
Ensure *connection complies to the integration.Integration type.
*/
var _ integration.Integration = (*connection)(nil)

/*
String returns the string representation of the Vault integration.
*/
func (conn *connection) String() string {
	return identifier
}

/*
Start does nothing since the Vault integration only exposes a client to communicate
with a Vault server.
*/
func (conn *connection) Start(ctx context.Context) error {
	return nil
}

/*
Close does nothing since the underlying Vault client doesn't need to be closed.
*/
func (conn *connection) Close(ctx context.Context) error {
	return nil
}
