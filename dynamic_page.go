package gojekyll

import (
	"bytes"
	"io"
	"io/ioutil"
	"path/filepath"
	"regexp"

	"github.com/osteele/gojekyll/helpers"
	"github.com/russross/blackfriday"

	yaml "gopkg.in/yaml.v2"
)

// DynamicPage is a static page, that includes frontmatter.
type DynamicPage struct {
	pageFields
	raw       []byte
	processed *[]byte
}

// Static returns a bool indicating that the page is a not static page.
func (page *DynamicPage) Static() bool { return false }

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
		raw:        data,
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

// Variables returns the attributes of the template page object.
func (page *DynamicPage) Variables() VariableMap {
	var (
		relpath = page.relpath
		ext     = filepath.Ext(relpath)
		root    = helpers.PathWithoutExtension(page.relpath)
		base    = filepath.Base(root)
		content = page.processed
	)

	if content == nil {
		content = &[]byte{}
	}

	data := VariableMap{
		"path":    relpath,
		"url":     page.Permalink(),
		"content": content,
		// TODO output

		// not documented, but present in both collection and non-collection pages
		"permalink": page.Permalink(),

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
		"relative_path": page.Path(),
		// TODO collection(name)

		// TODO undocumented; only present in collection pages:
		"ext": ext,
	}
	for k, v := range page.frontMatter {
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

// Context returns the local variables for template evaluation
func (page *DynamicPage) Context() VariableMap {
	return VariableMap{
		"page": page.Variables(),
		"site": page.site.Variables,
	}
}

// Output returns a bool indicating whether the page should be written.
func (page *DynamicPage) Output() bool {
	return page.pageFields.Output() && (page.collection == nil || page.collection.Output)
}

// Write applies Liquid and Markdown, as appropriate.
func (page *DynamicPage) Write(w io.Writer) (err error) {
	body, err := page.site.LiquidEngine().ParseAndRender(page.raw, page.Context())
	if err != nil {
		return helpers.PathError(err, "Liquid Error", page.Source())
	}

	if page.IsMarkdown() {
		body = blackfriday.MarkdownCommon(body)
		body, err = page.applyLayout(page.frontMatter, body)
		if err != nil {
			return
		}
	}

	if page.Site().IsSassPath(page.relpath) {
		return page.writeSass(w, body)
	}

	page.processed = &body
	_, err = w.Write(body)
	return
}
