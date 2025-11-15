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
	tocPatternBlock = regexp.MustCompile(`\{::\s*toc\s*\}`)
	// Match {:.no_toc} with optional whitespace - used to exclude headings from TOC
	noTocPattern = regexp.MustCompile(`\{:\s*\.no_toc\s*\}`)
	// Match <ul> list containing inline TOC marker {:toc} (Jekyll behavior)
	// Jekyll ONLY supports {:toc} (not {::toc}) in unordered lists (not ordered lists)
	// We use [^<]* to match only text before the marker, avoiding matching across elements
	tocUlListPatternInline = regexp.MustCompile(`<ul>\s*<li>[^<]*\{:\s*toc\s*\}\s*</li>\s*</ul>`)
)

// TOCOptions configures TOC generation behavior
type TOCOptions struct {
	MinLevel      int  // Minimum heading level to include (1-6)
	MaxLevel      int  // Maximum heading level to include (1-6)
	UseJekyllHTML bool // Use Jekyll-compatible HTML structure
}

func renderMarkdown(md []byte) ([]byte, error) {
	return renderMarkdownWithOptions(md, nil)
}

func renderMarkdownWithOptions(md []byte, opts *TOCOptions) ([]byte, error) {
	// Set default options if not provided
	if opts == nil {
		opts = &TOCOptions{
			MinLevel:      1,
			MaxLevel:      6,
			UseJekyllHTML: false,
		}
	}
	// Ensure valid level ranges
	if opts.MinLevel < 1 {
		opts.MinLevel = 1
	}
	if opts.MaxLevel > 6 {
		opts.MaxLevel = 6
	}
	if opts.MinLevel > opts.MaxLevel {
		opts.MinLevel = 1
		opts.MaxLevel = 6
	}

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
		html, err = processTOC(html, opts)
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

// search HTML for markdown attributes, and process if found
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
			shouldProcess, mode := hasMarkdownAttr(z)
			if mode != "" {
				// If we have a markdown attribute, always strip it from the output
				_, err := buf.Write(stripMarkdownAttr(z.Raw()))
				if err != nil {
					return nil, err
				}

				if shouldProcess {
					// Only process if the mode is one that enables processing
					if err := processInnerMarkdown(buf, z, mode); err != nil {
						return nil, err
					}
					// the above leaves z set to the end token
					// fall through to render it
				} else {
					// For markdown="0", just copy the content without processing
					if err := copyContent(buf, z); err != nil {
						return nil, err
					}
				}
				// fall through to write the end tag
			}
		}
		_, err := buf.Write(z.Raw())
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func hasMarkdownAttr(z *html.Tokenizer) (bool, string) {
	for {
		k, v, more := z.TagAttr()
		if string(k) == "markdown" {
			value := string(v)
			switch value {
			case "1", "block", "span":
				return true, value
			case "0":
				return false, value
			default:
				// Invalid or unknown markdown attribute value
				return false, ""
			}
		}
		if !more {
			return false, ""
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
func processTOC(content []byte, opts *TOCOptions) ([]byte, error) {
	// Generate the TOC HTML
	toc, err := generateTOC(content, opts)
	if err != nil {
		return nil, err
	}

	// First, replace unordered list elements containing {:toc} markers (Jekyll behavior)
	// This handles the pattern: <ul><li>text{:toc}</li></ul>
	// Jekyll ONLY supports this for {:toc} in <ul>, not {::toc} or <ol>
	result := tocUlListPatternInline.ReplaceAll(content, []byte(toc))

	// Then replace any remaining standalone TOC markers
	result = tocPatternInline.ReplaceAll(result, []byte(toc))
	result = tocPatternBlock.ReplaceAll(result, []byte(toc))

	// Remove no_toc markers from the final output
	result = noTocPattern.ReplaceAll(result, []byte(""))

	return result, nil
}

// generateTOC parses HTML content and creates a table of contents
func generateTOC(content []byte, opts *TOCOptions) (string, error) {
	// Parse the HTML document
	doc, err := html.Parse(bytes.NewReader(content))
	if err != nil {
		return "", err
	}

	// Extract headings
	headings := extractHeadings(doc)

	// Filter headings by level if opts is provided
	if opts != nil {
		filtered := make([]*TOCEntry, 0, len(headings))
		for _, h := range headings {
			if h.Level >= opts.MinLevel && h.Level <= opts.MaxLevel {
				filtered = append(filtered, h)
			}
		}
		headings = filtered
	}

	if len(headings) == 0 {
		emptyMsg := "<div class=\"toc\"><ul class=\"section-nav\"><li>No headings found</li></ul></div>"
		if opts != nil && opts.UseJekyllHTML {
			emptyMsg = "<ul id=\"markdown-toc\"><li>No headings found</li></ul>"
		}
		return emptyMsg, nil
	}

	// Create a nested TOC structure
	tocEntries := organizeTOCHierarchy(headings)

	// Render the TOC as HTML
	var buf bytes.Buffer
	if opts != nil && opts.UseJekyllHTML {
		buf.WriteString("<ul id=\"markdown-toc\">")
	} else {
		buf.WriteString("<div class=\"toc\"><ul class=\"section-nav\">")
	}
	renderTOCEntries(&buf, tocEntries)
	if opts != nil && opts.UseJekyllHTML {
		buf.WriteString("</ul>")
	} else {
		buf.WriteString("</ul></div>")
	}

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

// called once a markdown attribute is detected.
// Collects the HTML tokens into a string, applies markdown to them,
// and writes the result
func processInnerMarkdown(w io.Writer, z *html.Tokenizer, mode string) error {
	buf := new(bytes.Buffer)
	depth := 1
loop:
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			err := z.Err()
			if err == io.EOF {
				return utils.WrapError(err,
					"unexpected EOF while processing markdown attribute. "+
						"Common causes: unclosed HTML tags (use <br/> instead of <br>), "+
						"or mismatched opening/closing tags")
			}
			return err
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

	var html []byte
	var err error

	switch mode {
	case "span":
		// For span mode, process inline markdown only
		html, err = _renderMarkdownSpan(buf.Bytes())
	case "block", "1":
		// For block and 1 modes, process full markdown
		html, err = _renderMarkdown(buf.Bytes())
	default:
		// Should never happen as hasMarkdownAttr already filtered
		html = buf.Bytes()
	}

	if err != nil {
		return err
	}
	_, err = w.Write(html)
	return err
}

func _renderMarkdownSpan(md []byte) ([]byte, error) {
	// For span-level processing, we don't want to create block-level elements like paragraphs
	// Instead, we just want inline formatting (bold, italic, links, etc.)
	params := blackfriday.HTMLRendererParameters{
		Flags: blackfridayFlags,
	}
	renderer := blackfriday.NewHTMLRenderer(params)

	// Use only inline-level extensions for span mode
	inlineExtensions := blackfriday.NoIntraEmphasis |
		blackfriday.Autolink |
		blackfriday.Strikethrough |
		blackfriday.BackslashLineBreak

	// Process the content without creating paragraphs - we're handling inline elements
	content := bytes.TrimSpace(md)
	html := blackfriday.Run(
		content,
		blackfriday.WithRenderer(renderer),
		blackfriday.WithExtensions(inlineExtensions),
	)

	// Remove any potential wrapping paragraph tags that blackfriday might add
	html = bytes.TrimPrefix(html, []byte("<p>"))
	html = bytes.TrimSuffix(html, []byte("</p>\n"))

	return html, nil
}

// copyContent copies HTML content without processing markdown
func copyContent(w io.Writer, z *html.Tokenizer) error {
	buf := new(bytes.Buffer)
	depth := 1
loop:
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			err := z.Err()
			if err == io.EOF {
				return utils.WrapError(err,
					"unexpected EOF while processing markdown=\"0\" attribute. "+
						"Common causes: unclosed HTML tags (use <br/> instead of <br>), "+
						"or mismatched opening/closing tags")
			}
			return err
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
	_, err := w.Write(buf.Bytes())
	return err
}
