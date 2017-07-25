package pages

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"sync"
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

	Categories() []string
	Tags() []string
}

type page struct {
	file
	sync.Mutex
	firstLine int
	raw       []byte
	content   *[]byte
}

// Static is in the File interface.
func (p *page) Static() bool { return false }

func makePage(filename string, f file) (*page, error) {
	raw, lineNo, err := readFrontMatter(&f)
	if err != nil {
		return nil, err
	}
	p := page{
		file:      f,
		firstLine: lineNo,
		raw:       raw,
	}
	if err := p.setPermalink(); err != nil {
		return nil, err
	}
	return &p, nil
}

func (p *page) Reload() error {
	if err := p.file.Reload(); err != nil {
		return err
	}
	// FIXME use original defaults
	raw, lineNo, err := readFrontMatter(&p.file)
	if err != nil {
		return err
	}
	p.firstLine = lineNo
	p.raw = raw
	return nil
}

func readFrontMatter(f *file) (b []byte, lineNo int, err error) {
	b, err = ioutil.ReadFile(f.filename)
	if err != nil {
		return
	}
	lineNo = 1
	frontMatter, err := frontmatter.Read(&b, &lineNo)
	if err != nil {
		return
	}
	f.frontMatter = templates.MergeVariableMaps(f.frontMatter, frontMatter)
	return
}

func (p *page) FrontMatter() map[string]interface{} {
	return p.frontMatter
}

// Categories is in the Page interface
func (p *page) Categories() []string {
	return frontmatter.FrontMatter(p.frontMatter).SortedStringArray("categories")
}

// Tags is in the Page interface
func (p *page) Tags() []string {
	return frontmatter.FrontMatter(p.frontMatter).SortedStringArray("tags")
}

// TemplateContext returns the local variables for template evaluation
func (p *page) TemplateContext() map[string]interface{} {
	return map[string]interface{}{
		"page": p,
		"site": p.site,
	}
}

// PostDate is part of the Page interface.
// FIXME move this back to Page interface, or re-work this entirely.
func (f *file) PostDate() time.Time {
	switch value := f.frontMatter["date"].(type) {
	case time.Time:
		return value
	case string:
		t, err := evaluator.ParseDate(value)
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
	content := p.maybeContent(false)
	if content != nil {
		return content, nil
	}
	pipe := p.site.RenderingPipeline()
	buf := new(bytes.Buffer)
	b, err := pipe.Render(buf, p.raw, p.filename, p.firstLine, p.TemplateContext())
	if err != nil {
		return nil, err
	}
	p.SetContent(b)
	return b, nil
}

// retains its argument
func (p *page) SetContent(content []byte) {
	p.Lock()
	defer p.Unlock()
	p.content = &content
}
