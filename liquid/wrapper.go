package liquid

import (
	"bytes"
	"io"

	"github.com/acstech/liquid"
	"github.com/acstech/liquid/core"
)

// Engine is a configured liquid engine.
type Engine interface {
	ApplyTemplate([]byte, map[string]interface{}) ([]byte, error)
	ParseTemplate([]byte) (Template, error)
}

// Template is a liquid template.
type Template interface {
	Render(map[string]interface{}) ([]byte, error)
}

// LocalEngine is a wrapper around acstech/liquid.
type LocalEngine struct {
	config *core.Configuration
	lh     LinkHandler
}

// LocalTemplate is an wrapper around liquid.Template.
type LocalTemplate struct {
	e *LocalEngine
	t *liquid.Template
}

// NewLocalEngine creates a LocalEngine.
func NewLocalEngine() *LocalEngine {
	return &LocalEngine{}
}

// LinkHandler sets the link tag handler.
func (e *LocalEngine) LinkHandler(lh LinkHandler) {
	e.lh = lh
}

// IncludeHandler sets the include tag handler.
func (e *LocalEngine) IncludeHandler(ih func(string, io.Writer, map[string]interface{})) {
	e.config = liquid.Configure().IncludeHandler(ih)
}

// ParseTemplate is a wrapper for liquid.Parse.
func (e *LocalEngine) ParseTemplate(text []byte) (Template, error) {
	t, err := liquid.Parse(text, e.config)
	return &LocalTemplate{e, t}, err
}

// Render is a wrapper around liquid's template.Render that turns panics into errors.
func (t *LocalTemplate) Render(scope map[string]interface{}) (out []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				panic(r)
			}
		}
	}()
	SetLinkHandler(t.e.lh)
	writer := new(bytes.Buffer)
	t.t.Render(writer, scope)
	return writer.Bytes(), nil
}

// ApplyTemplate parses and then renders the template.
func (e *LocalEngine) ApplyTemplate(text []byte, scope map[string]interface{}) ([]byte, error) {
	t, err := e.ParseTemplate(text)
	if err != nil {
		return nil, err
	}
	return t.Render(scope)
}
