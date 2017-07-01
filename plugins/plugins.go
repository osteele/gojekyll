// Package plugins holds emulated Jekyll plugins.
//
// Unlike Jekyll, these are baked into the executable -- both because as of 2017.07 package "plugin' currently
// works only on Linux, but also because the gojekyll implementation is immature and any possible interfaces
// are far from baked.
package plugins

import (
	"fmt"
	"io"
	"strings"

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
	registerPlugin("jekyll-avatar", func(ctx PluginContext) error {
		ctx.TemplateEngine().DefineTag("avatar", avatarTag)
		return nil
	})

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

	// registerPlugin("jekyll-sitemap")
	// registerPlugin("jemoji")
}

// this template is from the plugin documentation
const avatarTemplate = `<img class="avatar avatar-small" src="https://avatars3.githubusercontent.com/{user}?v=3&amp;s=40" alt="{user}" srcset="https://avatars3.githubusercontent.com/{user}?v=3&amp;s=40 1x, https://avatars3.githubusercontent.com/{user}?v=3&amp;s=80 2x, https://avatars3.githubusercontent.com/{user}?v=3&amp;s=120 3x, https://avatars3.githubusercontent.com/{user}?v=3&amp;s=160 4x" width="40" height="40" />`

func avatarTag(_ string) (func(io.Writer, chunks.RenderContext) error, error) {
	return func(w io.Writer, ctx chunks.RenderContext) error {
		var (
			user string
			size = "40"
		)
		args, err := ctx.ParseTagArgs()
		fmt.Sprintln("args", args)
		if err != nil {
			fmt.Println("err", err)
			return err
		}
		for _, arg := range strings.Fields(args) {
			split := strings.SplitN(arg, "=", 2)
			if len(split) == 1 {
				split = []string{"user", arg}
			}
			switch split[0] {
			case "user":
				user = split[1]
			case "size":
				size = split[1]
			default:
				return fmt.Errorf("unknown avatar argument: %s", split[0])
			}
		}
		if user == "" {
			return fmt.Errorf("parse error in avatar tag parameters %s", args)
		}
		s := strings.Replace(avatarTemplate, "40", size, -1)
		s = strings.Replace(s, "{user}", user, -1)
		_, err = w.Write([]byte(s))
		return err
	}, nil
}
