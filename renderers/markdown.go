package renderers

import (
	"bytes"
	"io"
	"regexp"

	"github.com/osteele/gojekyll/utils"
	"github.com/russross/blackfriday"
	"golang.org/x/net/html"
)

const blackfridayFlags = 0 |
	blackfriday.HTML_USE_XHTML |
	blackfriday.HTML_USE_SMARTYPANTS |
	blackfriday.HTML_SMARTYPANTS_FRACTIONS |
	blackfriday.HTML_SMARTYPANTS_DASHES |
	blackfriday.HTML_SMARTYPANTS_LATEX_DASHES

const blackfridayExtensions = 0 |
	blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
	blackfriday.EXTENSION_TABLES |
	blackfriday.EXTENSION_FENCED_CODE |
	blackfriday.EXTENSION_AUTOLINK |
	blackfriday.EXTENSION_STRIKETHROUGH |
	blackfriday.EXTENSION_SPACE_HEADERS |
	blackfriday.EXTENSION_HEADER_IDS |
	blackfriday.EXTENSION_BACKSLASH_LINE_BREAK |
	blackfriday.EXTENSION_DEFINITION_LISTS |
	// added relative to commonExtensions
	blackfriday.EXTENSION_AUTO_HEADER_IDS

func renderMarkdown(md []byte) ([]byte, error) {
	renderer := blackfriday.HtmlRenderer(blackfridayFlags, "", "")
	html := blackfriday.MarkdownOptions(
		md,
		renderer,
		blackfriday.Options{Extensions: blackfridayExtensions},
	)
	html, err := renderInnerMarkdown(html)
	if err != nil {
		return nil, utils.WrapError(err, "markdown")
	}
	return html, nil
}

func _renderMarkdown(md []byte) ([]byte, error) {
	renderer := blackfriday.HtmlRenderer(blackfridayFlags, "", "")
	html := blackfriday.MarkdownOptions(
		md,
		renderer,
		blackfriday.Options{Extensions: blackfridayExtensions},
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

// Used inside markdown=1.
// TODO Instead of this approach, only count tags that match the start
// tag. For example, if <div markdown="1"> kicked off the inner markdown,
// count the div depth.
var notATagRE = regexp.MustCompile(`@|(https?|ftp):`)

// called once markdown="1" attribute is detected.
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
