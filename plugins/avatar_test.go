package plugins

import (
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
