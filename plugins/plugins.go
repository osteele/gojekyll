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

// PluginContext is the context for plugin initialization.
// Currently, the only thing a plugin can do is add filters and tags.
type PluginContext interface {
	TemplateEngine() liquid.Engine
}

// Site is the site interface that is available to a plugin.
type Site interface {
	AddDocument(pages.Document, bool)
	Config() *config.Config
	Pages() []pages.Page
}

// Plugin describes the hooks that a plugin can override.
type Plugin interface {
	PostRead(site Site) error
}

type plugin struct{}

func (p plugin) PostRead(site Site) error { return nil }

// Find looks up a plugin by name
func Find(name string) (Plugin, bool) {
	switch name {
	case "jekyll-redirect-from":
		return jekyllFeedPlugin{}, true
	default:
		return nil, false
	}
}

// Install installs a plugin from the plugin directory.
func Install(name string, ctx PluginContext) bool {
	p, found := directory[name]
	if p != nil {
		if err := p(ctx, pluginHelper{ctx, name}); err != nil {
			panic(err)
		}
	}
	return found
}

var directory = map[string]func(PluginContext, pluginHelper) error{}

// register installs a plugin in the plugin directory.
func register(name string, fn func(PluginContext, pluginHelper) error) {
	directory[name] = fn
}

func init() {
	register("jekyll-feed", func(ctx PluginContext, h pluginHelper) error {
		h.stubbed()
		h.tag("feed_meta", h.makeUnimplementedTag())
		return nil
	})

	register("jekyll-seo-tag", func(ctx PluginContext, h pluginHelper) error {
		h.stubbed()
		h.tag("seo", h.makeUnimplementedTag())
		return nil
	})

	// the following plugins are always active
	// no warning but effect; the server runs in this mode anyway
	register("jekyll-live-reload", func(ctx PluginContext, h pluginHelper) error {
		return nil
	})
	register("jekyll-sass-converter", func(ctx PluginContext, h pluginHelper) error {
		return nil
	})
}

type pluginHelper struct {
	PluginContext
	name string
}

func (h pluginHelper) stubbed() {
	fmt.Printf("warning: gojekyll does not emulate the %s plugin. Some tags have been stubbed to prevent errors.\n", h.name)
}

func (h pluginHelper) tag(name string, r liquid.Renderer) {
	h.TemplateEngine().RegisterTag(name, r)
}

func (h pluginHelper) makeUnimplementedTag() liquid.Renderer {
	warned := false
	return func(ctx render.Context) (string, error) {
		if !warned {
			fmt.Printf("The %q tag in the %q plugin has not been implemented.\n", ctx.TagName(), h.name)
			warned = true
		}
		return fmt.Sprintf(`<!-- unimplemented tag: %q -->`, ctx.TagName()), nil
	}
}
