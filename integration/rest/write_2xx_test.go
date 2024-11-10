package rest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.nunchi.studio/helix/errorstack"

	"github.com/stretchr/testify/assert"
)

type responseWithUser struct {
	Status string            `json:"status"`
	Error  *errorstack.Error `json:"error,omitempty"`
	Data   struct {
		User struct {
			ID string `json:"id"`
		} `json:"user"`
	} `json:"data,omitempty"`
}

type responseWithTemporalMetadata struct {
	Status       string            `json:"status"`
	Error        *errorstack.Error `json:"error,omitempty"`
	Metadatadata struct {
		Temporal struct {
			Workflow struct {
				ID  string `json:"id"`
				Run struct {
					ID string `json:"id"`
				} `json:"run"`
			} `json:"workflow"`
		} `json:"temporal"`
	} `json:"metadata,omitempty"`
}

func TestWriteOK(t *testing.T) {
	testcases := []struct {
		rw            *httptest.ResponseRecorder
		req           *http.Request
		withOnSuccess []WithOnSuccess
		expected      Response
	}{
		{
			rw:  httptest.NewRecorder(),
			req: nil,
			expected: Response{
				Status: http.StatusText(http.StatusOK),
			},
		},
		{
			rw:  httptest.NewRecorder(),
			req: nil,
			withOnSuccess: []WithOnSuccess{
				WithDataOnSuccess(map[string]map[string]string{
					"user": {
						"id": "fd97b96a-c399-4a2b-9c07-ff2f93ab97e0",
					},
				}),
			},
			expected: Response{
				Status: http.StatusText(http.StatusOK),
				Data: map[string]map[string]string{
					"user": {
						"id": "fd97b96a-c399-4a2b-9c07-ff2f93ab97e0",
					},
				},
			},
		},
	}

	for _, tc := range testcases {
		WriteOK[responseWithUser](tc.rw, tc.req, tc.withOnSuccess...)

		assert.Equal(t, http.StatusOK, tc.rw.Code)

		expected, _ := json.Marshal(tc.expected)
		actual := tc.rw.Body.Bytes()
		assert.JSONEq(t, string(expected), string(actual))
	}
}

func TestWriteCreated(t *testing.T) {
	testcases := []struct {
		rw            *httptest.ResponseRecorder
		req           *http.Request
		withOnSuccess []WithOnSuccess
		expected      Response
	}{
		{
			rw:            httptest.NewRecorder(),
			req:           nil,
			withOnSuccess: nil,
			expected: Response{
				Status: http.StatusText(http.StatusCreated),
			},
		},
		{
			rw:  httptest.NewRecorder(),
			req: nil,
			withOnSuccess: []WithOnSuccess{
				WithDataOnSuccess(map[string]map[string]string{
					"user": {
						"id": "fd97b96a-c399-4a2b-9c07-ff2f93ab97e0",
					},
				}),
			},
			expected: Response{
				Status: http.StatusText(http.StatusCreated),
				Data: map[string]map[string]string{
					"user": {
						"id": "fd97b96a-c399-4a2b-9c07-ff2f93ab97e0",
					},
				},
			},
		},
	}

	for _, tc := range testcases {
		WriteCreated[responseWithUser](tc.rw, tc.req, tc.withOnSuccess...)

		assert.Equal(t, http.StatusCreated, tc.rw.Code)

		expected, _ := json.Marshal(tc.expected)
		actual := tc.rw.Body.Bytes()
		assert.JSONEq(t, string(expected), string(actual))
	}
}

func TestWriteAccepted(t *testing.T) {
	testcases := []struct {
		rw            *httptest.ResponseRecorder
		req           *http.Request
		withOnSuccess []WithOnSuccess
		expected      Response
	}{
		{
			rw:            httptest.NewRecorder(),
			req:           nil,
			withOnSuccess: nil,
			expected: Response{
				Status: http.StatusText(http.StatusAccepted),
			},
		},
		{
			rw:  httptest.NewRecorder(),
			req: nil,
			withOnSuccess: []WithOnSuccess{
				WithMetadataOnSuccess(map[string]map[string]any{
					"temporal": {
						"workflow": map[string]any{
							"id": "fd54bb5b-5153-437e-ba8e-df1c4ad77706",
							"run": map[string]string{
								"id": "40d0fe4b-cab2-4b6f-a497-891545f241ae",
							},
						},
					},
				}),
			},
			expected: Response{
				Status: http.StatusText(http.StatusAccepted),
				Metadata: map[string]map[string]any{
					"temporal": {
						"workflow": map[string]any{
							"id": "fd54bb5b-5153-437e-ba8e-df1c4ad77706",
							"run": map[string]string{
								"id": "40d0fe4b-cab2-4b6f-a497-891545f241ae",
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range testcases {
		WriteAccepted[responseWithTemporalMetadata](tc.rw, tc.req, tc.withOnSuccess...)

		assert.Equal(t, http.StatusAccepted, tc.rw.Code)

		expected, _ := json.Marshal(tc.expected)
		actual := tc.rw.Body.Bytes()
		assert.JSONEq(t, string(expected), string(actual))
	}
}
