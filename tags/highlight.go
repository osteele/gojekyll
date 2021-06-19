package tags

import (
	"bytes"
	"fmt"
	"html"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/osteele/gojekyll/cache"
	"github.com/osteele/liquid/render"
)

const pygmentizeCmd = "pygmentize"

// warn once per execution, even on watch/rebuilds
var warnedMissingPygmentize = false
var highlightArgsRE = regexp.MustCompile(`^\s*(\S+)(\s+linenos)?\s*$`)

func highlightTag(rc render.Context) (string, error) {
	argStr, err := rc.ExpandTagArg()
	if err != nil {
		return "", err
	}
	args := highlightArgsRE.FindStringSubmatch(argStr)
	if args == nil {
		return "", fmt.Errorf("syntax error")
	}
	cmdArgs := []string{"-f", "html"}
	cmdArgs = append(cmdArgs, "-l"+args[1])
	if args[2] != "" {
		cmdArgs = append(cmdArgs, "-O", "linenos=1")
	}
	s, err := rc.InnerString()
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
		// This only works in go < 1.16:
		if e.Err == exec.ErrNotFound {
			r, err = maybeWarnMissingPygmentize(s, err)
		}
		// This is language-dependent, but works in go 1.16 too
		if strings.Contains(e.Err.Error(), "executable file not found") && e.Name == pygmentizeCmd {
			r, err = maybeWarnMissingPygmentize(s, err)
		}
	}
	// TODO: replace the test above by the following once support for go < 1.16
	// is dropped if pathErr, ok := err.(*fs.PathError); ok {
	//  if filepath.Base(pathErr.Path) == pygmentizeCmd {
	//      r, err = maybeWarnMissingPygmentize(s, err)
	//  }
	// }
	return r, err
}

func maybeWarnMissingPygmentize(s string, err error) (string, error) {
	r := `<code>` + html.EscapeString(s) + `</code>`
	if warnedMissingPygmentize {
		return r, nil
	}
	warnedMissingPygmentize = true
	_, err = fmt.Fprintf(os.Stderr,
		"Error: %s\n"+
			"Run `pip install Pygments` to install %s.\n"+
			"The {%% highlight %%} tag will use <code>…</code> instead.\n",
		err, pygmentizeCmd)
	return r, err
}
