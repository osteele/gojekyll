package liquid

import (
	"bytes"

	"github.com/acstech/liquid"
	"github.com/acstech/liquid/core"
)

// Template is an alias for liquid.Template.
type Template liquid.Template

// Parse is a wrapper for liquid.Parse.
func Parse(data []byte, config *core.Configuration) (*Template, error) {
	template, err := liquid.Parse(data, config)
	return (*Template)(template), err
}

// Render is a wrapper around liquid's template.Render that turns panics into errors
func Render(template *Template, variables map[string]interface{}) (bs []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				panic(r)
			}
		}
	}()
	writer := new(bytes.Buffer)
	(*liquid.Template)(template).Render(writer, variables)
	return writer.Bytes(), nil
}

// ParseAndApplyTemplate parses and then renders the template.
func ParseAndApplyTemplate(data []byte, variables map[string]interface{}, config *core.Configuration) ([]byte, error) {
	template, err := Parse(data, config)
	if err != nil {
		return nil, err
	}
	return Render(template, variables)
}
