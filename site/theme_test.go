package site

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/osteele/gojekyll/config"
	"github.com/stretchr/testify/require"
)

func TestReadThemeAssetsDoesNothingWithoutTheme(t *testing.T) {
	workingDir := t.TempDir()
	assets := filepath.Join(workingDir, "assets")
	require.NoError(t, os.Mkdir(assets, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(assets, "unrelated.txt"), []byte("unrelated"), 0o644))
	originalDir, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(workingDir))
	t.Cleanup(func() { require.NoError(t, os.Chdir(originalDir)) })

	s := New(config.Flags{})
	s.Routes = map[string]Document{}
	require.NoError(t, s.readThemeAssets())
	require.Empty(t, s.docs)
}
