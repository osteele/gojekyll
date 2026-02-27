package site

import (
	"fmt"
	"testing"
	"time"

	"github.com/osteele/gojekyll/config"
	"github.com/osteele/gojekyll/pages"
	"github.com/osteele/liquid/tags"
	"github.com/stretchr/testify/require"
)

func readTestSiteDrop(t *testing.T) map[string]interface{} {
	site, err := FromDirectory("testdata/site1", config.Flags{})
	require.NoError(t, err)
	require.NoError(t, site.Read())
	return site.ToLiquid().(tags.IterationKeyedMap)
}

// TODO test cases for collections, categories, tags, data

func TestSite_ToLiquid(t *testing.T) {
	drop := readTestSiteDrop(t)
	docs, ok := drop["documents"].([]Page)
	require.True(t, ok, fmt.Sprintf("documents has type %T", drop["documents"]))
	require.Len(t, docs, 3)
}

func TestSite_ToLiquid_time(t *testing.T) {
	drop := readTestSiteDrop(t)
	_, ok := drop["time"].(time.Time)
	require.True(t, ok)
	// TODO read time from config if present
}

func TestSite_ToLiquid_pages(t *testing.T) {
	drop := readTestSiteDrop(t)
	ps, ok := drop["pages"]
	require.True(t, ok, fmt.Sprintf("pages has type %T", drop["pages"]))
	require.Len(t, ps, 3) // includes main.scss for SCSS transpiler testing

	ps, ok = drop["html_pages"]
	require.True(t, ok, fmt.Sprintf("pages has type %T", drop["pages"]))
	require.Len(t, ps, 2)
}

func TestSite_ToLiquid_posts(t *testing.T) {
	drop := readTestSiteDrop(t)
	posts, ok := drop["posts"].([]Page)
	require.True(t, ok, fmt.Sprintf("posts has type %T", drop["posts"]))
	require.Len(t, posts, 1)
}

func TestSite_ToLiquid_related_posts(t *testing.T) {
	drop := readTestSiteDrop(t)
	posts, ok := drop["related_posts"].([]Page)
	require.True(t, ok, fmt.Sprintf("related_posts has type %T", drop["related_posts"]))
	require.Len(t, posts, 1)
}

func TestSite_readDataFiles_skips_directories(t *testing.T) {
	// Regression test: readDataFiles should skip directories and continue
	// reading subsequent files (previously used break instead of continue)
	site, err := FromDirectory("testdata/site1", config.Flags{})
	require.NoError(t, err)
	require.NoError(t, site.Read())

	// The _data dir has: alpha.json, subdir/, zulu.json
	// Both alpha and zulu should be loaded despite subdir/ between them
	require.Contains(t, site.data, "alpha", "data file before directory should be loaded")
	require.Contains(t, site.data, "zulu", "data file after directory should be loaded")
}

func TestSite_ToLiquid_tags_vs_categories(t *testing.T) {
	drop := readTestSiteDrop(t)

	// tags and categories should be distinct groupings
	tags, ok := drop["tags"].(map[string][]Page)
	require.True(t, ok, fmt.Sprintf("tags has type %T", drop["tags"]))

	categories, ok := drop["categories"].(map[string][]Page)
	require.True(t, ok, fmt.Sprintf("categories has type %T", drop["categories"]))

	// The test post has tags: ["go", "jekyll"] and categories: ["dev"]
	// tags should contain "go" and "jekyll" keys
	require.Contains(t, tags, "go", "tags should contain 'go'")
	require.Contains(t, tags, "jekyll", "tags should contain 'jekyll'")
	require.NotContains(t, tags, "dev", "tags should not contain category 'dev'")

	// categories should contain "dev" key
	require.Contains(t, categories, "dev", "categories should contain 'dev'")
	require.NotContains(t, categories, "go", "categories should not contain tag 'go'")
}

func TestSite_ToLiquid_static_files(t *testing.T) {
	drop := readTestSiteDrop(t)
	files, ok := drop["static_files"].([]*pages.StaticFile)
	require.True(t, ok, fmt.Sprintf("static_files has type %T", drop["static_files"]))
	require.Len(t, files, 1)

	f := files[0].ToLiquid().(tags.IterationKeyedMap)
	require.IsType(t, "", f["path"])
	require.IsType(t, time.Now(), f["modified_time"])
	require.Equal(t, ".html", f["extname"])
}
