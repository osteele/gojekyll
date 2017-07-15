package tags

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/osteele/liquid/render"
)

func (tc tagContext) includeTag(ctx render.Context) (string, error) {
	return includeFromDir(ctx, filepath.Join(tc.config.Source, tc.config.IncludesDir))
}

func (tc tagContext) includeRelativeTag(ctx render.Context) (string, error) {
	// TODO "Note that you cannot use the ../ syntax"
	return includeFromDir(ctx, path.Dir(ctx.SourceFile()))
}

func includeFromDir(ctx render.Context, dirname string) (string, error) {
	argsline, err := ctx.ExpandTagArg()
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
	include, err := args.EvalOptions(ctx)
	if err != nil {
		return "", err
	}
	filename := filepath.Join(dirname, args.Args[0])
	return ctx.RenderFile(filename, map[string]interface{}{"include": include})
}
