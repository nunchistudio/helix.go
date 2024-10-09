package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bunrouter"
)

func TestParamsFromContext(t *testing.T) {
	testcases := []struct {
		config     Config
		registered string
		requested  string
		expected   map[string]string
		found      bool
	}{
		{
			registered: "/hello",
			requested:  "/hello",
			expected:   nil,
			found:      false,
		},
		{
			registered: "/users/:username",
			requested:  "/users/nunchistudio",
			expected: map[string]string{
				"username": "nunchistudio",
			},
			found: true,
		},
	}

	for _, tc := range testcases {
		r, _ := New(tc.config)
		r.GET(tc.registered, func(rw http.ResponseWriter, req *http.Request) {
			params, found := ParamsFromContext(req.Context())

			assert.Equal(t, tc.expected, params)
			assert.Equal(t, tc.found, found)
		})

		rw := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, tc.requested, nil)

		rc := r.(*bunrouter.CompatRouter)
		rc.ServeHTTP(rw, req)
	}
}
