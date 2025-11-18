package plugins

import (
	"io"
	"testing"
	"time"

	"github.com/osteele/gojekyll/config"
	"github.com/osteele/gojekyll/pages"
	"github.com/osteele/liquid"
	"github.com/stretchr/testify/require"
)

// mockSite implements the Site interface for testing
type mockSite struct {
	cfg     *config.Config
	layouts map[string]bool
}

func (m *mockSite) AddHTMLPage(url string, tpl string, fm pages.FrontMatter) {}
func (m *mockSite) Config() *config.Config {
	if m.cfg == nil {
		m.cfg = &config.Config{}
	}
	return m.cfg
}
func (m *mockSite) TemplateEngine() *liquid.Engine { return nil }
func (m *mockSite) Pages() []Page                  { return nil }
func (m *mockSite) Posts() []Page                  { return nil }
func (m *mockSite) HasLayout(name string) bool {
	if m.layouts == nil {
		return false
	}
	return m.layouts[name]
}

// mockPage implements the Page interface for testing
type mockPage struct {
	fm     pages.FrontMatter
	isPost bool
	url    string
}

func (m *mockPage) FrontMatter() pages.FrontMatter {
	if m.fm == nil {
		m.fm = make(pages.FrontMatter)
	}
	return m.fm
}
func (m *mockPage) IsPost() bool          { return m.isPost }
func (m *mockPage) URL() string           { return m.url }
func (m *mockPage) IsStatic() bool        { return false }
func (m *mockPage) Published() bool       { return true }
func (m *mockPage) Source() string        { return "" }
func (m *mockPage) OutputExt() string     { return ".html" }
func (m *mockPage) Render() error         { return nil }
func (m *mockPage) SetContent(string)     {}
func (m *mockPage) PostDate() time.Time   { return time.Time{} }
func (m *mockPage) Categories() []string  { return nil }
func (m *mockPage) Tags() []string        { return nil }
func (m *mockPage) Write(io.Writer) error { return nil }
func (m *mockPage) Reload() error         { return nil }

func TestDefaultLayout_layoutNames(t *testing.T) {
	plugin := jekyllDefaultLayout{}

	tests := []struct {
		name             string
		availableLayouts []string
		expectedMap      map[string]string
	}{
		{
			name:             "only default layout",
			availableLayouts: []string{"default"},
			expectedMap: map[string]string{
				"default": "default",
				"post":    "default",
				"page":    "default",
				"home":    "default",
			},
		},
		{
			name:             "default and post layouts",
			availableLayouts: []string{"default", "post"},
			expectedMap: map[string]string{
				"default": "default",
				"post":    "post",
				"page":    "post",
				"home":    "post",
			},
		},
		{
			name:             "default, post, and page layouts",
			availableLayouts: []string{"default", "post", "page"},
			expectedMap: map[string]string{
				"default": "default",
				"post":    "post",
				"page":    "page",
				"home":    "page",
			},
		},
		{
			name:             "all layouts available",
			availableLayouts: []string{"default", "post", "page", "home"},
			expectedMap: map[string]string{
				"default": "default",
				"post":    "post",
				"page":    "page",
				"home":    "home",
			},
		},
		{
			name:             "only page layout (no default)",
			availableLayouts: []string{"page"},
			expectedMap: map[string]string{
				"default": "",
				"post":    "",
				"page":    "page",
				"home":    "page",
			},
		},
		{
			name:             "no layouts available",
			availableLayouts: []string{},
			expectedMap: map[string]string{
				"default": "",
				"post":    "",
				"page":    "",
				"home":    "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			site := &mockSite{
				layouts: make(map[string]bool),
			}
			for _, layout := range tt.availableLayouts {
				site.layouts[layout] = true
			}

			result := plugin.layoutNames(site)
			require.Equal(t, tt.expectedMap, result)
		})
	}
}

func TestDefaultLayout_PostInitPage(t *testing.T) {
	plugin := jekyllDefaultLayout{}

	tests := []struct {
		name             string
		page             *mockPage
		availableLayouts []string
		expectedLayout   interface{}
		description      string
	}{
		{
			name: "post without layout gets post layout",
			page: &mockPage{
				isPost: true,
				url:    "/2023/01/01/test-post.html",
			},
			availableLayouts: []string{"default", "post"},
			expectedLayout:   "post",
			description:      "Post should get 'post' layout",
		},
		{
			name: "post without layout falls back to default",
			page: &mockPage{
				isPost: true,
				url:    "/2023/01/01/test-post.html",
			},
			availableLayouts: []string{"default"},
			expectedLayout:   "default",
			description:      "Post should get 'default' layout when 'post' is not available",
		},
		{
			name: "home page without layout gets home layout",
			page: &mockPage{
				isPost: false,
				url:    "/",
			},
			availableLayouts: []string{"default", "home"},
			expectedLayout:   "home",
			description:      "Home page should get 'home' layout",
		},
		{
			name: "home page falls back to page layout",
			page: &mockPage{
				isPost: false,
				url:    "/",
			},
			availableLayouts: []string{"default", "page"},
			expectedLayout:   "page",
			description:      "Home page should get 'page' layout when 'home' is not available",
		},
		{
			name: "home page falls back to default layout",
			page: &mockPage{
				isPost: false,
				url:    "/",
			},
			availableLayouts: []string{"default"},
			expectedLayout:   "default",
			description:      "Home page should get 'default' layout when neither 'home' nor 'page' are available",
		},
		{
			name: "regular page without layout gets page layout",
			page: &mockPage{
				isPost: false,
				url:    "/about.html",
			},
			availableLayouts: []string{"default", "page"},
			expectedLayout:   "page",
			description:      "Regular page should get 'page' layout",
		},
		{
			name: "regular page falls back to default layout",
			page: &mockPage{
				isPost: false,
				url:    "/about.html",
			},
			availableLayouts: []string{"default"},
			expectedLayout:   "default",
			description:      "Regular page should get 'default' layout when 'page' is not available",
		},
		{
			name: "page with explicit layout is not modified",
			page: &mockPage{
				isPost: false,
				url:    "/about.html",
				fm:     pages.FrontMatter{"layout": "custom"},
			},
			availableLayouts: []string{"default", "page", "custom"},
			expectedLayout:   "custom",
			description:      "Page with explicit layout should keep its layout",
		},
		{
			name: "page with layout 'none' is not modified",
			page: &mockPage{
				isPost: false,
				url:    "/about.html",
				fm:     pages.FrontMatter{"layout": "none"},
			},
			availableLayouts: []string{"default", "page"},
			expectedLayout:   "none",
			description:      "Page with layout 'none' should keep it",
		},
		{
			name: "page with layout 'null' is not modified",
			page: &mockPage{
				isPost: false,
				url:    "/about.html",
				fm:     pages.FrontMatter{"layout": "null"},
			},
			availableLayouts: []string{"default", "page"},
			expectedLayout:   "null",
			description:      "Page with layout 'null' should keep it",
		},
		{
			name: "page with empty string layout is not modified",
			page: &mockPage{
				isPost: false,
				url:    "/about.html",
				fm:     pages.FrontMatter{"layout": ""},
			},
			availableLayouts: []string{"default", "page"},
			expectedLayout:   "",
			description:      "Page with empty string layout should keep it",
		},
		{
			name: "page without layout and no layouts available",
			page: &mockPage{
				isPost: false,
				url:    "/about.html",
			},
			availableLayouts: []string{},
			expectedLayout:   nil,
			description:      "Page should have no layout when none are available",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			site := &mockSite{
				layouts: make(map[string]bool),
			}
			for _, layout := range tt.availableLayouts {
				site.layouts[layout] = true
			}

			err := plugin.PostInitPage(site, tt.page)
			require.NoError(t, err)
			require.Equal(t, tt.expectedLayout, tt.page.FrontMatter()["layout"], tt.description)
		})
	}
}

func TestDefaultLayout_Integration(t *testing.T) {
	// Test that the plugin is properly registered
	plugin, found := Lookup("jekyll-default-layout")
	require.True(t, found, "jekyll-default-layout should be registered")
	require.NotNil(t, plugin, "plugin should not be nil")

	// Test that it implements the Plugin interface correctly
	_, ok := plugin.(jekyllDefaultLayout)
	require.True(t, ok, "plugin should be of type jekyllDefaultLayout")
}
