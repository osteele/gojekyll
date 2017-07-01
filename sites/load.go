package sites

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/osteele/gojekyll/collections"
	"github.com/osteele/gojekyll/config"
	"github.com/osteele/gojekyll/helpers"
	"github.com/osteele/gojekyll/pages"
)

// NewSiteFromDirectory reads the configuration file, if it exists.
func NewSiteFromDirectory(source string, flags config.Flags) (*Site, error) {
	s := NewSite(flags)
	configPath := filepath.Join(source, "_config.yml")
	bytes, err := ioutil.ReadFile(configPath)
	switch {
	case err != nil && os.IsNotExist(err):
		// ok
	case err != nil:
		return nil, err
	default:
		err = config.Unmarshal(bytes, &s.config)
		if err != nil {
			return nil, err
		}
		s.ConfigFile = &configPath
	}
	s.config.Source = source
	s.config.ApplyFlags(s.flags)
	return s, nil
}

// Load loads the site data and files. It doesn't load the configuration file; NewSiteFromDirectory did that.
func (s *Site) Load() error {
	if err := s.readFiles(); err != nil {
		return err
	}
	return s.readDataFiles()
}

// Reload reloads the config file and pages.
// If there's an error loading the config file, it has no effect.
func (s *Site) Reload() error {
	copy, err := NewSiteFromDirectory(s.SourceDir(), s.flags)
	if err != nil {
		return err
	}
	*s = *copy
	return s.Load()
}

// readFiles scans the source directory and creates pages and collections.
func (s *Site) readFiles() error {
	s.Routes = make(map[string]pages.Document)

	walkFn := func(filename string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relname, err := filepath.Rel(s.SourceDir(), filename)
		if err != nil {
			panic(err)
		}
		switch {
		case info.IsDir() && s.Exclude(relname):
			return filepath.SkipDir
		case info.IsDir(), s.Exclude(relname):
			return nil
		}
		defaultFrontmatter := s.config.GetFrontMatterDefaults(relname, "")
		p, err := pages.NewFile(filename, s, filepath.ToSlash(relname), defaultFrontmatter)
		if err != nil {
			return helpers.PathError(err, "read", filename)
		}
		s.AddDocument(p, true)
		return nil
	}

	if err := filepath.Walk(s.SourceDir(), walkFn); err != nil {
		return err
	}
	return s.ReadCollections()
}

// AddDocument adds a page to the site structures.
func (s *Site) AddDocument(p pages.Document, output bool) {
	if p.Published() || s.config.Unpublished {
		s.pages = append(s.pages, p)
		if output {
			s.Routes[p.Permalink()] = p
		}
	}
}

// ReadCollections reads the pages of the collections named in the site configuration.
// It adds each collection's pages to the site map, and creates a template site variable for each collection.
func (s *Site) ReadCollections() error {
	for name, data := range s.config.Collections {
		c := collections.NewCollection(name, data, s)
		s.Collections = append(s.Collections, c)
		if err := c.ReadPages(s.SourceDir(), s.config.GetFrontMatterDefaults); err != nil {
			return err
		}
		for _, p := range c.Pages() {
			s.AddDocument(p, c.Output())
		}
	}
	return nil
}
