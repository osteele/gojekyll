package renderers

import (
	"bytes"
	"fmt"
	stdhtml "html"
	"io"
	"regexp"
	"strings"

	blackfriday "github.com/danog/blackfriday/v2"
	"github.com/osteele/gojekyll/utils"
	"golang.org/x/net/html"
)

const blackfridayFlags = 0 |
	blackfriday.UseXHTML |
	blackfriday.Smartypants |
	blackfriday.SmartypantsFractions |
	blackfriday.SmartypantsDashes |
	blackfriday.SmartypantsLatexDashes |
	blackfriday.FootnoteReturnLinks

const blackfridayExtensions = 0 |
	blackfriday.NoIntraEmphasis |
	blackfriday.Tables |
	blackfriday.FencedCode |
	blackfriday.Autolink |
	blackfriday.Strikethrough |
	blackfriday.SpaceHeadings |
	blackfriday.HeadingIDs |
	blackfriday.BackslashLineBreak |
	blackfriday.DefinitionLists |
	blackfriday.NoEmptyLineBeforeBlock |
	// added relative to commonExtensions
	blackfriday.AutoHeadingIDs |
	blackfriday.Footnotes

// TOC marker patterns from kramdown
var (
	// Match {:toc} with optional whitespace
	tocPatternInline = regexp.MustCompile(`\{:\s*toc\s*\}`)
	// Match {::toc} with optional whitespace
	tocPatternBlock  = regexp.MustCompile(`\{::\s*toc\s*\}`)
	// Match {:.no_toc} with optional whitespace - used to exclude headings from TOC
	noTocPattern = regexp.MustCompile(`\{:\s*\.no_toc\s*\}`)
)

func renderMarkdown(md []byte) ([]byte, error) {
	params := blackfriday.HTMLRendererParameters{
		Flags: blackfridayFlags,
	}
	renderer := blackfriday.NewHTMLRenderer(params)
	html := blackfriday.Run(
		md,
		blackfriday.WithRenderer(renderer),
		blackfriday.WithExtensions(blackfridayExtensions),
	)
	html, err := renderInnerMarkdown(html)
	if err != nil {
		return nil, utils.WrapError(err, "markdown")
	}

	// Process TOC markers if they exist
	if tocPatternInline.Match(html) || tocPatternBlock.Match(html) {
		html, err = processTOC(html)
		if err != nil {
			return nil, utils.WrapError(err, "toc generation")
		}
	}
	return html, nil
}

func _renderMarkdown(md []byte) ([]byte, error) {
	params := blackfriday.HTMLRendererParameters{
		Flags: blackfridayFlags,
	}
	renderer := blackfriday.NewHTMLRenderer(params)
	html := blackfriday.Run(
		md,
		blackfriday.WithRenderer(renderer),
		blackfriday.WithExtensions(blackfridayExtensions),
	)
	return html, nil
}

// search HTML for markdown=1, and process if found
func renderInnerMarkdown(b []byte) ([]byte, error) {
	z := html.NewTokenizer(bytes.NewReader(b))
	buf := new(bytes.Buffer)
outer:
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			if z.Err() == io.EOF {
				break outer
			}
			return nil, z.Err()
		case html.StartTagToken:
			if hasMarkdownAttr(z) {
				_, err := buf.Write(stripMarkdownAttr(z.Raw()))
				if err != nil {
					return nil, err
				}
				if err := processInnerMarkdown(buf, z); err != nil {
					return nil, err
				}
				// the above leaves z set to the end token
				// fall through to render it
			}
		}
		_, err := buf.Write(z.Raw())
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func hasMarkdownAttr(z *html.Tokenizer) bool {
	for {
		k, v, more := z.TagAttr()
		switch {
		case string(k) == "markdown" && string(v) == "1":
			return true
		case !more:
			return false
		}
	}
}

var markdownAttrRE = regexp.MustCompile(`\s*markdown\s*=[^\s>]*\s*`)

// return the text of a start tag, w/out the markdown attribute
func stripMarkdownAttr(tag []byte) []byte {
	tag = markdownAttrRE.ReplaceAll(tag, []byte(" "))
	tag = bytes.Replace(tag, []byte(" >"), []byte(">"), 1)
	return tag
}

// TOCEntry represents a heading in the table of contents
type TOCEntry struct {
	ID       string
	Level    int
	Text     string
	Children []*TOCEntry
}

// processTOC parses HTML content and replaces TOC markers with generated table of contents
func processTOC(content []byte) ([]byte, error) {
	// Generate the TOC HTML
	toc, err := generateTOC(content)
	if err != nil {
		return nil, err
	}

	// Replace TOC markers with the generated TOC
	result := tocPatternInline.ReplaceAll(content, []byte(toc))
	result = tocPatternBlock.ReplaceAll(result, []byte(toc))

	// Remove no_toc markers from the final output
	result = noTocPattern.ReplaceAll(result, []byte(""))

	return result, nil
}

// generateTOC parses HTML content and creates a table of contents
func generateTOC(content []byte) (string, error) {
	// Parse the HTML document
	doc, err := html.Parse(bytes.NewReader(content))
	if err != nil {
		return "", err
	}

	// Extract headings
	headings := extractHeadings(doc)
	if len(headings) == 0 {
		return "<div class=\"toc\"><ul class=\"section-nav\"><li>No headings found</li></ul></div>", nil
	}

	// Create a nested TOC structure
	tocEntries := organizeTOCHierarchy(headings)

	// Render the TOC as HTML
	var buf bytes.Buffer
	buf.WriteString("<div class=\"toc\"><ul class=\"section-nav\">")
	renderTOCEntries(&buf, tocEntries)
	buf.WriteString("</ul></div>")

	return buf.String(), nil
}

// extractHeadings finds all heading elements (h1-h6) in the HTML document
func extractHeadings(n *html.Node) []*TOCEntry {
	var headings []*TOCEntry

	var extract func(*html.Node)
	extract = func(n *html.Node) {
		// Check if this is a heading element
		if n.Type == html.ElementNode && strings.HasPrefix(n.Data, "h") && len(n.Data) == 2 {
			// Parse the heading level (h1-h6)
			level := int(n.Data[1] - '0')
			if level >= 1 && level <= 6 {
				// Extract the heading ID
				id := ""
				for _, attr := range n.Attr {
					if attr.Key == "id" {
						id = attr.Val
						break
					}
				}

				// Extract the heading text and check for no_toc marker
				html := renderNodeToString(n)
				
				// Skip headings with the no_toc marker
				if noTocPattern.MatchString(html) {
					return
				}

				// Extract the heading text
				text := extractTextContent(n)

				// Create a TOC entry
				headings = append(headings, &TOCEntry{
					ID:    id,
					Level: level,
					Text:  text,
				})
			}
		}

		// Recursively process child nodes
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extract(c)
		}
	}

	extract(n)
	return headings
}

// renderNodeToString converts an HTML node to a string
func renderNodeToString(n *html.Node) string {
	var buf bytes.Buffer
	err := html.Render(&buf, n)
	if err != nil {
		return ""
	}
	return buf.String()
}

// extractTextContent gets the plain text from an HTML node
func extractTextContent(n *html.Node) string {
	var text string

	var extract func(*html.Node)
	extract = func(n *html.Node) {
		if n.Type == html.TextNode {
			text += n.Data
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extract(c)
		}
	}

	extract(n)
	return strings.TrimSpace(text)
}

// organizeTOCHierarchy organizes TOC entries into a hierarchical structure
func organizeTOCHierarchy(headings []*TOCEntry) []*TOCEntry {
	if len(headings) == 0 {
		return nil
	}

	// Find the minimum heading level to use as the top level
	minLevel := 6
	for _, h := range headings {
		if h.Level < minLevel {
			minLevel = h.Level
		}
	}

	// Create a root level to hold all entries
	var root []*TOCEntry
	var stack []*TOCEntry

	for _, h := range headings {
		// Create a new entry
		entry := &TOCEntry{
			ID:    h.ID,
			Level: h.Level,
			Text:  h.Text,
		}

		// Pop the stack until we find a parent with a lower level
		for len(stack) > 0 && stack[len(stack)-1].Level >= h.Level {
			stack = stack[:len(stack)-1]
		}

		// If the stack is empty, this is a top-level entry
		if len(stack) == 0 {
			root = append(root, entry)
		} else {
			// Add this entry as a child of the last item on the stack
			parent := stack[len(stack)-1]
			parent.Children = append(parent.Children, entry)
		}

		// Push this entry onto the stack
		stack = append(stack, entry)
	}

	return root
}

// renderTOCEntries renders TOC entries as HTML
func renderTOCEntries(buf *bytes.Buffer, entries []*TOCEntry) {
	for _, entry := range entries {
		buf.WriteString("<li>")

		// Add a link to the heading if it has an ID
		if entry.ID != "" {
			fmt.Fprintf(buf, "<a href=\"#%s\">%s</a>", entry.ID, stdhtml.EscapeString(entry.Text))
		} else {
			buf.WriteString(stdhtml.EscapeString(entry.Text))
		}

		// Recursively render children if any
		if len(entry.Children) > 0 {
			buf.WriteString("<ul>")
			renderTOCEntries(buf, entry.Children)
			buf.WriteString("</ul>")
		}

		buf.WriteString("</li>")
	}
}

// Used inside markdown=1.
// TODO Instead of this approach, only count tags that match the start
// tag. For example, if <div markdown="1"> kicked off the inner markdown,
// count the div depth.
var notATagRE = regexp.MustCompile(`@|(https?|ftp):`)

// Called once markdown="1" attribute is detected.
// Collects the HTML tokens into a string, applies markdown to them,
// and writes the result
func processInnerMarkdown(w io.Writer, z *html.Tokenizer) error {
	buf := new(bytes.Buffer)
	depth := 1
loop:
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return z.Err()
		case html.StartTagToken:
			if !notATagRE.Match(z.Raw()) {
				depth++
			}
		case html.EndTagToken:
			depth--
			if depth == 0 {
				break loop
			}
		}
		_, err := buf.Write(z.Raw())
		if err != nil {
			return err
		}
	}
	html, err := _renderMarkdown(buf.Bytes())
	if err != nil {
		return err
	}
	_, err = w.Write(html)
	return err
}
