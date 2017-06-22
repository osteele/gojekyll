package gojekyll

import (
	"bytes"
	"io"
	"io/ioutil"
	"path/filepath"
	"regexp"

	"github.com/osteele/gojekyll/helpers"
	"github.com/osteele/gojekyll/liquid"
	"github.com/osteele/gojekyll/templates"
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

func newDynamicPageFromFile(filename string, f pageFields) (*DynamicPage, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	frontMatter, err := readFrontMatter(&b)
	if err != nil {
		return nil, err
	}
	f.frontMatter = templates.MergeVariableMaps(f.frontMatter, frontMatter)
	return &DynamicPage{
		pageFields: f,
		raw:        b,
	}, nil
}

func readFrontMatter(sourcePtr *[]byte) (frontMatter templates.VariableMap, err error) {
	var (
		source = *sourcePtr
		start  = 0
	)
	// Replace Windows linefeeds. This allows the following regular expressions to work.
	source = bytes.Replace(source, []byte("\r\n"), []byte("\n"), -1)
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
func (p *DynamicPage) Variables() templates.VariableMap {
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

	data := templates.VariableMap{
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

// TemplateContext returns the local variables for template evaluation
func (p *DynamicPage) TemplateContext(ctx Context) templates.VariableMap {
	return templates.VariableMap{
		"page": p.Variables(),
		"site": ctx.SiteVariables(),
	}
}

// Output returns a bool indicating whether the page should be written.
func (p *DynamicPage) Output() bool {
	return p.pageFields.Output() && p.container.Output()
}

// Write applies Liquid and Markdown, as appropriate.
func (p *DynamicPage) Write(ctx Context, w io.Writer) error {
	if p.processed != nil {
		_, err := w.Write(*p.processed)
		return err
	}
	body, err := ctx.TemplateEngine().ParseAndRender(p.raw, p.TemplateContext(ctx))
	if err != nil {
		switch err := err.(type) {
		case *liquid.RenderError:
			if err.Filename == "" {
				err.Filename = p.filename
			}
			if rel, e := filepath.Rel(ctx.SourceDir(), err.Filename); e == nil {
				err.Filename = rel
			}
			return err
		default:
			return helpers.PathError(err, "Liquid Error", p.filename)
		}
	}

	if p.isMarkdown {
		body = blackfriday.MarkdownCommon(body)
	}
	body, err = p.applyLayout(ctx, p.frontMatter, body)
	if err != nil {
		return err
	}

	if ctx.IsSassPath(p.relpath) {
		return ctx.WriteSass(w, body)
	}

	p.processed = &body
	_, err = w.Write(body)
	return err
}

func (p *DynamicPage) applyLayout(ctx Context, frontMatter templates.VariableMap, body []byte) ([]byte, error) {
	for {
		name := frontMatter.String("layout", "")
		if name == "" {
			return body, nil
		}
		template, err := ctx.FindLayout(name, &frontMatter)
		if err != nil {
			return nil, err
		}
		vars := templates.MergeVariableMaps(p.TemplateContext(ctx), templates.VariableMap{
			"content": string(body),
			"layout":  frontMatter,
		})
		body, err = template.Render(vars)
		if err != nil {
			return nil, err
		}
	}
}
