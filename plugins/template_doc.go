package plugins

import (
	"bytes"
	"io"

	"github.com/osteele/gojekyll/pages"
	"github.com/osteele/liquid"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/html"
)

func newTemplateDoc(s Site, path, src string) pages.Document {
	tpl, err := s.TemplateEngine().ParseTemplate([]byte(src))
	if err != nil {
		panic(err)
	}
	return &templateDoc{pages.PageEmbed{Path: path}, s, tpl}
}

type templateDoc struct {
	pages.PageEmbed
	site Site
	tpl  *liquid.Template
}

func (d *templateDoc) Content() string {
	bindings := map[string]interface{}{"site": d.site}
	b, err := d.tpl.Render(bindings)
	if err != nil {
		panic(err)
	}
	m := minify.New()
	m.AddFunc("text/html", html.Minify)
	min := bytes.NewBuffer(make([]byte, 0, len(b)))
	if err := m.Minify("text/html", min, bytes.NewBuffer(b)); err != nil {
		panic(err)
	}
	return min.String()
}

func (d *templateDoc) Write(w io.Writer) error {
	_, err := io.WriteString(w, d.Content())
	return err
}
