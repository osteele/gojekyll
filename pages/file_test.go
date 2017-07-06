package pages

import (
	"path/filepath"
	"testing"

	"github.com/osteele/gojekyll/config"
	"github.com/stretchr/testify/require"
)

type containerFake struct {
	cfg    config.Config
	prefix string
}

func (c containerFake) Config() *config.Config    { return &c.cfg }
func (c containerFake) PathPrefix() string        { return c.prefix }
func (c containerFake) OutputExt(p string) string { return filepath.Ext(p) }

func TestPageCategories(t *testing.T) {
	require.Equal(t, []string{"a", "b"}, sortedStringValue("b a"))
	require.Equal(t, []string{"a", "b"}, sortedStringValue([]interface{}{"b", "a"}))
	require.Equal(t, []string{"a", "b"}, sortedStringValue([]string{"b", "a"}))
	require.Equal(t, []string{}, sortedStringValue(3))

	c := containerFake{config.Default(), ""}
	fm := map[string]interface{}{"categories": "b a"}
	f := file{container: c, frontMatter: fm}
	require.Equal(t, []string{"a", "b"}, f.Categories())
}
