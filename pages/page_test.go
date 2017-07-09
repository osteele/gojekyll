package pages

import (
	"bytes"
	"io"
	"testing"

	"github.com/osteele/gojekyll/config"
	"github.com/osteele/gojekyll/pipelines"
	"github.com/stretchr/testify/require"
)

type renderingContextFake struct {
	t   *testing.T
	cfg config.Config
}

func (c renderingContextFake) RenderingPipeline() pipelines.PipelineInterface { return c }
func (c renderingContextFake) Config() config.Config                          { return c.cfg }
func (c renderingContextFake) PathPrefix() string                             { return "." }
func (c renderingContextFake) OutputExt(string) string                        { return ".html" }
func (c renderingContextFake) Site() interface{}                              { return nil }
func (c renderingContextFake) ApplyLayout(layout string, src []byte, vars map[string]interface{}) ([]byte, error) {
	require.Equal(c.t, "layout1", layout)
	return nil, nil
}
func (c renderingContextFake) Render(w io.Writer, src []byte, filename string, lineNo int, vars map[string]interface{}) ([]byte, error) {
	require.Equal(c.t, "testdata/page_with_layout.md", filename)
	return nil, nil
}

func TestPageWrite(t *testing.T) {
	cfg := config.Default()
	p, err := NewFile("testdata/page_with_layout.md", containerFake{cfg, ""}, "page_with_layout.md", map[string]interface{}{})
	require.NoError(t, err)
	require.NotNil(t, p)
	buf := new(bytes.Buffer)
	require.NoError(t, p.Write(buf, renderingContextFake{t, cfg}))
}
