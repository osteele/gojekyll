package pages

import (
	"testing"

	"github.com/osteele/gojekyll/config"
	"github.com/stretchr/testify/require"
)

func TestPage_ToLiquid(t *testing.T) {
	cfg := config.Default()
	page, err := NewFile("testdata/excerpt.md", containerFake{cfg, ""}, "excerpt.md", map[string]interface{}{})
	require.NoError(t, err)
	drop := page.ToLiquid()
	excerpt := drop.(map[string]interface{})["excerpt"]
	// FIXME the following probably isn't right
	// TODO also test post-rendering.
	require.Equal(t, "First line.", excerpt)
}
