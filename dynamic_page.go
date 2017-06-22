package gojekyll

import (
	"bytes"
	"io"
	"io/ioutil"
	"path/filepath"
	"regexp"

	"github.com/osteele/gojekyll/helpers"
	"github.com/osteele/gojekyll/liquid"
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
func (p *DynamicPage) Static() bool { return false }

// NewDynamicPage reads the front matter from a file to create a new DynamicPage.
func NewDynamicPage(f pageFields) (p *DynamicPage, err error) {
	data, err := ioutil.ReadFile(filepath.Join(f.container.SourceDir(), f.relpath))
	if err != nil {
		return
	}
	// Replace Windows linefeeds. This allows regular expressions to work.
	data = bytes.Replace(data, []byte("\r\n"), []byte("\n"), -1)

	frontMatter, err := readFrontMatter(&data)
	if err != nil {
		return
	}
	f.frontMatter = MergeVariableMaps(f.frontMatter, frontMatter)
	return &DynamicPage{
		pageFields: f,
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
func (p *DynamicPage) Variables() VariableMap {
	var (
		relpath = p.relpath
		ext     = filepath.Ext(relpath)
		root    = helpers.TrimExt(p.relpath)
		base    = filepath.Base(root)
		content = p.processed
	)

	if content == nil {
		content = &[]byte{}
	}

	data := VariableMap{
		"path":    relpath,
		"url":     p.Permalink(),
		"content": content,
		// TODO output

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
		"categories": []string{},
		"tags":       []string{},

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

// Context returns the local variables for template evaluation
func (p *DynamicPage) Context() VariableMap {
	return VariableMap{
		"page": p.Variables(),
		"site": p.container.SiteVariables(),
	}
}

// Output returns a bool indicating whether the page should be written.
func (p *DynamicPage) Output() bool {
	return p.pageFields.Output() && (p.collection == nil || p.collection.Output)
}

// Write applies Liquid and Markdown, as appropriate.
func (p *DynamicPage) Write(w io.Writer) (err error) {
	if p.processed != nil {
		_, err = w.Write(*p.processed)
		return
	}
	body, err := p.container.LiquidEngine().ParseAndRender(p.raw, p.Context())
	if err != nil {
		switch err := err.(type) {
		case *liquid.RenderError:
			if err.Filename == "" {
				err.Filename = p.Source()
			}
			if rel, e := filepath.Rel(p.container.SourceDir(), err.Filename); e == nil {
				err.Filename = rel
			}
			return err
		default:
			return helpers.PathError(err, "Liquid Error", p.Source())
		}
	}

	if p.IsMarkdown() {
		body = blackfriday.MarkdownCommon(body)
	}
	body, err = p.applyLayout(p.frontMatter, body)
	if err != nil {
		return
	}

	if p.container.IsSassPath(p.relpath) {
		return p.writeSass(w, body)
	}

	p.processed = &body
	_, err = w.Write(body)
	return
}
