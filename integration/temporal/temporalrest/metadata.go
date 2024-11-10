package temporalrest

import (
	"go.temporal.io/sdk/client"
)

/*
Metadata holds public details about Temporal that shall be used in the "metadata"
object of a REST response.
*/
type Metadata struct {
	Workflow *MetadataWorkflow `json:"workflow,omitempty"`
}

/*
MetadataWorkflow holds public details about a Temporal workflow that shall be used
in the "metadata" object of a REST response.
*/
type MetadataWorkflow struct {
	ID  string       `json:"id"`
	Run *MetadataRun `json:"run"`
}

/*
MetadataRun holds public details about a Temporal run that shall be used in the
"metadata" object of a REST response.
*/
type MetadataRun struct {
	ID string `json:"id"`
}

/*
GetMetadata returns details about a Temporal workflow that shall be used in the
"metadata" object of a REST response.

Example:

	type CustomResponseMetadata struct {
	  Temporal *temporalrest.Metadata `json:"temporal,omitempty"`
	}

	func myHandlerFunc(rw http.ResponseWriter, req *http.Request) {
	  // ...

	  metadata := CustomResponseMetadata{
	    Temporal: temporalrest.GetMetadata(wr),
	  }

	  rest.WriteAccepted[CustomResponse](rw, req,
	    rest.WithMetadataOnSuccess[CustomResponseMetadata](metadata),
	  )
	}
*/
func GetMetadata(wr client.WorkflowRun) *Metadata {
	if wr == nil {
		return nil
	}

	return &Metadata{
		Workflow: &MetadataWorkflow{
			ID: wr.GetID(),
			Run: &MetadataRun{
				ID: wr.GetRunID(),
			},
		},
	}
}
