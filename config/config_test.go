package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	c := Default()
	require.Equal(t, ".", c.Source)
	require.Equal(t, "./_site", c.Destination)
	require.Equal(t, "_layouts", c.LayoutsDir)
}

func TestUnmarshal(t *testing.T) {
	c := Default()
	err := Unmarshal([]byte(`source: x`), &c)
	require.NoError(t, err)
	require.Equal(t, "x", c.Source)
	require.Equal(t, "./_site", c.Destination)
}

func TestIsMarkdown(t *testing.T) {
	c := Default()
	require.True(t, c.IsMarkdown("name.md"))
	require.True(t, c.IsMarkdown("name.markdown"))
	require.False(t, c.IsMarkdown("name.html"))
}
