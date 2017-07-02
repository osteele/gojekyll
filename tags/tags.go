package tags

import (
	"fmt"
	"io"
	"path"

	"github.com/osteele/gojekyll/config"
	"github.com/osteele/liquid"
	"github.com/osteele/liquid/chunks"
)

// A LinkTagHandler given an include tag file name returns a URL.
type LinkTagHandler func(string) (string, bool)

// AddJekyllTags adds the Jekyll tags to the Liquid engine.
func AddJekyllTags(e liquid.Engine, c config.Config, lh LinkTagHandler) {
	tc := tagContext{c, lh}
	e.RegisterBlock("highlight", highlightTag)
	e.RegisterTag("include", tc.includeTag)
	e.RegisterTag("link", tc.linkTag)
	e.RegisterTag("post_url", tc.postURLTag)

	// TODO unimplemented
	e.RegisterTag("post_url", CreateUnimplementedTag())
}

// tagContext provides the context to a tag renderer.
type tagContext struct {
	config config.Config
	lh     LinkTagHandler
}

// CreateUnimplementedTag creates a tag definition that prints a warning the first
// time it's rendered, and otherwise does nothing.
func CreateUnimplementedTag() liquid.TagDefinition {
	warned := false
	return func(_ io.Writer, ctx chunks.RenderContext) error {
		if !warned {
			fmt.Printf("The %q tag has not been implemented. It is being ignored.\n", ctx.TagName())
			warned = true
		}
		return nil
	}
}

func (tc tagContext) linkTag(w io.Writer, ctx chunks.RenderContext) error {
	filename := ctx.TagArgs()
	url, found := tc.lh(filename)
	if !found {
		return fmt.Errorf("missing link filename: %s", filename)
	}
	_, err := w.Write([]byte(url))
	return err
}

func (tc tagContext) postURLTag(w io.Writer, ctx chunks.RenderContext) error {
	var (
		filename = ctx.TagArgs()
		found    = false
		url      string
	)
	for _, ext := range append(tc.config.MarkdownExtensions(), ",") {
		url, found = tc.lh(path.Join("_posts", filename+ext))
		if found {
			break
		}
	}
	if !found {
		return fmt.Errorf("missing post_url filename: %s", filename)
	}
	_, err := w.Write([]byte(url))
	return err
}
