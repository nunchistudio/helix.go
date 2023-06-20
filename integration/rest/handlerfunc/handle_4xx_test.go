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

func TestBadRequest(t *testing.T) {
	testcases := []struct {
		rw          *httptest.ResponseRecorder
		req         *http.Request
		attachments []With
		expected    *rest.Response
	}{
		{
			rw:          httptest.NewRecorder(),
			req:         nil,
			attachments: nil,
			expected: &rest.Response{
				Status: http.StatusText(http.StatusBadRequest),
				Error: &errorstack.Error{
					Message: "Failed to validate request",
				},
			},
		},
		{
			rw:  httptest.NewRecorder(),
			req: nil,
			attachments: []With{
				WithValidations([]errorstack.Validation{
					{
						Message: "Property \"name\" is missing",
						Path:    []string{"request", "body", "event", "name"},
					},
				}),
			},
			expected: &rest.Response{
				Status: http.StatusText(http.StatusBadRequest),
				Error: &errorstack.Error{
					Message: "Failed to validate request",
					Validations: []errorstack.Validation{
						{
							Message: "Property \"name\" is missing",
							Path:    []string{"request", "body", "event", "name"},
						},
					},
				},
			},
		},
	}

	for _, tc := range testcases {
		BadRequest(tc.rw, tc.req, tc.attachments...)

		assert.Equal(t, http.StatusBadRequest, tc.rw.Code)

		var body *rest.Response
		json.Unmarshal(tc.rw.Body.Bytes(), &body)
		assert.Equal(t, tc.expected, body)
	}
}

func TestUnauthorized(t *testing.T) {
	testcases := []struct {
		rw          *httptest.ResponseRecorder
		req         *http.Request
		attachments []With
		expected    *rest.Response
	}{
		{
			rw:          httptest.NewRecorder(),
			req:         nil,
			attachments: nil,
			expected: &rest.Response{
				Status: http.StatusText(http.StatusUnauthorized),
				Error: &errorstack.Error{
					Message: "You are not authorized to perform this action",
				},
			},
		},
		{
			rw:  httptest.NewRecorder(),
			req: nil,
			attachments: []With{
				WithValidations([]errorstack.Validation{
					{
						Message: "Please contact the account owner",
					},
				}),
			},
			expected: &rest.Response{
				Status: http.StatusText(http.StatusUnauthorized),
				Error: &errorstack.Error{
					Message: "You are not authorized to perform this action",
					Validations: []errorstack.Validation{
						{
							Message: "Please contact the account owner",
						},
					},
				},
			},
		},
	}

	for _, tc := range testcases {
		Unauthorized(tc.rw, tc.req, tc.attachments...)

		assert.Equal(t, http.StatusUnauthorized, tc.rw.Code)

		var body *rest.Response
		json.Unmarshal(tc.rw.Body.Bytes(), &body)
		assert.Equal(t, tc.expected, body)
	}
}

func TestForbidden(t *testing.T) {
	testcases := []struct {
		rw          *httptest.ResponseRecorder
		req         *http.Request
		attachments []With
		expected    *rest.Response
	}{
		{
			rw:          httptest.NewRecorder(),
			req:         nil,
			attachments: nil,
			expected: &rest.Response{
				Status: http.StatusText(http.StatusForbidden),
				Error: &errorstack.Error{
					Message: "You don't have required permissions to perform this action",
				},
			},
		},
		{
			rw:  httptest.NewRecorder(),
			req: nil,
			attachments: []With{
				WithValidations([]errorstack.Validation{
					{
						Message: "Please contact the account owner",
					},
				}),
			},
			expected: &rest.Response{
				Status: http.StatusText(http.StatusForbidden),
				Error: &errorstack.Error{
					Message: "You don't have required permissions to perform this action",
					Validations: []errorstack.Validation{
						{
							Message: "Please contact the account owner",
						},
					},
				},
			},
		},
	}

	for _, tc := range testcases {
		Forbidden(tc.rw, tc.req, tc.attachments...)

		assert.Equal(t, http.StatusForbidden, tc.rw.Code)

		var body *rest.Response
		json.Unmarshal(tc.rw.Body.Bytes(), &body)
		assert.Equal(t, tc.expected, body)
	}
}

func TestNotFound(t *testing.T) {
	testcases := []struct {
		rw       *httptest.ResponseRecorder
		req      *http.Request
		expected *rest.Response
	}{
		{
			rw:  httptest.NewRecorder(),
			req: nil,
			expected: &rest.Response{
				Status: http.StatusText(http.StatusNotFound),
				Error: &errorstack.Error{
					Message: "Resource does not exist",
				},
			},
		},
	}

	for _, tc := range testcases {
		NotFound(tc.rw, tc.req)

		assert.Equal(t, http.StatusNotFound, tc.rw.Code)

		var body *rest.Response
		json.Unmarshal(tc.rw.Body.Bytes(), &body)
		assert.Equal(t, tc.expected, body)
	}
}

func TestMethodNotAllowed(t *testing.T) {
	testcases := []struct {
		rw       *httptest.ResponseRecorder
		req      *http.Request
		expected *rest.Response
	}{
		{
			rw:  httptest.NewRecorder(),
			req: nil,
			expected: &rest.Response{
				Status: http.StatusText(http.StatusMethodNotAllowed),
				Error: &errorstack.Error{
					Message: "Resource does not support this method",
				},
			},
		},
	}

	for _, tc := range testcases {
		MethodNotAllowed(tc.rw, tc.req)

		assert.Equal(t, http.StatusMethodNotAllowed, tc.rw.Code)

		var body *rest.Response
		json.Unmarshal(tc.rw.Body.Bytes(), &body)
		assert.Equal(t, tc.expected, body)
	}
}

func TestRequestEntityTooLarge(t *testing.T) {
	testcases := []struct {
		rw          *httptest.ResponseRecorder
		req         *http.Request
		attachments []With
		expected    *rest.Response
	}{
		{
			rw:          httptest.NewRecorder(),
			req:         nil,
			attachments: nil,
			expected: &rest.Response{
				Status: http.StatusText(http.StatusRequestEntityTooLarge),
				Error: &errorstack.Error{
					Message: "Can not process payload too large",
				},
			},
		},
		{
			rw:  httptest.NewRecorder(),
			req: nil,
			attachments: []With{
				WithValidations([]errorstack.Validation{
					{
						Message: "Payload must not be larger than 400kb",
					},
				}),
			},
			expected: &rest.Response{
				Status: http.StatusText(http.StatusRequestEntityTooLarge),
				Error: &errorstack.Error{
					Message: "Can not process payload too large",
					Validations: []errorstack.Validation{
						{
							Message: "Payload must not be larger than 400kb",
						},
					},
				},
			},
		},
	}

	for _, tc := range testcases {
		RequestEntityTooLarge(tc.rw, tc.req, tc.attachments...)

		assert.Equal(t, http.StatusRequestEntityTooLarge, tc.rw.Code)

		var body *rest.Response
		json.Unmarshal(tc.rw.Body.Bytes(), &body)
		assert.Equal(t, tc.expected, body)
	}
}

func TestTooManyRequests(t *testing.T) {
	testcases := []struct {
		rw          *httptest.ResponseRecorder
		req         *http.Request
		attachments []With
		expected    *rest.Response
	}{
		{
			rw:          httptest.NewRecorder(),
			req:         nil,
			attachments: nil,
			expected: &rest.Response{
				Status: http.StatusText(http.StatusTooManyRequests),
				Error: &errorstack.Error{
					Message: "Request-rate limit has been reached",
				},
			},
		},
		{
			rw:  httptest.NewRecorder(),
			req: nil,
			attachments: []With{
				WithValidations([]errorstack.Validation{
					{
						Message: "Subscription is limited to 100 requests per minute",
					},
				}),
			},
			expected: &rest.Response{
				Status: http.StatusText(http.StatusTooManyRequests),
				Error: &errorstack.Error{
					Message: "Request-rate limit has been reached",
					Validations: []errorstack.Validation{
						{
							Message: "Subscription is limited to 100 requests per minute",
						},
					},
				},
			},
		},
	}

	for _, tc := range testcases {
		TooManyRequests(tc.rw, tc.req, tc.attachments...)

		assert.Equal(t, http.StatusTooManyRequests, tc.rw.Code)

		var body *rest.Response
		json.Unmarshal(tc.rw.Body.Bytes(), &body)
		assert.Equal(t, tc.expected, body)
	}
}
