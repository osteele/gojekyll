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
	// Test markdown="1" (same as block mode)
	require.Contains(t, mustMarkdownString(`<div a=1 markdown="1">*b*</div>`), "<em>b</em>")
	require.Contains(t, mustMarkdownString(`<div a=1 markdown='1'>*b*</div>`), "<em>b</em>")
	require.Contains(t, mustMarkdownString(`<div a=1 markdown=1>*b*</div>`), "<em>b</em>")
	
	// Test markdown="block" (should be same as markdown="1")
	require.Contains(t, mustMarkdownString(`<div a=1 markdown="block">*b*</div>`), "<em>b</em>")
	
	// Test markdown="span" (no paragraphs, just inline elements)
	result := mustMarkdownString(`<div a=1 markdown="span">*b*</div>`)
	require.Contains(t, result, "<em>b</em>")
	require.NotContains(t, result, "<p><em>b</em></p>")
	
	// Test markdown="0" (no markdown processing)
	require.NotContains(t, mustMarkdownString(`<div a=1 markdown="0">*b*</div>`), "<em>")
	require.Contains(t, mustMarkdownString(`<div a=1 markdown="0">*b*</div>`), "*b*")
}

func TestRenderMarkdownWithHtml2(t *testing.T) {
	// No markdown attribute - content should not be processed
	require.Equal(t, "<p><div>*b*</div></p>\n", mustMarkdownString("<div>*b*</div>"))
	
	// Test autolink processing with different markdown modes
	require.Contains(t, mustMarkdownString(`<div markdown=1><user@example.com></div>`), `<a href="mailto:user@example.com">`)
	require.Contains(t, mustMarkdownString(`<div markdown="block"><user@example.com></div>`), `<a href="mailto:user@example.com">`)
	require.Contains(t, mustMarkdownString(`<div markdown="span"><user@example.com></div>`), `<a href="mailto:user@example.com">`)
	require.NotContains(t, mustMarkdownString(`<div markdown="0"><user@example.com></div>`), `<a href="mailto:user@example.com">`)
	
	// Test URL autolink processing with different markdown modes
	require.Contains(t, mustMarkdownString(`<div markdown=1><http://example.com></div>`), `<a href="http://example.com">`)
	require.Contains(t, mustMarkdownString(`<div markdown="block"><http://example.com></div>`), `<a href="http://example.com">`)
	require.Contains(t, mustMarkdownString(`<div markdown="span"><http://example.com></div>`), `<a href="http://example.com">`)
	require.NotContains(t, mustMarkdownString(`<div markdown="0"><http://example.com></div>`), `<a href="http://example.com">`)
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
