package service

import (
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
