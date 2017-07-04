package pages

import (
	"path/filepath"
	"testing"

	"github.com/osteele/gojekyll/config"
	"github.com/stretchr/testify/require"
)

type containerMock struct {
	c      config.Config
	prefix string
}

func (c containerMock) OutputExt(p string) string { return filepath.Ext(p) }

func (c containerMock) PathPrefix() string { return c.prefix }

func TestPageCategories(t *testing.T) {
	require.Equal(t, []string{"a", "b"}, sortedStringValue("b a"))
	require.Equal(t, []string{"a", "b"}, sortedStringValue([]interface{}{"b", "a"}))
	require.Equal(t, []string{"a", "b"}, sortedStringValue([]string{"b", "a"}))
	require.Equal(t, []string{}, sortedStringValue(3))

	c := containerMock{config.Default(), ""}
	fm := map[string]interface{}{"categories": "b a"}
	f := file{container: c, frontMatter: fm}
	require.Equal(t, []string{"a", "b"}, f.Categories())
}
