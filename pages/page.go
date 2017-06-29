package pages

import (
	"bytes"
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/osteele/gojekyll/helpers"
	"github.com/osteele/gojekyll/templates"
)

// Page is a post or collection page.
type Page struct {
	file
	raw       []byte
	processed *[]byte
}

// Static is in the File interface.
func (p *Page) Static() bool { return false }

func newPage(filename string, f file) (*Page, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	frontMatter, err := templates.ReadFrontMatter(&b)
	if err != nil {
		return nil, err
	}
	f.frontMatter = templates.MergeVariableMaps(f.frontMatter, frontMatter)
	return &Page{
		file: f,
		raw:  b,
	}, nil
}

// PageVariables returns the attributes of the template page object.
func (p *Page) PageVariables() templates.VariableMap {
	var (
		relpath = p.relpath
		ext     = filepath.Ext(relpath)
		root    = helpers.TrimExt(p.relpath)
		base    = filepath.Base(root)
	)

	data := templates.VariableMap{
		"path": relpath,
		"url":  p.Permalink(),
		// TODO output

		// not documented, but present in both collection and non-collection pages
		"permalink": p.Permalink(),

		// TODO only in non-collection pages:
		// TODO dir
		// TODO name
		// TODO next previous

		// TODO Documented as present in all pages, but de facto only defined for collection pages
		"id":    base,
		"title": base, // TODO capitalize
		// TODO date (of the collection?) 2017-06-15 07:44:21 -0400
		// TODO excerpt category? categories tags
		// TODO slug
		"categories": []string{},
		"tags":       []string{},

		// TODO Only present in collection pages https://jekyllrb.com/docs/collections/#documents
		"relative_path": p.Path(),
		// TODO collection(name)

		// TODO undocumented; only present in collection pages:
		"ext": ext,
	}
	for k, v := range p.frontMatter {
		switch k {
		// doc implies these aren't present, but they appear to be present in a collection page:
		// case "layout", "published":
		case "permalink":
		// omit this, in order to use the value above
		default:
			data[k] = v
		}
	}
	return data
}

// TemplateContext returns the local variables for template evaluation
func (p *Page) TemplateContext(rc RenderingContext) templates.VariableMap {
	return templates.VariableMap{
		"page": p.PageVariables(),
		"site": rc.SiteVariables(),
	}
}

// Write applies Liquid and Markdown, as appropriate.
func (p *Page) Write(rc RenderingContext, w io.Writer) error {
	rp := rc.RenderingPipeline()
	if p.processed != nil {
		_, err := w.Write(*p.processed)
		return err
	}
	b, err := rp.Render(w, p.raw, p.filename, p.TemplateContext(rc))
	if err != nil {
		return err
	}
	layout := p.frontMatter.String("layout", "")
	if layout != "" {
		b, err = rp.ApplyLayout(layout, b, p.TemplateContext(rc))
		if err != nil {
			return err
		}
	}
	_, err = w.Write(b)
	return err
}

// ComputeContent computes the page content.
func (p *Page) ComputeContent(rc RenderingContext) ([]byte, error) {
	if p.processed == nil {
		w := new(bytes.Buffer)
		if err := p.Write(rc, w); err != nil {
			return nil, err
		}
		b := w.Bytes()
		p.processed = &b
	}
	return *p.processed, nil
}
