package rest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.nunchi.studio/helix/errorstack"

	"github.com/stretchr/testify/assert"
)

func TestWriteInternalServerError(t *testing.T) {
	addFrenchLocalesForTest()
	reqWithLang, _ := http.NewRequest(http.MethodPost, "/anything", nil)
	reqWithLang.Header.Add("Accept-Language", "fr")

	testcases := []struct {
		rw          *httptest.ResponseRecorder
		req         *http.Request
		attachments []With
		expected    Response
	}{
		{
			rw:  httptest.NewRecorder(),
			req: nil,
			expected: Response{
				Status: http.StatusText(http.StatusInternalServerError),
				Error: &errorstack.Error{
					Message: "We have been notified of this unexpected internal error",
				},
			},
		},
		{
			rw:          httptest.NewRecorder(),
			req:         reqWithLang,
			attachments: nil,
			expected: Response{
				Status: http.StatusText(http.StatusInternalServerError),
				Error: &errorstack.Error{
					Message: "Nous avons été informés de cette erreur interne inattendue",
				},
			},
		},
		{
			rw:  httptest.NewRecorder(),
			req: nil,
			attachments: []With{
				WithMetadata(map[string]string{
					"anything": "value",
				}),
				WithData(map[string]string{
					"id": "Data object should be ignored",
				}),
			},
			expected: Response{
				Status: http.StatusText(http.StatusInternalServerError),
				Error: &errorstack.Error{
					Message: "We have been notified of this unexpected internal error",
				},
				Metadata: map[string]string{
					"anything": "value",
				},
			},
		},
	}

	for _, tc := range testcases {
		WriteInternalServerError(tc.rw, tc.req, tc.attachments...)

		assert.Equal(t, http.StatusInternalServerError, tc.rw.Code)

		expected, _ := json.Marshal(tc.expected)
		actual := tc.rw.Body.Bytes()
		assert.JSONEq(t, string(expected), string(actual))
	}
}

func TestWriteServiceUnavailable(t *testing.T) {
	addFrenchLocalesForTest()
	reqWithLang, _ := http.NewRequest(http.MethodPost, "/anything", nil)
	reqWithLang.Header.Add("Accept-Language", "fr")

	testcases := []struct {
		rw          *httptest.ResponseRecorder
		req         *http.Request
		attachments []With
		expected    Response
	}{
		{
			rw:  httptest.NewRecorder(),
			req: nil,
			expected: Response{
				Status: http.StatusText(http.StatusServiceUnavailable),
				Error: &errorstack.Error{
					Message: "Please try again in a few moments",
				},
			},
		},
		{
			rw:          httptest.NewRecorder(),
			req:         reqWithLang,
			attachments: nil,
			expected: Response{
				Status: http.StatusText(http.StatusServiceUnavailable),
				Error: &errorstack.Error{
					Message: "Veuillez réessayer dans quelques instants",
				},
			},
		},
		{
			rw:  httptest.NewRecorder(),
			req: nil,
			attachments: []With{
				WithMetadata(map[string]string{
					"anything": "value",
				}),
				WithData(map[string]string{
					"id": "Data object should be ignored",
				}),
			},
			expected: Response{
				Status: http.StatusText(http.StatusServiceUnavailable),
				Error: &errorstack.Error{
					Message: "Please try again in a few moments",
				},
				Metadata: map[string]string{
					"anything": "value",
				},
			},
		},
	}

	for _, tc := range testcases {
		WriteServiceUnavailable(tc.rw, tc.req, tc.attachments...)

		assert.Equal(t, http.StatusServiceUnavailable, tc.rw.Code)

		expected, _ := json.Marshal(tc.expected)
		actual := tc.rw.Body.Bytes()
		assert.JSONEq(t, string(expected), string(actual))
	}
}
