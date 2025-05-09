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
			name: "Basic TOC inline",
			html: `<h1 id="title">Title</h1>
<p>{:toc}</p>
<h2 id="section1">Section 1</h2>
<h3 id="subsection1">Subsection 1</h3>
<h2 id="section2">Section 2</h2>`,
			expected: `<h1 id="title">Title</h1>
<p><div class="toc"><ul class="section-nav"><li><a href="#title">Title</a><ul><li><a href="#section1">Section 1</a><ul><li><a href="#subsection1">Subsection 1</a></li></ul></li><li><a href="#section2">Section 2</a></li></ul></li></ul></div></p>
<h2 id="section1">Section 1</h2>
<h3 id="subsection1">Subsection 1</h3>
<h2 id="section2">Section 2</h2>`,
		},
		{
			name: "TOC inline with whitespace",
			html: `<h1 id="title">Title</h1>
<p>{: toc }</p>
<h2 id="section1">Section 1</h2>
<h3 id="subsection1">Subsection 1</h3>
<h2 id="section2">Section 2</h2>`,
			expected: `<h1 id="title">Title</h1>
<p><div class="toc"><ul class="section-nav"><li><a href="#title">Title</a><ul><li><a href="#section1">Section 1</a><ul><li><a href="#subsection1">Subsection 1</a></li></ul></li><li><a href="#section2">Section 2</a></li></ul></li></ul></div></p>
<h2 id="section1">Section 1</h2>
<h3 id="subsection1">Subsection 1</h3>
<h2 id="section2">Section 2</h2>`,
		},
		{
			name: "Basic TOC block",
			html: `<h1 id="title">Title</h1>
<p>{::toc}</p>
<h2 id="section1">Section 1</h2>
<h3 id="subsection1">Subsection 1</h3>
<h2 id="section2">Section 2</h2>`,
			expected: `<h1 id="title">Title</h1>
<p><div class="toc"><ul class="section-nav"><li><a href="#title">Title</a><ul><li><a href="#section1">Section 1</a><ul><li><a href="#subsection1">Subsection 1</a></li></ul></li><li><a href="#section2">Section 2</a></li></ul></li></ul></div></p>
<h2 id="section1">Section 1</h2>
<h3 id="subsection1">Subsection 1</h3>
<h2 id="section2">Section 2</h2>`,
		},
		{
			name: "TOC block with whitespace",
			html: `<h1 id="title">Title</h1>
<p>{:: toc }</p>
<h2 id="section1">Section 1</h2>
<h3 id="subsection1">Subsection 1</h3>
<h2 id="section2">Section 2</h2>`,
			expected: `<h1 id="title">Title</h1>
<p><div class="toc"><ul class="section-nav"><li><a href="#title">Title</a><ul><li><a href="#section1">Section 1</a><ul><li><a href="#subsection1">Subsection 1</a></li></ul></li><li><a href="#section2">Section 2</a></li></ul></li></ul></div></p>
<h2 id="section1">Section 1</h2>
<h3 id="subsection1">Subsection 1</h3>
<h2 id="section2">Section 2</h2>`,
		},
		{
			name: "No TOC marker",
			html: `<h1 id="title">Title</h1>
<h2 id="section1">Section 1</h2>
<h3 id="subsection1">Subsection 1</h3>
<h2 id="section2">Section 2</h2>`,
			expected: `<h1 id="title">Title</h1>
<h2 id="section1">Section 1</h2>
<h3 id="subsection1">Subsection 1</h3>
<h2 id="section2">Section 2</h2>`,
		},
		{
			name: "No headings",
			html: `<p>Just a paragraph</p>
<p>{:toc}</p>
<p>Another paragraph</p>`,
			expected: `<p>Just a paragraph</p>
<p><div class="toc"><ul class="section-nav"><li>No headings found</li></ul></div></p>
<p>Another paragraph</p>`,
		},
		{
			name: "With no_toc marker",
			html: `<h1 id="title">Title</h1>
<p>{:toc}</p>
<h2 id="section1">Section 1</h2>
<h3 id="subsection1">{:.no_toc}Subsection 1</h3>
<h2 id="section2">Section 2</h2>`,
			expected: `<h1 id="title">Title</h1>
<p><div class="toc"><ul class="section-nav"><li><a href="#title">Title</a><ul><li><a href="#section1">Section 1</a></li><li><a href="#section2">Section 2</a></li></ul></li></ul></div></p>
<h2 id="section1">Section 1</h2>
<h3 id="subsection1">Subsection 1</h3>
<h2 id="section2">Section 2</h2>`,
		},
		{
			name: "With no_toc marker and whitespace",
			html: `<h1 id="title">Title</h1>
<p>{:toc}</p>
<h2 id="section1">Section 1</h2>
<h3 id="subsection1">{: .no_toc }Subsection 1</h3>
<h2 id="section2">Section 2</h2>`,
			expected: `<h1 id="title">Title</h1>
<p><div class="toc"><ul class="section-nav"><li><a href="#title">Title</a><ul><li><a href="#section1">Section 1</a></li><li><a href="#section2">Section 2</a></li></ul></li></ul></div></p>
<h2 id="section1">Section 1</h2>
<h3 id="subsection1">Subsection 1</h3>
<h2 id="section2">Section 2</h2>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := processTOC([]byte(tt.html))
			if err != nil {
				t.Fatalf("Error processing HTML: %v", err)
			}

			if string(result) != tt.expected {
				t.Errorf("Expected:\n%s\n\nGot:\n%s", tt.expected, string(result))
			}
		})
	}
}

func TestMarkdownTOCIntegration(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		containsTOC bool
	}{
		{
			name: "Markdown with TOC inline",
			markdown: "# Title\n\n{:toc}\n\n## Section 1\n\nContent\n\n## Section 2\n\nMore content",
			containsTOC: true,
		},
		{
			name: "Markdown with TOC block",
			markdown: "# Title\n\n{::toc}\n\n## Section 1\n\nContent\n\n## Section 2\n\nMore content",
			containsTOC: true,
		},
		{
			name: "Markdown without TOC",
			markdown: "# Title\n\n## Section 1\n\nContent\n\n## Section 2\n\nMore content",
			containsTOC: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html, err := renderMarkdown([]byte(tt.markdown))
			if err != nil {
				t.Fatalf("Error rendering markdown: %v", err)
			}

			containsTOC := tocPatternInline.Match([]byte(tt.markdown)) || tocPatternBlock.Match([]byte(tt.markdown))
			containsTOCDiv := containsString(string(html), "<div class=\"toc\">")
			
			if containsTOC != tt.containsTOC {
				t.Errorf("Input markdown should %s contain TOC markers", map[bool]string{true: "", false: "not"}[tt.containsTOC])
			}
			
			if containsTOCDiv != tt.containsTOC {
				t.Errorf("Output HTML should %s contain TOC div", map[bool]string{true: "", false: "not"}[tt.containsTOC])
			}
		})
	}
}

func containsString(s, substr string) bool {
	return strings.Contains(s, substr)
}
