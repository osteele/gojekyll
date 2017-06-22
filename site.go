package gojekyll

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

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

	Collections []*Collection
	Variables   templates.VariableMap
	Paths       map[string]pages.Page // URL path -> Page

	config       SiteConfig
	liquidEngine liquid.Engine
	sassTempDir  string
}

// SourceDir returns the sites source directory.
func (s *Site) SourceDir() string { return s.Source }

// Output returns true, indicating that the files in the site should be written.
func (s *Site) Output() bool { return true }

// PathPrefix returns the relative directory prefix.
func (s *Site) PathPrefix() string { return "" }

// NewSite creates a new site record, initialized with the site defaults.
func NewSite() *Site {
	s := new(Site)
	if err := s.readConfigBytes([]byte(defaultSiteConfig)); err != nil {
		panic(err)
	}
	return s
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
		if err = s.readConfigBytes(bytes); err != nil {
			return nil, err
		}
		s.Source = filepath.Join(source, s.config.Source)
		s.ConfigFile = &configPath
	}
	s.Destination = filepath.Join(s.Source, s.config.Destination)
	return s, nil
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
	return nil
}

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
		defaults := s.GetFrontMatterDefaults(relname, "")
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
		c := NewCollection(s, name, data)
		s.Collections = append(s.Collections, c)
		if err := c.ReadPages(); err != nil {
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

func (s *Site) makeLocalLiquidEngine() liquid.Engine {
	engine := liquid.NewLocalWrapperEngine()
	engine.LinkTagHandler(s.RelPathURL)
	includeHandler := func(name string, w io.Writer, scope map[string]interface{}) error {
		filename := filepath.Join(s.Source, s.config.IncludesDir, name)
		template, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}
		text, err := engine.ParseAndRender(template, scope)
		if err != nil {
			return err
		}
		_, err = w.Write(text)
		return err
	}
	engine.IncludeHandler(includeHandler)
	return engine
}

func (s *Site) makeLiquidClient() (engine liquid.RemoteEngine, err error) {
	engine, err = liquid.NewRPCClientEngine(liquid.DefaultServer)
	if err != nil {
		return
	}
	urls := map[string]string{}
	for _, p := range s.Paths {
		urls[p.SiteRelPath()] = p.Permalink()
	}
	err = engine.FileURLMap(urls)
	if err != nil {
		return
	}
	err = engine.IncludeDirs([]string{filepath.Join(s.Source, s.config.IncludesDir)})
	return
}

func (s *Site) makeLiquidEngine() (liquid.Engine, error) {
	if s.UseRemoteLiquidEngine {
		return s.makeLiquidClient()
	}
	return s.makeLocalLiquidEngine(), nil
}

// TemplateEngine create a liquid engine configured to with include paths and link tag resolution
// for this site.
func (s *Site) TemplateEngine() liquid.Engine {
	return s.liquidEngine
}

// GetFrontMatterDefaults implements https://jekyllrb.com/docs/configuration/#front-matter-defaults
func (s *Site) GetFrontMatterDefaults(relpath, typename string) (m templates.VariableMap) {
	for _, entry := range s.config.Defaults {
		scope := &entry.Scope
		hasPrefix := strings.HasPrefix(relpath, scope.Path)
		hasType := scope.Type == "" || scope.Type == typename
		if hasPrefix && hasType {
			m = templates.MergeVariableMaps(m, entry.Values)
		}
	}
	return
}
