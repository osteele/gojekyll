package pages

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/osteele/gojekyll/frontmatter"
	"github.com/osteele/gojekyll/templates"
	"github.com/osteele/gojekyll/version"
	"github.com/osteele/liquid/evaluator"
)

// Page is a document with frontmatter.
type Page interface {
	Document
	// Render asks a page to compute its content.
	// This has the side effect of causing the content to subsequently appear in the drop.
	Render() error
	SetContent(string)
	FrontMatter() map[string]interface{}
	// PostDate returns the date computed from the filename or frontmatter.
	// It is an uncaught error to call this on a page that is not a Post.
	// TODO Should posts have their own interface?
	PostDate() time.Time

	Categories() []string
	Tags() []string
}

// PageEmbed can be embedded to give defaults for the Page interface.
type PageEmbed struct {
	Path string
}

// Permalink is in the pages.Page interface.
func (p *PageEmbed) Permalink() string { return p.Path }

// OutputExt is in the pages.Page interface.
func (p *PageEmbed) OutputExt() string { return path.Ext(p.Path) }

// SourcePath is in the pages.Page interface.
func (p *PageEmbed) SourcePath() string { return "" }

// Published is in the pages.Page interface.
func (p *PageEmbed) Published() bool { return true }

// Static is in the pages.Page interface.
func (p *PageEmbed) Static() bool { return false } // FIXME means different things to different callers

// Reload is in the pages.Page interface.
func (p *PageEmbed) Reload() error { return nil }

// A page is a concrete implementation of the Page interface.
type page struct {
	file
	firstLine int
	raw       []byte

	sync.RWMutex
	content      string
	contentError error
	contentOnce  sync.Once
	excerpt      interface{} // []byte or string, depending on rendering stage
	rendered     bool
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
	p.reset()
	return nil
}

func (p *page) reset() {
	p.contentOnce = sync.Once{}
	p.rendered = false
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
	env := os.Getenv("JEKYLL_ENV")
	if env == "" {
		env = "development"
	}
	return map[string]interface{}{
		"page": p,
		"site": p.site,
		"jekyll": map[string]string{
			"environment": env,
			"version":     fmt.Sprintf("%s (gojekyll)", version.Version)},
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
	err := p.Render()
	if err != nil {
		return err
	}
	p.RLock()
	defer p.RUnlock()
	content := p.content
	layout, ok := p.frontMatter["layout"].(string)
	if ok && layout != "" {
		rp := p.site.RendererManager()
		b, e := rp.ApplyLayout(layout, []byte(content), p.TemplateContext())
		if e != nil {
			return e
		}
		_, err = w.Write(b)
	} else {
		_, err = io.WriteString(w, content)
	}
	return err
}

// Content computes the page content.
func (p *page) Render() error {
	p.contentOnce.Do(func() {
		cn, ex, err := p.computeContent()
		p.Lock()
		defer p.Unlock()
		p.content = cn
		p.contentError = err
		p.excerpt = ex
		p.rendered = true
	})
	return p.contentError
}

func (p *page) SetContent(content string) {
	p.Lock()
	defer p.Unlock()
	p.content = content
	p.contentError = nil
}

func (p *page) computeContent() (cn string, ex string, err error) {
	pl := p.site.RendererManager()
	buf := new(bytes.Buffer)
	err = pl.Render(buf, p.raw, p.TemplateContext(), p.filename, p.firstLine)
	if err != nil {
		return
	}
	cn = buf.String()
	ex = cn
	pos := strings.Index(ex, p.site.Config().ExcerptSeparator)
	if pos >= 0 {
		ex = ex[:pos]
	}
	return
}

func (p *page) Excerpt() interface{} {
	if exc, ok := p.frontMatter["excerpt"]; ok {
		return exc
	}
	p.RLock()
	defer p.RUnlock()
	if p.rendered {
		return p.excerpt
	}
	return p.extractExcerpt()
}

func (p *page) extractExcerpt() []byte {
	raw := p.raw
	pos := bytes.Index(raw, []byte(p.site.Config().ExcerptSeparator))
	if pos >= 0 {
		return raw[:pos]
	}
	return raw
}
