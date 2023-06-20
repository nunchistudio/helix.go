package handlerfunc

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.nunchi.studio/helix/integration/rest"

	"github.com/stretchr/testify/assert"
)

func TestAccepted(t *testing.T) {
	testcases := []struct {
		rw       *httptest.ResponseRecorder
		req      *http.Request
		expected *rest.Response
	}{
		{
			rw:  httptest.NewRecorder(),
			req: nil,
			expected: &rest.Response{
				Status: http.StatusText(http.StatusAccepted),
			},
		},
	}

	for _, tc := range testcases {
		Accepted(tc.rw, tc.req)

		assert.Equal(t, http.StatusAccepted, tc.rw.Code)

		var body *rest.Response
		json.Unmarshal(tc.rw.Body.Bytes(), &body)
		assert.Equal(t, tc.expected, body)
	}
}
