package renderers

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/osteele/gojekyll/config"
	"github.com/osteele/liquid"
	"github.com/stretchr/testify/require"
)

func TestFindLayoutRejectsParentTraversal(t *testing.T) {
	cfg := config.Default()
	cfg.Source = t.TempDir()
	manager := Manager{cfg: cfg}

	_, err := manager.FindLayout(filepath.Join("..", "secret"), nil)
	require.ErrorContains(t, err, "escapes")
}

func TestFindLayoutCachesTemplateAndClonesFrontMatter(t *testing.T) {
	source := t.TempDir()
	layouts := filepath.Join(source, "_layouts")
	require.NoError(t, os.Mkdir(layouts, 0o755))
	filename := filepath.Join(layouts, "default.html")
	require.NoError(t, os.WriteFile(filename, []byte("---\nlayout: base\n---\nfirst"), 0o644))
	cfg := config.Default()
	cfg.Source = source
	manager := Manager{cfg: cfg, liquidEngine: liquid.NewEngine()}

	firstTemplate, err := manager.FindLayout("default", nil)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filename, []byte("---\nlayout: changed\n---\nsecond"), 0o644))
	var frontMatter map[string]interface{}
	secondTemplate, err := manager.FindLayout("default", &frontMatter)
	require.NoError(t, err)
	require.Same(t, firstTemplate, secondTemplate)
	require.Equal(t, "base", frontMatter["layout"])
	frontMatter["layout"] = "mutated"
	_, err = manager.FindLayout("default", &frontMatter)
	require.NoError(t, err)
	require.Equal(t, "base", frontMatter["layout"])
}
