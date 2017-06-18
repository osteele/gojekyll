package liquid

import (
	"bytes"
	"io"

	"github.com/acstech/liquid"
	"github.com/acstech/liquid/core"
)

// LocalWrapperEngine is a wrapper around acstech/liquid.
type LocalWrapperEngine struct {
	config      *core.Configuration
	linkHandler LinkHandler
}

type localTemplate struct {
	engine *LocalWrapperEngine
	lt     *liquid.Template
}

// NewLocalWrapperEngine creates a LocalEngine.
func NewLocalWrapperEngine() LocalEngine {
	return &LocalWrapperEngine{}
}

// LinkHandler sets the link tag handler.
func (engine *LocalWrapperEngine) LinkHandler(h LinkHandler) {
	engine.linkHandler = h
}

// IncludeHandler sets the include tag handler.
func (engine *LocalWrapperEngine) IncludeHandler(h func(string, io.Writer, map[string]interface{})) {
	engine.config = liquid.Configure().IncludeHandler(h)
}

// Parse is a wrapper for liquid.Parse.
func (engine *LocalWrapperEngine) Parse(text []byte) (Template, error) {
	template, err := liquid.Parse(text, engine.config)
	return &localTemplate{engine, template}, err
}

// ParseAndRender parses and then renders the template.
func (engine *LocalWrapperEngine) ParseAndRender(text []byte, scope map[string]interface{}) ([]byte, error) {
	template, err := engine.Parse(text)
	if err != nil {
		return nil, err
	}
	return template.Render(scope)
}

// Render is a wrapper around liquid's template.Render that turns panics into errors.
func (template *localTemplate) Render(scope map[string]interface{}) (out []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				panic(r)
			}
		}
	}()
	SetLinkHandler(template.engine.linkHandler)
	writer := new(bytes.Buffer)
	template.lt.Render(writer, scope)
	return writer.Bytes(), nil
}
