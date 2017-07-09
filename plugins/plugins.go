// Package plugins holds emulated Jekyll plugins.
//
// Unlike Jekyll, these are baked into the executable -- both because as of 2017.07 package "plugin' currently
// works only on Linux, but also because the gojekyll implementation is immature and any possible interfaces
// are far from baked.
package plugins

import (
	"fmt"

	"github.com/osteele/gojekyll/config"
	"github.com/osteele/gojekyll/pages"
	"github.com/osteele/liquid"
	"github.com/osteele/liquid/render"
)

// Site is the site interface that is available to a plugin.
type Site interface {
	AddDocument(pages.Document, bool)
	Config() *config.Config
	Pages() []pages.Page
}

// Plugin describes the hooks that a plugin can override.
type Plugin interface {
	ConfigureTemplateEngine(liquid.Engine) error
	Initialize(Site) error
	PostRead(site Site) error
}

type plugin struct{}

func (p plugin) ConfigureTemplateEngine(liquid.Engine) error { return nil }
func (p plugin) Initialize(Site) error                       { return nil }
func (p plugin) PostRead(Site) error                         { return nil }

// Lookup returns a plugin if it has been registered.
func Lookup(name string) (Plugin, bool) {
	p, found := directory[name]
	return p, found
}

// Install installs a plugin from the plugin directory.
func Install(names []string, site Site) {
	for _, name := range names {
		p, found := directory[name]
		if found {
			if err := p.Initialize(site); err != nil {
				panic(err)
			}
		} else {
			fmt.Printf("warning: gojekyll does not emulate the %s plugin.\n", name)
		}
	}
}

var directory = map[string]Plugin{}

// register installs a plugin in the plugin directory.
func register(name string, p Plugin) {
	directory[name] = p
}

func init() {
	register("jekyll-feed", jekyllFeedPlugin{})
	register("jekyll-seo-tag", jekyllSEOTagPlugin{})

	// the following plugins are always active
	// no warning but effect; the server runs in this mode anyway
	register("jekyll-live-reload", plugin{})
	register("jekyll-sass-converter", plugin{})
}

type jekyllFeedPlugin struct{ plugin }

func (p jekyllFeedPlugin) ConfigureTemplateEngine(e liquid.Engine) error {
	p.stubbed("jekyll-feed")
	e.RegisterTag("feed_meta", p.makeUnimplementedTag("jekyll-feed"))
	return nil
}

type jekyllSEOTagPlugin struct{ plugin }

func (p jekyllSEOTagPlugin) ConfigureTemplateEngine(e liquid.Engine) error {
	p.stubbed("jekyll-seo-tag")
	e.RegisterTag("seo", p.makeUnimplementedTag("jekyll-seo-tag"))
	return nil
}

func (p plugin) stubbed(name string) {
	fmt.Printf("warning: gojekyll does not emulate the %s plugin. Some tags have been stubbed to prevent errors.\n", name)
}

func (p plugin) makeUnimplementedTag(pluginName string) liquid.Renderer {
	warned := false
	return func(ctx render.Context) (string, error) {
		if !warned {
			fmt.Printf("The %q tag in the %q plugin has not been implemented.\n", ctx.TagName(), pluginName)
			warned = true
		}
		return fmt.Sprintf(`<!-- unimplemented tag: %q -->`, ctx.TagName()), nil
	}
}
