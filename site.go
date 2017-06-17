package gojekyll

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/acstech/liquid"
	"github.com/acstech/liquid/core"
	"github.com/osteele/gojekyll/helpers"
	liquidHelper "github.com/osteele/gojekyll/liquid"
)

// Site is a Jekyll site.
type Site struct {
	ConfigFile  *string
	Source      string
	Destination string

	Collections []*Collection
	Variables   VariableMap
	Paths       map[string]Page // URL path -> Page

	config              SiteConfig
	liquidConfiguration *core.Configuration
	sassTempDir         string
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
func (s *Site) Reload() (err error) {
	copy, err := NewSiteFromDirectory(s.Source)
	copy.Destination = s.Destination
	*s = *copy
	err = s.ReadFiles()
	if err != nil {
		return
	}
	return
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
func (s *Site) GetFileURL(path string) (string, bool) {
	for _, p := range s.Paths {
		if p.Path() == path {
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
			return err
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

// LiquidConfiguration configures the liquid tags with site-specific behavior.
func (s *Site) LiquidConfiguration() *core.Configuration {
	if s.liquidConfiguration != nil {
		return s.liquidConfiguration
	}
	liquidHelper.SetFilePathURLGetter(s.GetFileURL)
	includeHandler := func(name string, writer io.Writer, data map[string]interface{}) {
		name = strings.TrimLeft(strings.TrimRight(name, "}}"), "{{")
		filename := path.Join(s.Source, s.config.IncludesDir, name)
		template, err := liquid.ParseFile(filename, s.liquidConfiguration)
		if err != nil {
			panic(err)
		}
		template.Render(writer, data)
	}
	s.liquidConfiguration = liquid.Configure().IncludeHandler(includeHandler)
	return s.liquidConfiguration
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
