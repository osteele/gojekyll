package tags

import (
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/osteele/liquid/chunks"
)

// TODO string escapes
var includeLinePattern = regexp.MustCompile(`^\S+(?:\s+\S+=("[^"]+"|'[^']'|[^'"\s]+))*$`)
var includeParamPattern = regexp.MustCompile(`\b(\S+)=("[^"]+"|'[^']'|[^'"\s]+)(?:\s|$)`)

type includeArgSpec struct {
	value string
	eval  bool
}

type includeSpec struct {
	filename string
	args     map[string]includeArgSpec
}

func parseIncludeArgs(line string) (*includeSpec, error) {
	if !includeLinePattern.MatchString(line) {
		return nil, fmt.Errorf("parse error in include tag parameters")
	}
	spec := includeSpec{
		strings.Fields(line)[0],
		map[string]includeArgSpec{},
	}
	for _, m := range includeParamPattern.FindAllStringSubmatch(line, -1) {
		k, v, eval := m[1], m[2], true
		if strings.HasPrefix(v, `'`) || strings.HasPrefix(v, `"`) {
			v, eval = v[1:len(v)-1], false
		}
		spec.args[k] = includeArgSpec{v, eval}
	}
	return &spec, nil
}

func (spec *includeSpec) eval(ctx chunks.RenderContext) (map[string]interface{}, error) {
	include := map[string]interface{}{}
	for k, v := range spec.args {
		if v.eval {
			value, err := ctx.EvaluateString(v.value)
			if err != nil {
				return nil, err
			}
			include[k] = value
		} else {
			include[k] = v.value
		}
	}
	return include, nil
}

func (tc tagContext) includeTag(line string) (func(io.Writer, chunks.RenderContext) error, error) {
	spec, err := parseIncludeArgs(line)
	if err != nil {
		return nil, err
	}
	return func(w io.Writer, ctx chunks.RenderContext) error {
		params, err := spec.eval(ctx)
		if err != nil {
			return err
		}
		filename := filepath.Join(tc.config.Source, tc.config.IncludesDir, spec.filename)
		ctx2 := ctx.Clone()
		ctx2.UpdateBindings(map[string]interface{}{"include": params})
		return ctx2.RenderFile(w, filename)

	}, nil
}
