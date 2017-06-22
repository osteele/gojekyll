package pages

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"time"

	"github.com/osteele/gojekyll/helpers"
	"github.com/osteele/gojekyll/liquid"
	"github.com/osteele/gojekyll/templates"
)

var (
	frontMatterMatcher     = regexp.MustCompile(`(?s)^---\n(.+?\n)---\n`)
	emptyFontMatterMatcher = regexp.MustCompile(`(?s)^---\n+---\n`)
)

// Page is a Jekyll page.
type Page interface {
	// Paths
	SiteRelPath() string // relative to the site source directory
	Permalink() string   // relative URL path
	OutputExt() string

	// Output
	Published() bool
	Static() bool
	Output() bool
	Write(Context, io.Writer) error

	// Variables
	PageVariables() templates.VariableMap

	// internal
	initPermalink() error
}

// Context provides context information to a Page.
type Context interface {
	FindLayout(relname string, frontMatter *templates.VariableMap) (liquid.Template, error)
	IsMarkdown(filename string) bool
	IsSassPath(filename string) bool
	SassIncludePaths() []string
	SiteVariables() templates.VariableMap
	SourceDir() string
	TemplateEngine() liquid.Engine
	WriteSass(io.Writer, []byte) error
}

// Container is the Page container
type Container interface {
	PathPrefix() string
	Output() bool
}

// pageFields is embedded in StaticPage and DynamicPage
type pageFields struct {
	container          Container
	filename           string
	relpath            string // relative to site source directory
	outputExt          string
	permalink          string // cached permalink
	modTime            time.Time
	frontMatter templates.VariableMap
	isMarkdown         bool
}

func (p *pageFields) String() string {
	return fmt.Sprintf("%s{Path=%v, Permalink=%v}", reflect.TypeOf(p).Name(), p.relpath, p.permalink)
}

func (p *pageFields) Output() bool        { return p.Published() }
func (p *pageFields) OutputExt() string   { return p.outputExt }
func (p *pageFields) Path() string        { return p.relpath }
func (p *pageFields) Permalink() string   { return p.permalink }
func (p *pageFields) Published() bool     { return p.frontMatter.Bool("published", true) }
func (p *pageFields) SiteRelPath() string { return p.relpath }

// NewPageFromFile reads a Page from a file, using defaults as the default front matter.
func NewPageFromFile(ctx Context, c Container, filename string, relpath string, defaults templates.VariableMap) (Page, error) {
	magic, err := helpers.ReadFileMagic(filename)
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}

	fields := pageFields{
		container:          c,
		filename:           filename,
		frontMatter: defaults,
		modTime:            info.ModTime(),
		relpath:            relpath,
	}
	switch {
	case ctx.IsMarkdown(relpath):
		fields.isMarkdown = true
		fields.outputExt = ".html"
	case ctx.IsSassPath(relpath):
		fields.outputExt = ".css"
	default:
		fields.outputExt = filepath.Ext(relpath)
	}
	var p Page
	if string(magic) == "---\n" {
		p, err = newDynamicPageFromFile(filename, fields)
		if err != nil {
			return nil, err
		}
	} else {
		p = &StaticPage{fields}
	}
	// Compute this after creating the page, in order to pick up the front matter.
	err = p.initPermalink()
	if err != nil {
		return nil, err
	}
	return p, nil
}

// Variables returns the attributes of the template page object.
// See https://jekyllrb.com/docs/variables/#page-variables
func (p *pageFields) PageVariables() templates.VariableMap {
	var (
		relpath = "/" + filepath.ToSlash(p.relpath)
		base    = path.Base(relpath)
		ext     = path.Ext(relpath)
	)

	return templates.MergeVariableMaps(p.frontMatter, templates.VariableMap{
		"path":          relpath,
		"modified_time": p.modTime,
		"name":          base,
		"basename":      helpers.TrimExt(base),
		"extname":       ext,
	})
}
