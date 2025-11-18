package utils

import (
	"bytes"
	"io"

	"golang.org/x/net/html"
)

// ApplyToHTMLText applies a filter only to the text within an HTML document.
func ApplyToHTMLText(doc []byte, fn func(string) string) []byte {
	z := html.NewTokenizer(bytes.NewReader(doc))
	buf := new(bytes.Buffer)
	body := false
outer:
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			if z.Err() == io.EOF {
				break outer
			}
			panic(z.Err())
		case html.StartTagToken, html.EndTagToken:
			tn, _ := z.TagName()
			if string(tn) == "body" {
				body = tt == html.StartTagToken
			}
		case html.TextToken:
			if body {
				s := string(z.Text())
				_, err := buf.WriteString(fn(s))
				if err != nil {
					panic(err)
				}
				continue outer
			}
		}
		_, err := buf.Write(z.Raw())
		if err != nil {
			panic(err)
		}
	}
	return buf.Bytes()
}

// ProcessAnchorHrefs applies a filter to href attributes of anchor tags within an HTML document.
func ProcessAnchorHrefs(doc []byte, fn func(string) string) []byte {
	z := html.NewTokenizer(bytes.NewReader(doc))
	buf := new(bytes.Buffer)
outer:
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			if z.Err() == io.EOF {
				break outer
			}
			panic(z.Err())
		case html.StartTagToken, html.SelfClosingTagToken:
			// Get the full token first
			tag := z.Token()
			if tag.Data == "a" {
				modified := false
				// Process href attributes
				for i, attr := range tag.Attr {
					if attr.Key == "href" {
						newHref := fn(attr.Val)
						if newHref != attr.Val {
							tag.Attr[i].Val = newHref
							modified = true
						}
					}
				}
				if modified {
					// Write the modified tag
					_, err := buf.WriteString(tag.String())
					if err != nil {
						panic(err)
					}
					continue outer
				}
			}
			// Write unmodified non-anchor tags or anchor tags without href changes
			_, err := buf.WriteString(tag.String())
			if err != nil {
				panic(err)
			}
			continue outer
		}
		// Write the original token
		_, err := buf.Write(z.Raw())
		if err != nil {
			panic(err)
		}
	}
	return buf.Bytes()
}
