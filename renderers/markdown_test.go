package renderers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func renderMarkdownString(s string) string {
	return string(renderMarkdown([]byte(s)))
}

func TestRenderMarkdown(t *testing.T) {
	require.Equal(t, "<p><em>b</em></p>\n", renderMarkdownString("*b*"))
	require.Equal(t, "<div>*b*</div>\n", renderMarkdownString("<div>*b*</div>"))
	require.Equal(t, "<div a=1><p><em>b</em></p>\n</div>\n", renderMarkdownString(`<div a=1 markdown="1">*b*</div>`))
	require.Equal(t, "<div a=1><p><em>b</em></p>\n</div>\n", renderMarkdownString(`<div a=1 markdown='1'>*b*</div>`))
	require.Equal(t, "<div a=1><p><em>b</em></p>\n</div>\n", renderMarkdownString(`<div a=1 markdown=1>*b*</div>`))
}
