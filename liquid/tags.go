package liquid

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/osteele/liquid/chunks"
)

func (e *Wrapper) addJekyllTags() {
	e.engine.DefineTag("link", func(filename string) (func(io.Writer, chunks.Context) error, error) {
		return func(w io.Writer, _ chunks.Context) error {
			url, found := e.linkHandler(filename)
			if !found {
				return fmt.Errorf("missing link filename: %s", filename)
			}
			_, err := w.Write([]byte(url))
			return err
		}, nil
	})
	e.engine.DefineTag("include", func(line string) (func(io.Writer, chunks.Context) error, error) {
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
		return func(w io.Writer, ctx chunks.Context) error {
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
	})
}
