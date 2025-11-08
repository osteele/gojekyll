package site

import (
	"io"
	"testing"

	"github.com/osteele/gojekyll/config"
	"github.com/stretchr/testify/require"
)

// mockDocument implements Document interface for testing
type mockDocument struct {
	url       string
	published bool
}

func (d *mockDocument) URL() string           { return d.url }
func (d *mockDocument) Source() string        { return "" }
func (d *mockDocument) OutputExt() string     { return ".html" }
func (d *mockDocument) Published() bool       { return d.published }
func (d *mockDocument) IsStatic() bool        { return false }
func (d *mockDocument) Write(io.Writer) error { return nil }
func (d *mockDocument) Reload() error         { return nil }

func TestURLPageTrailingSlash(t *testing.T) {
	s := New(config.Flags{})
	s.Routes = make(map[string]Document)

	// Test trailing slash handling for directory-style permalinks
	aboutPage := &mockDocument{url: "/about/", published: true}
	s.Routes["/about/"] = aboutPage

	// Should find page with exact match
	p, found := s.URLPage("/about/")
	require.True(t, found, "Should find page with exact match including trailing slash")
	require.Equal(t, aboutPage, p)

	// Should find page when trailing slash is missing
	p, found = s.URLPage("/about")
	require.True(t, found, "Should find page without trailing slash when registered with trailing slash")
	require.Equal(t, aboutPage, p)

	// Test index.html fallback
	indexPage := &mockDocument{url: "/posts/index.html", published: true}
	s.Routes["/posts/index.html"] = indexPage

	p, found = s.URLPage("/posts")
	require.True(t, found, "Should find index.html when accessing directory")
	require.Equal(t, indexPage, p)

	// Test .html extension fallback
	contactPage := &mockDocument{url: "/contact.html", published: true}
	s.Routes["/contact.html"] = contactPage

	p, found = s.URLPage("/contact")
	require.True(t, found, "Should find .html file without extension")
	require.Equal(t, contactPage, p)

	// Test non-existent page
	_, found = s.URLPage("/nonexistent")
	require.False(t, found, "Should not find non-existent page")
}
