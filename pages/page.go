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
	raw     []byte
	content *[]byte
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
func (p *Page) PageVariables() map[string]interface{} {
	var (
		relpath = p.relpath
		ext     = filepath.Ext(relpath)
		root    = helpers.TrimExt(p.relpath)
		base    = filepath.Base(root)
	)

	data := map[string]interface{}{
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
func (p *Page) TemplateContext(rc RenderingContext) map[string]interface{} {
	return map[string]interface{}{
		"page": p.PageVariables(),
		"site": rc.SiteVariables(),
	}
}

// Write applies Liquid and Markdown, as appropriate.
func (p *Page) Write(rc RenderingContext, w io.Writer) error {
	rp := rc.RenderingPipeline()
	b, err := rp.Render(w, p.raw, p.filename, p.TemplateContext(rc))
	if err != nil {
		return err
	}
	layout := templates.VariableMap(p.frontMatter).String("layout", "")
	if layout != "" {
		b, err = rp.ApplyLayout(layout, b, p.TemplateContext(rc))
		if err != nil {
			return err
		}
	}
	_, err = w.Write(b)
	return err
}

// Content computes the page content.
func (p *Page) Content(rc RenderingContext) ([]byte, error) {
	if p.content == nil {
		// TODO DRY w/ Page.Write
		rp := rc.RenderingPipeline()
		buf := new(bytes.Buffer)
		b, err := rp.Render(buf, p.raw, p.filename, p.TemplateContext(rc))
		if err != nil {
			return nil, err
		}
		p.content = &b
	}
	return *p.content, nil
}
