package plugins

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/osteele/gojekyll/pages"
)

type jekyllRedirectFromPlugin struct{ plugin }

var redirectTemplate *template.Template

func init() {
	register("jekyll-redirect-from", jekyllRedirectFromPlugin{})
	tmpl, err := template.New("redirect_from").Parse(redirectFromTemplateSource)
	if err != nil {
		panic(err)
	}
	redirectTemplate = tmpl
}

func (p jekyllRedirectFromPlugin) PostRead(site Site) error {
	redirections := []pages.Document{}
	for _, p := range site.Pages() {
		rd, ok := p.FrontMatter()["redirect_from"]
		if ok {
			switch rd := rd.(type) {
			case string:
				siteurl := site.Config().AbsoluteURL
				baseurl := site.Config().BaseURL
				var p = redirectionDoc{From: rd, To: strings.Join([]string{siteurl, baseurl, p.Permalink()}, "")}
				redirections = append(redirections, &p)
			default:
				fmt.Printf("unimplemented redirect_from type: %T\n", rd)
			}
		}
		rd, ok = p.FrontMatter()["redirect_to"]
		if ok {
			switch rd := rd.(type) {
			case string:
				r := redirectionDoc{From: rd, To: p.Permalink()}
				p.SetContent(r.Content())
			default:
				fmt.Printf("unimplemented redirect_from type: %T\n", rd)
			}
		}
	}
	for _, p := range redirections {
		site.AddDocument(p, true)
	}
	return nil
}

type redirectionDoc struct {
	From string
	To   string
}

func (d *redirectionDoc) Permalink() string    { return d.From }
func (d *redirectionDoc) SourcePath() string   { return "" } // FIXME bad design
func (d *redirectionDoc) OutputExt() string    { return ".html" }
func (d *redirectionDoc) Published() bool      { return true }
func (d *redirectionDoc) Static() bool         { return false } // FIXME means different things to different callers
func (d *redirectionDoc) Categories() []string { return []string{} }
func (d *redirectionDoc) Tags() []string       { return []string{} }

func (d *redirectionDoc) Content() []byte {
	buf := new(bytes.Buffer)
	if err := redirectTemplate.Execute(buf, d); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func (d *redirectionDoc) Write(w io.Writer) error {
	return redirectTemplate.Execute(w, d)
}

// Adapted from https://github.com/jekyll/jekyll-redirect-from
const redirectFromTemplateSource = `<!DOCTYPE html>
<html lang="en-US">
  <meta charset="utf-8">
  <title>Redirecting…</title>
  <link rel="canonical" href="{{ .To }}">
  <meta http-equiv="refresh" content="0; url={{ .To }}">
  <meta name="robots" content="noindex">
  <h1>Redirecting…</h1>
  <a href="{{ .To }}">Click here if you are not redirected.</a>
  <script>location="{{ .To }}"</script>
</html>`
