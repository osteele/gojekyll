package sites

import (
	"io"
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
	Paths       map[string]pages.Page // URL path -> Page

	config       config.Config
	liquidEngine liquid.Engine
	sassTempDir  string
}

// SourceDir returns the sites source directory.
func (s *Site) SourceDir() string { return s.Source }

// DefaultPermalink returns the default Permalink for pages not in a collection.
func (s *Site) DefaultPermalink() string { return "/:path:output_ext" }

// Output returns true, indicating that the files in the site should be written.
func (s *Site) Output() bool { return true }

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
