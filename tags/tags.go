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

// AddJekyllTags adds the Jekyll tags to the Liquid engine.
func AddJekyllTags(e liquid.Engine, c config.Config, lh LinkTagHandler) {
	tc := tagContext{c, lh}
	e.DefineTag("link", tc.linkTag)
	e.DefineTag("include", tc.includeTag)

	// TODO unimplemented
	e.DefineTag("post_url", emptyTag)
	e.DefineStartTag("highlight", highlightTag)
}

type tagContext struct {
	config config.Config
	lh     LinkTagHandler
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
	// TODO this is disabled for performance; make it configurable instead.
	if true {
		_, err = w.Write([]byte(s))
		return err
	}
	cmd := exec.Command("pygmentize", cargs...) // nolint: gas
	cmd.Stdin = strings.NewReader(s)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (tc tagContext) linkTag(filename string) (func(io.Writer, chunks.RenderContext) error, error) {
	return func(w io.Writer, _ chunks.RenderContext) error {
		url, found := tc.lh(filename)
		if !found {
			return fmt.Errorf("missing link filename: %s", filename)
		}
		_, err := w.Write([]byte(url))
		return err
	}, nil
}
