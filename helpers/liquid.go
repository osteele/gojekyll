package helpers

import (
	"bytes"

	"github.com/acstech/liquid"
	"github.com/acstech/liquid/core"
)

// RenderTemplate is a wrapper around liquid template.Render that turns panics into errors
func RenderTemplate(template *liquid.Template, variables map[string]interface{}) (bs []byte, err error) {
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
	template.Render(writer, variables)
	return writer.Bytes(), nil
}

// ParseAndApplyTemplate parses and then renders the template.
func ParseAndApplyTemplate(bs []byte, variables map[string]interface{}, config *core.Configuration) ([]byte, error) {
	template, err := liquid.Parse(bs, config)
	if err != nil {
		return nil, err
	}
	return RenderTemplate(template, variables)
}
