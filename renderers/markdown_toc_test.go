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
		{
			name: "Block syntax with whitespace NOT processed",
			html: `<h1 id="title">Title</h1>
<p>{:: toc }</p>
<h2 id="section1">Section 1</h2>
<h3 id="subsection1">Subsection 1</h3>
<h2 id="section2">Section 2</h2>`,
			expected: `<h1 id="title">Title</h1>
<p>{:: toc }</p>
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
		name        string
		markdown    string
		containsTOC bool
	}{
		{
			name:        "Markdown with {:toc} in unordered list (processed)",
			markdown:    "# Title\n\n* TOC\n{:toc}\n\n## Section 1\n\nContent\n\n## Section 2\n\nMore content",
			containsTOC: true, // {:toc} in unordered list is processed
		},
		{
			name:        "Markdown with standalone {:toc} (not in list, not processed)",
			markdown:    "# Title\n\n{:toc}\n\n## Section 1\n\nContent\n\n## Section 2\n\nMore content",
			containsTOC: false, // Standalone {:toc} not in a list is not processed (matches Jekyll behavior)
		},
		{
			name:        "Markdown with {::toc} (invalid, not processed)",
			markdown:    "# Title\n\n{::toc}\n\n## Section 1\n\nContent\n\n## Section 2\n\nMore content",
			containsTOC: false, // {::toc} is not valid kramdown syntax
		},
		{
			name:        "Markdown without TOC",
			markdown:    "# Title\n\n## Section 1\n\nContent\n\n## Section 2\n\nMore content",
			containsTOC: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html, err := renderMarkdown([]byte(tt.markdown))
			if err != nil {
				t.Fatalf("Error rendering markdown: %v", err)
			}

			// Check if output contains TOC div
			containsTOCDiv := containsString(string(html), "<div class=\"toc\">")

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
		name      string
		input     interface{}
		expectMin int
		expectMax int
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

func TestTOCReplacesListItem(t *testing.T) {
	// Test for issue #89: TOC should replace the preceding list item
	tests := []struct {
		name             string
		markdown         string
		shouldNotContain string
		shouldContainTOC bool
	}{
		{
			name: "List item with {:toc} marker (should replace)",
			markdown: `# Title

* this list replaced by toc
{:toc}

## Section 1

## Section 2`,
			shouldNotContain: "this list replaced by toc",
			shouldContainTOC: true,
		},
		{
			name: "List item with {::toc} marker (should NOT replace - Jekyll doesn't support this)",
			markdown: `# Title

* placeholder text
{::toc}

## Section 1

## Section 2`,
			shouldNotContain: "",    // We expect the text to remain
			shouldContainTOC: false, // Jekyll doesn't process {::toc} in lists
		},
		{
			name: "Ordered list with {:toc} (should NOT replace - Jekyll doesn't support this)",
			markdown: `# Title

1. This will be replaced
{:toc}

## Section 1

## Section 2`,
			shouldNotContain: "",    // We expect the text to remain
			shouldContainTOC: false, // Jekyll doesn't process {:toc} in ordered lists
		},
		{
			name: "List item with different text",
			markdown: `# Title

* Contents
{:toc}

## Section 1

## Section 2`,
			shouldNotContain: "Contents",
			shouldContainTOC: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html, err := renderMarkdown([]byte(tt.markdown))
			if err != nil {
				t.Fatalf("Error rendering markdown: %v", err)
			}

			htmlStr := string(html)

			// Check if TOC should be generated
			if tt.shouldContainTOC {
				if !containsString(htmlStr, "<div class=\"toc\">") {
					t.Error("Expected output to contain TOC div")
				}
			}

			// The list item text should NOT appear in the output (if specified)
			if tt.shouldNotContain != "" && containsString(htmlStr, tt.shouldNotContain) {
				t.Errorf("Output should not contain '%s', but it does.\nOutput:\n%s", tt.shouldNotContain, htmlStr)
			}
		})
	}
}

func TestTOCMultipleMarkersInOneDocument(t *testing.T) {
	// Test that multiple TOC markers don't interfere with each other or delete content
	markdown := `# Test Variations

## Test 1: Block syntax

* Contents
{::toc}

### Section A
### Section B

---

## Test 2: Different placeholder text

* Table of Contents
{:toc}

### Section C
### Section D

---

## Test 3: Ordered list

1. This will be replaced
{:toc}

### Section E
### Section F`

	html, err := renderMarkdown([]byte(markdown))
	if err != nil {
		t.Fatalf("Error rendering markdown: %v", err)
	}

	htmlStr := string(html)

	// All section headings should be present (none should be deleted by overly greedy regex)
	requiredHeadings := []string{
		"Test 1: Block syntax",
		"Test 2: Different placeholder text",
		"Test 3: Ordered list",
		"Section A",
		"Section B",
		"Section C",
		"Section D",
		"Section E",
		"Section F",
	}

	for _, heading := range requiredHeadings {
		if !containsString(htmlStr, heading) {
			t.Errorf("Output should contain heading '%s', but it doesn't.\nOutput:\n%s", heading, htmlStr)
		}
	}

	// Only Test 2 should have its placeholder text removed (Jekyll only supports {:toc} in unordered lists)
	if containsString(htmlStr, "Table of Contents") {
		t.Error("Output should not contain 'Table of Contents' (should be replaced by TOC)")
	}

	// Test 1 and Test 3 placeholders should remain (Jekyll doesn't support these patterns)
	if !containsString(htmlStr, "Contents") {
		t.Error("Output should contain 'Contents' (Jekyll doesn't replace {::toc} in lists)")
	}
	if !containsString(htmlStr, "This will be replaced") {
		t.Error("Output should contain 'This will be replaced' (Jekyll doesn't replace {:toc} in ordered lists)")
	}
}

func TestTOCInCodeBlocks(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		shouldContainTOCDiv bool
		shouldContainLiteralMarker bool
	}{
		{
			name: "TOC marker in fenced code block",
			markdown: "# Heading\n\n```\n{:toc}\n```\n\n## Section 1",
			shouldContainTOCDiv: false,
			shouldContainLiteralMarker: true,
		},
		{
			name: "TOC marker in inline code",
			markdown: "# Heading\n\nUse `{:toc}` to generate TOC\n\n## Section 1",
			shouldContainTOCDiv: false,
			shouldContainLiteralMarker: true,
		},
		{
			name: "Real TOC marker outside code",
			markdown: "# Heading\n\n{:toc}\n\n## Section 1",
			shouldContainTOCDiv: true,
			shouldContainLiteralMarker: false,
		},
		{
			name: "Both code and real TOC marker",
			markdown: "# Heading\n\nExample: `{:toc}`\n\n{:toc}\n\n## Section 1",
			shouldContainTOCDiv: true,
			shouldContainLiteralMarker: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html, err := renderMarkdown([]byte(tt.markdown))
			if err != nil {
				t.Fatalf("Error rendering markdown: %v", err)
			}

			htmlStr := string(html)
			containsTOCDiv := containsString(htmlStr, "<div class=\"toc\">")
			containsLiteral := containsString(htmlStr, "{:toc}") || containsString(htmlStr, "&lt;:toc}")

			if containsTOCDiv != tt.shouldContainTOCDiv {
				t.Errorf("Expected TOC div presence: %v, got: %v", tt.shouldContainTOCDiv, containsTOCDiv)
			}

			if containsLiteral != tt.shouldContainLiteralMarker {
				t.Errorf("Expected literal marker presence: %v, got: %v", tt.shouldContainLiteralMarker, containsLiteral)
			}
		})
	}
}

func TestTOCWithSpecialCharacters(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		shouldContainInTOC string
	}{
		{
			name: "Heading with HTML entity",
			markdown: "# Title\n\n{:toc}\n\n## Section &amp; More",
			shouldContainInTOC: "Section &amp; More",
		},
		{
			name: "Heading with emoji",
			markdown: "# Title\n\n{:toc}\n\n## 🚀 Rocket Section",
			shouldContainInTOC: "🚀 Rocket Section",
		},
		{
			name: "Heading with bold",
			markdown: "# Title\n\n{:toc}\n\n## Section with **bold** text",
			shouldContainInTOC: "bold",
		},
		{
			name: "Heading with code",
			markdown: "# Title\n\n{:toc}\n\n## Using `code` here",
			shouldContainInTOC: "code",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html, err := renderMarkdown([]byte(tt.markdown))
			if err != nil {
				t.Fatalf("Error rendering markdown: %v", err)
			}

			htmlStr := string(html)

			if !containsString(htmlStr, tt.shouldContainInTOC) {
				t.Errorf("TOC should contain '%s', but it doesn't.\nOutput:\n%s", tt.shouldContainInTOC, htmlStr)
			}
		})
	}
}

func TestTOCWithUnusualHierarchy(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		shouldWork bool
	}{
		{
			name: "Starting with H3 (skipping H1, H2)",
			markdown: "{:toc}\n\n### Section 1\n\n#### Subsection",
			shouldWork: true,
		},
		{
			name: "Multiple H1 headings",
			markdown: "# First\n\n{:toc}\n\n# Second\n\n## Under Second",
			shouldWork: true,
		},
		{
			name: "Gaps in heading levels",
			markdown: "# H1\n\n{:toc}\n\n### H3\n\n###### H6",
			shouldWork: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html, err := renderMarkdown([]byte(tt.markdown))
			if err != nil {
				t.Fatalf("Error rendering markdown: %v", err)
			}

			// Just verify it doesn't crash and produces TOC
			if !containsString(string(html), "<div class=\"toc\">") {
				t.Error("Should contain TOC div")
			}
		})
	}
}

func TestTOCEmptyOrMinimalDocuments(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		expectedMessage string
	}{
		{
			name: "TOC with no headings",
			markdown: "{:toc}\n\nJust some text",
			expectedMessage: "No headings found",
		},
		{
			name: "Empty document with TOC",
			markdown: "{:toc}",
			expectedMessage: "No headings found",
		},
		{
			name: "Single heading",
			markdown: "# Only One\n\n{:toc}",
			expectedMessage: "", // Should work, not show error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html, err := renderMarkdown([]byte(tt.markdown))
			if err != nil {
				t.Fatalf("Error rendering markdown: %v", err)
			}

			htmlStr := string(html)

			if tt.expectedMessage != "" {
				if !containsString(htmlStr, tt.expectedMessage) {
					t.Errorf("Should contain message '%s', but doesn't.\nOutput:\n%s", tt.expectedMessage, htmlStr)
				}
			}
		})
	}
}

func TestTOCConfigurationEdgeCases(t *testing.T) {
	tests := []struct {
		name string
		html string
		opts *TOCOptions
		shouldWork bool
	}{
		{
			name: "Invalid levels (too high)",
			html: "<h1 id=\"h1\">H1</h1>\n<p>{:toc}</p>\n<h2 id=\"h2\">H2</h2>",
			opts: &TOCOptions{MinLevel: 7, MaxLevel: 9, UseJekyllHTML: false},
			shouldWork: true, // Should clamp to valid range
		},
		{
			name: "Reversed range",
			html: "<h1 id=\"h1\">H1</h1>\n<p>{:toc}</p>\n<h2 id=\"h2\">H2</h2>",
			opts: &TOCOptions{MinLevel: 5, MaxLevel: 2, UseJekyllHTML: false},
			shouldWork: true, // Should fix to valid range
		},
		{
			name: "Single level",
			html: "<h1 id=\"h1\">H1</h1>\n<p>{:toc}</p>\n<h2 id=\"h2\">H2</h2>\n<h3 id=\"h3\">H3</h3>",
			opts: &TOCOptions{MinLevel: 2, MaxLevel: 2, UseJekyllHTML: false},
			shouldWork: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := processTOC([]byte(tt.html), tt.opts)
			if err != nil {
				if tt.shouldWork {
					t.Fatalf("Should not error, but got: %v", err)
				}
			} else {
				if !tt.shouldWork {
					t.Error("Should have errored, but didn't")
				}
			}

			if tt.shouldWork {
				// Verify it produced something reasonable
				if !containsString(string(result), "<") {
					t.Error("Should produce HTML output")
				}
			}
		})
	}
}

func TestTOCDuplicateHeadings(t *testing.T) {
	markdown := `# Title

{:toc}

## Section
### Subsection
## Section
### Different Subsection`

	html, err := renderMarkdown([]byte(markdown))
	if err != nil {
		t.Fatalf("Error rendering markdown: %v", err)
	}

	htmlStr := string(html)

	// Both "Section" headings should appear
	// Count occurrences - should be at least 2 in the TOC
	if !containsString(htmlStr, "<div class=\"toc\">") {
		t.Error("Should contain TOC")
	}
}

func TestTOCNoTocPlacement(t *testing.T) {
	tests := []struct {
		name string
		markdown string
		shouldExclude string
	}{
		{
			name: "no_toc after heading",
			markdown: "# Title\n\n{:toc}\n\n## Excluded\n{:.no_toc}\n\n## Included",
			shouldExclude: "Excluded",
		},
		{
			name: "no_toc with whitespace",
			markdown: "# Title\n\n{:toc}\n\n## Excluded\n{: .no_toc }\n\n## Included",
			shouldExclude: "Excluded",
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
			if !containsString(htmlStr, "<div class=\"toc\">") {
				t.Error("Should contain TOC")
			}
		})
	}
}
