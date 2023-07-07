package nats

import (
	"context"

	"go.nunchi.studio/helix/integration"
)

/*
Ensure *connection complies to the integration.Integration type.
*/
var _ integration.Integration = (*connection)(nil)

/*
String returns the string representation of the NATS integration.
*/
func (conn *connection) String() string {
	return identifier
}

/*
Start does nothing since the NATS integration only exposes a client to communicate
with a NATS server.
*/
func (conn *connection) Start(ctx context.Context) error {
	return nil
}

/*
Close gracefully closes the connection with the NATS server.
*/
func (conn *connection) Close(ctx context.Context) error {
	conn.nats.Close()

	return nil
}

/*
Status indicates if the integration is able to connect to the NATS server or not.
Returns `200` if connection is in a proper state, `503` otherwise.
*/
func (conn *connection) Status(ctx context.Context) (int, error) {
	var status int = 503
	if conn.nats.Status().String() == "CONNECTED" {
		status = 200
	}

	return status, nil
}
