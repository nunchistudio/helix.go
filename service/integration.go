package service

import (
	"context"
	"sync"

	"go.nunchi.studio/helix/errorstack"
	"go.nunchi.studio/helix/integration"
)

/*
Attach allows to attach a third-party integration to a service. When attached,
the Init and Close methods of the integration are automatically called when the
service is initializing and stopping, so they shouldn't be called manually by
the clients.
*/
func Attach(inte integration.Integration) error {
	svc.mutex.Lock()
	defer svc.mutex.Unlock()

	stack := errorstack.New("Failed to attach integration")
	if svc.isInitialized {
		stack.WithValidations(errorstack.Validation{
			Message: "Service must not be initialized for attaching an integration",
		})

		return stack
	}

	if svc.isClosed {
		stack.WithValidations(errorstack.Validation{
			Message: "Service must not be closed for attaching an integration",
		})

		return stack
	}

	if inte == nil {
		stack.WithValidations(errorstack.Validation{
			Message: "Integration must not be nil",
		})

		return stack
	}

	if inte.String() == "" {
		stack.WithValidations(errorstack.Validation{
			Message: "Integration's name must be set and not be empty",
			Path:    []string{"integration.String()"},
		})

		return stack
	}

	svc.integrations = append(svc.integrations, inte)
	return nil
}

/*
Status executes a health check of each integration attached to the service, and
returns the highest HTTP status code returned. This means if all integrations are
healthy (status `200`) but one is temporarily unavailable (status `503`), the
status returned would be `503`.
*/
func Status(ctx context.Context) (int, error) {

	// Create a channel that will receive the HTTP status code of the health check
	// of each integration.
	chStatus := make(chan int, len(svc.integrations))
	chError := make(chan error, len(svc.integrations))

	// Go through each integration attached to the service, and execute the health
	// checks asynchronously. Write the status returned to the channel.
	var wg sync.WaitGroup
	for _, inte := range svc.integrations {
		inte := inte
		wg.Add(1)

		go func() {
			defer wg.Done()

			status, err := inte.Status(ctx)
			if err != nil {
				chError <- err
			}

			chStatus <- status
		}()
	}

	wg.Wait()
	close(chStatus)
	close(chError)

	// Define the highest status code returned, as it will be used as the main one
	// returned by this function.
	var max int = 200
	for status := range chStatus {
		if status > max {
			max = status
		}
	}

	// Build a list of returned errors, and returned the error stack if applicable.
	stack := errorstack.New("Service is not in a healthy state")
	for err := range chError {
		stack.WithChildren(err)
	}

	if stack.HasChildren() {
		return max, stack
	}

	return max, nil
}
