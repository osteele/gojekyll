package site

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/osteele/gojekyll/config"
	"github.com/stretchr/testify/require"
)

// func TestSite_WatchRebuild(t *testing.T) {

func TestSite_Reloaded(t *testing.T) {
	s0 := New(config.Flags{})
	s0.cfg.Incremental = true
	s1, _ := s0.Reloaded([]string{})
	require.Equal(t, s0, s1)

	s1, _ = s0.Reloaded([]string{"_config.yml"})
	require.NotEqual(t, s0, s1)
}

// func TestSite_processFilesEvent(t *testing.T) {
// func TestSite_rebuild(t *testing.T) {

func TestSite_RequiresFullReload(t *testing.T) {
	s := New(config.Flags{})
	require.False(t, s.RequiresFullReload([]string{}))
	require.True(t, s.RequiresFullReload([]string{"file.md"}))
	require.False(t, s.RequiresFullReload([]string{".git"}))
	// require.False(t, s.RequiresFullReload([]string{"_site"}))
	// require.False(t, s.RequiresFullReload([]string{"_site/index.html"}))

	s.cfg.Incremental = true
	require.False(t, s.RequiresFullReload([]string{}))
	require.True(t, s.RequiresFullReload([]string{"file.md"}))
	require.True(t, s.RequiresFullReload([]string{"_config.yml"}))
}

func TestIncrementalReload(t *testing.T) {
	source := t.TempDir()
	filename := filepath.Join(source, "index.md")
	require.NoError(t, os.WriteFile(filename, []byte("---\npermalink: /old/\n---\nold\n"), 0o644))
	s, err := FromDirectory(source, config.Flags{Incremental: boolPointer(true)})
	require.NoError(t, err)
	require.NoError(t, s.Read())

	require.NoError(t, os.WriteFile(filename, []byte("---\npermalink: /old/\n---\nnew\n"), 0o644))
	same, err := s.Reloaded([]string{"index.md"})
	require.NoError(t, err)
	require.Same(t, s, same)

	require.NoError(t, os.WriteFile(filename, []byte("---\npermalink: /new/\n---\nnew\n"), 0o644))
	reloaded, err := same.Reloaded([]string{"index.md"})
	require.NoError(t, err)
	require.NotSame(t, same, reloaded)
	require.Contains(t, reloaded.Routes, "/new/")
	require.NotContains(t, reloaded.Routes, "/old/")
}

func TestIncrementalReloadDetectsAddedAndDeletedFiles(t *testing.T) {
	source := t.TempDir()
	index := filepath.Join(source, "index.md")
	require.NoError(t, os.WriteFile(index, []byte("---\n---\nindex\n"), 0o644))
	s, err := FromDirectory(source, config.Flags{Incremental: boolPointer(true)})
	require.NoError(t, err)
	require.NoError(t, s.Read())

	added := filepath.Join(source, "added.md")
	require.NoError(t, os.WriteFile(added, []byte("---\n---\nadded\n"), 0o644))
	reloaded, err := s.Reloaded([]string{"added.md"})
	require.NoError(t, err)
	require.NotSame(t, s, reloaded)
	require.NotNil(t, reloaded.documentForRelativePath("added.md"))

	require.NoError(t, os.Remove(index))
	reloadedAgain, err := reloaded.Reloaded([]string{"index.md"})
	require.NoError(t, err)
	require.NotSame(t, reloaded, reloadedAgain)
	require.Nil(t, reloadedAgain.documentForRelativePath("index.md"))
}

func boolPointer(value bool) *bool { return &value }

// func TestSite_affectsBuildFilter(t *testing.T) {
// func TestSite_fileAffectsBuild(t *testing.T) {
// func TestSite_invalidatesDoc(t *testing.T) {
