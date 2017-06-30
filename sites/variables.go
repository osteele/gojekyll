package sites

import (
	"time"

	"github.com/osteele/gojekyll/templates"
)

// SiteVariables returns the site variable for template evaluation.
func (s *Site) SiteVariables() templates.VariableMap {
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
		// TODO pages, posts, related_posts, static_files, html_pages, html_files, collections, data, documents, categories.CATEGORY, tags.TAG
	})
	return s.setCollectionVariables(false)
}

func normalizeMaps(value interface{}) interface{} {
	return value
}

func (s *Site) setCollectionVariables(includeContent bool) error {
	for _, c := range s.Collections {
		v, err := c.TemplateVariable(s, includeContent)
		if err != nil {
			return err
		}
		s.siteVariables[c.Name] = v
	}
	return nil
}

func (s *Site) setCollectionContent() error {
	return s.setCollectionVariables(true)
}
