package renderers

import (
	"bytes"
	"fmt"
	stdhtml "html"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// TOC marker patterns from kramdown
// These are used for quick detection and text matching, not for HTML parsing
var (
	// Match {:toc} with optional whitespace - ONLY valid TOC syntax in kramdown
	// Note: {::toc} is NOT valid kramdown syntax and should not be processed
	tocPatternInline = regexp.MustCompile(`\{:\s*toc\s*\}`)
	// Match {:.no_toc} with optional whitespace - used to exclude headings from TOC
	noTocPattern = regexp.MustCompile(`\{:\s*\.no_toc\s*\}`)
)

// TOCOptions configures TOC generation behavior
type TOCOptions struct {
	MinLevel      int  // Minimum heading level to include (1-6)
	MaxLevel      int  // Maximum heading level to include (1-6)
	UseJekyllHTML bool // Use Jekyll-compatible HTML structure
}

// TOCEntry represents a heading in the table of contents
type TOCEntry struct {
	ID       string
	Level    int
	Text     string
	Children []*TOCEntry
}

// MarkerType identifies the context of a TOC marker in the HTML
type MarkerType int

const (
	MarkerStandalone MarkerType = iota
	MarkerInUnorderedList
	MarkerInOrderedList
	MarkerInCodeBlock
)

// MarkerContext describes a TOC marker's location and context in the DOM
type MarkerContext struct {
	Type       MarkerType
	Node       *html.Node  // The text node containing the marker
	ParentList *html.Node  // The <ul> or <ol> node if in a list
	MarkerText string      // The actual marker text: "{:toc}" or "{::toc}"
	IsBlock    bool        // true for {::toc}, false for {:toc}
}

// processTOC parses HTML content and replaces TOC markers with generated table of contents
// Uses DOM-based approach for robust handling of all marker contexts
// Note: Only {:toc} is valid kramdown syntax; {::toc} is not processed
func processTOC(content []byte, opts *TOCOptions) ([]byte, error) {
	// Quick check: if no TOC markers exist, skip processing
	if !tocPatternInline.Match(content) {
		// Still need to remove {:.no_toc} markers even if no TOC
		return noTocPattern.ReplaceAll(content, []byte("")), nil
	}

	// Parse the HTML into a DOM tree
	doc, err := html.Parse(bytes.NewReader(content))
	if err != nil {
		return nil, err
	}

	// Find all TOC markers and classify them by context
	markers := findTOCMarkersInDOM(doc)

	// If no markers found (e.g., all in code blocks), return original content
	if len(markers) == 0 {
		// Remove {:.no_toc} markers and return
		return noTocPattern.ReplaceAll(content, []byte("")), nil
	}

	// Generate the TOC HTML
	toc, err := generateTOC(content, opts)
	if err != nil {
		return nil, err
	}

	// Replace each marker in the DOM based on its context
	// Process in reverse order to avoid invalidating node references
	// when modifying the tree
	for i := len(markers) - 1; i >= 0; i-- {
		marker := markers[i]
		if err := replaceTOCMarkerInDOM(marker, toc); err != nil {
			return nil, fmt.Errorf("failed to replace TOC marker: %w", err)
		}
	}

	// Render the modified DOM back to HTML
	var buf bytes.Buffer
	if err := html.Render(&buf, doc); err != nil {
		return nil, err
	}

	// Extract the body content (html.Render wraps in <html><head></head><body>...</body></html>)
	result := extractBodyContent(buf.Bytes())

	// Remove {:.no_toc} markers from the final output
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
				htmlStr := renderNodeToString(n)

				// Check for {:.no_toc} marker in the heading's text content
				// The marker could be:
				// 1. In the heading text itself (before blackfriday strips it)
				// 2. In a following paragraph/text node (Jekyll behavior)
				if noTocPattern.MatchString(htmlStr) {
					return
				}

				// Also check the next sibling element for a no_toc marker
				// Jekyll allows the marker on the line after the heading
				// Skip over whitespace text nodes to find the next actual element
				nextElem := n.NextSibling
				for nextElem != nil {
					if nextElem.Type == html.ElementNode {
						// Found the next element, check if it contains {:.no_toc}
						nextHTML := renderNodeToString(nextElem)
						if noTocPattern.MatchString(nextHTML) {
							return
						}
						break
					} else if nextElem.Type == html.TextNode {
						// Skip whitespace-only text nodes
						text := strings.TrimSpace(nextElem.Data)
						if text != "" {
							// Found non-whitespace text, stop searching
							break
						}
					}
					nextElem = nextElem.NextSibling
				}

				// Extract the heading text (removing any remaining markers)
				text := extractTextContent(n)
				// Remove any {:.no_toc} markers from the text
				text = noTocPattern.ReplaceAllString(text, "")
				text = strings.TrimSpace(text)

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

// findTOCMarkersInDOM walks the HTML DOM tree to find all TOC markers and classify them by context
func findTOCMarkersInDOM(doc *html.Node) []*MarkerContext {
	var markers []*MarkerContext

	var walk func(*html.Node)
	walk = func(n *html.Node) {
		// Only process text nodes
		if n.Type == html.TextNode {
			text := n.Data

			// Check for TOC markers - only {:toc} is valid kramdown syntax
			if tocPatternInline.MatchString(text) {
				// Determine the marker type by examining parent nodes
				ctx := classifyMarkerContext(n, text)
				if ctx != nil {
					markers = append(markers, ctx)
				}
			}
		}

		// Recursively walk child nodes
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}

	walk(doc)
	return markers
}

// classifyMarkerContext examines parent nodes to determine the marker's context
// Note: Only {:toc} markers are passed here; {::toc} is not valid kramdown syntax
func classifyMarkerContext(textNode *html.Node, text string) *MarkerContext {
	// Walk up the tree to find significant parent elements
	var parentLI *html.Node
	var parentUL *html.Node
	var parentOL *html.Node

	for p := textNode.Parent; p != nil; p = p.Parent {
		if p.Type != html.ElementNode {
			continue
		}

		switch p.Data {
		case "pre", "code":
			// Don't process markers in code blocks
			return &MarkerContext{
				Type:       MarkerInCodeBlock,
				Node:       textNode,
				ParentList: nil,
				MarkerText: text,
				IsBlock:    false, // Only {:toc} is processed, which is inline
			}
		case "li":
			if parentLI == nil {
				parentLI = p
			}
		case "ul":
			if parentUL == nil {
				parentUL = p
			}
		case "ol":
			if parentOL == nil {
				parentOL = p
			}
		}
	}

	// Determine marker type based on context
	// Jekyll's TOC replacement rules (verified against Jekyll 4.4.1):
	// 1. {:toc} in <ul> where it's the only content -> replace entire <ul> with TOC
	// 2. {:toc} in <ol> -> leave as-is (Jekyll doesn't support this)
	// 3. {:toc} standalone (not in a list) -> replace marker with TOC

	// If we're in any list (UL or OL)
	if (parentUL != nil || parentOL != nil) && parentLI != nil {
		// {:toc} in unordered list - Jekyll's primary TOC pattern
		// Only replace the entire list if {:toc} is the only content
		if parentUL != nil && isOnlyContentInListItem(textNode, parentLI) {
			return &MarkerContext{
				Type:       MarkerInUnorderedList,
				Node:       textNode,
				ParentList: parentUL,
				MarkerText: text,
				IsBlock:    false,
			}
		}

		// {:toc} in ordered list OR {:toc} in UL but not the only content
		// Jekyll doesn't process these - leave as-is
		return &MarkerContext{
			Type:       MarkerInCodeBlock, // Use CodeBlock type to mean "don't process"
			Node:       textNode,
			ParentList: nil,
			MarkerText: text,
			IsBlock:    false,
		}
	}

	// Default: standalone marker (not in any list)
	// Jekyll supports standalone {:toc} markers when they appear as IAL on their own line,
	// but Blackfriday doesn't preserve this distinction - it renders {:toc} as plain text
	// in all contexts. We can't reliably distinguish between:
	// - {:toc} in heading text like "## Test with {:toc}" (should be literal)
	// - {:toc} as a standalone IAL marker (should be processed in Jekyll)
	// Therefore, we DON'T process standalone markers to avoid false positives.
	// This matches Jekyll's most common pattern (:{toc} in unordered lists).
	return &MarkerContext{
		Type:       MarkerInCodeBlock, // Use CodeBlock type to mean "don't process"
		Node:       textNode,
		ParentList: nil,
		MarkerText: text,
		IsBlock:    false,
	}
}

// isOnlyContentInListItem checks if the marker is the only significant content in the <li>
func isOnlyContentInListItem(textNode *html.Node, li *html.Node) bool {
	// Walk all children of the <li> and check if there's only whitespace + the marker
	for c := li.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			// Check if this is our marker node or just whitespace
			if c == textNode {
				continue
			}
			// If there's non-whitespace text, it's not the only content
			if strings.TrimSpace(c.Data) != "" {
				return false
			}
		} else if c.Type == html.ElementNode {
			// If there are other elements, need to check if they're empty or contain the marker
			if !isEmptyOrContainsNode(c, textNode) {
				return false
			}
		}
	}
	return true
}

// isEmptyOrContainsNode checks if an element is empty or contains the given node
func isEmptyOrContainsNode(elem *html.Node, target *html.Node) bool {
	if elem == target {
		return true
	}

	for c := elem.FirstChild; c != nil; c = c.NextSibling {
		if c == target {
			return true
		}
		if c.Type == html.ElementNode {
			if isEmptyOrContainsNode(c, target) {
				return true
			}
		} else if c.Type == html.TextNode {
			if strings.TrimSpace(c.Data) != "" && c != target {
				return false
			}
		}
	}

	return true
}

// replaceTOCMarkerInDOM replaces a TOC marker in the DOM with the generated TOC HTML
func replaceTOCMarkerInDOM(ctx *MarkerContext, tocHTML string) error {
	switch ctx.Type {
	case MarkerInCodeBlock:
		// Don't replace markers in code blocks - they should be displayed literally
		return nil

	case MarkerInUnorderedList:
		// Replace the entire <ul> parent with the TOC
		// This is the primary Jekyll-compatible TOC replacement pattern
		return replaceUnorderedListWithTOC(ctx, tocHTML)

	case MarkerInOrderedList, MarkerStandalone:
		// Replace just the marker text with the TOC
		// This handles:
		// - {:toc} and {::toc} outside of any list (MarkerStandalone)
		// - {:toc} in <ol> (MarkerInOrderedList)
		// - {::toc} in any list (classified as MarkerInOrderedList for simplicity)
		return replaceStandaloneMarkerWithTOC(ctx, tocHTML)

	default:
		return fmt.Errorf("unknown marker type: %d", ctx.Type)
	}
}

// replaceUnorderedListWithTOC replaces the entire <ul> element with the TOC HTML
func replaceUnorderedListWithTOC(ctx *MarkerContext, tocHTML string) error {
	if ctx.ParentList == nil {
		return fmt.Errorf("no parent list for unordered list marker")
	}

	// Parse the TOC HTML into nodes
	tocNodes, err := parseHTMLFragment(tocHTML)
	if err != nil {
		return err
	}

	// Replace the <ul> parent with the TOC nodes
	parent := ctx.ParentList.Parent
	if parent == nil {
		return fmt.Errorf("parent list has no parent")
	}

	// Insert TOC nodes before the <ul>
	for _, tocNode := range tocNodes {
		parent.InsertBefore(tocNode, ctx.ParentList)
	}

	// Remove the original <ul>
	parent.RemoveChild(ctx.ParentList)

	return nil
}

// replaceStandaloneMarkerWithTOC replaces just the marker text with the TOC HTML
func replaceStandaloneMarkerWithTOC(ctx *MarkerContext, tocHTML string) error {
	if ctx.Node == nil {
		return fmt.Errorf("no node for standalone marker")
	}

	// Parse the TOC HTML into nodes
	tocNodes, err := parseHTMLFragment(tocHTML)
	if err != nil {
		return err
	}

	// Remove the marker text from the text node
	// If the marker is the only content, replace the entire text node
	// Otherwise, split the text node and insert TOC in between
	markerText := ctx.MarkerText
	nodeText := ctx.Node.Data

	if strings.TrimSpace(nodeText) == strings.TrimSpace(markerText) {
		// Marker is the only content - replace the text node with TOC nodes
		parent := ctx.Node.Parent
		if parent == nil {
			return fmt.Errorf("text node has no parent")
		}

		// Insert TOC nodes before the text node
		for _, tocNode := range tocNodes {
			parent.InsertBefore(tocNode, ctx.Node)
		}

		// Remove the original text node
		parent.RemoveChild(ctx.Node)
	} else {
		// Marker is part of larger text - split the text node
		parent := ctx.Node.Parent
		if parent == nil {
			return fmt.Errorf("text node has no parent")
		}

		// Find the marker position
		markerIdx := strings.Index(nodeText, markerText)
		if markerIdx == -1 {
			// Try with regex to handle whitespace variations
			// Only {:toc} is processed (no {::toc} support)
			indices := tocPatternInline.FindStringIndex(nodeText)
			if indices != nil {
				markerIdx = indices[0]
			}
		}

		if markerIdx >= 0 {
			// Split text: before | marker | after
			before := nodeText[:markerIdx]
			after := nodeText[markerIdx+len(markerText):]

			// Replace current node with before text
			if before != "" {
				ctx.Node.Data = before
			} else {
				parent.RemoveChild(ctx.Node)
			}

			// Insert TOC nodes
			insertPoint := ctx.Node.NextSibling
			for _, tocNode := range tocNodes {
				parent.InsertBefore(tocNode, insertPoint)
			}

			// Insert after text if any
			if after != "" {
				afterNode := &html.Node{
					Type: html.TextNode,
					Data: after,
				}
				parent.InsertBefore(afterNode, insertPoint)
			}
		}
	}

	return nil
}
