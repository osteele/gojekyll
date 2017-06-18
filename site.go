package gojekyll

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/osteele/gojekyll/helpers"
	"github.com/osteele/gojekyll/liquid"
)

// Site is a Jekyll site.
type Site struct {
	ConfigFile  *string
	Source      string
	Destination string

	Collections []*Collection
	Variables   VariableMap
	Paths       map[string]Page // URL path -> Page

	config       SiteConfig
	liquidEngine liquid.Engine
	sassTempDir  string
}

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
	return s.ReadFiles()
}

// KeepFile returns a boolean indicating that clean should leave the file in the destination directory.
func (s *Site) KeepFile(path string) bool {
	// TODO
	return false
}

// FindPageByFilePath returns a Page or nil, referenced by relative path.
func (s *Site) FindPageByFilePath(relpath string) Page {
	for _, p := range s.Paths {
		if p.Path() == relpath {
			return p
		}
	}
	return nil
}

// GetFileURL returns the URL path given a file path, relative to the site source directory.
func (s *Site) GetFileURL(relpath string) (string, bool) {
	for _, p := range s.Paths {
		if p.Path() == relpath {
			return p.Permalink(), true
		}
	}
	return "", false
}

// PageForURL returns the page that will be served at URL
func (s *Site) PageForURL(urlpath string) (p Page, found bool) {
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

// ReadFiles scans the source directory and creates pages and collections.
func (s *Site) ReadFiles() error {
	s.Paths = make(map[string]Page)

	walkFn := func(name string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relname, err := filepath.Rel(s.Source, name)
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
		p, err := ReadPage(s, nil, relname, defaults)
		if err != nil {
			return helpers.PathError(err, "read", name)
		}
		if p.Published() {
			s.Paths[p.Permalink()] = p
		}
		return nil
	}

	if err := filepath.Walk(s.Source, walkFn); err != nil {
		return err
	}
	if err := s.ReadCollections(); err != nil {
		return err
	}
	s.initTemplateAttributes()
	return nil
}

func (s *Site) initTemplateAttributes() {
	// TODO site: {pages, posts, related_posts, static_files, html_pages, html_files, collections, data, documents, categories.CATEGORY, tags.TAG}
	s.Variables = MergeVariableMaps(s.Variables, VariableMap{
		"time": time.Now(),
	})
	for _, c := range s.Collections {
		s.Variables[c.Name] = c.PageTemplateObjects()
	}
}

func (s *Site) createLocalEngine() liquid.Engine {
	e := liquid.NewLocalWrapperEngine()
	e.LinkHandler(s.GetFileURL)
	includeHandler := func(name string, w io.Writer, scope map[string]interface{}) {
		name = strings.TrimLeft(strings.TrimRight(name, "}}"), "{{")
		filename := filepath.Join(s.Source, s.config.IncludesDir, name)
		template, err := ioutil.ReadFile(filename)
		if err != nil {
			panic(err)
		}
		text, err := e.ParseAndRender(template, scope)
		_, err = w.Write(text)
		if err != nil {
			panic(err)
		}
	}
	e.IncludeHandler(includeHandler)
	return e
}

func (s *Site) createRemoteEngine() liquid.Engine {
	e := liquid.NewRPCClientEngine(liquid.DefaultServer)
	m := map[string]string{}
	for _, p := range s.Paths {
		m[p.Path()] = p.Permalink()
	}
	e.FileUrlMap(m)
	e.IncludeDirs([]string{filepath.Join(s.Source, s.config.IncludesDir)})
	return e
}

// LiquidEngine create a liquid engine with site-specific behavior.
func (s *Site) LiquidEngine() liquid.Engine {
	if s.liquidEngine == nil {
		s.liquidEngine = s.createLocalEngine()
	}
	return s.liquidEngine
}

// GetFrontMatterDefaults implements https://jekyllrb.com/docs/configuration/#front-matter-defaults
func (s *Site) GetFrontMatterDefaults(relpath, typename string) (m VariableMap) {
	for _, entry := range s.config.Defaults {
		scope := &entry.Scope
		hasPrefix := strings.HasPrefix(relpath, scope.Path)
		hasType := scope.Type == "" || scope.Type == typename
		if hasPrefix && hasType {
			m = MergeVariableMaps(m, entry.Values)
		}
	}
	return
}
