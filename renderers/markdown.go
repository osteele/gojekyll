package renderers

import (
	"bytes"
	"io"
	"regexp"

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
			err := z.Err()
			if err == io.EOF {
				return utils.WrapError(err,
					"unexpected EOF while processing markdown=\"1\" attribute. "+
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
	html, err := _renderMarkdown(buf.Bytes())
	if err != nil {
		return err
	}
	_, err = w.Write(html)
	return err
}
