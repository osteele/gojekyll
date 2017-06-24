package sites

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/osteele/gojekyll/collections"
	"github.com/osteele/gojekyll/config"
	"github.com/osteele/gojekyll/helpers"
	"github.com/osteele/gojekyll/liquid"
	"github.com/osteele/gojekyll/pages"
	"github.com/osteele/gojekyll/templates"
)

// Site is a Jekyll site.
type Site struct {
	ConfigFile            *string
	Source                string
	Destination           string
	UseRemoteLiquidEngine bool

	Collections []*collections.Collection
	Variables   templates.VariableMap
	Paths       map[string]pages.Page // URL path -> Page, only for output pages

	config       config.Config
	liquidEngine liquid.Engine
	pages        []pages.Page // all pages, output or not
	sassTempDir  string
}

// Pages returns a slice of pages.
func (s *Site) Pages() []pages.Page { return s.pages }

// PathPrefix returns the relative directory prefix.
func (s *Site) PathPrefix() string { return "" }

// NewSite creates a new site record, initialized with the site defaults.
func NewSite() *Site {
	return &Site{config: config.Default()}
}

// NewSiteFromDirectory reads the configuration file, if it exists.
func NewSiteFromDirectory(source string) (*Site, error) {
	s := NewSite()
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
		s.Source = filepath.Join(source, s.config.Source)
		s.ConfigFile = &configPath
	}
	s.Destination = filepath.Join(s.Source, s.config.Destination)
	return s, nil
}

// KeepFile returns a boolean indicating that clean should leave the file in the destination directory.
func (s *Site) KeepFile(path string) bool {
	// TODO
	return false
}

// RelPathPage returns a Page, give a file path relative to site source directory.
func (s *Site) RelPathPage(relpath string) (pages.Page, bool) {
	for _, p := range s.Paths {
		if p.SiteRelPath() == relpath {
			return p, true
		}
	}
	return nil, false
}

// RelPathURL returns a page's relative URL, give a file path relative to the site source directory.
func (s *Site) RelPathURL(relpath string) (string, bool) {
	var url string
	p, found := s.RelPathPage(relpath)
	if found {
		url = p.Permalink()
	}
	return url, found
}

// URLPage returns the page that will be served at URL
func (s *Site) URLPage(urlpath string) (p pages.Page, found bool) {
	p, found = s.Paths[urlpath]
	if !found {
		p, found = s.Paths[filepath.Join(urlpath, "index.html")]
	}
	if !found {
		p, found = s.Paths[filepath.Join(urlpath, "index.htm")]
	}
	return
}

// Exclude returns a boolean indicating that the site excludes a file.
func (s *Site) Exclude(path string) bool {
	// TODO exclude based on glob, not exact match
	inclusionMap := helpers.StringArrayToMap(s.config.Include)
	exclusionMap := helpers.StringArrayToMap(s.config.Exclude)
	base := filepath.Base(path)
	switch {
	case inclusionMap[path]:
		return false
	case path == ".":
		return false
	case exclusionMap[path]:
		return true
	case strings.HasPrefix(base, "."), strings.HasPrefix(base, "_"):
		return true
	default:
		return false
	}
}

// LayoutsDir returns the path to the layouts directory.
func (s *Site) LayoutsDir() string {
	return filepath.Join(s.Source, s.config.LayoutsDir)
}
