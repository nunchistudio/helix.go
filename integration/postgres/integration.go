package postgres

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
String returns the string representation of the PostgreSQL integration.
*/
func (conn *connection) String() string {
	return identifier
}

/*
Start does nothing since the PostgreSQL integration only exposes a client to
communicate with a PostgreSQL server.
*/
func (conn *connection) Start(ctx context.Context) error {
	return nil
}

/*
Close tries to gracefully close the connection with the PostgreSQL server.
*/
func (conn *connection) Close(ctx context.Context) error {
	conn.client.Close()

	return nil
}

/*
Status indicates if the integration is able to ping the PostgreSQL server or not.
Returns `200` if connection is working, `503` otherwise.
*/
func (conn *connection) Status(ctx context.Context) (int, error) {
	stack := errorstack.New("Integration is not in a healthy state", errorstack.WithIntegration(identifier))

	err := conn.client.Ping(ctx)
	if err != nil {
		stack.WithValidations(errorstack.Validation{
			Message: err.Error(),
		})

		return 503, stack
	}

	return 200, nil
}
