package plugins

import (
	"bytes"
	"io"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// mathPlugin emulates kramdown's math engine for LaTeX math support.
// It converts $$...$$ delimiters to \(...\) for inline math and \[...\] for display math,
// matching the modern kramdown output format that works with MathJax v2, v3, and KaTeX.
type mathPlugin struct{ plugin }

var (
	// Match $$...$$ for both inline and display math
	mathDelimiterPattern = regexp.MustCompile(`\$\$([^$]+?)\$\$`)
)

// PostRender processes the HTML output to convert math delimiters.
func (p mathPlugin) PostRender(b []byte) ([]byte, error) {
	return processMathInHTML(b), nil
}

// processMathInHTML walks through HTML tokens and converts math expressions.
// It uses a more sophisticated approach than ApplyToHTMLText to detect display vs inline math.
func processMathInHTML(doc []byte) []byte {
	z := html.NewTokenizer(bytes.NewReader(doc))
	buf := new(bytes.Buffer)
	inBody := false
	hasBody := false
	inParagraph := false
	paragraphHTML := new(bytes.Buffer)

outer:
	for {
		tt := z.Next()
		raw := z.Raw()
		switch tt {
		case html.ErrorToken:
			if z.Err() == io.EOF {
				break outer
			}
			// On error, return original content
			return doc
		case html.StartTagToken:
			tn, _ := z.TagName()
			tagName := string(tn)
			if tagName == "body" {
				inBody = true
				hasBody = true
			} else if (inBody || !hasBody) && tagName == "p" {
				inParagraph = true
				paragraphHTML.Reset()
				// Write the opening <p> tag to main buffer
				_, _ = buf.Write(raw)
				continue outer
			} else if inParagraph {
				// Accumulate all HTML within the paragraph
				_, _ = paragraphHTML.Write(raw)
				continue outer
			}
		case html.EndTagToken:
			tn, _ := z.TagName()
			tagName := string(tn)
			if tagName == "body" {
				inBody = false
			} else if tagName == "p" && inParagraph {
				// Process the entire paragraph content at once
				processedContent := convertMathInHTML(paragraphHTML.String(), true)
				_, _ = buf.WriteString(processedContent)
				inParagraph = false
				// Write the closing </p> tag
				_, _ = buf.Write(raw)
				continue outer
			} else if inParagraph {
				// Accumulate all HTML within the paragraph
				_, _ = paragraphHTML.Write(raw)
				continue outer
			}
		case html.TextToken:
			if inParagraph {
				// Accumulate paragraph HTML including text
				_, _ = paragraphHTML.Write(raw)
				continue outer
			} else if inBody || !hasBody {
				// Process text outside paragraphs (shouldn't have display math)
				text := string(z.Text())
				processed := convertMathDelimiters(text, false)
				_, _ = buf.WriteString(processed)
				continue outer
			}
		default:
			if inParagraph {
				// Accumulate all tokens within paragraph
				_, _ = paragraphHTML.Write(raw)
				continue outer
			}
		}
		_, _ = buf.Write(raw)
	}
	return buf.Bytes()
}

// convertMathInHTML processes HTML content and converts math delimiters.
// This handles HTML with nested tags properly.
func convertMathInHTML(htmlContent string, allowDisplay bool) string {
	// First, extract all text content to determine context
	allText := extractAllText(htmlContent)

	z := html.NewTokenizer(bytes.NewReader([]byte(htmlContent)))
	buf := new(bytes.Buffer)

	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			if z.Err() == io.EOF {
				return buf.String()
			}
			return htmlContent
		case html.TextToken:
			text := string(z.Text())
			// Use the full paragraph text context to determine if math is standalone
			processed := convertMathDelimitersWithContext(text, allText, allowDisplay)
			_, _ = buf.WriteString(processed)
		default:
			_, _ = buf.Write(z.Raw())
		}
	}
}

// extractAllText extracts all text content from HTML, stripping tags.
func extractAllText(htmlContent string) string {
	z := html.NewTokenizer(bytes.NewReader([]byte(htmlContent)))
	var textParts []string

	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			break
		}
		if tt == html.TextToken {
			textParts = append(textParts, string(z.Text()))
		}
	}
	return strings.Join(textParts, "")
}

// convertMathDelimitersWithContext converts $$...$$ using context from the full text.
func convertMathDelimitersWithContext(text, fullText string, allowDisplay bool) string {
	if !strings.Contains(text, "$$") {
		return text
	}

	return mathDelimiterPattern.ReplaceAllStringFunc(text, func(match string) string {
		// Extract the content between $$...$$
		content := strings.TrimPrefix(match, "$$")
		content = strings.TrimSuffix(content, "$$")

		// Determine if this is display or inline math using the full text context
		// Display math indicators:
		// 1. Content contains newlines (multiline expressions are always display)
		// 2. The match is the only content in the full paragraph text (standalone)
		// 3. allowDisplay is true (we're in a paragraph context)
		isDisplay := allowDisplay && (
			strings.Contains(content, "\n") ||
				isStandaloneInText(fullText, match))

		if isDisplay {
			// Display math: \[...\]
			return "\\[" + content + "\\]"
		}
		// Inline math: \(...\)
		return "\\(" + content + "\\)"
	})
}

// convertMathDelimiters converts $$...$$ to LaTeX delimiters.
// If allowDisplay is true, it can detect and convert display math.
func convertMathDelimiters(text string, allowDisplay bool) string {
	if !strings.Contains(text, "$$") {
		return text
	}

	return mathDelimiterPattern.ReplaceAllStringFunc(text, func(match string) string {
		// Extract the content between $$...$$
		content := strings.TrimPrefix(match, "$$")
		content = strings.TrimSuffix(content, "$$")

		// Determine if this is display or inline math
		// Display math indicators:
		// 1. Content contains newlines (multiline expressions are always display)
		// 2. The match is the only content in the text (standalone in paragraph)
		// 3. allowDisplay is true (we're in a paragraph context)
		isDisplay := allowDisplay && (
			strings.Contains(content, "\n") ||
				isStandaloneInText(text, match))

		if isDisplay {
			// Display math: \[...\]
			return "\\[" + content + "\\]"
		}
		// Inline math: \(...\)
		return "\\(" + content + "\\)"
	})
}

// isStandaloneInText checks if the math expression is the only non-whitespace content in the text.
func isStandaloneInText(text, match string) bool {
	// Remove the match from the text
	textWithoutMatch := strings.Replace(text, match, "", 1)
	// If what remains is only whitespace, the match was standalone
	return len(strings.TrimSpace(textWithoutMatch)) == 0
}
