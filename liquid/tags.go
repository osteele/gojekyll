package liquid

import (
	"fmt"
	"io"

	"github.com/osteele/liquid/chunks"
)

func (e *LocalWrapperEngine) addJekyllTags() {
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
	e.engine.DefineTag("include", func(filename string) (func(io.Writer, chunks.Context) error, error) {
		return func(w io.Writer, ctx chunks.Context) error {
			return e.includeTagHandler(filename, w, ctx.GetVariableMap())
		}, nil
	})
}
