package tags

import (
	"fmt"
	"path/filepath"

	"github.com/osteele/gojekyll/utils"
	"github.com/osteele/liquid/render"
)

func (tc tagContext) includeTag(rc render.Context) (s string, err error) {
	for _, dir := range tc.includeDirs {
		s, err = includeFromDir(dir, rc)
		if err == nil {
			return
		}
	}
	return
}

func (tc tagContext) includeRelativeTag(rc render.Context) (string, error) {
	return includeFromDir(filepath.Dir(rc.SourceFile()), rc)
}

func includeFromDir(dir string, rc render.Context) (string, error) {
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
	filename, err := utils.JoinWithin(dir, args.Args[0])
	if err != nil {
		return "", err
	}
	return rc.RenderFile(filename, map[string]interface{}{"include": include})
}
