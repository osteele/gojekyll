package pages

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/osteele/gojekyll/frontmatter"
	"github.com/osteele/gojekyll/templates"
	"github.com/osteele/liquid/evaluator"
)

// Page is a document with frontmatter.
type Page interface {
	Document
	// Content asks a page to compute its content.
	// This has the side effect of causing the content to subsequently appear in the drop.
	Content() ([]byte, error)
	SetContent(content []byte)
	FrontMatter() map[string]interface{}
	// PostDate returns the date computed from the filename or frontmatter.
	// It is an uncaught error to call this on a page that is not a Post.
	// TODO Should posts have their own interface?
	PostDate() time.Time
}

type page struct {
	file
	firstLine int
	raw       []byte
	content   *[]byte
}

// Static is in the File interface.
func (p *page) Static() bool { return false }

func makePage(filename string, f file) (*page, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	lineNo := 1
	frontMatter, err := frontmatter.Read(&b, &lineNo)
	if err != nil {
		return nil, err
	}
	f.frontMatter = templates.MergeVariableMaps(f.frontMatter, frontMatter)
	p := page{
		file:      f,
		firstLine: lineNo,
		raw:       b,
	}
	if err = p.setPermalink(); err != nil {
		return nil, err
	}
	return &p, nil
}

// TemplateContext returns the local variables for template evaluation
func (p *page) TemplateContext() map[string]interface{} {
	return map[string]interface{}{
		"page": p,
		"site": p.site,
	}
}

func (p *page) FrontMatter() map[string]interface{} {
	return p.frontMatter
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
func (p *page) Write(w io.Writer) error {
	content, err := p.Content()
	if err != nil {
		return err
	}
	layout, ok := p.frontMatter["layout"].(string)
	if ok && layout != "" {
		rp := p.site.RenderingPipeline()
		content, err = rp.ApplyLayout(layout, content, p.TemplateContext())
		if err != nil {
			return err
		}
	}
	_, err = w.Write(content)
	return err
}

// Content computes the page content.
func (p *page) Content() ([]byte, error) {
	if p.content != nil {
		return *p.content, nil
	}
	pipe := p.site.RenderingPipeline()
	buf := new(bytes.Buffer)
	b, err := pipe.Render(buf, p.raw, p.filename, p.firstLine, p.TemplateContext())
	if err != nil {
		return nil, err
	}
	p.content = &b
	return b, nil
}

// retains content
func (p *page) SetContent(content []byte) {
	p.content = &content
}
