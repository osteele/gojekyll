package site

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/osteele/gojekyll/config"
	"github.com/stretchr/testify/require"
)

func TestCleanRejectsDestinationContainingSource(t *testing.T) {
	parent := t.TempDir()
	source := filepath.Join(parent, "source")
	require.NoError(t, os.Mkdir(source, 0o755))
	marker := filepath.Join(source, "keep.txt")
	require.NoError(t, os.WriteFile(marker, []byte("keep"), 0o644))

	s := New(config.Flags{})
	s.cfg.Source = source
	s.cfg.Destination = parent
	err := s.Clean()
	require.ErrorContains(t, err, "contains the source directory")
	_, statErr := os.Stat(marker)
	require.NoError(t, statErr)
}

func TestCleanRejectsSourceAsDestination(t *testing.T) {
	source := t.TempDir()
	marker := filepath.Join(source, "keep.txt")
	require.NoError(t, os.WriteFile(marker, []byte("keep"), 0o644))

	s := New(config.Flags{})
	s.cfg.Source = source
	s.cfg.Destination = "."
	require.Error(t, s.Clean())
	_, err := os.Stat(marker)
	require.NoError(t, err)
}

func TestCleanRejectsDestinationSymlinkAncestorContainingSource(t *testing.T) {
	parent := t.TempDir()
	source := filepath.Join(parent, "future-destination", "source")
	require.NoError(t, os.MkdirAll(source, 0o755))
	link := filepath.Join(source, "parent-link")
	if err := os.Symlink(parent, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	s := New(config.Flags{})
	s.cfg.Source = source
	s.cfg.Destination = filepath.Join("parent-link", "future-destination")
	require.ErrorContains(t, s.validateDestination(), "contains the source directory")
}

func TestReadSkipsDestinationDirectory(t *testing.T) {
	source := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(source, "index.md"), []byte("---\n---\nsource\n"), 0o644))
	destination := filepath.Join(source, "dist")
	require.NoError(t, os.Mkdir(destination, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(destination, "generated.html"), []byte("generated"), 0o644))

	s := New(config.Flags{})
	s.cfg.Source = source
	s.cfg.Destination = "dist"
	require.NoError(t, s.Read())
	require.Len(t, s.docs, 1)
	require.NotContains(t, s.Routes, "/dist/generated.html")
	_, err := s.Write()
	require.NoError(t, err)

	reloaded, err := FromDirectory(source, config.Flags{})
	require.NoError(t, err)
	reloaded.cfg.Destination = "dist"
	require.NoError(t, reloaded.Read())
	require.Len(t, reloaded.docs, 1)
	_, err = reloaded.Write()
	require.NoError(t, err)
}
