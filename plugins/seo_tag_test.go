package plugins

import (
	"testing"

	"github.com/osteele/gojekyll/config"
	"github.com/osteele/gojekyll/filters"
	"github.com/osteele/liquid"
	"github.com/osteele/liquid/tags"
	"github.com/stretchr/testify/require"
)

func TestSEOTag(t *testing.T) {
	engine := liquid.NewEngine()
	cfg := config.Default()
	cfg.BaseURL = "/my-baseurl"
	cfg.AbsoluteURL = "http://example.com"
	filters.AddJekyllFilters(engine, &cfg)
	plugins := []string{"jekyll-seo-tag"}
	_ = Install(plugins, siteFake{config.Default(), engine})
	require.NoError(t, directory[plugins[0]].ConfigureTemplateEngine(engine))
	bindings := liquid.Bindings{
		"site": tags.IterationKeyedMap{
			"title": "page title",
			"url":   "http://example.com/",
		},
		"page": tags.IterationKeyedMap{
			"title": "site title",
			"url":   "page",
		},
	}
	s, err := engine.ParseAndRenderString(`{% seo %}`, bindings)
	require.NoError(t, err)
	require.Contains(t, s, `<title>site title | page title</title>`)
}
