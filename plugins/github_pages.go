package plugins

import "github.com/osteele/gojekyll/utils"

func init() {
	register("github-pages", githubPagesPlugin{})
}

type githubPagesPlugin struct{ plugin }

func (p githubPagesPlugin) ModifyPluginList(names []string) []string {
	for _, name := range githubPagesPlugins {
		if !utils.StringArrayContains(names, name) {
			names = append(names, name)
		}
	}
	return names
}

var githubPagesPlugins = []string{
	"jekyll-avatar",
	"jekyll-coffeescript",
	"jekyll-default-layout",
	"jekyll-feed",
	"jekyll-gist",
	"jekyll-github-metadata",
	"jekyll-mentions",
	"jekyll-optional-front-matter",
	"jekyll-paginate",
	"jekyll-readme-index",
	"jekyll-redirect-from",
	"jekyll-relative-links",
	"jekyll-sass-converter",
	"jekyll-seo-tag",
	"jekyll-sitemap",
	"jekyll-titles-from-headings",
	"jemoji",
}
