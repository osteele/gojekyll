// Package plugins holds emulated Jekyll plugins.
//
// Unlike Jekyll, these are baked into the executable -- both because as of 2017.07 package "plugin' currently
// works only on Linux, but also because the gojekyll implementation is immature and any possible interfaces
// are far from baked.
package plugins

import (
	"fmt"

	"github.com/osteele/liquid"
	"github.com/osteele/liquid/render"
)

// PluginContext is the context for plugin initialization.
// Currently, the only thing a plugin can do is add filters and tags.
type PluginContext interface {
	TemplateEngine() liquid.Engine
}

// Install installs a plugin from the plugin directory.
func Install(name string, ctx PluginContext) bool {
	p, found := plugins[name]
	if p != nil {
		if err := p(ctx, pluginHelper{ctx, name}); err != nil {
			panic(err)
		}
	}
	return found
}

var plugins = map[string]func(PluginContext, pluginHelper) error{}

// registerPlugin installs a plugin in the plugin directory.
func registerPlugin(name string, fn func(PluginContext, pluginHelper) error) {
	plugins[name] = fn
}

func init() {
	registerPlugin("jekyll-feed", func(ctx PluginContext, h pluginHelper) error {
		h.stubbed()
		h.tag("feed_meta", h.makeUnimplementedTag())
		return nil
	})

	registerPlugin("jekyll-seo-tag", func(ctx PluginContext, h pluginHelper) error {
		h.stubbed()
		h.tag("seo", h.makeUnimplementedTag())
		return nil
	})

	// the following plugins are always active
	// no warning but effect; the server runs in this mode anyway
	registerPlugin("jekyll-live-reload", func(ctx PluginContext, h pluginHelper) error {
		return nil
	})
	registerPlugin("jekyll-sass-converter", func(ctx PluginContext, h pluginHelper) error {
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
		return "", nil
	}
}
