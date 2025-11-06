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
	{"base", "/:categories/:name:output_ext", "/a/b/base"},
	{"base", "none", "/a/b/base.html"},
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
	localLegacyDate := legacyTestDate.In(time.Local)
	legacyDateTests := []pathTest{
		{"base", "date", fmt.Sprintf("/a/b/%04d/%02d/%02d/base.html", localLegacyDate.Year(), localLegacyDate.Month(), localLegacyDate.Day())},
		{"base", "pretty", fmt.Sprintf("/a/b/%04d/%02d/%02d/base/", localLegacyDate.Year(), localLegacyDate.Month(), localLegacyDate.Day())},
		{"base", "ordinal", fmt.Sprintf("/a/b/%04d/%d/base.html", legacyTestDate.Year(), legacyTestDate.YearDay())}, // ordinal always used UTC yearDay
	}
	runLegacyTests(legacyDateTests, defaultConfig, defaultFM)

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
			pageDate:          time.Date(2023, 11, 21, 2, 0, 0, 0, time.UTC),
			expectedYear:      2023,
			expectedMonth:     time.November,
			expectedDay:       20,
		},
		{
			name:              "Europe/Berlin Timezone (Date change)",
			permalinkTimezone: "Europe/Berlin",
			// 2023-11-20 22:00:00 UTC is 2023-11-20 23:00:00 CET (same day)
			// Let's use a time that will be next day in Berlin
			// 2023-11-20 23:30:00 UTC is 2023-11-21 00:30:00 CET
			pageDate:          time.Date(2023, 11, 20, 23, 30, 0, 0, time.UTC),
			expectedYear:      2023,
			expectedMonth:     time.November,
			expectedDay:       21,
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
			fm := map[string]interface{}{
				"categories": "testcat",
				"title":      "testpage", // title is used in date style if slug/name not present
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
