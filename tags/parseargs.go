package tags

import (
	"fmt"
	"regexp"

	"github.com/osteele/liquid/render"
)

// TODO string escapes
var argPattern = regexp.MustCompile(`^([^=\s]+)(?:\s+|$)`)
var optionPattern = regexp.MustCompile(`^(\w+)=("[^"]*"|'[^']*'|[^'"\s]*)(?:\s+|$)`)

// ParsedArgs holds the parsed arguments from ParseArgs.
type ParsedArgs struct {
	Args    []string
	Options map[string]optionRecord
}

type optionRecord struct {
	value  string
	quoted bool
}

// ParseArgs parses a tag argument line {% include arg1 arg2 opt=a opt2='b' %}
func ParseArgs(argsline string) (*ParsedArgs, error) {
	args := ParsedArgs{
		[]string{},
		map[string]optionRecord{},
	}
	// Ranging over FindAllStringSubmatch would be better golf but got out of hand
	// maintenance-wise.
	for r, i := argsline, 0; len(r) > 0; r = r[i:] {
		am := argPattern.FindStringSubmatch(r)
		om := optionPattern.FindStringSubmatch(r)
		switch {
		case am != nil:
			args.Args = append(args.Args, am[1])
			i = len(am[0])
		case om != nil:
			k, v, quoted := om[1], om[2], false
			if v[0] == '\'' || v[0] == '"' {
				v, quoted = v[1:len(v)-1], true
			}
			args.Options[k] = optionRecord{v, quoted}
			i = len(om[0])
		default:
			return nil, fmt.Errorf("parse error in tag parameters %q", argsline)
		}
	}
	return &args, nil
}

// EvalOptions evaluates unquoted options.
func (r *ParsedArgs) EvalOptions(ctx render.Context) (map[string]interface{}, error) {
	options := map[string]interface{}{}
	for k, v := range r.Options {
		if v.quoted {
			options[k] = v.value
		} else {
			value, err := ctx.EvaluateString(v.value)
			if err != nil {
				return nil, err
			}
			options[k] = value
		}
	}
	return options, nil
}
