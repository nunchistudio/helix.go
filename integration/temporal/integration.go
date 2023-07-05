package temporal

import (
	"context"

	"go.nunchi.studio/helix/errorstack"
	"go.nunchi.studio/helix/integration"

	"go.temporal.io/sdk/worker"
)

/*
Ensure *connection complies to the integration.Integration type.
*/
var _ integration.Integration = (*connection)(nil)

/*
String returns the string representation of the Temporal integration.
*/
func (conn *connection) String() string {
	return identifier
}

/*
Start starts the Temporal worker, if applicable.
*/
func (conn *connection) Start(ctx context.Context) error {
	stack := errorstack.New("Failed to start worker", errorstack.WithIntegration(identifier))

	if conn.worker != nil {
		err := conn.worker.Run(worker.InterruptCh())
		if err != nil {
			stack.WithValidations(errorstack.Validation{
				Message: err.Error(),
			})

			return stack
		}
	}

	return nil
}

/*
Close gracefully stops the Temporal worker (if applicable) and closes the client's
connection with the server.
*/
func (conn *connection) Close(ctx context.Context) error {
	if conn.worker != nil {
		conn.worker.Stop()
	}

	conn.client.Close()
	return nil
}
