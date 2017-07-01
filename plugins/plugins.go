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
		if err := p(ctx); err != nil {
			panic(err)
		}
	}
	return found
}

var plugins = map[string]func(PluginContext) error{}

// registerPlugin installs a plugin in the plugin directory.
func registerPlugin(name string, fn func(PluginContext) error) {
	plugins[name] = fn
}

func warnUnimplemented(name string) {
	fmt.Printf("warning: gojekyll does not emulate the %s plugin. Some tags have been stubbed to prevent errors.\n", name)
}

func emptyTag(lexer string) (func(io.Writer, chunks.RenderContext) error, error) {
	return func(w io.Writer, _ chunks.RenderContext) error { return nil }, nil
}

func init() {
	registerPlugin("jekyll-feed", func(ctx PluginContext) error {
		warnUnimplemented("jekyll-feed")
		ctx.TemplateEngine().DefineTag("feed_meta", emptyTag)
		return nil
	})

	registerPlugin("jekyll-seo-tag", func(ctx PluginContext) error {
		warnUnimplemented("jekyll-seo-tag")
		ctx.TemplateEngine().DefineTag("seo", emptyTag)
		return nil
	})
}
