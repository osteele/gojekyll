package pages

import (
	"bytes"
	"io"
	"testing"

	"github.com/osteele/gojekyll/pipelines"
	"github.com/stretchr/testify/require"
)

type mockRenderingContext struct{ t *testing.T }

func (c mockRenderingContext) RenderingPipeline() pipelines.PipelineInterface { return c }
func (c mockRenderingContext) OutputExt(string) string                        { return ".html" }
func (c mockRenderingContext) Site() interface{}                              { return nil }
func (c mockRenderingContext) ApplyLayout(layout string, src []byte, vars map[string]interface{}) ([]byte, error) {
	require.Equal(c.t, "layout1", layout)
	return nil, nil
}
func (c mockRenderingContext) Render(w io.Writer, src []byte, filename string, vars map[string]interface{}) ([]byte, error) {
	require.Equal(c.t, "testdata/page_with_layout.md", filename)
	return nil, nil
}

func TestPageWrite(t *testing.T) {
	p, err := NewFile("testdata/page_with_layout.md", containerMock{}, "page_with_layout.md", map[string]interface{}{})
	require.NoError(t, err)
	require.NotNil(t, p)
	buf := new(bytes.Buffer)
	require.NoError(t, p.Write(buf, mockRenderingContext{t}))
}
