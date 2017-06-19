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
	ConfigFile            *string
	Source                string
	Destination           string
	UseRemoteLiquidEngine bool

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
func (site *Site) Reload() error {
	copy, err := NewSiteFromDirectory(site.Source)
	if err != nil {
		return err
	}
	copy.Destination = site.Destination
	*site = *copy
	site.sassTempDir = ""
	return nil
}

// Load loads the site data and files. It doesn't load the configuration file; NewSiteFromDirectory did that.
func (site *Site) Load() (err error) {
	err = site.ReadFiles()
	if err != nil {
		return
	}
	site.initSiteVariables()
	site.liquidEngine, err = site.makeLiquidEngine()
	return
}

// KeepFile returns a boolean indicating that clean should leave the file in the destination directory.
func (site *Site) KeepFile(path string) bool {
	// TODO
	return false
}

// FindPageByFilePath returns a Page or nil, referenced by relative path.
func (site *Site) FindPageByFilePath(relpath string) Page {
	for _, p := range site.Paths {
		if p.Path() == relpath {
			return p
		}
	}
	return nil
}

// GetFileURL returns the URL path given a file path, relative to the site source directory.
func (site *Site) GetFileURL(relpath string) (string, bool) {
	for _, p := range site.Paths {
		if p.Path() == relpath {
			return p.Permalink(), true
		}
	}
	return "", false
}

// PageForURL returns the page that will be served at URL
func (site *Site) PageForURL(urlpath string) (p Page, found bool) {
	p, found = site.Paths[urlpath]
	if !found {
		p, found = site.Paths[filepath.Join(urlpath, "index.html")]
	}
	if !found {
		p, found = site.Paths[filepath.Join(urlpath, "index.htm")]
	}
	return
}

// Exclude returns a boolean indicating that the site excludes a file.
func (site *Site) Exclude(path string) bool {
	// TODO exclude based on glob, not exact match
	inclusionMap := helpers.StringArrayToMap(site.config.Include)
	exclusionMap := helpers.StringArrayToMap(site.config.Exclude)
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
func (site *Site) LayoutsDir() string {
	return filepath.Join(site.Source, site.config.LayoutsDir)
}

// ReadFiles scans the source directory and creates pages and collections.
func (site *Site) ReadFiles() error {
	site.Paths = make(map[string]Page)

	walkFn := func(name string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relname, err := filepath.Rel(site.Source, name)
		if err != nil {
			panic(err)
		}
		switch {
		case info.IsDir() && site.Exclude(relname):
			return filepath.SkipDir
		case info.IsDir(), site.Exclude(relname):
			return nil
		}
		defaults := site.GetFrontMatterDefaults(relname, "")
		p, err := ReadPage(site, nil, relname, defaults)
		if err != nil {
			return helpers.PathError(err, "read", name)
		}
		if p.Published() {
			site.Paths[p.Permalink()] = p
		}
		return nil
	}

	if err := filepath.Walk(site.Source, walkFn); err != nil {
		return err
	}
	return site.ReadCollections()
}

func (site *Site) initSiteVariables() {
	site.Variables = MergeVariableMaps(site.Variables, VariableMap{
		// TODO read time from _config, if it's available
		"time": time.Now(),
		// TODO pages, posts, related_posts, static_files, html_pages, html_files, collections, data, documents, categories.CATEGORY, tags.TAG
	})
	for _, c := range site.Collections {
		site.Variables[c.Name] = c.PageTemplateObjects()
	}
}

func (site *Site) makeLocalLiquidEngine() liquid.Engine {
	engine := liquid.NewLocalWrapperEngine()
	engine.LinkTagHandler(site.GetFileURL)
	includeHandler := func(name string, w io.Writer, scope map[string]interface{}) error {
		filename := filepath.Join(site.Source, site.config.IncludesDir, name)
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

func (site *Site) makeLiquidClient() (engine liquid.RemoteEngine, err error) {
	engine, err = liquid.NewRPCClientEngine(liquid.DefaultServer)
	if err != nil {
		return
	}
	urls := map[string]string{}
	for _, p := range site.Paths {
		urls[p.Path()] = p.Permalink()
	}
	err = engine.FileURLMap(urls)
	if err != nil {
		return
	}
	err = engine.IncludeDirs([]string{filepath.Join(site.Source, site.config.IncludesDir)})
	return
}

func (site *Site) makeLiquidEngine() (liquid.Engine, error) {
	if site.UseRemoteLiquidEngine {
		return site.makeLiquidClient()
	}
	return site.makeLocalLiquidEngine(), nil
}

// LiquidEngine create a liquid engine configured to with include paths and link tag resolution
// for this site.
func (site *Site) LiquidEngine() liquid.Engine {
	return site.liquidEngine
}

// GetFrontMatterDefaults implements https://jekyllrb.com/docs/configuration/#front-matter-defaults
func (site *Site) GetFrontMatterDefaults(relpath, typename string) (m VariableMap) {
	for _, entry := range site.config.Defaults {
		scope := &entry.Scope
		hasPrefix := strings.HasPrefix(relpath, scope.Path)
		hasType := scope.Type == "" || scope.Type == typename
		if hasPrefix && hasType {
			m = MergeVariableMaps(m, entry.Values)
		}
	}
	return
}
