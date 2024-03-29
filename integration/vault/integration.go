package vault

import (
	"context"

	"go.nunchi.studio/helix/errorstack"
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

/*
Status indicates if the integration is able to connect to the Vault server or not.
Returns `200` if connection is working, `503` if Vault is sealed, is not
initialized, or if an error occurred.
*/
func (conn *connection) Status(ctx context.Context) (int, error) {
	stack := errorstack.New("Integration is not in a healthy state", errorstack.WithIntegration(identifier))

	res, err := conn.client.Sys().Health()
	if err != nil {
		stack.WithValidations(errorstack.Validation{
			Message: err.Error(),
		})

		return 503, stack
	}

	if !res.Initialized || res.Sealed {
		return 503, nil
	}

	return 200, nil
}
