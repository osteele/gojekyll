package plugins

import (
	"bytes"
	"fmt"
	"io"
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
	ps := site.Pages()
	newPages, err := p.processRedirectFrom(site, ps)
	if err != nil {
		return err
	}
	if err := p.processRedirectTo(site, ps); err != nil {
		return err
	}
	for _, r := range newPages {
		site.AddDocument(r, true)
	}
	return nil
}

func (p jekyllRedirectFromPlugin) processRedirectFrom(site Site, ps []pages.Page) ([]pages.Document, error) {
	var (
		cfg          = site.Config()
		siteurl      = cfg.AbsoluteURL
		baseurl      = cfg.BaseURL
		prefix       = siteurl + baseurl
		redirections = []pages.Document{}
	)
	addRedirectFrom := func(from string, to pages.Page) {
		r := redirectionDoc{pages.PageEmbed{Path: from}, prefix + to.URL()}
		redirections = append(redirections, &r)
	}
	for _, p := range ps {
		sources, err := getStringArray(p, "redirect_from")
		if err != nil {
			return nil, err
		}
		for _, from := range sources {
			addRedirectFrom(from, p)
		}
	}
	return redirections, nil
}

func (p jekyllRedirectFromPlugin) processRedirectTo(site Site, ps []pages.Page) error {
	for _, p := range ps {
		sources, err := getStringArray(p, "redirect_to")
		if err != nil {
			return err
		}
		if len(sources) > 0 {
			r := redirectionDoc{pages.PageEmbed{Path: p.URL()}, sources[0]}
			p.SetContent(r.Content())
		}
	}
	return nil
}

func getStringArray(p pages.Page, fieldName string) (out []string, err error) {
	if value, ok := p.FrontMatter()[fieldName]; ok {
		switch value := value.(type) {
		case []string:
			out = value
		case []interface{}:
			out = make([]string, len(value))
			for i, item := range value {
				out[i] = fmt.Sprintf("%s", item)
			}
		case string:
			out = []string{value}
		default:
			err = fmt.Errorf("unimplemented redirect_from type %T", value)
		}
	}
	return
}

type redirectionDoc struct {
	pages.PageEmbed
	To string
}

func (d *redirectionDoc) Content() string {
	buf := new(bytes.Buffer)
	if err := redirectTemplate.Execute(buf, d); err != nil {
		panic(err)
	}
	return buf.String()
}

func (d *redirectionDoc) Write(w io.Writer) error {
	return redirectTemplate.Execute(w, d)
}

// Adapted from https://github.com/jekyll/jekyll-redirect-from
const redirectFromTemplateSource = `<!DOCTYPE html>
<html lang="en-US">
  <meta charset="utf-8">
  <title>Redirecting…</title>
  <link rel="canonical" href="{{.To}}">
  <meta http-equiv="refresh" content="0; url={{.To}}">
  <meta name="robots" content="noindex">
  <h1>Redirecting…</h1>
  <a href="{{.To}}">Click here if you are not redirected.</a>
  <script>location="{{.To}}"</script>
</html>`
