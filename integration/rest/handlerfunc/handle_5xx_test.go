package handlerfunc

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.nunchi.studio/helix/errorstack"
	"go.nunchi.studio/helix/integration/rest"

	"github.com/stretchr/testify/assert"
)

func TestInternalServerError(t *testing.T) {
	testcases := []struct {
		rw       *httptest.ResponseRecorder
		req      *http.Request
		expected *rest.Response
	}{
		{
			rw:  httptest.NewRecorder(),
			req: nil,
			expected: &rest.Response{
				Status: http.StatusText(http.StatusInternalServerError),
				Error: &errorstack.Error{
					Message: "We have been notified of this unexpected internal error",
				},
			},
		},
	}

	for _, tc := range testcases {
		InternalServerError(tc.rw, tc.req)

		assert.Equal(t, http.StatusInternalServerError, tc.rw.Code)

		var body *rest.Response
		json.Unmarshal(tc.rw.Body.Bytes(), &body)
		assert.Equal(t, tc.expected, body)
	}
}

func TestServiceUnavailable(t *testing.T) {
	testcases := []struct {
		rw       *httptest.ResponseRecorder
		req      *http.Request
		expected *rest.Response
	}{
		{
			rw:  httptest.NewRecorder(),
			req: nil,
			expected: &rest.Response{
				Status: http.StatusText(http.StatusServiceUnavailable),
				Error: &errorstack.Error{
					Message: "Please try again in a few moments",
				},
			},
		},
	}

	for _, tc := range testcases {
		ServiceUnavailable(tc.rw, tc.req)

		assert.Equal(t, http.StatusServiceUnavailable, tc.rw.Code)

		var body *rest.Response
		json.Unmarshal(tc.rw.Body.Bytes(), &body)
		assert.Equal(t, tc.expected, body)
	}
}
