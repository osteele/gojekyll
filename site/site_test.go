package site

import (
	"testing"

	"github.com/osteele/gojekyll/config"
	"github.com/stretchr/testify/require"
)

func TestKeepFile(t *testing.T) {
	s := New(config.Flags{})
	require.Equal(t, "", s.PathPrefix())
	require.False(t, s.KeepFile("random"))
	require.True(t, s.KeepFile(".git"))
	require.True(t, s.KeepFile(".svn"))
}

func TestExclude(t *testing.T) {
	s := New(config.Flags{})
	s.cfg.Exclude = append(s.cfg.Exclude, "exclude/")
	s.cfg.Include = append(s.cfg.Include, ".include/")
	require.False(t, s.Exclude("."))
	require.True(t, s.Exclude(".git"))
	require.True(t, s.Exclude(".dir"))
	require.True(t, s.Exclude(".dir/file"))
	require.False(t, s.Exclude(".htaccess"))
	require.False(t, s.Exclude("dir"))
	require.False(t, s.Exclude("dir/file"))
	require.True(t, s.Exclude("dir/.file"))
	require.True(t, s.Exclude("dir/#file"))
	require.True(t, s.Exclude("dir/~file"))
	require.True(t, s.Exclude("dir/file~"))
	require.True(t, s.Exclude("dir/subdir/.file"))
	require.False(t, s.Exclude(".include/file"))
	require.True(t, s.Exclude("exclude/file"))
	require.False(t, s.Exclude("_posts"))
	require.False(t, s.Exclude("_posts/file"))
	require.True(t, s.Exclude("_posts/_file"))
	require.True(t, s.Exclude("_posts/_dir/file"))

	// The following aren't documented but are evident
	// TODO submit a doc PR to Jekyll
	require.True(t, s.Exclude("#file"))
	require.True(t, s.Exclude("~file"))
	require.True(t, s.Exclude("file~"))
}
