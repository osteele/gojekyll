package collections

import (
	"fmt"
	"io"
	"testing"

	"github.com/osteele/gojekyll/liquid"
	"github.com/osteele/gojekyll/templates"
	"github.com/stretchr/testify/require"
)

var tests = []struct{ in, out string }{
	{"pre</head>post", "pre:insertion:</head>post"},
	{"pre:insertion:</head>post", "pre:insertion:</head>post"},
	{"post", ":insertion:post"},
}

type MockContext struct{}

func (c MockContext) FindLayout(_ string, _ *templates.VariableMap) (liquid.Template, error) {
	return nil, fmt.Errorf("unimplemented")
}
func (c MockContext) IsMarkdown(_ string) bool             { return true }
func (c MockContext) IsSassPath(_ string) bool             { return true }
func (c MockContext) SassIncludePaths() []string           { return []string{} }
func (c MockContext) SiteVariables() templates.VariableMap { return templates.VariableMap{} }
func (c MockContext) TemplateEngine() liquid.Engine        { return nil }
func (c MockContext) WriteSass(io.Writer, []byte) error    { return nil }

func TestNewCollection(t *testing.T) {
	ctx := MockContext{}

	c1 := NewCollection(ctx, "c", templates.VariableMap{"output": true})
	require.Equal(t, true, c1.Output())
	require.Equal(t, "_c/", c1.PathPrefix())

	c2 := NewCollection(ctx, "c", templates.VariableMap{})
	require.Equal(t, false, c2.Output())
}

func TestPermalinkPattern(t *testing.T) {
	ctx := MockContext{}

	c1 := NewCollection(ctx, "c", templates.VariableMap{})
	require.Contains(t, c1.PermalinkPattern(), ":collection")

	c2 := NewCollection(ctx, "c", templates.VariableMap{"permalink": "out"})
	require.Equal(t, "out", c2.PermalinkPattern())

	c3 := NewCollection(ctx, "posts", templates.VariableMap{})
	require.Contains(t, c3.PermalinkPattern(), "/:year/:month/:day/:title")
}
