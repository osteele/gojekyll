package renderers

import (
	"path/filepath"
	"testing"

	"github.com/osteele/gojekyll/config"
	"github.com/stretchr/testify/require"
)

func TestFindLayoutRejectsParentTraversal(t *testing.T) {
	cfg := config.Default()
	cfg.Source = t.TempDir()
	manager := Manager{cfg: cfg}

	_, err := manager.FindLayout(filepath.Join("..", "secret"), nil)
	require.ErrorContains(t, err, "escapes")
}
