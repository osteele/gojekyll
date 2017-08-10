package pages

import (
	"testing"
	"time"

	"github.com/osteele/gojekyll/config"
	"github.com/osteele/liquid"
	"github.com/stretchr/testify/require"
)

func TestStaticFile_ToLiquid(t *testing.T) {
	site := siteFake{t, config.Default()}
	page, err := NewFile(site, "testdata/static.html", "static.html", map[string]interface{}{})
	require.NoError(t, err)
	drop := page.(liquid.Drop).ToLiquid().(map[string]interface{})

	require.Equal(t, "static", drop["basename"])
	require.Equal(t, "static.html", drop["name"])
	require.Equal(t, "/static.html", drop["path"])
	require.Equal(t, ".html", drop["extname"])
	require.IsType(t, time.Now(), drop["modified_time"])
}

func TestPage_ToLiquid(t *testing.T) {
	site := siteFake{t, config.Default()}
	p, err := NewFile(site, "testdata/excerpt.md", "excerpt.md", map[string]interface{}{})
	require.NoError(t, err)
	_, err = p.(Page).Content()
	require.NoError(t, err)
	drop := p.(liquid.Drop).ToLiquid()
	excerpt := drop.(map[string]interface{})["excerpt"]
	ex, ok := excerpt.(string)
	require.True(t, ok)
	require.Equal(t, "First line.", ex)
}
