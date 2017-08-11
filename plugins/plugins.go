// Package plugins holds emulated Jekyll plugins.
//
// Unlike Jekyll, these are baked into the executable -- both because package "plugin'
// works only on Linux (as of 2017.07); and because the gojekyll implementation is immature and any possible interfaces
// are far from baked.
package plugins

import (
	"fmt"
	"regexp"
	"sort"

	"github.com/kyokomi/emoji"
	"github.com/osteele/gojekyll/config"
	"github.com/osteele/gojekyll/pages"
	"github.com/osteele/gojekyll/utils"
	"github.com/osteele/liquid"
)

// Plugin describes the hooks that a plugin can override.
type Plugin interface {
	Initialize(Site) error
	ConfigureTemplateEngine(*liquid.Engine) error
	ModifySiteDrop(Site, map[string]interface{}) error
	PostRead(Site) error
	PostRender([]byte) ([]byte, error)
}

// Site is the site interface that is available to a plugin.
type Site interface {
	AddDocument(pages.Document, bool)
	Config() *config.Config
	TemplateEngine() *liquid.Engine
	Pages() []pages.Page
}

// Lookup returns a plugin if it has been registered.
func Lookup(name string) (Plugin, bool) {
	p, found := directory[name]
	return p, found
}

// Install installs a registered plugin.
func Install(names []string, site Site) {
	for _, name := range names {
		if p, found := directory[name]; found {
			if err := p.Initialize(site); err != nil {
				panic(err)
			}
		} else {
			fmt.Printf("warning: gojekyll does not emulate the %s plugin.\n", name)
		}
	}
}

// Names returns a sorted list of names of registered plugins.
func Names() []string {
	var names []string
	for name := range directory {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Embed plugin to implement defaults implementations of the Plugin interface.
//
// This is internal until better baked.
type plugin struct{}

func (p plugin) Initialize(Site) error                             { return nil }
func (p plugin) ConfigureTemplateEngine(*liquid.Engine) error      { return nil }
func (p plugin) ModifySiteDrop(Site, map[string]interface{}) error { return nil }
func (p plugin) PostRead(Site) error                               { return nil }
func (p plugin) PostRender(b []byte) ([]byte, error)               { return b, nil }

var directory = map[string]Plugin{}

// register installs a plugin in the plugin directory.
//
// This is internal until better baked.
func register(name string, p Plugin) {
	directory[name] = p
}

// Add the built-in plugins defined in this file.
// More extensive plugins are defined and registered in own files.
func init() {
	register("jemoji", jemojiPlugin{})
	register("jekyll-mentions", jekyllMentionsPlugin{})
	register("jekyll-optional-front-matter", jekyllOptionalFrontMatterPlugin{})

	// Gojekyll behaves as though the following plugins are always loaded.
	// Define them here so we don't see warnings that they aren't defined.
	register("jekyll-live-reload", plugin{})
	register("jekyll-sass-converter", plugin{})
}

// Some small plugins are below. More involved plugins are in separate files.

// jemojiPlugin emulates the jekyll-jemoji plugin.
type jemojiPlugin struct{ plugin }

func (p jemojiPlugin) PostRender(b []byte) ([]byte, error) {
	return utils.ApplyToHTMLText(b, func(s string) string {
		return emoji.Sprint(s)
	}), nil
}

// jekyllMentionsPlugin emulates the jekyll-mentions plugin.
type jekyllMentionsPlugin struct{ plugin }

var mentionPattern = regexp.MustCompile(`@(\w+)`)

func (p jekyllMentionsPlugin) PostRender(b []byte) ([]byte, error) {
	return utils.ApplyToHTMLText(b, func(s string) string {
		return mentionPattern.ReplaceAllString(s, `<a href="https://github.com/$1" class="user-mention">@$1</a>`)
	}), nil
}

// jekyllOptionalFrontMatterPlugin emulates the jekyll-optional-front-matter plugin.
type jekyllOptionalFrontMatterPlugin struct{ plugin }

var requireFrontMatterExclude = []string{
	"README",
	"LICENSE",
	"LICENCE",
	"COPYING",
	"CODE_OF_CONDUCT",
	"CONTRIBUTING",
	"ISSUE_TEMPLATE",
	"PULL_REQUEST_TEMPLATE",
}

func (p jekyllOptionalFrontMatterPlugin) Initialize(s Site) error {
	m := map[string]bool{}
	for _, k := range requireFrontMatterExclude {
		m[k] = true
	}
	s.Config().RequireFrontMatter = false
	s.Config().RequireFrontMatterExclude = m
	return nil
}

// helpers

// func (p plugin) stubbed(name string) {
// 	fmt.Printf("warning: gojekyll does not emulate the %s plugin. Some tags have been stubbed to prevent errors.\n", name)
// }

// func (p plugin) makeUnimplementedTag(pluginName string) liquid.Renderer {
// 	warned := false
// 	return func(ctx render.Context) (string, error) {
// 		if !warned {
// 			fmt.Printf("The %q tag in the %q plugin has not been implemented.\n", ctx.TagName(), pluginName)
// 			warned = true
// 		}
// 		return fmt.Sprintf(`<!-- unimplemented tag: %q -->`, ctx.TagName()), nil
// 	}
// }
