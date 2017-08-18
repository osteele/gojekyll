package pages

import (
	"bytes"
	"os"
	"testing"

	"github.com/osteele/gojekyll/config"
	"github.com/stretchr/testify/require"
)

func TestPage_TemplateContext(t *testing.T) {
	s := siteFake{t, config.Default()}
	f := file{site: s}
	p := page{file: f}
	os.Setenv("JEKYLL_ENV", "") // nolint: errcheck
	tc := p.TemplateContext()
	j := tc["jekyll"].(map[string]string)
	require.Equal(t, "development", j["environment"])
	require.Contains(t, j["version"], "gojekyll")

	os.Setenv("JEKYLL_ENV", "production") // nolint: errcheck
	tc = p.TemplateContext()
	j = tc["jekyll"].(map[string]string)
	require.Equal(t, "production", j["environment"])
}

func TestPage_Categories(t *testing.T) {
	s := siteFake{t, config.Default()}
	fm := map[string]interface{}{"categories": "b a"}
	f := file{site: s, frontMatter: fm}
	p := page{file: f}
	require.Equal(t, []string{"a", "b"}, p.Categories())
}

func TestPage_Write(t *testing.T) {
	cfg := config.Default()
	p, err := NewFile(siteFake{t, cfg}, "testdata/page_with_layout.md", "page_with_layout.md", map[string]interface{}{})
	require.NoError(t, err)
	require.NotNil(t, p)
	buf := new(bytes.Buffer)
	require.NoError(t, p.Write(buf))
}
