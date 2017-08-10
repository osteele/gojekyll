package site

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/osteele/gojekyll/collection"
	"github.com/osteele/gojekyll/config"
	"github.com/osteele/gojekyll/pages"
	"github.com/osteele/gojekyll/plugins"
	"github.com/osteele/gojekyll/utils"
)

// FromDirectory reads the configuration file, if it exists.
func FromDirectory(dir string, flags config.Flags) (*Site, error) {
	s := New(flags)
	if err := s.config.FromDirectory(dir); err != nil {
		return nil, err
	}
	s.config.ApplyFlags(s.flags)
	return s, nil
}

// Read loads the site data and files.
func (s *Site) Read() error {
	s.Routes = make(map[string]pages.Document)
	plugins.Install(s.config.Plugins, s)
	if err := s.findTheme(); err != nil {
		return err
	}
	if err := s.readDataFiles(); err != nil {
		return err
	}
	if err := s.readThemeAssets(); err != nil {
		return err
	}
	if err := s.readFiles(s.SourceDir(), s.SourceDir()); err != nil {
		return err
	}
	if err := s.ReadCollections(); err != nil {
		return err
	}
	if err := s.initializeRenderingPipeline(); err != nil {
		return err
	}
	return s.runHooks(func(p plugins.Plugin) error { return p.PostRead(s) })
}

// readFiles scans the source directory and creates pages and collection.
func (s *Site) readFiles(dir, base string) error {
	return filepath.Walk(dir, func(filename string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel := utils.MustRel(base, filename)
		switch {
		case info.IsDir() && s.Exclude(rel):
			return filepath.SkipDir
		case info.IsDir():
			return nil
		case s.Exclude(rel):
			return nil
		case strings.HasPrefix(rel, "_"):
			return nil
		}
		defaultFrontmatter := s.config.GetFrontMatterDefaults("", rel)
		d, err := pages.NewFile(s, filename, filepath.ToSlash(rel), defaultFrontmatter)
		if err != nil {
			return utils.WrapPathError(err, filename)
		}
		s.AddDocument(d, true)
		if p, ok := d.(pages.Page); ok {
			s.nonCollectionPages = append(s.nonCollectionPages, p)
		}
		return nil
	})
}

// AddDocument adds a document to the site's fields.
// It ignores unpublished documents unless config.Unpublished is true.
func (s *Site) AddDocument(d pages.Document, output bool) {
	if d.Published() || s.config.Unpublished {
		s.docs = append(s.docs, d)
		if output {
			s.Routes[d.Permalink()] = d
		}
	}
}

// ReadCollections reads the pages of the collections named in the site configuration.
// It adds each collection's pages to the site map, and creates a template site variable for each collection.
func (s *Site) ReadCollections() error {
	for name, data := range s.config.Collections {
		c := collection.New(s, name, data)
		s.Collections = append(s.Collections, c)
		if err := c.ReadPages(); err != nil {
			return err
		}
		for _, p := range c.Pages() {
			s.AddDocument(p, c.Output())
		}
	}
	return nil
}
