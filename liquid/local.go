package liquid

import (
	"fmt"
	"io"

	"github.com/osteele/liquid"
	"github.com/osteele/liquid/chunks"
)

// LocalWrapperEngine is a wrapper around osteele/liquid.
type LocalWrapperEngine struct {
	engine            liquid.Engine
	linkHandler       LinkTagHandler
	includeTagHandler IncludeTagHandler
}

// NewLocalWrapperEngine creates a LocalEngine.
func NewLocalWrapperEngine() LocalEngine {
	e := &LocalWrapperEngine{engine: liquid.NewEngine()}
	AddStandardFilters(e)
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
	return e
}

// LinkTagHandler sets the link tag handler.
func (e *LocalWrapperEngine) LinkTagHandler(h LinkTagHandler) {
	e.linkHandler = h
}

// IncludeHandler sets the include tag handler.
func (e *LocalWrapperEngine) IncludeHandler(h IncludeTagHandler) {
	e.includeTagHandler = h
}

// Parse is a wrapper for liquid.Parse.
func (e *LocalWrapperEngine) Parse(source []byte) (Template, error) {
	// fmt.Println("parse", string(source))
	t, err := e.engine.ParseTemplate(source)
	// return &localTemplate{t}, err
	return t, err
}

// ParseAndRender parses and then renders the template.
func (e *LocalWrapperEngine) ParseAndRender(source []byte, scope map[string]interface{}) ([]byte, error) {
	t, err := e.Parse(source)
	if err != nil {
		return nil, err
	}
	// fmt.Println("render", t)
	return t.Render(scope)
}
