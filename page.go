package gojekyll

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"time"

	"github.com/osteele/gojekyll/helpers"
	"github.com/osteele/gojekyll/liquid"
)

var (
	frontMatterMatcher     = regexp.MustCompile(`(?s)^---\n(.+?\n)---\n`)
	emptyFontMatterMatcher = regexp.MustCompile(`(?s)^---\n+---\n`)
)

// Page is a Jekyll page.
type Page interface {
	// Paths
	Path() string
	Permalink() string
	OutputExt() string
	Source() string // Source returns the file path of the page source.

	// Output
	Published() bool
	Static() bool
	Output() bool
	Write(Context, io.Writer) error

	// Variables
	Variables() VariableMap

	// internal
	initPermalink() error
}

// Context provides context information to a Page.
type Context interface {
	FindLayout(relname string, frontMatter *VariableMap) (liquid.Template, error)
	IsMarkdown(filename string) bool
	IsSassPath(filename string) bool
	SassIncludePaths() []string
	SiteVariables() VariableMap
	SourceDir() string
	TemplateEngine() liquid.Engine
	WriteSass(io.Writer, []byte) error
}

// pageFields is embedded in StaticPage and DynamicPage
type pageFields struct {
	relpath     string // relative to site source, e.g. "_post/base.ext"
	filename    string
	outputExt   string
	permalink   string // cached permalink
	modTime     time.Time
	frontMatter VariableMap // page front matter, merged with defaults
	collection  *Collection
	isMarkdown  bool
}

func (p *pageFields) String() string {
	return fmt.Sprintf("%s{Path=%v, Permalink=%v}", reflect.TypeOf(p).Name(), p.relpath, p.permalink)
}

func (p *pageFields) Path() string      { return p.relpath }
func (p *pageFields) Output() bool      { return p.Published() }
func (p *pageFields) Permalink() string { return p.permalink }
func (p *pageFields) Published() bool   { return p.frontMatter.Bool("published", true) }
func (p *pageFields) OutputExt() string { return p.outputExt }
func (p *pageFields) Source() string    { return p.filename }

// func (p *pageFields) IsMarkdown() bool { return p.isMarkdown }

// NewPageFromFile reads a Page from a file, using defaults as the default front matter.
func NewPageFromFile(ctx Context, collection *Collection, relpath string, defaults VariableMap) (p Page, err error) {
	filename := filepath.Join(ctx.SourceDir(), relpath)
	magic, err := helpers.ReadFileMagic(filename)
	if err != nil {
		return
	}
	info, err := os.Stat(filename)
	if err != nil {
		return
	}

	fields := pageFields{
		collection:  collection,
		modTime:     info.ModTime(),
		relpath:     relpath,
		filename:    filename,
		frontMatter: defaults,
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
	if string(magic) == "---\n" {
		p, err = newDynamicPageFromFile(filename, fields)
		if err != nil {
			return
		}
	} else {
		p = &StaticPage{fields}
	}
	// Compute this after creating the page, in order to pick up the front matter.
	err = p.initPermalink()
	if err != nil {
		return
	}
	return
}

// Variables returns the attributes of the template page object.
// See https://jekyllrb.com/docs/variables/#page-variables
func (p *pageFields) Variables() VariableMap {
	var (
		relpath = "/" + filepath.ToSlash(p.relpath)
		base    = path.Base(relpath)
		ext     = path.Ext(relpath)
	)

	return VariableMap{
		"path":          relpath,
		"modified_time": p.modTime,
		"name":          base,
		"basename":      helpers.TrimExt(base),
		"extname":       ext,
	}
}

// StaticPage is a static page.
type StaticPage struct {
	pageFields
}

// Static returns a bool indicating that the page is a static page.
func (p *StaticPage) Static() bool { return true }

// Variables returns metadata for use in the representation of the page as a collection item
func (p *StaticPage) Variables() VariableMap {
	return MergeVariableMaps(p.frontMatter, p.pageFields.Variables())
}

func (p *StaticPage) Write(_ Context, w io.Writer) error {
	b, err := ioutil.ReadFile(p.Source())
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}
