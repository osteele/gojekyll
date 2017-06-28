package liquid

import (
	"github.com/osteele/liquid"
)

// LocalWrapperEngine is a wrapper around osteele/liquid.
type LocalWrapperEngine struct {
	engine               liquid.Engine
	linkHandler          LinkTagHandler
	includeTagHandler    IncludeTagHandler
	BaseURL, AbsoluteURL string
}

// NewLocalWrapperEngine creates a LocalEngine.
func NewLocalWrapperEngine() *LocalWrapperEngine {
	e := &LocalWrapperEngine{engine: liquid.NewEngine()}
	e.addJekyllFilters()
	e.addJekyllTags()
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
	return e.engine.ParseTemplate(source)
}

// ParseAndRender parses and then renders the template.
func (e *LocalWrapperEngine) ParseAndRender(source []byte, scope map[string]interface{}) ([]byte, error) {
	t, err := e.Parse(source)
	if err != nil {
		return nil, err
	}
	return t.Render(scope)
}
