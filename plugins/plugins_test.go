package plugins

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/osteele/liquid"
	"github.com/stretchr/testify/require"
)

type train struct{ e liquid.Engine }

func (t train) TemplateEngine() liquid.Engine { return t.e }

func TestAvatarTag(t *testing.T) {
	engine := liquid.NewEngine()
	bindings := map[string]interface{}{"user": "osteele"}
	Install("jekyll-avatar", train{engine})

	s, err := engine.ParseAndRenderString(`{% avatar osteele %}`, bindings)
	require.NoError(t, err)
	re := regexp.MustCompile(`<img class="avatar.*avatar.*usercontent\.com/osteele\b`)
	require.True(t, re.MatchString(s))

	s, err = engine.ParseAndRenderString(`{% avatar user='osteele' %}`, bindings)
	require.NoError(t, err)
	require.True(t, re.MatchString(s))

	s, err = engine.ParseAndRenderString(`{% avatar user=user %}`, bindings)
	require.NoError(t, err)
	require.True(t, re.MatchString(s))

	s, err = engine.ParseAndRenderString(`{% avatar user=user size=20 %}`, bindings)
	require.NoError(t, err)
	require.Contains(t, s, "20")
}

func TestGistTag(t *testing.T) {
	engine := liquid.NewEngine()
	bindings := map[string]interface{}{}
	Install("jekyll-gist", train{engine})

	s, err := engine.ParseAndRenderString(`{% gist parkr/931c1c8d465a04042403 %}`, bindings)
	require.NoError(t, err)
	re := regexp.MustCompile(`<script.*>\s*</script>`)
	fmt.Println(s)
	require.Contains(t, s, `src=https://gist.github.com/parkr/931c1c8d465a04042403.js`)
	require.True(t, re.MatchString(s))
}
