package site

import (
	"time"

	"github.com/osteele/gojekyll/templates"
	"github.com/osteele/liquid/generics"
)

// ToLiquid returns the site variable for template evaluation.
func (s *Site) ToLiquid() interface{} {
	if len(s.drop) == 0 {
		s.initializeDrop()
	}
	return s.drop
}

// MarshalYAML is part of the yaml.Marshaler interface
// The variables subcommand uses this.
func (s *Site) MarshalYAML() (interface{}, error) {
	return s.ToLiquid(), nil
}

func (s *Site) initializeDrop() {
	vars := templates.MergeVariableMaps(s.config.Variables, map[string]interface{}{
		"data": s.data,
		// "collections": s.computeCollections(), // generics.MustConvert(s.config.Collections, reflect.TypeOf([]interface{}{})),
		// TODO read time from _config, if it's available
		"time": time.Now(),
		// TODO pages, static_files, html_pages, html_files, documents, tags.TAG
	})
	collections := []interface{}{}
	for _, c := range s.Collections {
		vars[c.Name] = c.Pages()
		collections = append(collections, c.ToLiquid())
	}
	generics.SortByProperty(collections, "label", true)
	vars["collections"] = collections
	s.drop = vars
	s.setPostVariables()
}

func (s *Site) setPageContent() error {
	for _, c := range s.Collections {
		if err := c.SetPageContent(s); err != nil {
			return err
		}
	}
	return nil
}