package rest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.nunchi.studio/helix/errorstack"

	"github.com/stretchr/testify/assert"
)

func TestWriteBadRequest(t *testing.T) {
	addFrenchLocalesForTest()
	reqWithLang, _ := http.NewRequest(http.MethodPost, "/anything", nil)
	reqWithLang.Header.Add("Accept-Language", "fr")

	testcases := []struct {
		rw          *httptest.ResponseRecorder
		req         *http.Request
		withOnError []WithOnError
		expected    Response
	}{
		{
			rw:          httptest.NewRecorder(),
			req:         nil,
			withOnError: nil,
			expected: Response{
				Status: http.StatusText(http.StatusBadRequest),
				Error: &errorstack.Error{
					Message: "Failed to validate request",
				},
			},
		},
		{
			rw:          httptest.NewRecorder(),
			req:         reqWithLang,
			withOnError: nil,
			expected: Response{
				Status: http.StatusText(http.StatusBadRequest),
				Error: &errorstack.Error{
					Message: "Échec de la validation de la requête",
				},
			},
		},
		{
			rw:  httptest.NewRecorder(),
			req: nil,
			withOnError: []WithOnError{
				WithMetadataOnError(map[string]string{
					"anything": "value",
				}),
			},
			expected: Response{
				Status: http.StatusText(http.StatusBadRequest),
				Error: &errorstack.Error{
					Message: "Failed to validate request",
				},
				Metadata: map[string]string{
					"anything": "value",
				},
			},
		},
		{
			rw:  httptest.NewRecorder(),
			req: nil,
			withOnError: []WithOnError{
				WithValidationsOnError([]errorstack.Validation{
					{
						Message: "Property \"name\" is missing",
						Path:    []string{"request", "body", "event", "name"},
					},
				}),
			},
			expected: Response{
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
		WriteBadRequest[Response](tc.rw, tc.req, tc.withOnError...)

		assert.Equal(t, http.StatusBadRequest, tc.rw.Code)

		expected, _ := json.Marshal(tc.expected)
		actual := tc.rw.Body.Bytes()
		assert.JSONEq(t, string(expected), string(actual))
	}
}

func TestWriteUnauthorized(t *testing.T) {
	addFrenchLocalesForTest()
	reqWithLang, _ := http.NewRequest(http.MethodPost, "/anything", nil)
	reqWithLang.Header.Add("Accept-Language", "fr")

	testcases := []struct {
		rw          *httptest.ResponseRecorder
		req         *http.Request
		withOnError []WithOnError
		expected    Response
	}{
		{
			rw:          httptest.NewRecorder(),
			req:         nil,
			withOnError: nil,
			expected: Response{
				Status: http.StatusText(http.StatusUnauthorized),
				Error: &errorstack.Error{
					Message: "You are not authorized to perform this action",
				},
			},
		},
		{
			rw:          httptest.NewRecorder(),
			req:         reqWithLang,
			withOnError: nil,
			expected: Response{
				Status: http.StatusText(http.StatusUnauthorized),
				Error: &errorstack.Error{
					Message: "Vous n'êtes pas autorisé à effectuer cette action",
				},
			},
		},
		{
			rw:  httptest.NewRecorder(),
			req: nil,
			withOnError: []WithOnError{
				WithMetadataOnError(map[string]string{
					"anything": "value",
				}),
			},
			expected: Response{
				Status: http.StatusText(http.StatusUnauthorized),
				Error: &errorstack.Error{
					Message: "You are not authorized to perform this action",
				},
				Metadata: map[string]string{
					"anything": "value",
				},
			},
		},
		{
			rw:  httptest.NewRecorder(),
			req: nil,
			withOnError: []WithOnError{
				WithValidationsOnError([]errorstack.Validation{
					{
						Message: "Please contact the account owner",
					},
				}),
			},
			expected: Response{
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
		WriteUnauthorized[Response](tc.rw, tc.req, tc.withOnError...)

		assert.Equal(t, http.StatusUnauthorized, tc.rw.Code)

		expected, _ := json.Marshal(tc.expected)
		actual := tc.rw.Body.Bytes()
		assert.JSONEq(t, string(expected), string(actual))
	}
}

func TestWritePaymentRequired(t *testing.T) {
	addFrenchLocalesForTest()
	reqWithLang, _ := http.NewRequest(http.MethodPost, "/anything", nil)
	reqWithLang.Header.Add("Accept-Language", "fr")

	testcases := []struct {
		rw          *httptest.ResponseRecorder
		req         *http.Request
		withOnError []WithOnError
		expected    Response
	}{
		{
			rw:          httptest.NewRecorder(),
			req:         nil,
			withOnError: nil,
			expected: Response{
				Status: http.StatusText(http.StatusPaymentRequired),
				Error: &errorstack.Error{
					Message: "Request failed because payment is required",
				},
			},
		},
		{
			rw:          httptest.NewRecorder(),
			req:         reqWithLang,
			withOnError: nil,
			expected: Response{
				Status: http.StatusText(http.StatusPaymentRequired),
				Error: &errorstack.Error{
					Message: "La requête a échoué car un paiement est requis",
				},
			},
		},
		{
			rw:  httptest.NewRecorder(),
			req: nil,
			withOnError: []WithOnError{
				WithMetadataOnError(map[string]string{
					"anything": "value",
				}),
			},
			expected: Response{
				Status: http.StatusText(http.StatusPaymentRequired),
				Error: &errorstack.Error{
					Message: "Request failed because payment is required",
				},
				Metadata: map[string]string{
					"anything": "value",
				},
			},
		},
		{
			rw:  httptest.NewRecorder(),
			req: nil,
			withOnError: []WithOnError{
				WithValidationsOnError([]errorstack.Validation{
					{
						Message: "Please contact the account owner",
					},
				}),
			},
			expected: Response{
				Status: http.StatusText(http.StatusPaymentRequired),
				Error: &errorstack.Error{
					Message: "Request failed because payment is required",
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
		WritePaymentRequired[Response](tc.rw, tc.req, tc.withOnError...)

		assert.Equal(t, http.StatusPaymentRequired, tc.rw.Code)

		expected, _ := json.Marshal(tc.expected)
		actual := tc.rw.Body.Bytes()
		assert.JSONEq(t, string(expected), string(actual))
	}
}

func TestWriteForbidden(t *testing.T) {
	addFrenchLocalesForTest()
	reqWithLang, _ := http.NewRequest(http.MethodPost, "/anything", nil)
	reqWithLang.Header.Add("Accept-Language", "fr")

	testcases := []struct {
		rw          *httptest.ResponseRecorder
		req         *http.Request
		withOnError []WithOnError
		expected    Response
	}{
		{
			rw:          httptest.NewRecorder(),
			req:         nil,
			withOnError: nil,
			expected: Response{
				Status: http.StatusText(http.StatusForbidden),
				Error: &errorstack.Error{
					Message: "You don't have required permissions to perform this action",
				},
			},
		},
		{
			rw:          httptest.NewRecorder(),
			req:         reqWithLang,
			withOnError: nil,
			expected: Response{
				Status: http.StatusText(http.StatusForbidden),
				Error: &errorstack.Error{
					Message: "Vous n'avez pas les permissions requises pour effectuer cette action",
				},
			},
		},
		{
			rw:  httptest.NewRecorder(),
			req: nil,
			withOnError: []WithOnError{
				WithMetadataOnError(map[string]string{
					"anything": "value",
				}),
			},
			expected: Response{
				Status: http.StatusText(http.StatusForbidden),
				Error: &errorstack.Error{
					Message: "You don't have required permissions to perform this action",
				},
				Metadata: map[string]string{
					"anything": "value",
				},
			},
		},
		{
			rw:  httptest.NewRecorder(),
			req: nil,
			withOnError: []WithOnError{
				WithValidationsOnError([]errorstack.Validation{
					{
						Message: "Please contact the account owner",
					},
				}),
			},
			expected: Response{
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
		WriteForbidden[Response](tc.rw, tc.req, tc.withOnError...)

		assert.Equal(t, http.StatusForbidden, tc.rw.Code)

		expected, _ := json.Marshal(tc.expected)
		actual := tc.rw.Body.Bytes()
		assert.JSONEq(t, string(expected), string(actual))
	}
}

func TestWriteNotFound(t *testing.T) {
	addFrenchLocalesForTest()
	reqWithLang, _ := http.NewRequest(http.MethodPost, "/anything", nil)
	reqWithLang.Header.Add("Accept-Language", "fr")

	testcases := []struct {
		rw          *httptest.ResponseRecorder
		req         *http.Request
		withOnError []WithOnError
		expected    Response
	}{
		{
			rw:  httptest.NewRecorder(),
			req: nil,
			expected: Response{
				Status: http.StatusText(http.StatusNotFound),
				Error: &errorstack.Error{
					Message: "Resource does not exist",
				},
			},
		},
		{
			rw:          httptest.NewRecorder(),
			req:         reqWithLang,
			withOnError: nil,
			expected: Response{
				Status: http.StatusText(http.StatusNotFound),
				Error: &errorstack.Error{
					Message: "La ressource n'existe pas",
				},
			},
		},
		{
			rw:  httptest.NewRecorder(),
			req: nil,
			withOnError: []WithOnError{
				WithMetadataOnError(map[string]string{
					"anything": "value",
				}),
			},
			expected: Response{
				Status: http.StatusText(http.StatusNotFound),
				Error: &errorstack.Error{
					Message: "Resource does not exist",
				},
				Metadata: map[string]string{
					"anything": "value",
				},
			},
		},
	}

	for _, tc := range testcases {
		WriteNotFound[Response](tc.rw, tc.req, tc.withOnError...)

		assert.Equal(t, http.StatusNotFound, tc.rw.Code)

		expected, _ := json.Marshal(tc.expected)
		actual := tc.rw.Body.Bytes()
		assert.JSONEq(t, string(expected), string(actual))
	}
}

func TestWriteMethodNotAllowed(t *testing.T) {
	addFrenchLocalesForTest()
	reqWithLang, _ := http.NewRequest(http.MethodPost, "/anything", nil)
	reqWithLang.Header.Add("Accept-Language", "fr")

	testcases := []struct {
		rw          *httptest.ResponseRecorder
		req         *http.Request
		withOnError []WithOnError
		expected    Response
	}{
		{
			rw:  httptest.NewRecorder(),
			req: nil,
			expected: Response{
				Status: http.StatusText(http.StatusMethodNotAllowed),
				Error: &errorstack.Error{
					Message: "Resource does not support this method",
				},
			},
		},
		{
			rw:          httptest.NewRecorder(),
			req:         reqWithLang,
			withOnError: nil,
			expected: Response{
				Status: http.StatusText(http.StatusMethodNotAllowed),
				Error: &errorstack.Error{
					Message: "La ressource ne supporte pas cette méthode",
				},
			},
		},
		{
			rw:  httptest.NewRecorder(),
			req: nil,
			withOnError: []WithOnError{
				WithMetadataOnError(map[string]string{
					"anything": "value",
				}),
			},
			expected: Response{
				Status: http.StatusText(http.StatusMethodNotAllowed),
				Error: &errorstack.Error{
					Message: "Resource does not support this method",
				},
				Metadata: map[string]string{
					"anything": "value",
				},
			},
		},
	}

	for _, tc := range testcases {
		WriteMethodNotAllowed[Response](tc.rw, tc.req, tc.withOnError...)

		assert.Equal(t, http.StatusMethodNotAllowed, tc.rw.Code)

		expected, _ := json.Marshal(tc.expected)
		actual := tc.rw.Body.Bytes()
		assert.JSONEq(t, string(expected), string(actual))
	}
}

func TestWriteConflict(t *testing.T) {
	addFrenchLocalesForTest()
	reqWithLang, _ := http.NewRequest(http.MethodPost, "/anything", nil)
	reqWithLang.Header.Add("Accept-Language", "fr")

	testcases := []struct {
		rw          *httptest.ResponseRecorder
		req         *http.Request
		withOnError []WithOnError
		expected    Response
	}{
		{
			rw:  httptest.NewRecorder(),
			req: nil,
			expected: Response{
				Status: http.StatusText(http.StatusConflict),
				Error: &errorstack.Error{
					Message: "Failed to process target resource because of conflict",
				},
			},
		},
		{
			rw:          httptest.NewRecorder(),
			req:         reqWithLang,
			withOnError: nil,
			expected: Response{
				Status: http.StatusText(http.StatusConflict),
				Error: &errorstack.Error{
					Message: "Échec du traitement de la requête en raison d'un conflit",
				},
			},
		},
		{
			rw:  httptest.NewRecorder(),
			req: nil,
			withOnError: []WithOnError{
				WithMetadataOnError(map[string]string{
					"anything": "value",
				}),
			},
			expected: Response{
				Status: http.StatusText(http.StatusConflict),
				Error: &errorstack.Error{
					Message: "Failed to process target resource because of conflict",
				},
				Metadata: map[string]string{
					"anything": "value",
				},
			},
		},
	}

	for _, tc := range testcases {
		WriteConflict[Response](tc.rw, tc.req, tc.withOnError...)

		assert.Equal(t, http.StatusConflict, tc.rw.Code)

		expected, _ := json.Marshal(tc.expected)
		actual := tc.rw.Body.Bytes()
		assert.JSONEq(t, string(expected), string(actual))
	}
}

func TestWriteRequestEntityTooLarge(t *testing.T) {
	addFrenchLocalesForTest()
	reqWithLang, _ := http.NewRequest(http.MethodPost, "/anything", nil)
	reqWithLang.Header.Add("Accept-Language", "fr")

	testcases := []struct {
		rw          *httptest.ResponseRecorder
		req         *http.Request
		withOnError []WithOnError
		expected    Response
	}{
		{
			rw:          httptest.NewRecorder(),
			req:         nil,
			withOnError: nil,
			expected: Response{
				Status: http.StatusText(http.StatusRequestEntityTooLarge),
				Error: &errorstack.Error{
					Message: "Can not process payload too large",
				},
			},
		},
		{
			rw:          httptest.NewRecorder(),
			req:         reqWithLang,
			withOnError: nil,
			expected: Response{
				Status: http.StatusText(http.StatusRequestEntityTooLarge),
				Error: &errorstack.Error{
					Message: "Impossible de traiter une requête avec un payload trop large",
				},
			},
		},
		{
			rw:  httptest.NewRecorder(),
			req: nil,
			withOnError: []WithOnError{
				WithMetadataOnError(map[string]string{
					"anything": "value",
				}),
			},
			expected: Response{
				Status: http.StatusText(http.StatusRequestEntityTooLarge),
				Error: &errorstack.Error{
					Message: "Can not process payload too large",
				},
				Metadata: map[string]string{
					"anything": "value",
				},
			},
		},
		{
			rw:  httptest.NewRecorder(),
			req: nil,
			withOnError: []WithOnError{
				WithValidationsOnError([]errorstack.Validation{
					{
						Message: "Payload must not be larger than 400kb",
					},
				}),
			},
			expected: Response{
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
		WriteRequestEntityTooLarge[Response](tc.rw, tc.req, tc.withOnError...)

		assert.Equal(t, http.StatusRequestEntityTooLarge, tc.rw.Code)

		expected, _ := json.Marshal(tc.expected)
		actual := tc.rw.Body.Bytes()
		assert.JSONEq(t, string(expected), string(actual))
	}
}

func TestWriteTooManyRequests(t *testing.T) {
	addFrenchLocalesForTest()
	reqWithLang, _ := http.NewRequest(http.MethodPost, "/anything", nil)
	reqWithLang.Header.Add("Accept-Language", "fr")

	testcases := []struct {
		rw          *httptest.ResponseRecorder
		req         *http.Request
		withOnError []WithOnError
		expected    Response
	}{
		{
			rw:          httptest.NewRecorder(),
			req:         nil,
			withOnError: nil,
			expected: Response{
				Status: http.StatusText(http.StatusTooManyRequests),
				Error: &errorstack.Error{
					Message: "Request-rate limit has been reached",
				},
			},
		},
		{
			rw:          httptest.NewRecorder(),
			req:         reqWithLang,
			withOnError: nil,
			expected: Response{
				Status: http.StatusText(http.StatusTooManyRequests),
				Error: &errorstack.Error{
					Message: "La limite du taux de requêtes a été atteinte",
				},
			},
		},
		{
			rw:  httptest.NewRecorder(),
			req: nil,
			withOnError: []WithOnError{
				WithMetadataOnError(map[string]string{
					"anything": "value",
				}),
			},
			expected: Response{
				Status: http.StatusText(http.StatusTooManyRequests),
				Error: &errorstack.Error{
					Message: "Request-rate limit has been reached",
				},
				Metadata: map[string]string{
					"anything": "value",
				},
			},
		},
		{
			rw:  httptest.NewRecorder(),
			req: nil,
			withOnError: []WithOnError{
				WithValidationsOnError([]errorstack.Validation{
					{
						Message: "Subscription is limited to 100 requests per minute",
					},
				}),
			},
			expected: Response{
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
		WriteTooManyRequests[Response](tc.rw, tc.req, tc.withOnError...)

		assert.Equal(t, http.StatusTooManyRequests, tc.rw.Code)

		expected, _ := json.Marshal(tc.expected)
		actual := tc.rw.Body.Bytes()
		assert.JSONEq(t, string(expected), string(actual))
	}
}
