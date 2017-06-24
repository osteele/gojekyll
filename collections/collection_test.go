package collections

import (
	"path"
	"testing"

	"github.com/osteele/gojekyll/templates"
	"github.com/stretchr/testify/require"
)

var tests = []struct{ in, out string }{
	{"pre</head>post", "pre:insertion:</head>post"},
	{"pre:insertion:</head>post", "pre:insertion:</head>post"},
	{"post", ":insertion:post"},
}

type MockContainer struct{}

func (c MockContainer) PathPrefix() string               { return "" }
func (c MockContainer) OutputExt(filename string) string { return path.Ext(filename) }

// func (c MockPipeline) Render(_ io.Writer, _ []byte, _ string, _ templates.VariableMap) ([]byte, error) {
// 	return nil, fmt.Errorf("unimplemented")
// }

// func (c MockPipeline) ApplyLayout(_ string, _ []byte, _ templates.VariableMap) ([]byte, error) {
// 	return nil, fmt.Errorf("unimplemented")
// }

// func (c MockPipeline) SiteVariables() templates.VariableMap { return templates.VariableMap{} }

func TestNewCollection(t *testing.T) {
	ctx := MockContainer{}

	c1 := NewCollection("c", templates.VariableMap{"output": true}, ctx)
	require.Equal(t, true, c1.Output())
	require.Equal(t, "_c/", c1.PathPrefix())

	c2 := NewCollection("c", templates.VariableMap{}, ctx)
	require.Equal(t, false, c2.Output())
}

func TestPermalinkPattern(t *testing.T) {
	ctx := MockContainer{}

	c1 := NewCollection("c", templates.VariableMap{}, ctx)
	require.Contains(t, c1.PermalinkPattern(), ":collection")

	c2 := NewCollection("c", templates.VariableMap{"permalink": "out"}, ctx)
	require.Equal(t, "out", c2.PermalinkPattern())

	c3 := NewCollection("posts", templates.VariableMap{}, ctx)
	require.Contains(t, c3.PermalinkPattern(), "/:year/:month/:day/:title")
}
