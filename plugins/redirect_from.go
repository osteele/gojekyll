package plugins

import (
	"bytes"
	"text/template"
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

func (p jekyllRedirectFromPlugin) PostReadSite(s Site) error {
	ps := s.Pages()
	addRedirects, err := p.processRedirectFrom(s, ps)
	if err != nil {
		return err
	}
	if err := p.processRedirectTo(s, ps); err != nil {
		return err
	}
	addRedirects()
	return nil
}

func (p jekyllRedirectFromPlugin) processRedirectFrom(s Site, ps []Page) (func(), error) {
	var (
		cfg       = s.Config()
		siteurl   = cfg.AbsoluteURL
		baseurl   = cfg.BaseURL
		prefix    = siteurl + baseurl
		redirects = []func(){}
	)
	addRedirectFrom := func(from string, to Page) {
		f := func() {
			s.AddHTMLPage(from, createRedirectionHTML(prefix+to.URL()), nil)
		}
		redirects = append(redirects, f)
	}
	for _, p := range ps {
		for _, from := range p.FrontMatter().StringArray("redirect_from") {
			addRedirectFrom(from, p)
		}
	}
	return func() {
		for _, f := range redirects {
			f()
		}
	}, nil
}

func (p jekyllRedirectFromPlugin) processRedirectTo(_ Site, ps []Page) error {
	for _, p := range ps {
		sources := p.FrontMatter().StringArray("redirect_to")
		if len(sources) > 0 {
			p.SetContent(createRedirectionHTML(sources[0]))
		}
	}
	return nil
}

func createRedirectionHTML(to string) string {
	r := redirection{to}
	buf := new(bytes.Buffer)
	if err := redirectTemplate.Execute(buf, r); err != nil {
		panic(err)
	}
	return buf.String()
}

type redirection struct {
	To string
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
