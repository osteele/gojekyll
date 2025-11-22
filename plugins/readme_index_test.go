package plugins

import (
	"testing"

	"github.com/osteele/gojekyll/pages"
	"github.com/stretchr/testify/require"
)

func TestIsReadmePage(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		expected bool
	}{
		{"root README.md", "/path/to/site/README.md", true},
		{"root README.markdown", "/path/to/site/README.markdown", true},
		{"nested README.md", "/path/to/site/foo/README.md", true},
		{"nested README.markdown", "/path/to/site/foo/bar/README.markdown", true},
		{"case insensitive", "/path/to/site/readme.md", true},
		{"not a README", "/path/to/site/index.md", false},
		{"empty source", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock page with source in frontmatter
			fm := make(pages.FrontMatter)
			page := &mockPage{fm: fm}
			// Override the Source method by creating a custom page wrapper
			testPage := &testPageWithSource{
				Page:   page,
				source: tt.source,
			}
			result := isReadmePage(testPage)
			require.Equal(t, tt.expected, result)
		})
	}
}

// testPageWithSource wraps mockPage to provide a custom Source method
type testPageWithSource struct {
	Page
	source string
}

func (p *testPageWithSource) Source() string {
	return p.source
}

func TestCalculateIndexURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{"root README", "/README.html", "/"},
		{"one level deep", "/foo/README.html", "/foo/"},
		{"two levels deep", "/foo/bar/README.html", "/foo/bar/"},
		{"already has slash", "/foo/", "/foo/"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateIndexURL(tt.url)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestReadmeIndexPlugin_PostInitPage(t *testing.T) {
	plugin := jekyllReadmeIndexPlugin{}
	site := &mockSite{}

	tests := []struct {
		name            string
		source          string
		url             string
		expectedChanged bool
		expectedURL     string
	}{
		{
			name:            "root README",
			source:          "/path/to/site/README.md",
			url:             "/README.html",
			expectedChanged: true,
			expectedURL:     "/",
		},
		{
			name:            "nested README",
			source:          "/path/to/site/foo/README.md",
			url:             "/foo/README.html",
			expectedChanged: true,
			expectedURL:     "/foo/",
		},
		{
			name:            "non-README page",
			source:          "/path/to/site/index.md",
			url:             "/index.html",
			expectedChanged: false,
			expectedURL:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fm := make(pages.FrontMatter)
			page := &mockPage{
				url: tt.url,
				fm:  fm,
			}
			testPage := &testPageWithSource{
				Page:   page,
				source: tt.source,
			}

			err := plugin.PostInitPage(site, testPage)
			require.NoError(t, err)

			if tt.expectedChanged {
				permalink, hasPermalink := fm["permalink"]
				require.True(t, hasPermalink, "permalink should be set in frontmatter")
				require.Equal(t, tt.expectedURL, permalink)
			} else {
				_, hasPermalink := fm["permalink"]
				require.False(t, hasPermalink, "permalink should not be set for non-README pages")
			}
		})
	}
}
