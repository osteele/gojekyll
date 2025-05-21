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
				continue // Skip the buf.Write(z.Raw()) below since we've already written
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
	_, err := w.Write(buf.Bytes())
	return err
}
