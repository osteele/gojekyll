package site

import (
	"os"
	"path/filepath"
	"sort"
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
	if err := s.cfg.FromDirectory(dir); err != nil {
		return nil, utils.WrapError(err, "reading site")
	}
	s.cfg.ApplyFlags(s.flags)
	return s, nil
}

// Read loads the site data and files.
func (s *Site) Read() error {
	if err := s.installPlugins(); err != nil {
		return utils.WrapError(err, "initializing plugins")
	}
	s.Routes = make(map[string]Document)
	if err := s.findTheme(); err != nil {
		return utils.WrapError(err, "finding theme")
	}
	if err := s.readDataFiles(); err != nil {
		return utils.WrapError(err, "reading data files")
	}
	if err := s.readThemeAssets(); err != nil {
		return utils.WrapError(err, "reading theme assets")
	}
	if err := s.readFiles(s.SourceDir(), s.SourceDir()); err != nil {
		return utils.WrapError(err, "reading files")
	}
	if err := s.ReadCollections(); err != nil {
		return utils.WrapError(err, "reading collections")
	}
	if err := s.initializeRenderers(); err != nil {
		return utils.WrapError(err, "initializing renderers")
	}
	for _, p := range s.Pages() {
		err := s.runHooks(func(h plugins.Plugin) error {
			return h.PostInitPage(s, p)
		})
		if err != nil {
			return err
		}
	}
	return s.runHooks(func(p plugins.Plugin) error { return p.PostReadSite(s) })
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
		defaultFrontmatter := s.cfg.GetFrontMatterDefaults("", rel)
		d, err := pages.NewFile(s, filename, filepath.ToSlash(rel), defaultFrontmatter)
		if err != nil {
			return utils.WrapPathError(err, filename)
		}
		s.AddDocument(d, true)
		if p, ok := d.(Page); ok {
			s.nonCollectionPages = append(s.nonCollectionPages, p)
		}
		return nil
	})
}

// AddDocument adds a document to the site's fields.
// It ignores unpublished documents unless config.Unpublished is true.
func (s *Site) AddDocument(d Document, output bool) {
	if d.Published() || s.cfg.Unpublished {
		s.docs = append(s.docs, d)
		if output {
			s.Routes[d.URL()] = d
		}
	}
}

// ReadCollections reads the pages of the collections named in the site configuration.
// It adds each collection's pages to the site map, and creates a template site variable for each collection.
func (s *Site) ReadCollections() (err error) {
	var cols []*collection.Collection
	for name, data := range s.cfg.Collections {
		c := collection.New(s, name, data)
		cols = append(cols, c)
		err = c.ReadPages()
		if err != nil {
			break
		}
		for _, p := range c.Pages() {
			s.AddDocument(p, c.Output())
		}
	}
	sort.Slice(cols, func(i, j int) bool {
		return cols[i].Name < cols[j].Name
	})
	s.Collections = cols
	return nil
}
