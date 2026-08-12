package renderers

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRootedTemplateStoreReadsWithinRoot(t *testing.T) {
	root := t.TempDir()
	filename := filepath.Join(root, "include.html")
	require.NoError(t, os.WriteFile(filename, []byte("included"), 0o644))
	store, err := newRootedTemplateStore(root)
	require.NoError(t, err)

	contents, err := store.ReadTemplate(filename)
	require.NoError(t, err)
	require.Equal(t, []byte("included"), contents)
}

func TestRootedTemplateStoreRejectsPathsOutsideRoot(t *testing.T) {
	root := t.TempDir()
	outside := filepath.Join(t.TempDir(), "outside.html")
	require.NoError(t, os.WriteFile(outside, []byte("outside"), 0o644))
	store, err := newRootedTemplateStore(root)
	require.NoError(t, err)

	_, err = store.ReadTemplate(outside)
	require.Error(t, err)
}

func TestRootedTemplateStoreRejectsSymlinkEscapes(t *testing.T) {
	root := t.TempDir()
	outside := filepath.Join(t.TempDir(), "outside.html")
	require.NoError(t, os.WriteFile(outside, []byte("outside"), 0o644))
	link := filepath.Join(root, "link.html")
	if err := os.Symlink(outside, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	store, err := newRootedTemplateStore(root)
	require.NoError(t, err)

	_, err = store.ReadTemplate(link)
	require.Error(t, err)
}
