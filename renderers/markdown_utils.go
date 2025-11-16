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
