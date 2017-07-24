package tags

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/osteele/liquid/render"
)

func (tc tagContext) includeTag(rc render.Context) (string, error) {
	return includeFromDir(rc, filepath.Join(tc.cfg.Source, tc.cfg.IncludesDir))
}

func (tc tagContext) includeRelativeTag(rc render.Context) (string, error) {
	// TODO "Note that you cannot use the ../ syntax"
	return includeFromDir(rc, path.Dir(rc.SourceFile()))
}

func includeFromDir(rc render.Context, dirname string) (string, error) {
	argsline, err := rc.ExpandTagArg()
	if err != nil {
		return "", err
	}
	args, err := ParseArgs(argsline)
	if err != nil {
		return "", err
	}
	if len(args.Args) != 1 {
		return "", fmt.Errorf("parse error")
	}
	include, err := args.EvalOptions(rc)
	if err != nil {
		return "", err
	}
	filename := filepath.Join(dirname, args.Args[0])
	return rc.RenderFile(filename, map[string]interface{}{"include": include})
}
