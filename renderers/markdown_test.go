package renderers

import (
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRenderMarkdown(t *testing.T) {
	require.Equal(t, "<p><em>b</em></p>\n", mustMarkdownString("*b*"))
	require.Equal(t, "<div>*b*</div>\n", mustMarkdownString("<div>*b*</div>"))
	require.Equal(t, "<div a=1><p><em>b</em></p>\n</div>\n", mustMarkdownString(`<div a=1 markdown="1">*b*</div>`))
	require.Equal(t, "<div a=1><p><em>b</em></p>\n</div>\n", mustMarkdownString(`<div a=1 markdown='1'>*b*</div>`))
	require.Equal(t, "<div a=1><p><em>b</em></p>\n</div>\n", mustMarkdownString(`<div a=1 markdown=1>*b*</div>`))

	_, err := renderMarkdownString(`<div a=1 markdown=1><p></div>`)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "EOF")

	require.Contains(t, mustMarkdownString(`<div markdown=1><user@example.com></div>`), `<a href="mailto:user@example.com">user@example.com</a>`)
	require.Contains(t, mustMarkdownString(`<div markdown=1><http://example.com></div>`), `<a href="http://example.com">http://example.com</a>`)
}

func mustMarkdownString(md string) string {
	s, err := renderMarkdown([]byte(md))
	if err != nil {
		log.Fatal(err)
	}
	return string(s)
}

func renderMarkdownString(md string) (string, error) {
	s, err := renderMarkdown([]byte(md))
	if err != nil {
		return "", err
	}
	return string(s), err
}
