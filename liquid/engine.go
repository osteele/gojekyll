package liquid

import (
	"io"

	"github.com/osteele/liquid"
	"github.com/osteele/liquid/chunks"
)

// Engine is a configured liquid engine.
type Engine interface {
	Parse([]byte) (liquid.Template, error)
	ParseAndRender([]byte, map[string]interface{}) ([]byte, error)
	DefineTag(string, func(string) (func(io.Writer, chunks.Context) error, error))
}

// Wrapper is a wrapper around the Liquid engine.
type Wrapper struct {
	engine               liquid.Engine
	linkHandler          LinkTagHandler
	includeTagHandler    IncludeTagHandler
	BaseURL, AbsoluteURL string
}

// IncludeTagHandler resolves the filename in a Liquid include tag into the expanded content
// of the included file.
type IncludeTagHandler func(string, io.Writer, map[string]interface{}) error

// A LinkTagHandler given an include tag file name returns a URL.
type LinkTagHandler func(string) (string, bool)

// NewEngine creates a LocalEngine.
func NewEngine() *Wrapper {
	e := &Wrapper{engine: liquid.NewEngine()}
	e.addJekyllFilters()
	e.addJekyllTags()
	return e
}

// DefineTag is in the Engine interface.
func (e *Wrapper) DefineTag(name string, f func(string) (func(io.Writer, chunks.Context) error, error)) {
	e.engine.DefineTag(name, f)
}

// LinkTagHandler sets the link tag handler.
func (e *Wrapper) LinkTagHandler(h LinkTagHandler) {
	e.linkHandler = h
}

// IncludeHandler sets the include tag handler.
func (e *Wrapper) IncludeHandler(h IncludeTagHandler) {
	e.includeTagHandler = h
}

// Parse is a wrapper for liquid.Parse.
func (e *Wrapper) Parse(source []byte) (liquid.Template, error) {
	return e.engine.ParseTemplate(source)
}

// ParseAndRender parses and then renders the template.
func (e *Wrapper) ParseAndRender(source []byte, scope map[string]interface{}) ([]byte, error) {
	t, err := e.Parse(source)
	if err != nil {
		return nil, err
	}
	return t.Render(scope)
}
