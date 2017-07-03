package collections

import (
	"testing"

	"github.com/osteele/gojekyll/config"
	"github.com/stretchr/testify/require"
)

type siteMock struct{ c config.Config }

func (c siteMock) Config() *config.Config    { return &c.c }
func (c siteMock) OutputExt(s string) string { return "" }
func (c siteMock) Site() interface{}         { return c }

func TestNewCollection(t *testing.T) {
	site := siteMock{config.Default()}

	c1 := NewCollection(site, "c", map[string]interface{}{"output": true})
	require.Equal(t, true, c1.Output())
	require.Equal(t, "_c/", c1.PathPrefix())

	c2 := NewCollection(site, "c", map[string]interface{}{})
	require.Equal(t, false, c2.Output())
}

func TestPermalinkPattern(t *testing.T) {
	site := siteMock{config.Default()}

	c1 := NewCollection(site, "c", map[string]interface{}{})
	require.Contains(t, c1.PermalinkPattern(), ":collection")

	c2 := NewCollection(site, "c", map[string]interface{}{"permalink": "out"})
	require.Equal(t, "out", c2.PermalinkPattern())

	c3 := NewCollection(site, "posts", map[string]interface{}{})
	require.Contains(t, c3.PermalinkPattern(), "/:year/:month/:day/:title")
}

func TestReadPosts(t *testing.T) {
	site := siteMock{config.FromString("source: testdata")}
	c := NewCollection(site, "posts", map[string]interface{}{})
	c.ReadPages()
	require.Len(t, c.Pages(), 1)

	site = siteMock{config.FromString("source: testdata\nunpublished: true")}
	c = NewCollection(site, "posts", map[string]interface{}{})
	c.ReadPages()
	require.Len(t, c.Pages(), 2)

	site = siteMock{config.FromString("source: testdata\nfuture: true")}
	c = NewCollection(site, "posts", map[string]interface{}{})
	c.ReadPages()
	require.Len(t, c.Pages(), 2)
}
