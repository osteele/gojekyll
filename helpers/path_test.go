package helpers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func timeMustParse(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestFilenameDate(t *testing.T) {
	d, found := FilenameDate("2017-07-02-post.html")
	require.True(t, found)
	require.Equal(t, timeMustParse("2017-07-02T00:00:00Z"), d)

	d, found = FilenameDate("not-post.html")
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
