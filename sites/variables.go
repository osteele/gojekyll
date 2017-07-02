package sites

import (
	"fmt"
	"time"

	"github.com/osteele/gojekyll/templates"
	"github.com/osteele/liquid/generics"
)

// SiteVariables returns the site variable for template evaluation.
func (s *Site) SiteVariables() map[string]interface{} {
	if len(s.siteVariables) == 0 {
		if err := s.initializeSiteVariables(); err != nil {
			panic(err)
		}
	}
	return s.siteVariables
}

func (s *Site) initializeSiteVariables() error {
	s.siteVariables = templates.MergeVariableMaps(s.config.Variables, map[string]interface{}{
		"data": s.data,
		// "collections": s.computeCollections(), // generics.MustConvert(s.config.Collections, reflect.TypeOf([]interface{}{})),
		// TODO read time from _config, if it's available
		"time": time.Now(),
		// TODO pages, static_files, html_pages, html_files, documents, tags.TAG
	})
	return s.setCollectionVariables(false)
}

// set site[collection.name] for each collection.
func (s *Site) setCollectionVariables(includeContent bool) error {
	cv := []interface{}{}
	for _, c := range s.Collections {
		pages, err := c.TemplateVariable(s, includeContent)
		if err != nil {
			return err
		}
		cv = append(cv, c.TemplateObject(pages))
		s.siteVariables[c.Name] = pages
		if c.IsPostsCollection() {
			s.setPostVariables(pages)
		}
	}
	generics.SortByProperty(cv, "label", true)
	s.siteVariables["collections"] = cv
	return nil
}

func (s *Site) setPostVariables(pages []interface{}) {
	var (
		related    = pages
		categories = map[string][]interface{}{}
		tags       = map[string][]interface{}{}
	)
	if len(related) > 10 {
		related = related[:10]
	}
	for _, p := range pages {
		b := p.(map[string]interface{})
		switch cs := b["categories"].(type) {
		case []interface{}:
			for _, c := range cs {
				key := fmt.Sprint(c)
				ps, found := categories[key]
				if !found {
					ps = []interface{}{}
				}
				categories[key] = append(ps, p)
			}
		}
	}
	s.siteVariables["categories"] = categories
	s.siteVariables["tags"] = tags
	s.siteVariables["related_posts"] = related
}

func (s *Site) setCollectionContent() error {
	return s.setCollectionVariables(true)
}
