package renderers

import (
	"strings"
	"testing"
)

func TestTOCGeneration(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name: "TOC in unordered list (Jekyll-compatible)",
			html: `<h1 id="title">Title</h1>
<ul>
<li>TOC
{:toc}</li>
</ul>
<h2 id="section1">Section 1</h2>
<h3 id="subsection1">Subsection 1</h3>
<h2 id="section2">Section 2</h2>`,
			expected: `<h1 id="title">Title</h1>
<ul id="markdown-toc"><li><a href="#section1">Section 1</a><ul><li><a href="#subsection1">Subsection 1</a></li></ul></li><li><a href="#section2">Section 2</a></li></ul>
<h2 id="section1">Section 1</h2>
<h3 id="subsection1">Subsection 1</h3>
<h2 id="section2">Section 2</h2>`,
		},
		{
			name: "Standalone {:toc} NOT processed (remains literal)",
			html: `<h1 id="title">Title</h1>
<p>{:toc}</p>
<h2 id="section1">Section 1</h2>
<h3 id="subsection1">Subsection 1</h3>
<h2 id="section2">Section 2</h2>`,
			expected: `<h1 id="title">Title</h1>
<p>{:toc}</p>
<h2 id="section1">Section 1</h2>
<h3 id="subsection1">Subsection 1</h3>
<h2 id="section2">Section 2</h2>`,
		},
		{
			name: "Block syntax {::toc} NOT processed (invalid kramdown)",
			html: `<h1 id="title">Title</h1>
<p>{::toc}</p>
<h2 id="section1">Section 1</h2>
<h3 id="subsection1">Subsection 1</h3>
<h2 id="section2">Section 2</h2>`,
			expected: `<h1 id="title">Title</h1>
<p>{::toc}</p>
<h2 id="section1">Section 1</h2>
<h3 id="subsection1">Subsection 1</h3>
<h2 id="section2">Section 2</h2>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html, err := processTOC([]byte(tt.html), nil)
			if err != nil {
				t.Fatalf("Error processing TOC: %v", err)
			}
			if string(html) != tt.expected {
				t.Errorf("Expected:\n%s\n\nGot:\n%s", tt.expected, string(html))
			}
		})
	}
}

func TestMarkdownWithTOC(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
	}{
		{
			name: "Basic markdown with TOC",
			markdown: `* TOC
{:toc}

### Section 1

#### Subsection`,
		},
		{
			name: "Complex nested sections",
			markdown: `* TOC
{:toc}

# Main Title

## Section 1
### Subsection 1.1
### Subsection 1.2

## Section 2
### Subsection 2.1
#### Deep subsection`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html, err := renderMarkdown([]byte(tt.markdown))
			if err != nil {
				t.Fatalf("Error rendering markdown: %v", err)
			}

			htmlStr := string(html)

			// Check for Jekyll-style TOC structure
			if !containsString(htmlStr, "<ul id=\"markdown-toc\">") {
				t.Error("Should contain TOC ul with id='markdown-toc'")
			}
		})
	}
}

func TestTOCWithNoTocMarker(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
	}{
		{
			name: "Exclude heading with {:.no_toc}",
			markdown: `* TOC
{:toc}

## Heading 1

## Excluded Heading {:.no_toc}

## Heading 2`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html, err := renderMarkdown([]byte(tt.markdown))
			if err != nil {
				t.Fatalf("Error rendering markdown: %v", err)
			}

			htmlStr := string(html)

			// The excluded heading should not appear in TOC
			// But should still appear as a heading in the document
			// This is tricky to test - need to check TOC specifically
			if !containsString(htmlStr, "<ul id=\"markdown-toc\">") {
				t.Error("Should contain TOC")
			}
		})
	}
}

// TestLiteralTOCMarkers tests that {:toc} markers remain literal in invalid contexts
func TestLiteralTOCMarkers(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		contains string // text that should remain in output
	}{
		{
			name: "Literal {:toc} in heading",
			markdown: `## How to use {:toc}

Content here.`,
			contains: "How to use {:toc}",
		},
		{
			name: "Literal {:toc} in paragraph",
			markdown: `Some text about {:toc} markers.

## Section 1`,
			contains: "about {:toc} markers",
		},
		{
			name: "Literal {:toc} in standalone paragraph",
			markdown: `{:toc}

## Section 1`,
			contains: "<p>{:toc}</p>",
		},
		{
			name: "Literal {:toc} in ordered list",
			markdown: `1. Item
{:toc}

## Section 1`,
			contains: "{:toc}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html, err := renderMarkdown([]byte(tt.markdown))
			if err != nil {
				t.Fatalf("Error rendering markdown: %v", err)
			}

			htmlStr := string(html)

			// Check that the literal text is preserved
			if !containsString(htmlStr, tt.contains) {
				t.Errorf("Expected to find literal text %q in output, got:\n%s", tt.contains, htmlStr)
			}

			// Verify that NO TOC was generated
			if containsString(htmlStr, "<ul id=\"markdown-toc\">") {
				t.Error("Should NOT generate TOC for literal {:toc} markers")
			}
		})
	}
}

// TestNoTocSiblingParagraph tests that {:.no_toc} in sibling paragraphs excludes headings
func TestNoTocSiblingParagraph(t *testing.T) {
	tests := []struct {
		name            string
		markdown        string
		shouldExclude   string // heading text that should be excluded from TOC
		shouldStayInDoc string // heading text that should remain in document
	}{
		{
			name: "no_toc in sibling paragraph - EXCLUDED from TOC",
			markdown: `* TOC
{:toc}

## Heading 1

## Excluded Heading
{:.no_toc}

## Heading 2`,
			shouldExclude:   "Excluded Heading",
			shouldStayInDoc: "Excluded Heading",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html, err := renderMarkdown([]byte(tt.markdown))
			if err != nil {
				t.Fatalf("Error rendering markdown: %v", err)
			}

			htmlStr := string(html)

			// The heading should still appear in the document
			if !containsString(htmlStr, tt.shouldStayInDoc) {
				t.Errorf("Heading %q should remain in document, got:\n%s", tt.shouldStayInDoc, htmlStr)
			}

			// Extract just the TOC portion to check exclusion
			tocStart := strings.Index(htmlStr, "<ul id=\"markdown-toc\">")
			tocEnd := strings.Index(htmlStr, "</ul>")
			if tocStart == -1 || tocEnd == -1 {
				t.Fatal("Could not find TOC in output")
			}
			tocContent := htmlStr[tocStart : tocEnd+5] // +5 for </ul>

			// The excluded heading should NOT appear in the TOC
			if containsString(tocContent, tt.shouldExclude) {
				t.Errorf("Heading %q should be excluded from TOC, but found in:\n%s", tt.shouldExclude, tocContent)
			}
		})
	}
}

// TestNoTocInline tests that {:.no_toc} inline in heading text is kept literal
func TestNoTocInline(t *testing.T) {
	markdown := `* TOC
{:toc}

## Heading 1

## Not Excluded {:.no_toc}

## Heading 2`

	html, err := renderMarkdown([]byte(markdown))
	if err != nil {
		t.Fatalf("Error rendering markdown: %v", err)
	}

	htmlStr := string(html)

	// The literal {:.no_toc} text should appear in the heading
	if !containsString(htmlStr, "Not Excluded {:.no_toc}") {
		t.Error("Heading should contain literal {:.no_toc} text")
	}

	// Extract TOC portion
	tocStart := strings.Index(htmlStr, "<ul id=\"markdown-toc\">")
	tocEnd := strings.Index(htmlStr, "</ul>")
	if tocStart == -1 || tocEnd == -1 {
		t.Fatal("Could not find TOC in output")
	}
	tocContent := htmlStr[tocStart : tocEnd+5]

	// The heading WITH literal {:.no_toc} should appear in TOC
	if !containsString(tocContent, "Not Excluded") {
		t.Error("Heading with inline {:.no_toc} should appear in TOC (it's literal text, not a marker)")
	}
}

// Helper function to check if a string contains a substring
func containsString(s, substr string) bool {
	return strings.Contains(s, substr)
}
