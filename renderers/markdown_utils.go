package renderers

import (
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

// parseHTMLFragment parses an HTML fragment string into DOM nodes
func parseHTMLFragment(htmlStr string) ([]*html.Node, error) {
	// Wrap in a container to parse as a fragment
	wrapped := "<div>" + htmlStr + "</div>"
	doc, err := html.Parse(strings.NewReader(wrapped))
	if err != nil {
		return nil, err
	}

	// Find the <div> container (html > head > div or html > body > div)
	var container *html.Node
	var findDiv func(*html.Node)
	findDiv = func(n *html.Node) {
		if container != nil {
			return
		}
		if n.Type == html.ElementNode && n.Data == "div" {
			container = n
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findDiv(c)
		}
	}
	findDiv(doc)

	if container == nil {
		return nil, fmt.Errorf("could not find container div")
	}

	// Extract all children of the container
	var nodes []*html.Node
	for c := container.FirstChild; c != nil; c = c.NextSibling {
		// Clone the node to detach it from the parsed tree
		nodes = append(nodes, cloneNode(c))
	}

	return nodes, nil
}

// cloneNode creates a deep copy of a node and its children
func cloneNode(n *html.Node) *html.Node {
	clone := &html.Node{
		Type:      n.Type,
		DataAtom:  n.DataAtom,
		Data:      n.Data,
		Namespace: n.Namespace,
		Attr:      make([]html.Attribute, len(n.Attr)),
	}
	copy(clone.Attr, n.Attr)

	// Clone children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		childClone := cloneNode(c)
		clone.AppendChild(childClone)
	}

	return clone
}

// extractBodyContent extracts content from within <body> tags
// html.Render wraps content in <html><head></head><body>...</body></html>
func extractBodyContent(htmlBytes []byte) []byte {
	// Find <body> and </body> tags
	bodyStart := bytes.Index(htmlBytes, []byte("<body>"))
	bodyEnd := bytes.Index(htmlBytes, []byte("</body>"))

	if bodyStart == -1 || bodyEnd == -1 {
		// Fallback: return original if body tags not found
		return htmlBytes
	}

	// Extract content between <body> and </body>
	bodyStart += len("<body>")
	return htmlBytes[bodyStart:bodyEnd]
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

// hasNoTocSibling checks if a heading has a sibling paragraph containing {:.no_toc}
// Kramdown syntax allows IAL markers on the line after a heading, which Blackfriday
// renders as a sibling <p> element
// If found, the sibling paragraph is removed from the DOM (matching Kramdown behavior)
func hasNoTocSibling(heading *html.Node) bool {
	// Walk forward through siblings, skipping whitespace text nodes
	for sibling := heading.NextSibling; sibling != nil; sibling = sibling.NextSibling {
		// Skip whitespace-only text nodes
		if sibling.Type == html.TextNode {
			if strings.TrimSpace(sibling.Data) == "" {
				continue
			}
			// Non-whitespace text means no more siblings to check
			break
		}

		// Check if this is a <p> element
		if sibling.Type == html.ElementNode && sibling.Data == "p" {
			// Check if the paragraph contains {:.no_toc}
			pText := extractTextContent(sibling)
			if strings.TrimSpace(pText) == "{:.no_toc}" {
				// Remove the IAL marker paragraph (matching Kramdown behavior)
				if sibling.Parent != nil {
					sibling.Parent.RemoveChild(sibling)
				}
				return true
			}
			// Found a non-marker paragraph, stop checking
			break
		}

		// Found a non-paragraph element, stop checking
		if sibling.Type == html.ElementNode {
			break
		}
	}

	return false
}
