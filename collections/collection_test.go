package collections

import (
	"path"
	"testing"

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

// func (c MockPipeline) Render(_ io.Writer, _ []byte, _ string, _ map[string]interface{}) ([]byte, error) {
// 	return nil, fmt.Errorf("unimplemented")
// }

// func (c MockPipeline) ApplyLayout(_ string, _ []byte, _ map[string]interface{}) ([]byte, error) {
// 	return nil, fmt.Errorf("unimplemented")
// }

// func (c MockPipeline) SiteVariables() map[string]interface{} { return map[string]interface{}{} }

func TestNewCollection(t *testing.T) {
	ctx := MockContainer{}

	c1 := NewCollection("c", map[string]interface{}{"output": true}, ctx)
	require.Equal(t, true, c1.Output())
	require.Equal(t, "_c/", c1.PathPrefix())

	c2 := NewCollection("c", map[string]interface{}{}, ctx)
	require.Equal(t, false, c2.Output())
}

func TestPermalinkPattern(t *testing.T) {
	ctx := MockContainer{}

	c1 := NewCollection("c", map[string]interface{}{}, ctx)
	require.Contains(t, c1.PermalinkPattern(), ":collection")

	c2 := NewCollection("c", map[string]interface{}{"permalink": "out"}, ctx)
	require.Equal(t, "out", c2.PermalinkPattern())

	c3 := NewCollection("posts", map[string]interface{}{}, ctx)
	require.Contains(t, c3.PermalinkPattern(), "/:year/:month/:day/:title")
}
