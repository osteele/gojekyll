package site

import (
	"github.com/osteele/gojekyll/collection"
	"github.com/osteele/gojekyll/pages"
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
		ps         = c.Pages()
		related    = ps
		categories = map[string][]pages.Page{}
		tags       = map[string][]pages.Page{}
	)
	if len(related) > 10 {
		related = related[:10]
	}
	for _, p := range ps {
		for _, k := range p.Categories() {
			ps, found := categories[k]
			if !found {
				ps = []pages.Page{}
			}
			categories[k] = append(ps, p)
		}
	}
	s.drop["categories"] = categories
	s.drop["tags"] = tags
	s.drop["related_posts"] = related
}
