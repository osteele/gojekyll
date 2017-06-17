package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsMarkdown(t *testing.T) {
	site := NewSite()
	require.True(t, site.IsMarkdown("name.md"))
	require.True(t, site.IsMarkdown("name.markdown"))
	require.False(t, site.IsMarkdown("name.html"))
}
