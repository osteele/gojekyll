package plugins

import (
	"fmt"
	"io"
	"strings"

	"github.com/osteele/gojekyll/liquid"
	"github.com/osteele/liquid/chunks"
)

type PluginContext interface {
	TemplateEngine() liquid.Engine
}

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

func registerPlugin(name string, fn func(PluginContext) error) {
	plugins[name] = fn
}

func warnUnimplemented(name string) {
	fmt.Printf("warning: gojekyll does not emulate the %s plugin. Some tags have been stubbed to prevent errors.\n", name)
}

func emptyTag(lexer string) (func(io.Writer, chunks.Context) error, error) {
	return func(w io.Writer, _ chunks.Context) error { return nil }, nil
}

func init() {
	registerPlugin("jekyll-avatar", func(ctx PluginContext) error {
		warnUnimplemented("jekyll-avatar")
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

const avatarTemplate = `<img class="avatar avatar-small" src="https://avatars3.githubusercontent.com/{username}?v=3&amp;s=40" alt="{username}" srcset="https://avatars3.githubusercontent.com/{username}?v=3&amp;s=40 1x, https://avatars3.githubusercontent.com/{username}?v=3&amp;s=80 2x, https://avatars3.githubusercontent.com/{username}?v=3&amp;s=120 3x, https://avatars3.githubusercontent.com/{username}?v=3&amp;s=160 4x" width="40" height="40" />`

func avatarTag(filename string) (func(io.Writer, chunks.Context) error, error) {
	username := "osteele" // TODO replace with real name
	size := 40
	return func(w io.Writer, _ chunks.Context) error {
		s := strings.Replace(avatarTemplate, "40", fmt.Sprintf("%s", size), -1)
		s = strings.Replace(s, "{username}", username, -1)
		_, err := w.Write([]byte(s))
		return err
	}, nil
}
