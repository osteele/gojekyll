package site

import (
	"testing"

	"github.com/osteele/gojekyll/config"
	"github.com/stretchr/testify/require"
)

//func TestSite_WatchRebuild(t *testing.T) {

func TestSite_Reloaded(t *testing.T) {
	s0 := New(config.Flags{})
	s0.cfg.Incremental = true
	s1, _ := s0.Reloaded([]string{})
	require.Equal(t, s0, s1)

	s1, _ = s0.Reloaded([]string{"_config.yml"})
	require.NotEqual(t, s0, s1)
}

//func TestSite_processFilesEvent(t *testing.T) {
//func TestSite_rebuild(t *testing.T) {

func TestSite_RequiresFullReload(t *testing.T) {
	s := New(config.Flags{})
	require.False(t, s.RequiresFullReload([]string{}))
	require.True(t, s.RequiresFullReload([]string{"file.md"}))
	require.False(t, s.RequiresFullReload([]string{".git"}))
	// require.False(t, s.RequiresFullReload([]string{"_site"}))
	// require.False(t, s.RequiresFullReload([]string{"_site/index.html"}))

	s.cfg.Incremental = true
	require.False(t, s.RequiresFullReload([]string{}))
	require.False(t, s.RequiresFullReload([]string{"file.md"}))
	require.True(t, s.RequiresFullReload([]string{"_config.yml"}))
}

//func TestSite_affectsBuildFilter(t *testing.T) {
//func TestSite_fileAffectsBuild(t *testing.T) {
//func TestSite_invalidatesDoc(t *testing.T) {
