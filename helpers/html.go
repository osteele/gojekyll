package helpers

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
				s := (string(z.Text()))
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
