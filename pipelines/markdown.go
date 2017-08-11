package pipelines

import (
	"bytes"
	"io"
	"regexp"

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

func renderMarkdown(md []byte) []byte {
	renderer := blackfriday.HtmlRenderer(blackfridayFlags, "", "")
	html := blackfriday.MarkdownOptions(md, renderer, blackfriday.Options{
		Extensions: blackfridayExtensions})
	html, err := renderInnerMarkdown(html)
	if err != nil {
		panic(err)
	}
	return html
}

var markdownAttrRE = regexp.MustCompile(`\s*markdown\s*=\s*("1"|'1'|1)\s*`)

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
			if hasAttr(z) {
				tag := markdownAttrRE.ReplaceAll(z.Raw(), []byte(" "))
				tag = bytes.Replace(tag, []byte(" >"), []byte(">"), 1)
				_, err := buf.Write(tag)
				if err != nil {
					return nil, err
				}
				if err := processInnerMarkdown(buf, z); err != nil {
					return nil, err
				}
				// the above leaves z set to the end token
			}
		}
		_, err := buf.Write(z.Raw())
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func hasAttr(z *html.Tokenizer) bool {
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
			depth++
		case html.EndTagToken:
			depth--
			if depth == 0 {
				break loop
			}
		}
		_, err := buf.Write(z.Raw())
		if err != nil {
			panic(err)
		}
	}
	html := renderMarkdown(buf.Bytes())
	_, err := w.Write(html)
	return err
}
