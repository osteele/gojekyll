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
	"strings"

	"github.com/acstech/liquid"
	"github.com/russross/blackfriday"
	yaml "gopkg.in/yaml.v2"
)

var (
	frontmatterMatcher             = regexp.MustCompile(`(?s)^---\n(.+?\n)---\n`)
	templateVariableMatcher        = regexp.MustCompile(`:\w+\b`)
	nonAlphanumericSequenceMatcher = regexp.MustCompile(`[^[:alnum:]]+`)
)

// Page is a Jekyll page.
type Page interface {
	Path() string
	Source() string
	Static() bool
	Published() bool
	Permalink() string
	TemplateObject() VariableMap
	Write(io.Writer) error
	DebugVariables() VariableMap
}

type pageFields struct {
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

// StaticPage is a static page.
type StaticPage struct {
	pageFields
}

// DynamicPage is a static page, that includes frontmatter.
type DynamicPage struct {
	pageFields
	Content []byte
}

// Static returns a bool indicatingthat the page is a not static page.
func (p *DynamicPage) Static() bool { return false }

// Static returns a bool indicating that the page is a static page.
func (p *StaticPage) Static() bool { return true }

// TemplateObject returns metadata for use in the representation of the page as a collection item
func (p *StaticPage) TemplateObject() VariableMap {
	return mergeVariableMaps(p.frontMatter, p.pageFields.TemplateObject())
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

// TemplateVariables returns the local variables for template evaluation
func (p *DynamicPage) TemplateVariables() VariableMap {
	return VariableMap{
		"page": p.TemplateObject(),
		"site": site.Variables,
	}
}

// DebugVariables returns a map that's useful to present during diagnostics.
// For a static page, this is just the page's template object attributes.
func (p *pageFields) DebugVariables() VariableMap {
	return p.TemplateObject()
}

// DebugVariables returns a map that's useful to present during diagnostics.
// For a dynamic page, this is the local variable map that is used for template evaluation.
func (p *DynamicPage) DebugVariables() VariableMap {
	return p.TemplateVariables()
}

// ReadPage reads a Page from a file, using defaults as the default front matter.
func ReadPage(path string, defaults VariableMap) (p Page, err error) {
	// TODO don't read, parse binary files
	source, err := ioutil.ReadFile(filepath.Join(site.Source, path))
	if err != nil {
		return nil, err
	}

	if match := frontmatterMatcher.FindSubmatchIndex(source); match != nil {
		p, err = makeDynamicPage(path, defaults, source, match)
	} else {
		p, err = makeStaticPage(path, defaults)
	}
	return
}

func makeDynamicPage(path string, defaults VariableMap, source []byte, match []int) (*DynamicPage, error) {
	// TODO only prepend newlines if it's markdown
	body := append(
		regexp.MustCompile(`[^\n\r]+`).ReplaceAllLiteral(source[:match[1]], []byte{}),
		source[match[1]:]...)

	frontMatter := VariableMap{}
	if err := yaml.Unmarshal(source[match[2]:match[3]], &frontMatter); err != nil {
		err := &os.PathError{Op: "read frontmatter", Path: path, Err: err}
		return nil, err
	}
	frontMatter = mergeVariableMaps(defaults, frontMatter)

	pattern := frontMatter.String("permalink", ":path")
	permalink, err := expandPermalinkPattern(pattern, path, frontMatter)
	if err != nil {
		return nil, err
	}

	p := &DynamicPage{
		pageFields: pageFields{
			path:        path,
			permalink:   permalink,
			frontMatter: frontMatter,
		},
		Content: body,
	}
	return p, nil
}

func makeStaticPage(path string, frontMatter VariableMap) (*StaticPage, error) {
	permalink := "/" + path // TODO resolve same as for dynamic page
	p := &StaticPage{
		pageFields: pageFields{
			path:        path,
			permalink:   permalink,
			frontMatter: frontMatter,
		},
	}
	return p, nil
}

// Source returns the file path of the page source.
func (p *pageFields) Source() string {
	return filepath.Join(site.Source, p.path)
}

func (p *StaticPage) Write(w io.Writer) error {
	source, err := ioutil.ReadFile(p.Source())
	if err != nil {
		return err
	}
	_, err = w.Write(source)
	return err
}

// Write applies Liquid and Markdown, as appropriate.
func (p *DynamicPage) Write(w io.Writer) error {
	parsingLiquid := true
	defer func() {
		if parsingLiquid {
			fmt.Println("While processing", p.Source())
		}
	}()
	template, err := liquid.Parse(p.Content, nil)
	parsingLiquid = false
	if err != nil {
		err := &os.PathError{Op: "Liquid Error", Path: p.Source(), Err: err}
		return err
	}
	writer := new(bytes.Buffer)
	template.Render(writer, p.TemplateVariables())
	body := writer.Bytes()

	if isMarkdown(p.path) {
		body = blackfriday.MarkdownCommon(body)
	}

	_, err = w.Write(body)
	return err
}

func isMarkdown(path string) bool {
	ext := filepath.Ext(path)
	return site.MarkdownExtensions()[strings.TrimLeft(ext, ".")]
}

// replace each sequence of non-alphanumerics by a single hyphen
func hyphenateNonAlphaSequence(s string) string {
	return nonAlphanumericSequenceMatcher.ReplaceAllString(s, "-")
}
