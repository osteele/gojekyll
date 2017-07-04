package tags

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/osteele/liquid/render"
)

func highlightTag(ctx render.Context) (string, error) {
	args, err := ctx.ParseTagArgs()
	if err != nil {
		return "", err
	}
	cargs := []string{"-f", "html"}
	if args != "" {
		cargs = append(cargs, "-l"+args)
	}
	s, err := ctx.InnerString()
	if err != nil {
		return "", err
	}
	return withFileCache(fmt.Sprintf("pygments %s", args), s, func() (string, error) {
		buf := new(bytes.Buffer)
		cmd := exec.Command("pygmentize", cargs...) // nolint: gas
		cmd.Stdin = strings.NewReader(s)
		cmd.Stdout = buf
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return "", err
		}
		return buf.String(), nil
	})
}
