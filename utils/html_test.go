package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProcessAnchorHrefs(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		{
			name:     "modifies href attribute",
			input:    []byte(`<html><body><a href="test.md">Link</a></body></html>`),
			expected: []byte(`<html><body><a href="modified">Link</a></body></html>`),
		},
		{
			name:     "preserves other attributes",
			input:    []byte(`<html><body><a href="test.md" class="link">Link</a></body></html>`),
			expected: []byte(`<html><body><a href="modified" class="link">Link</a></body></html>`),
		},
		{
			name:     "handles multiple links",
			input:    []byte(`<html><body><a href="one.md">One</a> <a href="two.md">Two</a></body></html>`),
			expected: []byte(`<html><body><a href="modified">One</a> <a href="modified">Two</a></body></html>`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProcessAnchorHrefs(tt.input, func(href string) string {
				if href == "test.md" || href == "one.md" || href == "two.md" {
					return "modified"
				}
				return href
			})
			require.Equal(t, string(tt.expected), string(result))
		})
	}
}
