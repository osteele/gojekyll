package sites

import (
	"time"

	"github.com/osteele/gojekyll/templates"
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
	s.siteVariables = templates.MergeVariableMaps(s.config.Variables, templates.VariableMap{
		"data": s.data,
		// TODO read time from _config, if it's available
		"time": time.Now(),
		// TODO pages, static_files, html_pages, html_files, documents, categories.CATEGORY, tags.TAG
	})
	return s.setCollectionVariables(false)
}

func (s *Site) setCollectionVariables(includeContent bool) error {
	for _, c := range s.Collections {
		pages, err := c.TemplateVariable(s, includeContent)
		if err != nil {
			return err
		}
		s.siteVariables[c.Name] = pages
		if c.IsPostsCollection() {
			related := pages
			if len(related) > 10 {
				related = related[:10]
			}
			s.siteVariables["related_posts"] = related
		}
	}
	// Set these here instead of initializeSiteVariables so that they're
	// re-generated once page.content has been rendered. Obviously
	// this method has the wrong name.
	return nil
}

func (s *Site) setCollectionContent() error {
	return s.setCollectionVariables(true)
}
