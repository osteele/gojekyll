package site

import (
	"fmt"
	"testing"
	"time"

	"github.com/osteele/gojekyll/config"
	"github.com/osteele/gojekyll/pages"
	"github.com/stretchr/testify/require"
)

func readTestSiteDrop(t *testing.T) map[string]interface{} {
	site, err := FromDirectory("testdata/site1", config.Flags{})
	require.NoError(t, err)
	require.NoError(t, site.Load())
	return site.ToLiquid().(map[string]interface{})
}

// TODO test cases for collections, categories, tags, data

func TestSite_ToLiquid(t *testing.T) {
	drop := readTestSiteDrop(t)
	docs, isTime := drop["documents"].([]pages.Document)
	require.True(t, isTime, fmt.Sprintf("documents has type %T", drop["documents"]))
	require.Len(t, docs, 4)
}
func TestSite_ToLiquid_time(t *testing.T) {
	drop := readTestSiteDrop(t)

	_, ok := drop["time"].(time.Time)
	require.True(t, ok)

	// TODO read time from config if present
}

func TestSite_ToLiquid_pages(t *testing.T) {
	drop := readTestSiteDrop(t)
	pages, ok := drop["pages"].([]pages.Page)
	require.True(t, ok, fmt.Sprintf("pages has type %T", drop["pages"]))
	require.Len(t, pages, 3)
}

func TestSite_ToLiquid_posts(t *testing.T) {
	drop := readTestSiteDrop(t)
	posts, ok := drop["posts"].([]pages.Page)
	require.True(t, ok, fmt.Sprintf("posts has type %T", drop["posts"]))
	require.Len(t, posts, 1)
}

func TestSite_ToLiquid_related_posts(t *testing.T) {
	drop := readTestSiteDrop(t)
	posts, ok := drop["related_posts"].([]pages.Page)
	require.True(t, ok, fmt.Sprintf("related_posts has type %T", drop["related_posts"]))
	require.Len(t, posts, 1)
}

func TestSite_ToLiquid_static_files(t *testing.T) {
	drop := readTestSiteDrop(t)
	files, ok := drop["static_files"].([]*pages.StaticFile)
	require.True(t, ok, fmt.Sprintf("static_files has type %T", drop["static_files"]))
	require.Len(t, files, 1)

	// TODO move this test to pages package
	f := files[0].ToLiquid().(map[string]interface{})
	require.Equal(t, "static.html", f["path"])
	_, isTime := f["modified_time"].(time.Time)
	require.True(t, isTime)
	require.Equal(t, ".html", f["extname"])
}
