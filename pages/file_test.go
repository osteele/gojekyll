package pages

import (
	"io"
	"path/filepath"
	"testing"

	"github.com/osteele/gojekyll/config"
	"github.com/osteele/gojekyll/pipelines"
	"github.com/osteele/liquid"
	"github.com/stretchr/testify/require"
)

type siteFake struct {
	t   *testing.T
	cfg config.Config
}

func (s siteFake) Config() *config.Config                         { return &s.cfg }
func (s siteFake) RelativePath(p string) string                   { return p }
func (s siteFake) RenderingPipeline() pipelines.PipelineInterface { return &pipelineFake{s.t} }
func (s siteFake) OutputExt(p string) string                      { return filepath.Ext(p) }

type pipelineFake struct{ t *testing.T }

func (p pipelineFake) OutputExt(string) string { return ".html" }
func (p pipelineFake) ApplyLayout(layout string, src []byte, vars liquid.Bindings) ([]byte, error) {
	require.Equal(p.t, "layout1", layout)
	return nil, nil
}
func (p pipelineFake) Render(w io.Writer, src []byte, vars liquid.Bindings, filename string, lineNo int) error {
	_, err := io.WriteString(w, "rendered: ")
	if err != nil {
		return err
	}
	_, err = w.Write(src)
	return err
}

func (p pipelineFake) RenderTemplate(src []byte, vars liquid.Bindings, filename string, lineNo int) ([]byte, error) {
	return append([]byte("rendered: "), src...), nil
}
