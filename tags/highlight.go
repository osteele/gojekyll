package tags

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/osteele/liquid/chunks"
)

func highlightTag(w io.Writer, ctx chunks.RenderContext) error {
	args, err := ctx.ParseTagArgs()
	if err != nil {
		return err
	}
	cargs := []string{"-f", "html"}
	if args != "" {
		cargs = append(cargs, "-l"+args)
	}
	s, err := ctx.InnerString()
	if err != nil {
		return err
	}
	return withFileCache(w, fmt.Sprintf("pygments %s", args), s, func(w io.Writer) error {
		cmd := exec.Command("pygmentize", cargs...) // nolint: gas
		cmd.Stdin = strings.NewReader(s)
		cmd.Stdout = w
		cmd.Stderr = os.Stderr
		return cmd.Run()
	})
}
