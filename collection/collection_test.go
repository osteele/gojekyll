package collection

import (
	"testing"

	"github.com/osteele/gojekyll/config"
	"github.com/osteele/gojekyll/pipelines"
	"github.com/stretchr/testify/require"
)

type siteFake struct{ c config.Config }

func (s siteFake) Config() *config.Config                         { return &s.c }
func (s siteFake) Exclude(string) bool                            { return false }
func (s siteFake) OutputExt(string) string                        { return "" }
func (s siteFake) RelativePath(string) string                     { panic("unimplemented") }
func (s siteFake) RenderingPipeline() pipelines.PipelineInterface { panic("unimplemented") }

func TestNewCollection(t *testing.T) {
	site := siteFake{config.Default()}

	c1 := New(site, "c", map[string]interface{}{"output": true})
	require.Equal(t, true, c1.Output())
	require.Equal(t, "_c/", c1.PathPrefix())

	c2 := New(site, "c", map[string]interface{}{})
	require.Equal(t, false, c2.Output())
}

func TestPermalinkPattern(t *testing.T) {
	site := siteFake{config.Default()}

	c1 := New(site, "c", map[string]interface{}{})
	require.Contains(t, c1.PermalinkPattern(), ":collection")

	c2 := New(site, "c", map[string]interface{}{"permalink": "out"})
	require.Equal(t, "out", c2.PermalinkPattern())

	c3 := New(site, "posts", map[string]interface{}{})
	require.Contains(t, c3.PermalinkPattern(), "/:year/:month/:day/:title")
}

func Test_ReadPages(t *testing.T) {
	site := siteFake{config.FromString("source: testdata")}
	c := New(site, "posts", map[string]interface{}{})
	require.NoError(t, c.ReadPages())
	require.Len(t, c.Pages(), 1)

	site = siteFake{config.FromString("source: testdata\nunpublished: true")}
	c = New(site, "posts", map[string]interface{}{})
	require.NoError(t, c.ReadPages())
	require.Len(t, c.Pages(), 2)

	site = siteFake{config.FromString("source: testdata\nfuture: true")}
	c = New(site, "posts", map[string]interface{}{})
	require.NoError(t, c.ReadPages())
	require.Len(t, c.Pages(), 2)

	pages := c.Pages()
	require.Equal(t, nil, pages[0].FrontMatter()["previous"])
	require.Equal(t, pages[1], pages[0].FrontMatter()["next"])
	require.Equal(t, pages[0], pages[1].FrontMatter()["previous"])
	require.Equal(t, nil, pages[1].FrontMatter()["next"])
}
