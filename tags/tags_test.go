package tags

import (
	"testing"

	"github.com/osteele/gojekyll/config"
	"github.com/osteele/liquid"
	"github.com/stretchr/testify/require"
)

func TestTags(t *testing.T) {
	engine := liquid.NewEngine()
	cfg := config.Default()
	AddJekyllTags(engine, &cfg, []string{}, func(s string) (string, bool) {
		if s == "_posts/2017-07-04-test.md" {
			return "post.html", true
		}
		return "", false
	})

	s, err := engine.ParseAndRenderString(`{% post_url 2017-07-04-test.md %}`, liquid.Bindings{})
	require.NoError(t, err)
	require.Equal(t, "post.html", s)
}
