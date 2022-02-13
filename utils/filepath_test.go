package utils

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func timeMustParse(s string) time.Time {
	// ParseInLocation is necessary on Windows. time.Parse defaults to
	// time.Location("") which doesn't match the times returned by the tested
	// functions.
	t, err := time.ParseInLocation(time.RFC3339, s, time.Local)
	if err != nil {
		panic(err)
	}
	return t
}

func TestMustAbs(t *testing.T) {
	require.True(t, strings.HasPrefix(MustAbs("."), "/"))
}

func TestParseFilenameDate(t *testing.T) {
	os.Setenv("TZ", "America/New_York") // nolint: errcheck
	d, title, found := ParseFilenameDateTitle("2017-07-02-post.html")
	require.True(t, found)
	require.Equal(t, "Post", title)
	require.Equal(t, timeMustParse("2017-07-02T00:00:00-04:00"), d)

	_, _, found = ParseFilenameDateTitle("not-post.html")
	require.False(t, found)
}

func TestTrimExt(t *testing.T) {
	require.Equal(t, "/a/b", TrimExt("/a/b.c"))
	require.Equal(t, "/a/b", TrimExt("/a/b"))
}

func TestURLPathClean(t *testing.T) {
	require.Equal(t, "/a/b", URLPathClean("/a/b"))
	require.Equal(t, "/a/b/", URLPathClean("/a/b/"))
	require.Equal(t, "/a/b", URLPathClean("/a//b"))
	require.Equal(t, "/b", URLPathClean("/a/../b"))
	require.Equal(t, "/", URLPathClean("/"))
}
