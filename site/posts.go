package site

import (
	"github.com/osteele/gojekyll/collection"
)

func (s *Site) findPostCollection() *collection.Collection {
	for _, c := range s.Collections {
		if c.Name == "posts" {
			return c
		}
	}
	return nil
}

func (s *Site) setPostVariables() {
	c := s.findPostCollection()
	if c == nil {
		return
	}
	var (
		ps      = c.Pages()
		related = ps
	)
	if len(related) > 10 {
		related = related[:10]
	}
	s.drop["categories"] = s.groupPagesBy(func(p Page) []string { return p.Categories() })
	s.drop["tags"] = s.groupPagesBy(func(p Page) []string { return p.Tags() })
	s.drop["related_posts"] = related
}

func (s *Site) groupPagesBy(getter func(Page) []string) map[string][]Page {
	categories := map[string][]Page{}
	for _, p := range s.Pages() {
		for _, k := range p.Categories() {
			ps, found := categories[k]
			if !found {
				ps = []Page{}
			}
			categories[k] = append(ps, p)
		}
	}
	return categories
}
