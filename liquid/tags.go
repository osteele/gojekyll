package liquid

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/osteele/liquid/chunks"
)

func (e *Wrapper) addJekyllTags() {
	e.engine.DefineTag("link", e.linkTag)
	e.engine.DefineTag("include", e.includeTag)

	// TODO unimplemented
	e.engine.DefineTag("post_url", emptyTag)
	e.engine.DefineStartTag("highlight", highlightTag)
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

func (e *Wrapper) linkTag(filename string) (func(io.Writer, chunks.RenderContext) error, error) {
	return func(w io.Writer, _ chunks.RenderContext) error {
		url, found := e.linkHandler(filename)
		if !found {
			return fmt.Errorf("missing link filename: %s", filename)
		}
		_, err := w.Write([]byte(url))
		return err
	}, nil
}

func (e *Wrapper) includeTag(line string) (func(io.Writer, chunks.RenderContext) error, error) {
	// TODO string escapes
	includeLinePattern := regexp.MustCompile(`^\S+(?:\s+\S+=("[^"]+"|'[^']'|[^'"\s]+))*$`)
	includeParamPattern := regexp.MustCompile(`\b(\S+)=("[^"]+"|'[^']'|[^'"\s]+)(?:\s|$)`)
	if !includeLinePattern.MatchString(line) {
		return nil, fmt.Errorf("parse error in include tag parameters")
	}
	filename := strings.Fields(line)[0]
	type paramSpec struct {
		value string
		eval  bool
	}
	params := map[string]paramSpec{}
	for _, m := range includeParamPattern.FindAllStringSubmatch(line, -1) {
		k, v, eval := m[1], m[2], true
		if strings.HasPrefix(v, `'`) || strings.HasPrefix(v, `"`) {
			v, eval = v[1:len(v)-1], false
		}
		params[k] = paramSpec{v, eval}
	}
	return func(w io.Writer, ctx chunks.RenderContext) error {
		include := map[string]interface{}{}
		for k, v := range params {
			if v.eval {
				value, err := ctx.EvaluateString(v.value)
				if err != nil {
					return err
				}
				include[k] = value
			} else {
				include[k] = v.value
			}
		}
		bindings := map[string]interface{}{}
		for k, v := range ctx.GetVariableMap() {
			bindings[k] = v
		}
		bindings["include"] = include
		return e.includeTagHandler(filename, w, bindings)
	}, nil
}
