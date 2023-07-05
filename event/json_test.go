package event

import (
	"encoding/json"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventFromJSON(t *testing.T) {
	testcases := []struct {
		input    json.RawMessage
		expected Event
		success  bool
	}{
		{
			input:    []byte(`not a valid JSON input`),
			expected: Event{},
			success:  false,
		},
		{
			input:    []byte(`{ "_no_event_key": {} }`),
			expected: Event{},
			success:  false,
		},
		{
			input: []byte(`{ "event": {
        "invalid_key": true
      } }`),
			expected: Event{},
			success:  true,
		},
		{
			input: []byte(`{ "event": {
        "name": "testing",
        "params": {
          "filters": ["a", "b"],
          "query": ["search_query"]
        }
      } }`),
			expected: Event{
				Name: "testing",
				Params: url.Values{
					"filters": []string{"a", "b"},
					"query":   []string{"search_query"},
				},
			},
			success: true,
		},
	}

	for _, tc := range testcases {
		actual, ok := EventFromJSON(tc.input)

		assert.Equal(t, tc.expected, actual)
		assert.Equal(t, tc.success, ok)
	}
}
