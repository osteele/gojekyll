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

var collectionTests = []pathTest{
	{"/a/b/c.d", "/prefix/:collection/post", "/prefix/c/post"},
	{"/a/b/c.d", "/prefix:path/post", "/prefix/a/b/c/post"},
}

// testPermalinkPatternHelper creates a page with a given config and computes its permalink.
func testPermalinkPatternHelper(t *testing.T, cfg config.Config, pageDate time.Time, pattern, path string, data map[string]interface{}) (string, error) {
	s := siteFake{t, cfg}
	fm := frontmatter.Merge(data, FrontMatter{"permalink": pattern})
	ext := filepath.Ext(path)
	switch ext {
	case ".md", ".markdown":
		ext = ".html"
	}
	f := file{site: s, relPath: path, fm: fm, outputExt: ext}
	p := page{file: f}
	p.modTime = pageDate // This is used by p.PostDate() if no date in frontmatter
	return p.computePermalink(p.permalinkVariables())
}

func TestExpandPermalinkPattern(t *testing.T) {
	// Default frontmatter for most tests
	defaultFM := map[string]interface{}{
		"categories": "b a",
	}

	// Legacy test date
	legacyTestDate, err := time.Parse(time.RFC3339, "2006-02-03T15:04:05Z")
	require.NoError(t, err)

	runLegacyTests := func(tests []pathTest, cfg config.Config, fm map[string]interface{}) {
		for i, test := range tests {
			t.Run(fmt.Sprintf("legacy/%s", test.pattern), func(t *testing.T) {
				p, err := testPermalinkPatternHelper(t, cfg, legacyTestDate, test.pattern, test.path, fm)
				require.NoError(t, err)
				require.Equalf(t, test.out, p, "%d: pattern=%s", i+1, test.pattern)
			})
		}
	}

	defaultConfig := config.Default()

	// Run non-date-dependent tests
	runLegacyTests(staticTests, defaultConfig, defaultFM)

	// Date-dependent tests (legacy behavior - uses local time)
	// These tests use posts to ensure date/category placeholders work
	postFM := map[string]interface{}{
		"categories": "b a",
		"collection": "posts",
	}
	localLegacyDate := legacyTestDate.In(time.Local)
	legacyDateTests := []pathTest{
		{"base", "date", fmt.Sprintf("/a/b/%04d/%02d/%02d/base.html", localLegacyDate.Year(), localLegacyDate.Month(), localLegacyDate.Day())},
		{"base", "pretty", fmt.Sprintf("/a/b/%04d/%02d/%02d/base/", localLegacyDate.Year(), localLegacyDate.Month(), localLegacyDate.Day())},
		{"base", "ordinal", fmt.Sprintf("/a/b/%04d/%d/base.html", legacyTestDate.Year(), legacyTestDate.YearDay())}, // ordinal always used UTC yearDay
	}
	runLegacyTests(legacyDateTests, defaultConfig, postFM)

	// Collection tests (legacy)
	collectionFM := map[string]interface{}{"categories": "b a", "collection": "c"}
	runLegacyTests(collectionTests, defaultConfig, collectionFM)

	t.Run("legacy/invalid template variable", func(t *testing.T) {
		p, err := testPermalinkPatternHelper(t, defaultConfig, legacyTestDate, "/:invalid", "/a/b/base.html", defaultFM)
		require.Error(t, err)
		require.Zero(t, p)
	})

	// --- Timezone-aware tests ---
	// Test date in UTC. Using a time that will fall on different dates in different timezones.
	// Example: 2023-11-20 22:00:00 UTC is:
	// - 2023-11-20 in America/Los_Angeles (PST, UTC-8)
	// - 2023-11-21 in Europe/Berlin (CET, UTC+1)
	permalinkTestDateUTC := time.Date(2023, 11, 20, 22, 0, 0, 0, time.UTC)

	// Pre-load locations to ensure tests don't fail if system doesn't have them
	locNewYork, err := time.LoadLocation("America/New_York")
	require.NoError(t, err, "Failed to load America/New_York timezone for testing")
	// locLosAngeles, err := time.LoadLocation("America/Los_Angeles")
	// require.NoError(t, err, "Failed to load America/Los_Angeles timezone for testing")

	type permalinkTimezoneTestCase struct {
		name              string
		permalinkTimezone string
		pageDate          time.Time
		expectedYear      int
		expectedMonth     time.Month
		expectedDay       int
	}

	timezoneTestCases := []permalinkTimezoneTestCase{
		{
			name:              "No Timezone (Fallback to Local)",
			permalinkTimezone: "", // Explicitly empty
			pageDate:          permalinkTestDateUTC,
			expectedYear:      permalinkTestDateUTC.In(time.Local).Year(),
			expectedMonth:     permalinkTestDateUTC.In(time.Local).Month(),
			expectedDay:       permalinkTestDateUTC.In(time.Local).Day(),
		},
		{
			name:              "UTC Timezone",
			permalinkTimezone: "UTC",
			pageDate:          permalinkTestDateUTC,
			expectedYear:      2023,
			expectedMonth:     time.November,
			expectedDay:       20, // 2023-11-20 22:00:00 UTC
		},
		{
			name:              "America/New_York Timezone",
			permalinkTimezone: "America/New_York",
			pageDate:          permalinkTestDateUTC, // 2023-11-20 22:00:00 UTC is 2023-11-20 17:00:00 ET
			expectedYear:      2023,
			expectedMonth:     time.November,
			expectedDay:       20,
		},
		{
			name:              "America/New_York Timezone (Date change)",
			permalinkTimezone: "America/New_York",
			// 2023-11-21 02:00:00 UTC is 2023-11-20 21:00:00 ET (previous day in ET)
			pageDate:      time.Date(2023, 11, 21, 2, 0, 0, 0, time.UTC),
			expectedYear:  2023,
			expectedMonth: time.November,
			expectedDay:   20,
		},
		{
			name:              "Europe/Berlin Timezone (Date change)",
			permalinkTimezone: "Europe/Berlin",
			// 2023-11-20 22:00:00 UTC is 2023-11-20 23:00:00 CET (same day)
			// Let's use a time that will be next day in Berlin
			// 2023-11-20 23:30:00 UTC is 2023-11-21 00:30:00 CET
			pageDate:      time.Date(2023, 11, 20, 23, 30, 0, 0, time.UTC),
			expectedYear:  2023,
			expectedMonth: time.November,
			expectedDay:   21,
		},
		{
			name:              "Invalid Timezone (Fallback to Local)",
			permalinkTimezone: "Invalid/Timezone",
			pageDate:          permalinkTestDateUTC,
			expectedYear:      permalinkTestDateUTC.In(time.Local).Year(),
			expectedMonth:     permalinkTestDateUTC.In(time.Local).Month(),
			expectedDay:       permalinkTestDateUTC.In(time.Local).Day(),
		},
	}

	for _, tc := range timezoneTestCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := config.Default()
			cfg.PermalinkTimezone = tc.permalinkTimezone

			// For these tests, we use a simple path and the "date" permalink style
			// to make it easy to check the date components.
			// Categories "testcat" and title "testpage" are arbitrary.
			// IMPORTANT: Mark as a post so date/category placeholders are used
			fm := map[string]interface{}{
				"categories": "testcat",
				"title":      "testpage", // title is used in date style if slug/name not present
				"collection": "posts",    // Mark as post so date placeholders work
			}
			path := "/testcat/testpage.md" // Path provides categories if not in FM.
			pattern := "date"              // /:categories/:year/:month/:day/:title.html

			actualPermalink, err := testPermalinkPatternHelper(t, cfg, tc.pageDate, pattern, path, fm)
			require.NoError(t, err)

			// Construct expected permalink
			// Note: categories in permalink are slugified and ordered. "testcat" -> "testcat"
			// title in permalink is slugified. "testpage" -> "testpage"
			expectedPermalink := fmt.Sprintf("/testcat/%04d/%02d/%02d/testpage.html", tc.expectedYear, tc.expectedMonth, tc.expectedDay)

			// If categories were more complex, e.g., "b a", they'd be "/a/b/..."
			// If title had spaces/special chars, it would be slugified.

			require.Equal(t, expectedPermalink, actualPermalink)

			// Additionally, test with "pretty" style to ensure trailing slash behavior is maintained
			patternPretty := "pretty" // /:categories/:year/:month/:day/:title/
			actualPermalinkPretty, err := testPermalinkPatternHelper(t, cfg, tc.pageDate, patternPretty, path, fm)
			require.NoError(t, err)
			expectedPermalinkPretty := fmt.Sprintf("/testcat/%04d/%02d/%02d/testpage/", tc.expectedYear, tc.expectedMonth, tc.expectedDay)
			require.Equal(t, expectedPermalinkPretty, actualPermalinkPretty)

		})
	}
	// Ensure that the NewYork location loaded correctly for subsequent tests if any
	_ = locNewYork
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
		{
			name:            "custom permalink pattern for page (issue #81)",
			globalPermalink: "/blog/:slug/",
			pagePath:        "/index.html",
			frontMatter:     map[string]interface{}{"title": "Home"},
			expected:        "/index.html", // Custom patterns don't apply to pages, use default
		},
		{
			name:            "custom permalink pattern for post",
			globalPermalink: "/blog/:slug/",
			pagePath:        "/_posts/2006-02-03-hello.html",
			frontMatter:     map[string]interface{}{"title": "Hello", "collection": "posts"},
			expected:        "/blog/2006-02-03-hello/", // Custom patterns apply to posts (slug from filename)
		},
		{
			name:            "custom permalink with :path for page",
			globalPermalink: "/custom/:path/",
			pagePath:        "/about.html",
			frontMatter:     map[string]interface{}{},
			expected:        "/about.html", // Custom patterns don't apply to pages
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
