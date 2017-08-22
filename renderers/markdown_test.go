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
