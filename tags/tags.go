package tags

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/osteele/gojekyll/config"
	"github.com/osteele/liquid"
	"github.com/osteele/liquid/chunks"
)

// A LinkTagHandler given an include tag file name returns a URL.
type LinkTagHandler func(string) (string, bool)

func AddJekyllTags(engine liquid.Engine, config config.Config, linkHandler LinkTagHandler) {
	tc := tagContext{config, linkHandler}
	engine.DefineTag("link", tc.linkTag)
	engine.DefineTag("include", tc.includeTag)

	// TODO unimplemented
	engine.DefineTag("post_url", emptyTag)
	engine.DefineStartTag("highlight", highlightTag)
}

type tagContext struct {
	config      config.Config
	linkHandler LinkTagHandler
}

func emptyTag(_ string) (func(io.Writer, chunks.RenderContext) error, error) {
	return func(w io.Writer, _ chunks.RenderContext) error { return nil }, nil
}

func highlightTag(w io.Writer, ctx chunks.RenderContext) error {
	args, err := ctx.ParseTagArgs()
	if err != nil {
		return err
	}
	cargs := []string{}
	if args != "" {
		cargs = append(cargs, "-l"+args)
	}
	s, err := ctx.InnerString()
	if err != nil {
		return err
	}
	if true {
		_, err = w.Write([]byte(s))
		return err
	}
	cmd := exec.Command("pygmentize", cargs...)
	cmd.Stdin = strings.NewReader(s)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (tc tagContext) linkTag(filename string) (func(io.Writer, chunks.RenderContext) error, error) {
	return func(w io.Writer, _ chunks.RenderContext) error {
		url, found := tc.linkHandler(filename)
		if !found {
			return fmt.Errorf("missing link filename: %s", filename)
		}
		_, err := w.Write([]byte(url))
		return err
	}, nil
}
