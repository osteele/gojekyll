// Package plugins holds emulated Jekyll plugins.
//
// Unlike Jekyll, these are baked into the executable -- both because as of 2017.07 package "plugin' currently
// works only on Linux, but also because the gojekyll implementation is immature and any possible interfaces
// are far from baked.
package plugins

import (
	"fmt"
	"io"

	"github.com/osteele/liquid"
	"github.com/osteele/liquid/chunks"
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
		if err := p(ctx, pluginHelper{name}); err != nil {
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
		ctx.TemplateEngine().DefineTag("feed_meta", h.makeUnimplementedTag())
		return nil
	})

	registerPlugin("jekyll-seo-tag", func(ctx PluginContext, h pluginHelper) error {
		h.stubbed()
		ctx.TemplateEngine().DefineTag("seo", h.makeUnimplementedTag())
		return nil
	})
}

type pluginHelper struct{ name string }

func (h pluginHelper) stubbed() {
	fmt.Printf("warning: gojekyll does not emulate the %s plugin. Some tags have been stubbed to prevent errors.\n", h.name)
}

func (h pluginHelper) makeUnimplementedTag() liquid.TagDefinition {
	warned := false
	return func(_ io.Writer, ctx chunks.RenderContext) error {
		if !warned {
			fmt.Printf("The %q tag in the %q plugin has not been implemented.\n", ctx.TagName(), h.name)
			warned = true
		}
		return nil
	}
}
