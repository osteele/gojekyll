package pages

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/osteele/gojekyll/config"
	"github.com/osteele/gojekyll/frontmatter"
	"github.com/stretchr/testify/require"
)

type pathTest struct{ path, pattern, out string }

// Non-date-dependent tests
var staticTests = []pathTest{
	{"/a/b/base.html", "/out:output_ext", "/out.html"},
	{"/a/b/base.md", "/out:output_ext", "/out.html"},
	{"/a/b/base.markdown", "/out:output_ext", "/out.html"},
	{"/a/b/base.html", "/:path/out:output_ext", "/a/b/base/out.html"},
	{"/a/b/base.html", "/prefix/:name", "/prefix/base"},
	{"/a/b/base.html", "/prefix/:path/post", "/prefix/a/b/base/post"},
	{"/a/b/base.html", "/prefix/:title", "/prefix/base"},
	{"/a/b/base.html", "/prefix/:slug", "/prefix/base"},
	{"base", "/:categories/:name:output_ext", "/base"}, // categories ignored for non-posts
	{"base", "none", "/base.html"},                     // categories ignored for non-posts
}

// Date-dependent tests will be generated dynamically based on the test date
// This approach allows tests to pass in any time zone while we investigate the proper
// time zone handling in permalinks. See: https://github.com/osteele/gojekyll/issues/63

var collectionTests = []pathTest{
	{"/a/b/c.d", "/prefix/:collection/post", "/prefix/c/post"},
	{"/a/b/c.d", "/prefix:path/post", "/prefix/a/b/c/post"},
}

func TestExpandPermalinkPattern(t *testing.T) {
	var (
		s = siteFake{t, config.Default()}
		d = map[string]interface{}{
			"categories": "b a",
		}
	)

	// Create a test date in UTC - this is the reference date for all tests
	testDate, err := time.Parse(time.RFC3339, "2006-02-03T15:04:05Z")
	require.NoError(t, err)

	testPermalinkPattern := func(pattern, path string, data map[string]interface{}) (string, error) {
		fm := frontmatter.Merge(data, FrontMatter{"permalink": pattern})
		ext := filepath.Ext(path)
		switch ext {
		case ".md", ".markdown":
			ext = ".html"
		}
		f := file{site: s, relPath: path, fm: fm, outputExt: ext}
		p := page{file: f}
		// Use the same test date that we use for generating expectations
		p.modTime = testDate
		return p.computePermalink(p.permalinkVariables())
	}

	runTests := func(tests []pathTest) {
		for i, test := range tests {
			t.Run(test.pattern, func(t *testing.T) {
				p, err := testPermalinkPattern(test.pattern, test.path, d)
				require.NoError(t, err)
				require.Equalf(t, test.out, p, "%d: pattern=%s", i+1, test.pattern)
			})
		}
	}

	// Generate date-dependent tests with expected values
	// NOTE: These are pages (not posts), so date/category placeholders are ignored per Jekyll behavior
	dateTests := []pathTest{
		{"base", "date", "/base.html"},    // dates/categories ignored for non-posts
		{"base", "pretty", "/base/"},      // dates/categories ignored for non-posts
		{"base", "ordinal", "/base.html"}, // dates/categories ignored for non-posts
	}

	// Run the non-date-dependent tests
	runTests(staticTests)

	// Run the date-dependent tests
	runTests(dateTests)

	s = siteFake{t, config.Default()}
	d["collection"] = "c"
	runTests(collectionTests)

	t.Run("invalid template variable", func(t *testing.T) {
		p, err := testPermalinkPattern("/:invalid", "/a/b/base.html", d)
		require.Error(t, err)
		require.Zero(t, p)
	})
}

func TestPostPermalinkPatterns(t *testing.T) {
	// Test that posts correctly use date and category placeholders
	var (
		s = siteFake{t, config.Default()}
		d = map[string]interface{}{
			"categories": "blog tech",
			"collection": "posts", // Mark as post
			"title":      "My Post",
		}
	)

	testDate, err := time.Parse(time.RFC3339, "2006-02-03T15:04:05Z")
	require.NoError(t, err)
	localDate := testDate.In(time.Local)

	testPermalinkPattern := func(pattern, path string, data map[string]interface{}) (string, error) {
		fm := frontmatter.Merge(data, FrontMatter{"permalink": pattern})
		ext := filepath.Ext(path)
		switch ext {
		case ".md", ".markdown":
			ext = ".html"
		}
		f := file{site: s, relPath: path, fm: fm, outputExt: ext}
		p := page{file: f}
		p.modTime = testDate
		return p.computePermalink(p.permalinkVariables())
	}

	tests := []struct {
		pattern  string
		expected string
	}{
		{"date", fmt.Sprintf("/blog/tech/%04d/%02d/%02d/my-post.html", localDate.Year(), localDate.Month(), localDate.Day())},
		{"pretty", fmt.Sprintf("/blog/tech/%04d/%02d/%02d/my-post/", localDate.Year(), localDate.Month(), localDate.Day())},
		{"ordinal", fmt.Sprintf("/blog/tech/%04d/%d/my-post.html", testDate.Year(), testDate.YearDay())},
		{"none", "/blog/tech/my-post.html"},
		{"/:categories/:year/:month/:title/", fmt.Sprintf("/blog/tech/%04d/%02d/my-post/", localDate.Year(), localDate.Month())},
	}

	for _, test := range tests {
		t.Run(test.pattern, func(t *testing.T) {
			p, err := testPermalinkPattern(test.pattern, "/_posts/2006-02-03-my-post.md", d)
			require.NoError(t, err)
			require.Equal(t, test.expected, p)
		})
	}
}

func TestPagePermalinkEdgeCases(t *testing.T) {
	// Test edge cases for non-post permalink handling
	var (
		s = siteFake{t, config.Default()}
		d = map[string]interface{}{
			"title": "Test Page",
		}
	)

	testDate, err := time.Parse(time.RFC3339, "2006-02-03T15:04:05Z")
	require.NoError(t, err)

	testPermalinkPattern := func(pattern, path string, data map[string]interface{}) (string, error) {
		fm := frontmatter.Merge(data, FrontMatter{"permalink": pattern})
		ext := filepath.Ext(path)
		switch ext {
		case ".md", ".markdown":
			ext = ".html"
		}
		f := file{site: s, relPath: path, fm: fm, outputExt: ext}
		p := page{file: f}
		p.modTime = testDate
		return p.computePermalink(p.permalinkVariables())
	}

	tests := []struct {
		name     string
		pattern  string
		path     string
		expected string
	}{
		// Complex patterns with multiple placeholders
		{"complex with categories", "/:categories/:year/:month/:day/:title/", "/test.md", "/test-page/"},
		{"categories at end", "/blog/:categories", "/test.md", "/blog"},
		{"categories in middle", "/prefix/:categories/suffix/:title", "/test.md", "/prefix/suffix/test-page"},
		
		// Date placeholders in various positions
		{"year only", "/:year/:title", "/test.md", "/test-page"},
		{"date at end", "/blog/:title/:year/:month/:day", "/test.md", "/blog/test-page"},
		{"mixed dates", "/:i_month/:short_year/:title/:y_day", "/test.md", "/test-page"},
		
		// Edge cases for cleanup
		{"multiple slashes", "/:categories//:year///:title", "/test.md", "/test-page"},
		{"trailing dates", "/blog/:title/:year/", "/test.md", "/blog/test-page"},
		
		// Patterns that become empty or minimal
		{"only categories", ":categories", "/test.md", "/test-page"},
		{"only dates", ":year/:month/:day", "/test.md", "/test-page"},
		{"dates and categories", ":categories/:year/:month/:day", "/test.md", "/test-page"},
		
		// Edge case specifically mentioned in PR review
		{"categories with colon after", ":categories:slug", "/test.md", "/test"},
		{"categories with multiple colons", "/prefix:categories:year:title", "/test.md", "/prefixtest-page"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p, err := testPermalinkPattern(test.pattern, test.path, d)
			require.NoError(t, err)
			require.Equal(t, test.expected, p, "Pattern: %s", test.pattern)
		})
	}
}

func TestGlobalPermalinkConfiguration(t *testing.T) {
	testDate, err := time.Parse(time.RFC3339, "2006-02-03T15:04:05Z")
	require.NoError(t, err)
	localDate := testDate.In(time.Local)

	tests := []struct {
		name            string
		globalPermalink string
		pagePath        string
		frontMatter     map[string]interface{}
		expected        string
	}{
		{
			name:            "pretty permalink for regular page",
			globalPermalink: "pretty",
			pagePath:        "/bread.html",
			frontMatter:     map[string]interface{}{"title": "Bread Page"},
			expected:        "/bread-page/", // Jekyll ignores dates/categories for pages
		},
		{
			name:            "date permalink for regular page",
			globalPermalink: "date",
			pagePath:        "/about.html",
			frontMatter:     map[string]interface{}{"title": "About"},
			expected:        "/about.html", // Date placeholders ignored for pages
		},
		{
			name:            "none permalink for regular page",
			globalPermalink: "none",
			pagePath:        "/contact.html",
			frontMatter:     map[string]interface{}{"title": "Contact"},
			expected:        "/contact.html",
		},
		{
			name:            "pretty permalink for post",
			globalPermalink: "pretty",
			pagePath:        "/_posts/2006-02-03-hello.html",
			frontMatter:     map[string]interface{}{"title": "Hello World", "collection": "posts"},
			expected:        fmt.Sprintf("/%04d/%02d/%02d/hello-world/", localDate.Year(), localDate.Month(), localDate.Day()),
		},
		{
			name:            "date permalink for post",
			globalPermalink: "date",
			pagePath:        "/_posts/2006-02-03-hello.html",
			frontMatter:     map[string]interface{}{"title": "Hello World", "collection": "posts"},
			expected:        fmt.Sprintf("/%04d/%02d/%02d/hello-world.html", localDate.Year(), localDate.Month(), localDate.Day()),
		},
		{
			name:            "front matter overrides global",
			globalPermalink: "pretty",
			pagePath:        "/special.html",
			frontMatter:     map[string]interface{}{"permalink": "/custom/path/"},
			expected:        "/custom/path/",
		},
		{
			name:            "no global permalink uses default",
			globalPermalink: "",
			pagePath:        "/default.html",
			frontMatter:     map[string]interface{}{"title": "Default"},
			expected:        "/default.html",
		},
		{
			name:            "collection document with pretty permalink",
			globalPermalink: "pretty",
			pagePath:        "/_authors/john.html",
			frontMatter:     map[string]interface{}{"title": "John Doe", "collection": "authors"},
			expected:        "/john-doe/", // Date/categories ignored for non-post collections
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.Default()
			cfg.Permalink = tt.globalPermalink
			s := siteFake{t, cfg}

			p := page{
				file: file{
					site:      s,
					relPath:   tt.pagePath,
					outputExt: ".html",
					fm:        tt.frontMatter,
					modTime:   testDate,
				},
			}

			err := p.setPermalink()
			require.NoError(t, err)
			require.Equal(t, tt.expected, p.URL())
		})
	}
}
