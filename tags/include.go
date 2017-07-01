package tags

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/osteele/liquid/chunks"
)

func (tc tagContext) includeTag(w io.Writer, ctx chunks.RenderContext) error {
	args, err := ParseArgs(ctx.TagArgs())
	if err != nil {
		return err
	}
	if len(args.Args) != 1 {
		return fmt.Errorf("parse error")
	}
	include, err := args.EvalOptions(ctx)
	if err != nil {
		return err
	}
	filename := filepath.Join(tc.config.Source, tc.config.IncludesDir, args.Args[0])
	ctx2 := ctx.Clone()
	ctx2.UpdateBindings(map[string]interface{}{"include": include})
	return ctx2.RenderFile(w, filename)

}
