package renderers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRenderMarkdown(t *testing.T) {
	require.Equal(t, "<p><em>b</em></p>\n", mustMarkdownString("*b*"))
}

func TestRenderMarkdownWithHtml1(t *testing.T) {
	// Test markdown="1" (same as block mode) - using block-level HTML
	require.Contains(t, mustMarkdownString("\n<div markdown=\"1\">\n*b*\n</div>\n"), "<em>b</em>")
	require.Contains(t, mustMarkdownString("\n<div markdown='1'>\n*b*\n</div>\n"), "<em>b</em>")
	require.Contains(t, mustMarkdownString("\n<div markdown=1>\n*b*\n</div>\n"), "<em>b</em>")

	// Test markdown="block" (should be same as markdown="1")
	require.Contains(t, mustMarkdownString("\n<div markdown=\"block\">\n*b*\n</div>\n"), "<em>b</em>")

	// Test markdown="span" (no paragraphs, just inline elements)
	result := mustMarkdownString("\n<div markdown=\"span\">\n*b*\n</div>\n")
	require.Contains(t, result, "<em>b</em>")
	require.NotContains(t, result, "<p><em>b</em></p>")

	// Test markdown="0" (no markdown processing)
	require.NotContains(t, mustMarkdownString("\n<div markdown=\"0\">\n*b*\n</div>\n"), "<em>")
	require.Contains(t, mustMarkdownString("\n<div markdown=\"0\">\n*b*\n</div>\n"), "*b*")
}

func TestRenderMarkdownWithHtml2(t *testing.T) {
	// No markdown attribute with block-level HTML - content should not be processed
	result := mustMarkdownString("\n<div>\n*b*\n</div>\n")
	require.NotContains(t, result, "<em>")
	require.Contains(t, result, "*b*")

	// Test autolink processing with different markdown modes (block-level HTML)
	require.Contains(t, mustMarkdownString("\n<div markdown=1>\n<user@example.com>\n</div>\n"), `<a href="mailto:user@example.com">`)
	require.Contains(t, mustMarkdownString("\n<div markdown=\"block\">\n<user@example.com>\n</div>\n"), `<a href="mailto:user@example.com">`)
	require.Contains(t, mustMarkdownString("\n<div markdown=\"span\">\n<user@example.com>\n</div>\n"), `<a href="mailto:user@example.com">`)

	emailResult := mustMarkdownString("\n<div markdown=\"0\">\n<user@example.com>\n</div>\n")
	require.NotContains(t, emailResult, `<a href="mailto:user@example.com">`)
	require.Contains(t, emailResult, "user@example.com")

	// Test URL autolink processing with different markdown modes (block-level HTML)
	require.Contains(t, mustMarkdownString("\n<div markdown=1>\n<http://example.com>\n</div>\n"), `<a href="http://example.com">`)
	require.Contains(t, mustMarkdownString("\n<div markdown=\"block\">\n<http://example.com>\n</div>\n"), `<a href="http://example.com">`)
	require.Contains(t, mustMarkdownString("\n<div markdown=\"span\">\n<http://example.com>\n</div>\n"), `<a href="http://example.com">`)

	urlResult := mustMarkdownString("\n<div markdown=\"0\">\n<http://example.com>\n</div>\n")
	require.NotContains(t, urlResult, `<a href="http://example.com">`)
	require.Contains(t, urlResult, "http://example.com")
}

func mustMarkdownString(md string) string {
	s, err := renderMarkdown([]byte(md))
	if err != nil {
		panic(err)
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
