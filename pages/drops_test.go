package pages

import (
	"fmt"
	"testing"
	"time"

	"github.com/osteele/gojekyll/config"
	"github.com/osteele/liquid"
	"github.com/osteele/liquid/tags"
	"github.com/stretchr/testify/require"
)

func TestStaticFile_ToLiquid(t *testing.T) {
	site := siteFake{t, config.Default()}
	page, err := NewFile(site, "testdata/static.html", "static.html", map[string]interface{}{})
	require.NoError(t, err)
	drop := page.(liquid.Drop).ToLiquid().(tags.IterationKeyedMap)

	require.Equal(t, "static", drop["basename"])
	require.Equal(t, "static.html", drop["name"])
	require.Equal(t, "/static.html", drop["path"])
	require.Equal(t, ".html", drop["extname"])
	require.IsType(t, time.Now(), drop["modified_time"])
}

func TestPage_ToLiquid_excerpt(t *testing.T) {
	site := siteFake{t, config.Default()}
	p, err := NewFile(site, "testdata/excerpt.md", "excerpt.md", map[string]interface{}{})
	require.NoError(t, err)

	t.Run("before render", func(t *testing.T) {
		drop := p.(liquid.Drop).ToLiquid()
		excerpt := drop.(tags.IterationKeyedMap)["excerpt"]
		require.Equal(t, "First line.", fmt.Sprintf("%s", excerpt))
	})

	t.Run("after render", func(t *testing.T) {
		require.NoError(t, p.(renderer).Render())
		drop := p.(liquid.Drop).ToLiquid()
		excerpt := drop.(tags.IterationKeyedMap)["excerpt"]
		require.Equal(t, "rendered: First line.", fmt.Sprintf("%s", excerpt))
	})
}

type renderer interface {
	Render() error
}
