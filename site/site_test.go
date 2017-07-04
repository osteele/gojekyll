package site

import (
	"testing"

	"github.com/osteele/gojekyll/config"
	"github.com/stretchr/testify/require"
)

func TestIsMarkdown(t *testing.T) {
	s := New(config.Flags{})
	require.Equal(t, "", s.PathPrefix())
	require.False(t, s.KeepFile("random"))
	require.True(t, s.KeepFile(".git"))
	require.True(t, s.KeepFile(".svn"))
}
