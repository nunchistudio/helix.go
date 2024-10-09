package rest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.nunchi.studio/helix/errorstack"

	"github.com/stretchr/testify/assert"
)

func TestWriteOK(t *testing.T) {
	testcases := []struct {
		rw          *httptest.ResponseRecorder
		req         *http.Request
		attachments []With
		expected    Response
	}{
		{
			rw:          httptest.NewRecorder(),
			req:         nil,
			attachments: nil,
			expected: Response{
				Status: http.StatusText(http.StatusOK),
			},
		},
		{
			rw:  httptest.NewRecorder(),
			req: nil,
			attachments: []With{
				WithData(map[string]map[string]string{
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
		{
			rw:  httptest.NewRecorder(),
			req: nil,
			attachments: []With{
				WithErrorMessage("This error message should be ignored"),
				WithValidations([]errorstack.Validation{
					{
						Message: "This validation should be ignored",
					},
				}),
			},
			expected: Response{
				Status: http.StatusText(http.StatusOK),
			},
		},
	}

	for _, tc := range testcases {
		WriteOK(tc.rw, tc.req, tc.attachments...)

		assert.Equal(t, http.StatusOK, tc.rw.Code)

		expected, _ := json.Marshal(tc.expected)
		actual := tc.rw.Body.Bytes()
		assert.JSONEq(t, string(expected), string(actual))
	}
}

func TestWriteCreated(t *testing.T) {
	testcases := []struct {
		rw          *httptest.ResponseRecorder
		req         *http.Request
		attachments []With
		expected    Response
	}{
		{
			rw:          httptest.NewRecorder(),
			req:         nil,
			attachments: nil,
			expected: Response{
				Status: http.StatusText(http.StatusCreated),
			},
		},
		{
			rw:  httptest.NewRecorder(),
			req: nil,
			attachments: []With{
				WithData(map[string]map[string]string{
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
		{
			rw:  httptest.NewRecorder(),
			req: nil,
			attachments: []With{
				WithErrorMessage("This error message should be ignored"),
				WithValidations([]errorstack.Validation{
					{
						Message: "This validation should be ignored",
					},
				}),
			},
			expected: Response{
				Status: http.StatusText(http.StatusCreated),
			},
		},
	}

	for _, tc := range testcases {
		WriteCreated(tc.rw, tc.req, tc.attachments...)

		assert.Equal(t, http.StatusCreated, tc.rw.Code)

		expected, _ := json.Marshal(tc.expected)
		actual := tc.rw.Body.Bytes()
		assert.JSONEq(t, string(expected), string(actual))
	}
}

func TestWriteAccepted(t *testing.T) {
	testcases := []struct {
		rw          *httptest.ResponseRecorder
		req         *http.Request
		attachments []With
		expected    Response
	}{
		{
			rw:          httptest.NewRecorder(),
			req:         nil,
			attachments: nil,
			expected: Response{
				Status: http.StatusText(http.StatusAccepted),
			},
		},
		{
			rw:  httptest.NewRecorder(),
			req: nil,
			attachments: []With{
				WithMetadata(map[string]map[string]string{
					"temporal": {
						"workflow_id": "fd54bb5b-5153-437e-ba8e-df1c4ad77706",
						"run_id":      "40d0fe4b-cab2-4b6f-a497-891545f241ae",
					},
				}),
			},
			expected: Response{
				Status: http.StatusText(http.StatusAccepted),
				Metadata: map[string]map[string]string{
					"temporal": {
						"workflow_id": "fd54bb5b-5153-437e-ba8e-df1c4ad77706",
						"run_id":      "40d0fe4b-cab2-4b6f-a497-891545f241ae",
					},
				},
			},
		},
		{
			rw:  httptest.NewRecorder(),
			req: nil,
			attachments: []With{
				WithErrorMessage("This error message should be ignored"),
				WithValidations([]errorstack.Validation{
					{
						Message: "This validation should be ignored",
					},
				}),
			},
			expected: Response{
				Status: http.StatusText(http.StatusAccepted),
			},
		},
	}

	for _, tc := range testcases {
		WriteAccepted(tc.rw, tc.req, tc.attachments...)

		assert.Equal(t, http.StatusAccepted, tc.rw.Code)

		expected, _ := json.Marshal(tc.expected)
		actual := tc.rw.Body.Bytes()
		assert.JSONEq(t, string(expected), string(actual))
	}
}
