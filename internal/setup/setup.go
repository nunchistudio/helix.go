package setup

import (
	"go.nunchi.studio/helix/internal/cloudprovider"
	"go.nunchi.studio/helix/internal/cloudprovider/kubernetes"
	"go.nunchi.studio/helix/internal/cloudprovider/nomad"
	"go.nunchi.studio/helix/internal/cloudprovider/qovery"
	"go.nunchi.studio/helix/internal/cloudprovider/render"
	"go.nunchi.studio/helix/internal/cloudprovider/unknown"
)

/*
init ensures helix.go global environment is properly setup: cloud provider is
mandatory for logger and tracer, which are required for a service to work as
expected.
*/
func init() {
	if cloudprovider.Detected == nil {
		cloudproviders := []cloudprovider.CloudProvider{
			qovery.Get(),
			kubernetes.Get(),
			nomad.Get(),
			render.Get(),
			unknown.Get(),
		}

		for _, orch := range cloudproviders {
			if orch != nil {
				cloudprovider.Detected = orch
				break
			}
		}
	}
}
