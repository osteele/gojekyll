package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/osteele/gojekyll/helpers"
	"github.com/russross/blackfriday"

	yaml "gopkg.in/yaml.v2"
)

// DynamicPage is a static page, that includes frontmatter.
type DynamicPage struct {
	pageFields
	Content []byte
}

// Static returns a bool indicating that the page is a not static page.
func (p *DynamicPage) Static() bool { return false }

// NewDynamicPage reads the front matter from a file to create a new DynamicPage.
func NewDynamicPage(fields pageFields) (p *DynamicPage, err error) {
	data, err := ioutil.ReadFile(filepath.Join(fields.site.Source, fields.relpath))
	if err != nil {
		return
	}
	// Replace Windows linefeeds. This allows regular expressions to work.
	data = bytes.Replace(data, []byte("\r\n"), []byte("\n"), -1)

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
	// TODO find a less hack-ey solution
	*sourcePtr = append(
		regexp.MustCompile(`[^\n\r]+`).ReplaceAllLiteral(source[:start], []byte{}),
		source[start:]...)
	return
}

// TemplateObject returns the attributes of the template page object.
func (p *DynamicPage) TemplateObject() VariableMap {
	var (
		relpath = p.relpath
		ext     = filepath.Ext(relpath)
		root    = helpers.PathWithoutExtension(p.relpath)
		base    = filepath.Base(root)
	)

	data := VariableMap{
		"path": relpath,
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

// Write applies Liquid and Markdown, as appropriate.
func (p *DynamicPage) Write(w io.Writer) (err error) {
	config := p.site.LiquidConfiguration()
	body, err := helpers.ParseAndApplyTemplate(p.Content, p.TemplateVariables(), config)
	if err != nil {
		err = &os.PathError{Op: "Liquid Error", Path: p.Source(), Err: err}
		return
	}

	if p.Site().IsMarkdown(p.relpath) {
		body = blackfriday.MarkdownCommon(body)
		body, err = p.applyLayout(p.frontMatter, body, config)
		if err != nil {
			return
		}
	}

	if p.Site().IsSassPath(p.relpath) {
		return p.writeSass(w, body)
	}

	_, err = w.Write(body)
	return
}
