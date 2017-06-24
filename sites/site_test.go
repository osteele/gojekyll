package sites

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsMarkdown(t *testing.T) {
	s := NewSite()
	require.Equal(t, "", s.PathPrefix())
	require.False(t, s.KeepFile("random"))
}
