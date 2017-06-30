package sites

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/osteele/gojekyll/collections"
	"github.com/osteele/gojekyll/config"
	"github.com/osteele/gojekyll/helpers"
	"github.com/osteele/gojekyll/pages"
	"github.com/osteele/gojekyll/pipelines"
	"github.com/osteele/gojekyll/plugins"
	"github.com/osteele/gojekyll/templates"
	"github.com/osteele/liquid"
)

// Site is a Jekyll site.
type Site struct {
	ConfigFile  *string
	Source      string
	Destination string

	Collections []*collections.Collection
	// Variables   templates.VariableMap
	Routes map[string]pages.Document // URL path -> Page, only for output pages

	config           config.Config
	data             map[string]interface{}
	pipeline         *pipelines.Pipeline
	pages            []pages.Document // all pages, output or not
	preparedToRender bool
	siteVariables    templates.VariableMap
}

// OutputPages returns a list of output pages.
func (s *Site) OutputPages() []pages.Document {
	out := make([]pages.Document, 0, len(s.Routes))
	for _, p := range s.Routes {
		out = append(out, p)
	}
	return out
}

// Pages returns all the pages, output or not.
func (s *Site) Pages() []pages.Document { return s.pages }

// PathPrefix is in the page.Container interface.
func (s *Site) PathPrefix() string { return "" }

// NewSite creates a new site record, initialized with the site defaults.
func NewSite() *Site {
	return &Site{config: config.Default()}
}

// SetAbsoluteURL overrides the loaded configuration.
// The server uses this.
func (s *Site) SetAbsoluteURL(url string) {
	s.config.AbsoluteURL = url
	s.config.Variables["url"] = url
	if s.siteVariables != nil {
		s.siteVariables["url"] = url
	}
}

// FilenameURLs returns a map of relative filenames to URL paths
func (s *Site) FilenameURLs() map[string]string {
	urls := map[string]string{}
	for _, page := range s.Pages() {
		urls[page.SiteRelPath()] = page.Permalink()
	}
	return urls
}

// KeepFile returns a boolean indicating that clean should leave the file in the destination directory.
func (s *Site) KeepFile(path string) bool {
	// TODO
	return false
}

// RelPathPage returns a Page, give a file path relative to site source directory.
func (s *Site) RelPathPage(relpath string) (pages.Document, bool) {
	for _, p := range s.Routes {
		if p.SiteRelPath() == relpath {
			return p, true
		}
	}
	return nil, false
}

// RelativeFilenameToURL returns a page's relative URL, give a file path relative to the site source directory.
func (s *Site) RelativeFilenameToURL(relpath string) (string, bool) {
	var url string
	p, found := s.RelPathPage(relpath)
	if found {
		url = p.Permalink()
	}
	return url, found
}

// RenderingPipeline returns the rendering pipeline.
func (s *Site) RenderingPipeline() pipelines.PipelineInterface {
	if s.pipeline == nil {
		panic(fmt.Errorf("uninitialized rendering pipeline"))
	}
	return s.pipeline
}

type pluginContext struct {
	engine liquid.Engine
}

// Engine is in the PluginContext interface.
func (c pluginContext) TemplateEngine() liquid.Engine { return c.engine }

// initializeRenderingPipeline initializes the rendering pipeline
func (s *Site) initializeRenderingPipeline() (err error) {
	options := pipelines.PipelineOptions{
		SourceDir:             s.Source,
		RelativeFilenameToURL: s.RelativeFilenameToURL,
	}
	s.pipeline, err = pipelines.NewPipeline(s.config, options)
	ctx := pluginContext{s.pipeline.TemplateEngine()}
	for _, name := range s.config.Plugins {
		if !plugins.Install(name, ctx) {
			fmt.Printf("warning: gojekyll does not emulate the %s plugin.\n", name)
		}
	}
	return
}

// OutputExt is in the page.Container interface.
func (s *Site) OutputExt(pathname string) string {
	return s.config.OutputExt(pathname)
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
