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
			result, err := processTOC([]byte(tt.html), nil)
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

func TestTOCLevelsFiltering(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		opts     *TOCOptions
		expected string
	}{
		{
			name: "Filter to H2-H3 only",
			html: `<h1 id="title">Title</h1>
<p>{:toc}</p>
<h2 id="section1">Section 1</h2>
<h3 id="subsection1">Subsection 1</h3>
<h4 id="subsubsection1">SubSubsection 1</h4>
<h2 id="section2">Section 2</h2>`,
			opts: &TOCOptions{MinLevel: 2, MaxLevel: 3, UseJekyllHTML: false},
			expected: `<h1 id="title">Title</h1>
<p><div class="toc"><ul class="section-nav"><li><a href="#section1">Section 1</a><ul><li><a href="#subsection1">Subsection 1</a></li></ul></li><li><a href="#section2">Section 2</a></li></ul></div></p>
<h2 id="section1">Section 1</h2>
<h3 id="subsection1">Subsection 1</h3>
<h4 id="subsubsection1">SubSubsection 1</h4>
<h2 id="section2">Section 2</h2>`,
		},
		{
			name: "Filter to H1-H2 only",
			html: `<h1 id="title">Title</h1>
<p>{:toc}</p>
<h2 id="section1">Section 1</h2>
<h3 id="subsection1">Subsection 1</h3>
<h2 id="section2">Section 2</h2>`,
			opts: &TOCOptions{MinLevel: 1, MaxLevel: 2, UseJekyllHTML: false},
			expected: `<h1 id="title">Title</h1>
<p><div class="toc"><ul class="section-nav"><li><a href="#title">Title</a><ul><li><a href="#section1">Section 1</a></li><li><a href="#section2">Section 2</a></li></ul></li></ul></div></p>
<h2 id="section1">Section 1</h2>
<h3 id="subsection1">Subsection 1</h3>
<h2 id="section2">Section 2</h2>`,
		},
		{
			name: "Filter removes all headings",
			html: `<h1 id="title">Title</h1>
<p>{:toc}</p>
<h2 id="section1">Section 1</h2>`,
			opts: &TOCOptions{MinLevel: 3, MaxLevel: 4, UseJekyllHTML: false},
			expected: `<h1 id="title">Title</h1>
<p><div class="toc"><ul class="section-nav"><li>No headings found</li></ul></div></p>
<h2 id="section1">Section 1</h2>`,
		},
		{
			name: "Jekyll-compatible HTML structure",
			html: `<h1 id="title">Title</h1>
<p>{:toc}</p>
<h2 id="section1">Section 1</h2>
<h3 id="subsection1">Subsection 1</h3>`,
			opts: &TOCOptions{MinLevel: 1, MaxLevel: 6, UseJekyllHTML: true},
			expected: `<h1 id="title">Title</h1>
<p><ul id="markdown-toc"><li><a href="#title">Title</a><ul><li><a href="#section1">Section 1</a><ul><li><a href="#subsection1">Subsection 1</a></li></ul></li></ul></li></ul></p>
<h2 id="section1">Section 1</h2>
<h3 id="subsection1">Subsection 1</h3>`,
		},
		{
			name: "Jekyll HTML with toc_levels filtering",
			html: `<h1 id="title">Title</h1>
<p>{:toc}</p>
<h2 id="section1">Section 1</h2>
<h3 id="subsection1">Subsection 1</h3>
<h2 id="section2">Section 2</h2>`,
			opts: &TOCOptions{MinLevel: 2, MaxLevel: 2, UseJekyllHTML: true},
			expected: `<h1 id="title">Title</h1>
<p><ul id="markdown-toc"><li><a href="#section1">Section 1</a></li><li><a href="#section2">Section 2</a></li></ul></p>
<h2 id="section1">Section 1</h2>
<h3 id="subsection1">Subsection 1</h3>
<h2 id="section2">Section 2</h2>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := processTOC([]byte(tt.html), tt.opts)
			if err != nil {
				t.Fatalf("Error processing HTML: %v", err)
			}

			if string(result) != tt.expected {
				t.Errorf("Expected:\n%s\n\nGot:\n%s", tt.expected, string(result))
			}
		})
	}
}

func TestTOCLevelsParsing(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expectMin   int
		expectMax   int
	}{
		{
			name:      "String range 1..6",
			input:     "1..6",
			expectMin: 1,
			expectMax: 6,
		},
		{
			name:      "String range 2..3",
			input:     "2..3",
			expectMin: 2,
			expectMax: 3,
		},
		{
			name:      "String range 2..4",
			input:     "2..4",
			expectMin: 2,
			expectMax: 4,
		},
		{
			name:      "Array format",
			input:     []interface{}{2, 3, 4},
			expectMin: 2,
			expectMax: 4,
		},
		{
			name:      "Single level array",
			input:     []interface{}{3},
			expectMin: 3,
			expectMax: 3,
		},
		{
			name:      "Invalid input returns defaults",
			input:     "invalid",
			expectMin: 1,
			expectMax: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			min, max := parseTOCLevels(tt.input)
			if min != tt.expectMin || max != tt.expectMax {
				t.Errorf("Expected (%d, %d), got (%d, %d)", tt.expectMin, tt.expectMax, min, max)
			}
		})
	}
}
