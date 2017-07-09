package plugins

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/osteele/gojekyll/pages"
)

type jekyllFeedPlugin struct{ plugin }

var redirectTemplate *template.Template

func init() {
	register("jekyll-redirect-from", func(ctx PluginContext, h pluginHelper) error {
		return nil
	})
	tmpl, err := template.New("redirect_from").Parse(redirectFromText)
	if err != nil {
		panic(err)
	}
	redirectTemplate = tmpl
}

type redirector struct {
	From string
	To   string
}

func (p *redirector) Permalink() string    { return p.From }
func (p *redirector) SourcePath() string   { return "" } // FIXME bad design
func (p *redirector) OutputExt() string    { return ".html" }
func (p *redirector) Published() bool      { return true }
func (p *redirector) Static() bool         { return false } // FIXME means different things to different callers
func (p *redirector) Categories() []string { return []string{} }
func (p *redirector) Tags() []string       { return []string{} }

func (p *redirector) Content() []byte {
	buf := new(bytes.Buffer)
	if err := redirectTemplate.Execute(buf, p); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func (p *redirector) Write(w io.Writer, c pages.RenderingContext) error {
	return redirectTemplate.Execute(w, p)
}

func (p jekyllFeedPlugin) PostRead(site Site) error {
	redirections := []pages.Document{}
	for _, p := range site.Pages() {
		rd, ok := p.FrontMatter()["redirect_from"]
		if ok {
			switch rd := rd.(type) {
			case string:
				siteurl := site.Config().AbsoluteURL
				baseurl := site.Config().BaseURL
				var p = redirector{From: rd, To: strings.Join([]string{siteurl, baseurl, p.Permalink()}, "")}
				redirections = append(redirections, &p)
			default:
				fmt.Printf("unimplemented redirect_from type: %T\n", rd)
			}
		}
		rd, ok = p.FrontMatter()["redirect_to"]
		if ok {
			switch rd := rd.(type) {
			case string:
				r := redirector{From: rd, To: p.Permalink()}
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

// Adapted from https://github.com/jekyll/jekyll-redirect-from
var redirectFromText = `<!DOCTYPE html>
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
