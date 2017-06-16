package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"

	. "github.com/osteele/gojekyll/helpers"

	yaml "gopkg.in/yaml.v2"

	"github.com/acstech/liquid"
	"github.com/russross/blackfriday"
)

var (
	frontMatterMatcher     = regexp.MustCompile(`(?s)^---\n(.+?\n)---\n`)
	emptyFontMatterMatcher = regexp.MustCompile(`(?s)^---\n+---\n`)
)

// Page is a Jekyll page.
type Page interface {
	Path() string
	Site() *Site
	Source() string
	Static() bool
	Published() bool
	Permalink() string
	TemplateObject() VariableMap
	Write(io.Writer) error
	DebugVariables() VariableMap

	initPermalink() error
}

type pageFields struct {
	site        *Site
	path        string // this is the relative path
	permalink   string
	frontMatter VariableMap
}

func (p *pageFields) String() string {
	return fmt.Sprintf("%s{Path=%v, Permalink=%v}",
		reflect.TypeOf(p).Name(), p.path, p.permalink)
}

func (p *pageFields) Path() string      { return p.path }
func (p *pageFields) Permalink() string { return p.permalink }
func (p *pageFields) Published() bool {
	return p.frontMatter.Bool("published", true)
}
func (p *pageFields) Site() *Site { return p.site }

// ReadPage reads a Page from a file, using defaults as the default front matter.
func ReadPage(site *Site, rel string, defaults VariableMap) (p Page, err error) {
	magic, err := ReadFileMagic(filepath.Join(site.Source, rel))
	if err != nil {
		return
	}

	fields := pageFields{site: site, path: rel, frontMatter: defaults}
	if string(magic) == "---\n" {
		p, err = NewDynamicPage(fields)
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

func (p *StaticPage) Write(w io.Writer) error {
	source, err := ioutil.ReadFile(p.Source())
	if err != nil {
		return err
	}
	_, err = w.Write(source)
	return err
}

// TemplateObject returns the attributes of the template page object.
// See https://jekyllrb.com/docs/variables/#page-variables
func (p *pageFields) TemplateObject() VariableMap {
	var (
		path = "/" + p.path
		base = filepath.Base(path)
		ext  = filepath.Ext(path)
	)

	return VariableMap{
		"path":          path,
		"modified_time": 0, // TODO
		"name":          base,
		"basename":      base[:len(base)-len(ext)],
		"extname":       ext,
	}
}

// DebugVariables returns a map that's useful to present during diagnostics.
// For a static page, this is just the page's template object attributes.
func (p *pageFields) DebugVariables() VariableMap {
	return p.TemplateObject()
}

// Source returns the file path of the page source.
func (p *pageFields) Source() string {
	return filepath.Join(p.site.Source, p.path)
}

// StaticPage is a static page.
type StaticPage struct {
	pageFields
}

// Static returns a bool indicating that the page is a static page.
func (p *StaticPage) Static() bool { return true }

// TemplateObject returns metadata for use in the representation of the page as a collection item
func (p *StaticPage) TemplateObject() VariableMap {
	return MergeVariableMaps(p.frontMatter, p.pageFields.TemplateObject())
}

// DynamicPage is a static page, that includes frontmatter.
type DynamicPage struct {
	pageFields
	Content []byte
}

// Static returns a bool indicating that the page is a not static page.
func (p *DynamicPage) Static() bool { return false }

// NewDynamicPage reads the front matter from a file to create a new DynamicPage.
func NewDynamicPage(fields pageFields) (p *DynamicPage, err error) {
	data, err := ioutil.ReadFile(filepath.Join(fields.site.Source, fields.path))
	if err != nil {
		return
	}
	data = bytes.Replace(data, []byte("\r"), []byte("\n"), -1)

	frontMatter, err := readFrontMatter(&data)
	if err != nil {
		return
	}
	fields.frontMatter = MergeVariableMaps(fields.frontMatter, frontMatter)
	return &DynamicPage{
		pageFields: fields,
		Content:    data,
	}, nil
}

func readFrontMatter(sourcePtr *[]byte) (frontMatter VariableMap, err error) {
	var (
		source = *sourcePtr
		start  = 0
	)
	if match := frontMatterMatcher.FindSubmatchIndex(source); match != nil {
		start = match[1]
		if err = yaml.Unmarshal(source[match[2]:match[3]], &frontMatter); err != nil {
			return
		}
	} else if match := emptyFontMatterMatcher.FindSubmatchIndex(source); match != nil {
		start = match[1]
	}
	// This fixes the line numbers for template errors
	// TODO find a less hacky solution
	*sourcePtr = append(
		regexp.MustCompile(`[^\n\r]+`).ReplaceAllLiteral(source[:start], []byte{}),
		source[start:]...)
	return
}

// TemplateObject returns the attributes of the template page object.
func (p *DynamicPage) TemplateObject() VariableMap {
	var (
		path = p.path
		ext  = filepath.Ext(path)
		root = p.path[:len(path)-len(ext)]
		base = filepath.Base(root)
	)

	data := VariableMap{
		"path": p.path,
		"url":  p.Permalink(),
		// TODO content output

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

// TemplateVariables returns the local variables for template evaluation
func (p *DynamicPage) TemplateVariables() VariableMap {
	return VariableMap{
		"page": p.TemplateObject(),
		"site": p.site.Variables,
	}
}

// DebugVariables returns a map that's useful to present during diagnostics.
// For a dynamic page, this is the local variable map that is used for template evaluation.
func (p *DynamicPage) DebugVariables() VariableMap {
	return p.TemplateVariables()
}

// renderTemplate is a wrapper around liquid template.Render that turns panics into errors
func renderTemplate(template *liquid.Template, variables VariableMap) (bs []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				panic(r)
			}
		}
	}()
	writer := new(bytes.Buffer)
	template.Render(writer, variables)
	return writer.Bytes(), nil
}

// applyTemplate parses and then renders the template.
func parseAndApplyTemplate(bs []byte, variables VariableMap) ([]byte, error) {
	template, err := liquid.Parse(bs, nil)
	if err != nil {
		return nil, err
	}
	return renderTemplate(template, variables)
}

// Write applies Liquid and Markdown, as appropriate.
func (p *DynamicPage) Write(w io.Writer) (err error) {
	body, err := parseAndApplyTemplate(p.Content, p.TemplateVariables())
	if err != nil {
		err = &os.PathError{Op: "Liquid Error", Path: p.Source(), Err: err}
		return
	}

	if p.Site().IsMarkdown(p.path) {
		body = blackfriday.MarkdownCommon(body)
		body, err = p.applyLayout(p.frontMatter, body)
		if err != nil {
			return
		}
	}

	if p.Site().IsSassPath(p.path) {
		return p.writeSass(w, body)
	}

	_, err = w.Write(body)
	return
}
