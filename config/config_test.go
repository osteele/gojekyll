package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSourceDir(t *testing.T) {
	c := Default()
	require.True(t, strings.HasPrefix(c.SourceDir(), "/"))
}
func TestDefaultConfig(t *testing.T) {
	c := Default()
	require.Equal(t, ".", c.Source)
	require.Equal(t, "./_site", c.Destination)
	require.Equal(t, "_layouts", c.LayoutsDir)
}

func TestPlugins(t *testing.T) {
	c := Default()
	Unmarshal([]byte(`plugins: ['a']`), &c)
	require.Equal(t, []string{"a"}, c.Plugins)

	c = Default()
	Unmarshal([]byte(`gems: ['a']`), &c)
	require.Equal(t, []string{"a"}, c.Plugins)
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
