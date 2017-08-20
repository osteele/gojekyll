package pages

import (
	"io"
	"testing"

	"github.com/osteele/gojekyll/config"
	"github.com/osteele/gojekyll/renderers"
	"github.com/osteele/liquid"
	"github.com/stretchr/testify/require"
)

type siteFake struct {
	t   *testing.T
	cfg config.Config
}

func (s siteFake) Config() *config.Config               { return &s.cfg }
func (s siteFake) RelativePath(p string) string         { return p }
func (s siteFake) RendererManager() renderers.Renderers { return &renderManagerFake{s.t} }

type renderManagerFake struct{ t *testing.T }

func (rm renderManagerFake) ApplyLayout(layout string, src []byte, vars liquid.Bindings) ([]byte, error) {
	require.Equal(rm.t, "layout1", layout)
	return nil, nil
}
func (rm renderManagerFake) Render(w io.Writer, src []byte, vars liquid.Bindings, filename string, lineNo int) error {
	_, err := io.WriteString(w, "rendered: ")
	if err != nil {
		return err
	}
	_, err = w.Write(src)
	return err
}

func (rm renderManagerFake) RenderTemplate(src []byte, vars liquid.Bindings, filename string, lineNo int) ([]byte, error) {
	return append([]byte("rendered: "), src...), nil
}
