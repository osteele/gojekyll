package pages

import (
	"testing"

	"github.com/osteele/gojekyll/config"
	"github.com/osteele/liquid"
	"github.com/stretchr/testify/require"
)

func TestPage_ToLiquid(t *testing.T) {
	cfg := config.Default()
	page, err := NewFile(siteFake{t, cfg}, "testdata/excerpt.md", "excerpt.md", map[string]interface{}{})
	require.NoError(t, err)
	drop := page.(liquid.Drop).ToLiquid()
	excerpt := drop.(map[string]interface{})["excerpt"]
	// FIXME the following probably isn't right
	// TODO also test post-rendering.
	require.Equal(t, "First line.", excerpt)
}
