package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMustAbs(t *testing.T) {
	abs := MustAbs(".")
	require.True(t, filepath.IsAbs(abs) || strings.HasPrefix(abs, "/"))
}

func TestParseFilenameDate(t *testing.T) {
	os.Setenv("TZ", "America/New_York") // nolint: errcheck
	d, title, found := ParseFilenameDateTitle("2017-07-02-post.html")
	require.True(t, found)
	require.Equal(t, "Post", title)
	// The date should be 2017-07-02 at midnight in the local timezone
	expected := time.Date(2017, 7, 2, 0, 0, 0, 0, time.Local)
	require.Equal(t, expected, d)

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

func TestJoinWithin(t *testing.T) {
	root := t.TempDir()
	inside, err := JoinWithin(root, filepath.Join("nested", "file.html"))
	require.NoError(t, err)
	require.Equal(t, filepath.Join(root, "nested", "file.html"), inside)

	_, err = JoinWithin(root, filepath.Join("..", "secret.html"))
	require.ErrorContains(t, err, "escapes")
	_, err = JoinWithin(root, filepath.Join(root, "secret.html"))
	require.ErrorContains(t, err, "absolute path")
}

func TestJoinWithinRejectsSymlinkEscape(t *testing.T) {
	parent := t.TempDir()
	root := filepath.Join(parent, "root")
	outside := filepath.Join(parent, "outside")
	require.NoError(t, os.Mkdir(root, 0o755))
	require.NoError(t, os.Mkdir(outside, 0o755))
	if err := os.Symlink(outside, filepath.Join(root, "link")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	_, err := JoinWithin(root, "link")
	require.ErrorContains(t, err, "symlink")
}
