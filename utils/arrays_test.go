package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSearchStrings(t *testing.T) {
	tests := []struct {
		a      []string
		s      string
		expect bool
	}{
		{[]string{}, "a", false},
		{[]string{"a", "b"}, "a", true},
		{[]string{"a", "b"}, "b", true},
		{[]string{"a", "b"}, "c", false},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("SearchStrings(%v, %v)", tt.a, tt.s), func(t *testing.T) {
			require.Equal(t, tt.expect, SearchStrings(tt.a, tt.s))
		})
	}
}
