package pages

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/osteele/gojekyll/templates"
	"github.com/osteele/liquid/evaluator"
)

// Page is a document with frontmatter.
type Page interface {
	Document
	// Content asks a page to compute its content.
	// This has the side effect of causing the content to subsequently appear in the drop.
	Content(rc RenderingContext) ([]byte, error)
	// PostDate returns the date computed from the filename or frontmatter.
	// It is an uncaught error to call this on a page that is not a Post.
	// TODO Should posts have their own interface?
	PostDate() time.Time
}

type page struct {
	file
	raw     []byte
	content *[]byte
}

// Static is in the File interface.
func (p *page) Static() bool { return false }

func newPage(filename string, f file) (*page, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	frontMatter, err := templates.ReadFrontMatter(&b)
	if err != nil {
		return nil, err
	}
	f.frontMatter = templates.MergeVariableMaps(f.frontMatter, frontMatter)
	return &page{
		file: f,
		raw:  b,
	}, nil
}

// TemplateContext returns the local variables for template evaluation
func (p *page) TemplateContext(rc RenderingContext) map[string]interface{} {
	return map[string]interface{}{
		"page": p,
		"site": rc.Site(),
	}
}

// PostDate is part of the Page interface.
func (p *page) PostDate() time.Time {
	switch value := p.frontMatter["date"].(type) {
	case time.Time:
		return value
	case string:
		t, err := evaluator.ParseTime(value)
		if err == nil {
			return t
		}
	default:
		panic(fmt.Sprintf("expected a date %v", value))
	}
	panic("read posts should have set this")
}

// Write applies Liquid and Markdown, as appropriate.
func (p *page) Write(w io.Writer, rc RenderingContext) error {
	content, err := p.Content(rc)
	if err != nil {
		return err
	}
	layout, ok := p.frontMatter["layout"].(string)
	if ok && layout != "" {
		rp := rc.RenderingPipeline()
		content, err = rp.ApplyLayout(layout, content, p.TemplateContext(rc))
		if err != nil {
			return err
		}
	}
	_, err = w.Write(content)
	return err
}

// Content computes the page content.
func (p *page) Content(rc RenderingContext) ([]byte, error) {
	if p.content == nil {
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
