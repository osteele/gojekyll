package pages

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/osteele/gojekyll/helpers"
	"github.com/osteele/gojekyll/templates"
)

var (
	frontMatterMatcher     = regexp.MustCompile(`(?s)^---\n(.+?\n)---\n`)
	emptyFontMatterMatcher = regexp.MustCompile(`(?s)^---\n+---\n`)
)

// pageFields is embedded in StaticPage and DynamicPage
type pageFields struct {
	container   Container
	filename    string // target os filepath
	relpath     string // slash-separated path relative to site or container source
	outputExt   string
	permalink   string // cached permalink
	fileModTime time.Time
	frontMatter templates.VariableMap
}

func (p *pageFields) String() string {
	return fmt.Sprintf("%s{Path=%v, Permalink=%v}", reflect.TypeOf(p).Name(), p.relpath, p.permalink)
}

func (p *pageFields) OutputExt() string   { return p.outputExt }
func (p *pageFields) Path() string        { return p.relpath }
func (p *pageFields) Permalink() string   { return p.permalink }
func (p *pageFields) Published() bool     { return p.frontMatter.Bool("published", true) }
func (p *pageFields) SiteRelPath() string { return p.relpath }

// NewPageFromFile reads a Page from a file, using defaults as the default front matter.
func NewPageFromFile(ctx RenderingContext, c Container, filename string, relpath string, defaults templates.VariableMap) (Page, error) {
	magic, err := helpers.ReadFileMagic(filename)
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}

	fields := pageFields{
		container:   c,
		filename:    filename,
		frontMatter: defaults,
		fileModTime: info.ModTime(),
		relpath:     relpath,
		outputExt:   ctx.OutputExt(relpath),
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
		"modified_time": p.fileModTime,
		"name":          base,
		"basename":      helpers.TrimExt(base),
		"extname":       ext,
	})
}

func (p *pageFields) categories() []string {
	if v, found := p.frontMatter["categories"]; found {
		switch v := v.(type) {
		case string:
			return strings.Fields(v)
		case []interface{}:
			sl := make([]string, len(v))
			for i, s := range v {
				switch s := s.(type) {
				case fmt.Stringer:
					sl[i] = s.String()
				default:
					sl[i] = fmt.Sprint(s)
				}
			}
			return sl
		default:
			fmt.Printf("%T", v)
			panic("unimplemented")
		}
	}
	return []string{p.frontMatter.String("category", "")}
}
