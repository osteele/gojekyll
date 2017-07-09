package plugins

import (
	"regexp"
	"testing"

	"github.com/osteele/gojekyll/config"
	"github.com/osteele/liquid"
	"github.com/stretchr/testify/require"
)

func TestGistTag(t *testing.T) {
	engine := liquid.NewEngine()
	plugins := []string{"jekyll-gist"}
	Install(plugins, siteFake{config.Default(), engine})
	require.NoError(t, directory[plugins[0]].ConfigureTemplateEngine(engine))
	bindings := map[string]interface{}{}

	s, err := engine.ParseAndRenderString(`{% gist parkr/931c1c8d465a04042403 %}`, bindings)
	require.NoError(t, err)
	re := regexp.MustCompile(`<script.*>\s*</script>`)
	require.Contains(t, s, `src=https://gist.github.com/parkr/931c1c8d465a04042403.js`)
	require.True(t, re.MatchString(s))
}
