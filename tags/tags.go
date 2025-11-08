package tags

import (
	"fmt"
	"path"

	"github.com/osteele/gojekyll/config"
	"github.com/osteele/gojekyll/logger"
	"github.com/osteele/liquid"
	"github.com/osteele/liquid/render"
)

// A LinkTagHandler given an include tag file name returns a URL.
type LinkTagHandler func(string) (string, bool)

// AddJekyllTags adds the Jekyll tags to the Liquid engine.
func AddJekyllTags(e *liquid.Engine, c *config.Config, includeDirs []string, lh LinkTagHandler) {
	tc := tagContext{c, includeDirs, lh}
	e.RegisterBlock("highlight", highlightTag)
	e.RegisterTag("include", tc.includeTag)
	e.RegisterTag("include_relative", tc.includeRelativeTag)
	e.RegisterTag("link", tc.linkTag)
	e.RegisterTag("post_url", tc.postURLTag)
}

// tagContext provides the context to a tag renderer.
type tagContext struct {
	cfg         *config.Config
	includeDirs []string
	lh          LinkTagHandler
}

// CreateUnimplementedTag creates a tag definition that prints a warning the first
// time it's rendered, and otherwise does nothing.
func CreateUnimplementedTag() liquid.Renderer {
	warned := false
	log := logger.Default()
	return func(rc render.Context) (string, error) {
		if !warned {
			log.Warn("The %q tag has not been implemented. It is being ignored.", rc.TagName())
			warned = true
		}
		return "", nil
	}
}

func (tc tagContext) linkTag(rc render.Context) (string, error) {
	filename := rc.TagArgs()
	url, found := tc.lh(filename)
	if !found {
		return "", fmt.Errorf("missing link filename: %s", filename)
	}
	return url, nil
}

func (tc tagContext) postURLTag(rc render.Context) (string, error) {
	var (
		filename = rc.TagArgs()
		found    = false
		url      string
	)
	for _, ext := range append(tc.cfg.MarkdownExtensions(), "") {
		url, found = tc.lh(path.Join("_posts", filename+ext))
		if found {
			break
		}
	}
	if !found {
		return "", fmt.Errorf("missing post_url filename: %s", filename)
	}
	return url, nil
}
