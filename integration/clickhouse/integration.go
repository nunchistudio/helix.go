package clickhouse

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
String returns the string representation of the ClickHouse integration.
*/
func (conn *connection) String() string {
	return identifier
}

/*
Start does nothing since the ClickHouse integration only exposes a client to
communicate with a ClickHouse server.
*/
func (conn *connection) Start(ctx context.Context) error {
	return nil
}

/*
Close tries to gracefully close the connection with the ClickHouse server.
*/
func (conn *connection) Close(ctx context.Context) error {
	stack := errorstack.New("Failed to gracefully close connection with database", errorstack.WithIntegration(identifier))

	err := conn.client.Close()
	if err != nil {
		stack.WithValidations(errorstack.Validation{
			Message: normalizeErrorMessage(err),
		})

		return stack
	}

	return nil
}
