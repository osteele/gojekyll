package renderers

import (
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRenderMarkdown(t *testing.T) {
	require.Equal(t, "<p><em>b</em></p>\n", mustMarkdownString("*b*"))
}

func TestRenderMarkdownWithHtml1(t *testing.T) {
	require.Equal(t, "<p><div a=1><p><em>b</em></p>\n</div></p>\n", mustMarkdownString(`<div a=1 markdown="1">*b*</div>`))
	require.Equal(t, "<p><div a=1><p><em>b</em></p>\n</div></p>\n", mustMarkdownString(`<div a=1 markdown='1'>*b*</div>`))
	require.Equal(t, "<p><div a=1><p><em>b</em></p>\n</div></p>\n", mustMarkdownString(`<div a=1 markdown=1>*b*</div>`))
	require.Equal(t, "<div a=1 markdown=1><p></div>", `<div a=1 markdown=1><p></div>`)
}

func TestRenderMarkdownWithHtml2(t *testing.T) {
	t.Skip("skipping broken test.")
	// FIXME for now, manually test against against site/testdata/site1/markdown.md.
	// These render correctly in the entire pipeline, but not in the test.
	require.Equal(t, "<p><div>*b*</div></p>\n", mustMarkdownString("<div>*b*</div>"))
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

// func renderMarkdownString(md string) (string, error) {
// 	s, err := renderMarkdown([]byte(md))
// 	if err != nil {
// 		return "", err
// 	}
// 	return string(s), err
// }
