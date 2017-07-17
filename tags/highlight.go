package tags

import (
	"bytes"
	"fmt"
	"html"
	"os"
	"os/exec"
	"strings"

	"github.com/osteele/gojekyll/cache"
	"github.com/osteele/liquid/render"
)

const pygmentizeCmd = "pygmentize"

// warn once per execution, even on watch/rebuilds
var warnedMissingPygmentize = false

func highlightTag(ctx render.Context) (string, error) {
	args, err := ctx.ExpandTagArg()
	if err != nil {
		return "", err
	}
	cmdArgs := []string{"-f", "html"}
	if args != "" {
		cmdArgs = append(cmdArgs, "-l"+args)
	}
	s, err := ctx.InnerString()
	if err != nil {
		return "", err
	}
	r, err := cache.WithFile(fmt.Sprintf("pygments %s", args), s, func() (string, error) {
		buf := new(bytes.Buffer)
		cmd := exec.Command(pygmentizeCmd, cmdArgs...) // nolint: gas
		cmd.Stdin = strings.NewReader(s)
		cmd.Stdout = buf
		cmd.Stderr = os.Stderr
		if e := cmd.Run(); e != nil {
			return "", e
		}
		return buf.String(), nil
	})
	if e, ok := err.(*exec.Error); ok {
		if e.Err == exec.ErrNotFound {
			r, err = `<code>`+html.EscapeString(s)+`</code>`, nil
			if !warnedMissingPygmentize {
				warnedMissingPygmentize = true
				_, err = fmt.Fprintf(os.Stdout, "%s\nThe {%% highlight %%} tag will use <code>…</code> instead\n", err)
			}
		}
	}
	return r, err
}
