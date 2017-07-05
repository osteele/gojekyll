package site

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/osteele/gojekyll/collection"
	"github.com/osteele/gojekyll/config"
	"github.com/osteele/gojekyll/helpers"
	"github.com/osteele/gojekyll/pages"
)

// FromDirectory reads the configuration file, if it exists.
func FromDirectory(source string, flags config.Flags) (*Site, error) {
	s := New(flags)
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
// It returns a copy.
// If there's an error loading the config file, it has no effect.
func (s *Site) Reload() (*Site, error) {
	copy, err := FromDirectory(s.SourceDir(), s.flags)
	if err != nil {
		return nil, err
	}
	return copy, copy.Load()
}

// readFiles scans the source directory and creates pages and collection.
func (s *Site) readFiles() error {
	s.Routes = make(map[string]pages.Document)

	walkFn := func(filename string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relname := helpers.MustRel(s.SourceDir(), filename)
		switch {
		case info.IsDir() && s.Exclude(relname):
			return filepath.SkipDir
		case info.IsDir(), s.Exclude(relname):
			return nil
		}
		defaultFrontmatter := s.config.GetFrontMatterDefaults("", relname)
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
		s.docs = append(s.docs, p)
		if output {
			s.Routes[p.Permalink()] = p
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
