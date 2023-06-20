package setup

import (
	"go.nunchi.studio/helix/internal/orchestrator"
	"go.nunchi.studio/helix/internal/orchestrator/kubernetes"
	"go.nunchi.studio/helix/internal/orchestrator/nomad"
	"go.nunchi.studio/helix/internal/orchestrator/unknown"
)

/*
init ensures helix.go global environment is properly setup: orhcestrator is
mandatory for logger and tracer, which are required for a service to work as
expected.
*/
func init() {
	if orchestrator.Detected == nil {
		orchestrators := []orchestrator.Orchestrator{
			kubernetes.Get(),
			nomad.Get(),
			unknown.Get(),
		}

		for _, orch := range orchestrators {
			if orch != nil {
				orchestrator.Detected = orch
				break
			}
		}
	}
}
