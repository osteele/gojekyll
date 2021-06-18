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
	require.Len(t, ps, 2)

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
