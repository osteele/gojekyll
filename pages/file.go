package pages

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/osteele/gojekyll/helpers"
	"github.com/osteele/gojekyll/templates"
)

// file is embedded in StaticPage and DynamicPage
type file struct {
	container   Container
	filename    string // target os filepath
	relpath     string // slash-separated path relative to site or container source
	outputExt   string
	permalink   string // cached permalink
	fileModTime time.Time
	frontMatter templates.VariableMap
}

func (p *file) String() string {
	return fmt.Sprintf("%s{Path=%v, Permalink=%v}", reflect.TypeOf(p).Name(), p.relpath, p.permalink)
}

func (p *file) OutputExt() string   { return p.outputExt }
func (p *file) Path() string        { return p.relpath }
func (p *file) Permalink() string   { return p.permalink }
func (p *file) Published() bool     { return p.frontMatter.Bool("published", true) }
func (p *file) SiteRelPath() string { return p.relpath }

// NewFile creates a Post or StaticFile.
func NewFile(filename string, c Container, relpath string, defaults templates.VariableMap) (Document, error) {
	magic, err := helpers.ReadFileMagic(filename)
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}

	fields := file{
		container:   c,
		filename:    filename,
		frontMatter: defaults,
		fileModTime: info.ModTime(),
		relpath:     relpath,
		outputExt:   c.OutputExt(relpath),
	}
	var p Document
	if string(magic) == "---\n" {
		p, err = newPage(filename, fields)
		if err != nil {
			return nil, err
		}
	} else {
		p = &StaticFile{fields}
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
func (p *file) PageVariables() templates.VariableMap {
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

func (p *file) categories() []string {
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
