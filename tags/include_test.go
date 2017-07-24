package tags

import (
	"strings"
	"testing"

	"github.com/osteele/gojekyll/config"
	"github.com/osteele/liquid"
	"github.com/stretchr/testify/require"
)

func TestIncludeTag(t *testing.T) {
	engine := liquid.NewEngine()
	cfg := config.Default()
	cfg.Source = "testdata"
	AddJekyllTags(engine, &cfg, []string{"testdata/_includes"}, func(s string) (string, bool) {
		if s == "_posts/2017-07-04-test.md" {
			return "post.html", true
		}
		return "", false
	})
	bindings := map[string]interface{}{}

	s, err := engine.ParseAndRenderString(`{% include include_target.html %}`, bindings)
	require.NoError(t, err)
	require.Equal(t, "include target", strings.TrimSpace(s))

	// TODO {% include {{ page.my_variable }} %}
	// TODO {% include note.html content="This is my sample note." %}
}

func TestIncludeRelativeTag(t *testing.T) {
	engine := liquid.NewEngine()
	cfg := config.Default()
	AddJekyllTags(engine, &cfg, []string{}, func(s string) (string, bool) { return "", false })
	bindings := map[string]interface{}{}

	path := "testdata/dir/include_relative_source.md"
	tpl, err := engine.ParseTemplateLocation([]byte(`{% include_relative subdir/include_relative.html %}`), path, 1)
	require.NoError(t, err)
	s, err := tpl.Render(bindings)
	require.NoError(t, err)
	require.Equal(t, "include_relative target", strings.TrimSpace(string(s)))
}
