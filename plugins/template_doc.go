package plugins

import (
	"io"
	"path"

	"github.com/osteele/liquid"
)

type templateDoc struct {
	site Site
	path string
	tpl  *liquid.Template
}

func (d *templateDoc) Permalink() string  { return "/" + d.path }
func (d *templateDoc) SourcePath() string { return "" }
func (d *templateDoc) OutputExt() string  { return path.Ext(d.path) }
func (d *templateDoc) Published() bool    { return true }
func (d *templateDoc) Static() bool       { return false } // FIXME means different things to different callers

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
