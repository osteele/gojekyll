package sites

import (
	"os"
	"path/filepath"

	"github.com/osteele/gojekyll/collections"
	"github.com/osteele/gojekyll/helpers"
	"github.com/osteele/gojekyll/pages"
)

// Load loads the site data and files. It doesn't load the configuration file; NewSiteFromDirectory did that.
func (s *Site) Load() (err error) {
	err = s.readFiles()
	if err != nil {
		return
	}
	err = s.initSiteVariables()
	if err != nil {
		return
	}
	s.liquidEngine, err = s.makeLiquidEngine()
	return
}

// Reload reloads the config file and pages.
// If there's an error loading the config file, it has no effect.
func (s *Site) Reload() error {
	copy, err := NewSiteFromDirectory(s.Source)
	if err != nil {
		return err
	}
	copy.Destination = s.Destination
	*s = *copy
	s.sassTempDir = ""
	return s.Load()
}

// readFiles scans the source directory and creates pages and collections.
func (s *Site) readFiles() error {
	s.Paths = make(map[string]pages.Page)

	walkFn := func(filename string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relname, err := filepath.Rel(s.Source, filename)
		if err != nil {
			panic(err)
		}
		switch {
		case info.IsDir() && s.Exclude(relname):
			return filepath.SkipDir
		case info.IsDir(), s.Exclude(relname):
			return nil
		}
		defaults := s.config.GetFrontMatterDefaults(relname, "")
		p, err := pages.NewPageFromFile(s, s, filename, relname, defaults)
		if err != nil {
			return helpers.PathError(err, "read", filename)
		}
		if p.Published() {
			s.Paths[p.Permalink()] = p
		}
		return nil
	}

	if err := filepath.Walk(s.Source, walkFn); err != nil {
		return err
	}
	return s.ReadCollections()
}

// ReadCollections reads the pages of the collections named in the site configuration.
// It adds each collection's pages to the site map, and creates a template site variable for each collection.
func (s *Site) ReadCollections() error {
	for name, data := range s.config.Collections {
		c := collections.NewCollection(s, name, data)
		s.Collections = append(s.Collections, c)
		if err := c.ReadPages(s, s.Source, s.config.GetFrontMatterDefaults); err != nil {
			return err
		}
		for _, p := range c.Pages() {
			if p.Published() {
				s.Paths[p.Permalink()] = p
			}
		}
	}
	return nil
}
