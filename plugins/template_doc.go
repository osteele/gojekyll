package plugins

import (
	"io"

	"github.com/osteele/gojekyll/pages"
	"github.com/osteele/liquid"
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

func (d *templateDoc) Content() []byte {
	bindings := map[string]interface{}{"site": d.site}
	b, err := d.tpl.Render(bindings)
	if err != nil {
		panic(err)
	}
	return b
}

func (d *templateDoc) Write(w io.Writer) error {
	_, err := w.Write(d.Content())
	return err
}
