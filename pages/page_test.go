package pages

import (
	"bytes"
	"testing"

	"github.com/osteele/gojekyll/config"
	"github.com/stretchr/testify/require"
)

// func (c pipelineFake) RenderingPipeline() pipelines.PipelineInterface { return c }
// func (c pipelineFake) Config() config.Config                          { return c.cfg }
// func (c pipelineFake) PathPrefix() string                             { return "." }
// func (c pipelineFake) Site() interface{}                              { return nil }

func TestPageWrite(t *testing.T) {
	cfg := config.Default()
	p, err := NewFile(siteFake{t, cfg}, "testdata/page_with_layout.md", "page_with_layout.md", map[string]interface{}{})
	require.NoError(t, err)
	require.NotNil(t, p)
	buf := new(bytes.Buffer)
	require.NoError(t, p.Write(buf))
}
