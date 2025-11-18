package plugins

import (
	"testing"

	"github.com/osteele/gojekyll/config"
	"github.com/osteele/gojekyll/pages"
	"github.com/osteele/liquid"
	"github.com/stretchr/testify/require"
)

type relativeLinksTestSite struct {
	c     config.Config
	pages map[string]string
}

func (s relativeLinksTestSite) AddHTMLPage(string, string, pages.FrontMatter) {}
func (s relativeLinksTestSite) Config() *config.Config                        { return &s.c }
func (s relativeLinksTestSite) HasLayout(string) bool                         { return true }
func (s relativeLinksTestSite) Pages() []Page                                 { return nil }
func (s relativeLinksTestSite) Posts() []Page                                 { return nil }
func (s relativeLinksTestSite) TemplateEngine() *liquid.Engine                { return nil }
func (s relativeLinksTestSite) FilenameURLPath(path string) (string, bool) {
	if url, found := s.pages[path]; found {
		return url, true
	}
	return "", false
}

func TestRelativeLinksPlugin(t *testing.T) {
	cfg := config.Default()
	site := relativeLinksTestSite{
		c: cfg,
		pages: map[string]string{
			"about.md":                   "/about/",
			"docs/guide.md":              "/docs/guide.html",
			"docs/installation.md":       "/docs/installation/",
			"_posts/2023-01-01-hello.md": "/2023/01/01/hello.html",
		},
	}

	plugin := &jekyllRelativeLinksPlugin{}
	err := plugin.AfterInitSite(site)
	require.NoError(t, err)

	tests := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		{
			name:     "converts simple markdown link",
			input:    []byte(`<html><body><a href="about.md">About</a></body></html>`),
			expected: []byte(`<html><body><a href="/about/">About</a></body></html>`),
		},
		{
			name:     "converts markdown link with path",
			input:    []byte(`<html><body><a href="docs/guide.md">Guide</a></body></html>`),
			expected: []byte(`<html><body><a href="/docs/guide.html">Guide</a></body></html>`),
		},
		{
			name:     "preserves fragment in converted link",
			input:    []byte(`<html><body><a href="about.md#contact">Contact</a></body></html>`),
			expected: []byte(`<html><body><a href="/about/#contact">Contact</a></body></html>`),
		},
		{
			name:     "preserves query string in converted link",
			input:    []byte(`<html><body><a href="docs/guide.md?version=2">Guide v2</a></body></html>`),
			expected: []byte(`<html><body><a href="/docs/guide.html?version=2">Guide v2</a></body></html>`),
		},
		{
			name:     "preserves query string and fragment",
			input:    []byte(`<html><body><a href="about.md?foo=bar#section">About</a></body></html>`),
			expected: []byte(`<html><body><a href="/about/?foo=bar#section">About</a></body></html>`),
		},
		{
			name:     "ignores absolute URLs",
			input:    []byte(`<html><body><a href="https://example.com/page.md">External</a></body></html>`),
			expected: []byte(`<html><body><a href="https://example.com/page.md">External</a></body></html>`),
		},
		{
			name:     "ignores protocol-relative URLs",
			input:    []byte(`<html><body><a href="//example.com/page.md">External</a></body></html>`),
			expected: []byte(`<html><body><a href="//example.com/page.md">External</a></body></html>`),
		},
		{
			name:     "ignores anchor-only links",
			input:    []byte(`<html><body><a href="#section">Jump</a></body></html>`),
			expected: []byte(`<html><body><a href="#section">Jump</a></body></html>`),
		},
		{
			name:     "ignores mailto links",
			input:    []byte(`<html><body><a href="mailto:test@example.com">Email</a></body></html>`),
			expected: []byte(`<html><body><a href="mailto:test@example.com">Email</a></body></html>`),
		},
		{
			name:     "ignores non-markdown files",
			input:    []byte(`<html><body><a href="image.png">Image</a></body></html>`),
			expected: []byte(`<html><body><a href="image.png">Image</a></body></html>`),
		},
		{
			name:     "ignores non-existent markdown files",
			input:    []byte(`<html><body><a href="nonexistent.md">Missing</a></body></html>`),
			expected: []byte(`<html><body><a href="nonexistent.md">Missing</a></body></html>`),
		},
		{
			name:     "converts multiple links",
			input:    []byte(`<html><body><a href="about.md">About</a> and <a href="docs/guide.md">Guide</a></body></html>`),
			expected: []byte(`<html><body><a href="/about/">About</a> and <a href="/docs/guide.html">Guide</a></body></html>`),
		},
		{
			name:     "handles mixed link types",
			input:    []byte(`<html><body><a href="about.md">About</a> <a href="https://example.com">External</a> <a href="#top">Top</a></body></html>`),
			expected: []byte(`<html><body><a href="/about/">About</a> <a href="https://example.com">External</a> <a href="#top">Top</a></body></html>`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.PostRender(tt.input)
			require.NoError(t, err)
			require.Equal(t, string(tt.expected), string(result))
		})
	}
}

func TestIsMarkdownFile(t *testing.T) {
	cfg := config.Default()
	site := relativeLinksTestSite{
		c:     cfg,
		pages: map[string]string{},
	}

	plugin := &jekyllRelativeLinksPlugin{}
	err := plugin.AfterInitSite(site)
	require.NoError(t, err)

	tests := []struct {
		path     string
		expected bool
	}{
		{"file.md", true},
		{"file.markdown", true},
		{"file.mkd", true},
		{"file.mkdn", true},
		{"file.mkdown", true},
		{"file.MD", true},
		{"file.html", false},
		{"file.txt", false},
		{"file.png", false},
		{"file", false},
		{"file.md?query=1", true},
		{"file.md#anchor", true},
		{"file.html?query=1", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := plugin.isMarkdownFile(tt.path)
			require.Equal(t, tt.expected, result)
		})
	}
}
