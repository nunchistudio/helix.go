package event

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/baggage"
)

func TestInjectEventPageToFlatMap(t *testing.T) {
	testcases := []struct {
		input    Page
		expected map[string]string
	}{
		{
			input:    Page{},
			expected: map[string]string{},
		},
		{
			input: Page{
				Path:     "page_path_test",
				Referrer: "page_referrer_test",
				Search:   "page_search_test",
				Title:    "page_title_test",
				URL:      "page_url_test",
			},
			expected: map[string]string{
				"event.page.path":     "page_path_test",
				"event.page.referrer": "page_referrer_test",
				"event.page.search":   "page_search_test",
				"event.page.title":    "page_title_test",
				"event.page.url":      "page_url_test",
			},
		},
	}

	for _, tc := range testcases {
		var actual = make(map[string]string)
		for k, v := range tc.expected {
			actual[k] = v
		}

		injectEventPageToFlatMap(tc.input, tc.expected)

		assert.Equal(t, tc.expected, actual)
	}
}

func TestApplyEventPageFromBaggageMember(t *testing.T) {
	testcases := []struct {
		input    func() baggage.Member
		expected *Event
	}{
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.page.unknown", "anything")
				return m
			},
			expected: &Event{
				Page: Page{},
			},
		},
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.page.path", "page_path_test")
				return m
			},
			expected: &Event{
				Page: Page{
					Path: "page_path_test",
				},
			},
		},
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.page.referrer", "page_referrer_test")
				return m
			},
			expected: &Event{
				Page: Page{
					Referrer: "page_referrer_test",
				},
			},
		},
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.page.search", "page_search_test")
				return m
			},
			expected: &Event{
				Page: Page{
					Search: "page_search_test",
				},
			},
		},
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.page.title", "page_title_test")
				return m
			},
			expected: &Event{
				Page: Page{
					Title: "page_title_test",
				},
			},
		},
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.page.url", "page_url_test")
				return m
			},
			expected: &Event{
				Page: Page{
					URL: "page_url_test",
				},
			},
		},
	}

	for _, tc := range testcases {
		actual := new(Event)
		applyEventPageFromBaggageMember(tc.input(), actual)

		assert.Equal(t, tc.expected, actual)
	}
}
