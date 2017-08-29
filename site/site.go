package site

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/osteele/gojekyll/collection"
	"github.com/osteele/gojekyll/config"
	"github.com/osteele/gojekyll/pages"
	"github.com/osteele/gojekyll/plugins"
	"github.com/osteele/gojekyll/renderers"
	"github.com/osteele/gojekyll/utils"
	"github.com/osteele/liquid"
)

// Site is a Jekyll site.
type Site struct {
	Collections []*collection.Collection
	Routes      map[string]pages.Document // URL path -> Document; only for output pages

	cfg      config.Config
	data     map[string]interface{} // from _data files
	flags    config.Flags           // command-line flags, override config files
	themeDir string                 // absolute path to theme directory

	docs               []pages.Document // all documents, whether or not they are output
	nonCollectionPages []pages.Page

	renderer   *renderers.Manager
	renderOnce sync.Once

	drop     map[string]interface{} // cached drop value
	dropOnce sync.Once
}

// SourceDir returns the site source directory.
func (s *Site) SourceDir() string { return s.cfg.Source }

// DestDir returns the site destination directory.
func (s *Site) DestDir() string {
	if filepath.IsAbs(s.cfg.Destination) {
		return s.cfg.Destination
	}
	return filepath.Join(s.cfg.Source, s.cfg.Destination)
}

// OutputDocs returns a list of output pages.
func (s *Site) OutputDocs() []pages.Document {
	out := make([]pages.Document, 0, len(s.Routes))
	for _, p := range s.Routes {
		out = append(out, p)
	}
	return out
}

// Pages returns all the pages, output or not.
func (s *Site) Pages() (out []pages.Page) {
	for _, d := range s.docs {
		if p, ok := d.(pages.Page); ok {
			out = append(out, p)
		}
	}
	return
}

// Posts is part of the plugins.Site interface.
func (s *Site) Posts() []pages.Page {
	for _, c := range s.Collections {
		if c.Name == "posts" {
			return c.Pages()
		}
	}
	return nil
}

// AbsDir is in the collection.Site interface.
func (s *Site) AbsDir() string {
	d, err := filepath.Abs(s.SourceDir())
	if err != nil {
		panic(err)
	}
	return d
}

// Config is in the collection.Site interface.
func (s *Site) Config() *config.Config {
	return &s.cfg
}

func (s *Site) runHooks(h func(plugins.Plugin) error) error {
	for _, name := range s.cfg.Plugins {
		p, ok := plugins.Lookup(name)
		if ok {
			if err := h(p); err != nil {
				return err
			}
		}
	}
	return nil
}

// Site is in the pages.RenderingContext interface.
func (s *Site) Site() interface{} {
	return s
}

// PathPrefix is in the page.Container interface.
func (s *Site) PathPrefix() string { return "" }

// New creates a new site record, initialized with the site defaults.
func New(flags config.Flags) *Site {
	s := &Site{cfg: config.Default(), flags: flags}
	s.cfg.ApplyFlags(flags)
	return s
}

// SetAbsoluteURL overrides the loaded configuration.
// The server uses this.
func (s *Site) SetAbsoluteURL(url string) {
	s.cfg.AbsoluteURL = url
	s.cfg.Set("url", url)
	if s.drop != nil {
		s.drop["url"] = url
	}
}

// FilenameURLs returns a map of site-relative pathnames to URL paths
func (s *Site) FilenameURLs() map[string]string {
	urls := map[string]string{}
	for _, page := range s.Pages() {
		urls[utils.MustRel(s.SourceDir(), page.SourcePath())] = page.Permalink()
	}
	return urls
}

// KeepFile returns a boolean indicating that clean should leave the file in the destination directory.
func (s *Site) KeepFile(filename string) bool {
	return utils.SearchStrings(s.cfg.KeepFiles, filename)
}

// FilePathPage returns a Page, give a file path relative to site source directory.
func (s *Site) FilePathPage(rel string) (pages.Document, bool) {
	// This looks wasteful. If it shows up as a hotspot, you know what to do.
	for _, p := range s.Routes {
		if p.SourcePath() != "" {
			if r, err := filepath.Rel(s.SourceDir(), p.SourcePath()); err == nil {
				if r == rel {
					return p, true
				}
			}
		}
	}
	return nil, false
}

// FilenameURLPath returns a page's URL path, give a relative file path relative to the site source directory.
func (s *Site) FilenameURLPath(relpath string) (string, bool) {
	if p, found := s.FilePathPage(relpath); found {
		return p.Permalink(), true
	}
	return "", false
}

// RendererManager returns the rendering manager.
func (s *Site) RendererManager() renderers.Renderers {
	if s.renderer == nil {
		panic(fmt.Errorf("uninitialized rendering manager"))
	}
	return s.renderer
}

// TemplateEngine is part of the plugins.Site interface.
func (s *Site) TemplateEngine() *liquid.Engine {
	return s.renderer.TemplateEngine()
}

// initializeRenderers initializes the rendering manager
func (s *Site) initializeRenderers() (err error) {
	options := renderers.Options{
		RelativeFilenameToURL: s.FilenameURLPath,
		ThemeDir:              s.themeDir,
	}
	s.renderer, err = renderers.New(s.cfg, options)
	if err != nil {
		return err
	}
	engine := s.renderer.TemplateEngine()
	return s.runHooks(func(p plugins.Plugin) error {
		return p.ConfigureTemplateEngine(engine)
	})
}

// RelativePath is in the page.Container interface.
func (s *Site) RelativePath(path string) string {
	if s.themeDir != "" {
		if rel, err := filepath.Rel(s.themeDir, path); err == nil {
			return rel
		}
	}
	return utils.MustRel(s.cfg.Source, path)
}

// URLPage returns the page that will be served at URL
func (s *Site) URLPage(urlpath string) (p pages.Document, found bool) {
	p, found = s.Routes[urlpath]
	if !found {
		p, found = s.Routes[filepath.Join(urlpath, "index.html")]
	}
	if !found {
		p, found = s.Routes[filepath.Join(urlpath, "index.htm")]
	}
	return
}
