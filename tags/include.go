package tags

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/osteele/liquid/chunks"
)

func (tc tagContext) includeTag(argsline string) (func(io.Writer, chunks.RenderContext) error, error) {
	args, err := ParseArgs(argsline)
	if err != nil {
		return nil, err
	}
	if len(args.Args) != 1 {
		return nil, fmt.Errorf("parse error")
	}
	return func(w io.Writer, ctx chunks.RenderContext) error {
		include, err := args.EvalOptions(ctx)
		if err != nil {
			return err
		}
		filename := filepath.Join(tc.config.Source, tc.config.IncludesDir, args.Args[0])
		ctx2 := ctx.Clone()
		ctx2.UpdateBindings(map[string]interface{}{"include": include})
		return ctx2.RenderFile(w, filename)

	}, nil
}
