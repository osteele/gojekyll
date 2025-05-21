package plugins

import (
	"fmt"
	"strings"
)

type paginatePlugin struct{ plugin }

func init() {
	register("jekyll-paginate", &paginatePlugin{})
}

func (p *paginatePlugin) ModifyRenderContext(s Site, m map[string]interface{}) error {
	m["paginator"] = createPaginator(1, 5, s.Posts())
	return nil
}

func (p *paginatePlugin) PostReadSite(s Site) error {
	// s.AddDocument(newTemplateDoc(s, "/sitemap.xml", sitemapTemplateSource), true)
	// s.AddDocument(newTemplateDoc(s, "/robots.txt", `Sitemap: {{ "sitemap.xml" | absolute_url }}`), true)
	// iterate over pagesets
	return nil
}

func createPaginator(n, perPage int, posts []Page) map[string]interface{} {
	pageCount := (len(posts) + perPage - 1) / perPage
	paginatePath := "/blog/page:num/"
	pagePath := func(n int) string {
		return strings.ReplaceAll(paginatePath, ":num", fmt.Sprint(n))
	}
	m := map[string]interface{}{
		"page":               n,
		"per_page":           perPage,
		"posts":              posts,
		"total_posts":        len(posts),
		"total_pages":        pageCount,
		"previous_page":      nil,
		"previous_page_path": nil,
		"next_page":          nil,
		"next_page_path":     nil,
	}
	if n > 1 {
		m["previous_page"] = n - 1
		m["previous_page_path"] = pagePath(n - 1)
	}
	if n < pageCount {
		m["next_page"] = n + 1
		m["next_page_path"] = pagePath(n + 1)
	}
	return m
}
