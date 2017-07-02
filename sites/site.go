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
	"github.com/osteele/liquid"
)

// Site is a Jekyll site.
type Site struct {
	ConfigFile  *string
	Collections []*collections.Collection
	// Variables   map[string]interface{}
	Routes map[string]pages.Document // URL path -> Page, only for output pages

	config           config.Config
	data             map[string]interface{}
	flags            config.Flags
	pipeline         *pipelines.Pipeline
	docs             []pages.Document // all documents, whether or not they are output
	preparedToRender bool
	siteVariables    map[string]interface{}
}

// SourceDir returns the site source directory.
func (s *Site) SourceDir() string { return s.config.Source }

// DestDir returns the site destination directory.
func (s *Site) DestDir() string {
	if filepath.IsAbs(s.config.Destination) {
		return s.config.Destination
	}
	return filepath.Join(s.config.Source, s.config.Destination)
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
func (s *Site) Pages() []pages.Document { return s.docs }

// AbsDir is in the page.Container interface.
func (s *Site) AbsDir() string {
	d, err := filepath.Abs(s.SourceDir())
	if err != nil {
		panic(err)
	}
	return d
}

// Config is in the page.Container interface.
func (s *Site) Config() config.Config {
	return s.config
}

// PathPrefix is in the page.Container interface.
func (s *Site) PathPrefix() string { return "" }

// NewSite creates a new site record, initialized with the site defaults.
func NewSite(flags config.Flags) *Site {
	s := &Site{config: config.Default(), flags: flags}
	s.config.ApplyFlags(flags)
	return s
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
func (s *Site) KeepFile(filename string) bool {
	return helpers.SearchStrings(s.config.KeepFiles, filename)
}

// FilePathPage returns a Page, give a file path relative to site source directory.
func (s *Site) FilePathPage(relpath string) (pages.Document, bool) {
	for _, p := range s.Routes {
		if p.SiteRelPath() == relpath {
			return p, true
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
		RelativeFilenameToURL: s.FilenameURLPath,
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
