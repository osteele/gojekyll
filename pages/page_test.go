package pages

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/osteele/gojekyll/config"
	"github.com/osteele/gojekyll/frontmatter"
	"github.com/osteele/gojekyll/utils"
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
	fm := frontmatter.FrontMatter{"categories": "b a"}
	f := file{site: s, frontMatter: fm}
	p := page{file: f}
	require.Equal(t, []string{"a", "b"}, p.Categories())
}

func TestPage_Write(t *testing.T) {
	t.Run("rendering", func(t *testing.T) {
		p := requirePageFromFile(t, "page_with_layout.md")
		buf := new(bytes.Buffer)
		require.NoError(t, p.Write(buf))
		require.Contains(t, buf.String(), "page with layout")
	})

	t.Run("rendering error", func(t *testing.T) {
		p := requirePageFromFile(t, "liquid_error.md")
		err := p.Write(ioutil.Discard)
		require.NotNil(t, err)
		require.Contains(t, err.Error(), "render error")
		pe, ok := err.(utils.PathError)
		require.True(t, ok)
		require.Equal(t, "testdata/liquid_error.md", pe.Path())
	})
}

func fakePageFromFile(t *testing.T, file string) (Document, error) {
	return NewFile(
		siteFake{t, config.Default()},
		filepath.Join("testdata", file),
		file,
		frontmatter.FrontMatter{},
	)
}

func requirePageFromFile(t *testing.T, file string) Document {
	p, err := fakePageFromFile(t, file)
	require.NoError(t, err)
	require.NotNil(t, p)
	return p
}
