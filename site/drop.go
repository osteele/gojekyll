package site

import (
	"log"
	"time"

	"github.com/osteele/gojekyll/pages"
	"github.com/osteele/gojekyll/plugins"
	"github.com/osteele/gojekyll/templates"
	"github.com/osteele/liquid/evaluator"
)

// ToLiquid returns the site variable for template evaluation.
func (s *Site) ToLiquid() interface{} {
	s.Do(func() {
		if err := s.initializeDrop(); err != nil {
			log.Fatalf("ToLiquid failed: %s\n", err)
		}
	})
	return s.drop
}

// MarshalYAML is part of the yaml.Marshaler interface
// The variables subcommand uses this.
func (s *Site) MarshalYAML() (interface{}, error) {
	return s.ToLiquid(), nil
}

func (s *Site) initializeDrop() error {
	drop := templates.MergeVariableMaps(s.config.Variables, map[string]interface{}{
		"data":         s.data,
		"documents":    s.docs,
		"html_files":   s.htmlFiles(),
		"html_pages":   s.htmlPages(),
		"pages":        s.Pages(),
		"static_files": s.staticFiles(),
		// TODO read time from _config, if it's available
		"time": time.Now(),
	})
	collections := []interface{}{}
	for _, c := range s.Collections {
		drop[c.Name] = c.Pages()
		collections = append(collections, c.ToLiquid())
	}
	evaluator.SortByProperty(collections, "label", true)
	drop["collections"] = collections
	s.drop = drop
	s.setPostVariables()
	return s.runHooks(func(h plugins.Plugin) error {
		return h.ModifySiteDrop(s, drop)
	})
}

func (s *Site) setPageContent() error {
	for _, c := range s.Collections {
		if err := c.SetPageContent(); err != nil {
			return err
		}
	}
	return nil
}

// The following functions are only used in the drop, therefore they're
// non-public and they're listed here.
//
// Since the drop is cached, there's no effort to cache these too.

func (s *Site) htmlFiles() (out []*pages.StaticFile) {
	for _, p := range s.staticFiles() {
		if p.OutputExt() == ".html" {
			out = append(out, p)
		}
	}
	return
}

func (s *Site) htmlPages() (out []pages.Page) {
	for _, p := range s.Pages() {
		if p.OutputExt() == ".html" {
			out = append(out, p)
		}
	}
	return
}

func (s *Site) staticFiles() (out []*pages.StaticFile) {
	for _, d := range s.docs {
		if sd, ok := d.(*pages.StaticFile); ok {
			out = append(out, sd)
		}
	}
	return
}
