package site

import (
	"fmt"
	"time"

	"github.com/osteele/gojekyll/pages"
	"github.com/osteele/gojekyll/plugins"
	"github.com/osteele/gojekyll/templates"
	"github.com/osteele/liquid"
)

// ToLiquid returns the site variable for template evaluation.
func (s *Site) ToLiquid() interface{} {
	s.dropOnce.Do(func() {
		s.dropErr = s.initializeDrop()
	})
	if s.dropErr != nil {
		panic(fmt.Sprintf("site drop initialization failed: %s", s.dropErr))
	}
	return liquid.IterationKeyedMap(s.drop)
}

func (s *Site) initializeDrop() error {
	var docs []Page
	for _, c := range s.Collections {
		docs = append(docs, c.Pages()...)
	}
	drop := templates.MergeVariableMaps(s.cfg.Variables(), map[string]interface{}{
		"collections":  s.collectionDrops(),
		"data":         s.data,
		"documents":    docs,
		"html_files":   s.htmlFiles(),
		"html_pages":   s.htmlPages(),
		"pages":        s.nonCollectionPages,
		"static_files": s.staticFiles(),
		// TODO read time from _config, if it's available
		"time": time.Now(),
	})
	for _, c := range s.Collections {
		drop[c.Name] = c.Pages()
	}
	s.drop = drop
	s.setPostVariables()
	return s.runHooks(func(h plugins.Plugin) error {
		return h.ModifySiteDrop(s, drop)
	})
}

// The following functions are only used in the drop.
//
// Since the drop is cached, there's no effort to cache these too.

func (s *Site) collectionDrops() []interface{} {
	var drops []interface{}
	for _, c := range s.Collections {
		drops = append(drops, c.ToLiquid())
	}
	return drops
}

func (s *Site) htmlFiles() (result []*pages.StaticFile) {
	for _, p := range s.staticFiles() {
		if p.OutputExt() == ".html" {
			result = append(result, p)
		}
	}
	return
}

func (s *Site) htmlPages() (result []Page) {
	for _, p := range s.nonCollectionPages {
		if p.OutputExt() == ".html" {
			result = append(result, p)
		}
	}
	return
}

func (s *Site) staticFiles() (result []*pages.StaticFile) {
	for _, d := range s.docs {
		if sd, ok := d.(*pages.StaticFile); ok {
			result = append(result, sd)
		}
	}
	return
}
