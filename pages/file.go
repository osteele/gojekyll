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

// file is embedded in StaticFile and page
type file struct {
	container   Container
	filename    string // target os filepath
	relpath     string // slash-separated path relative to site or container source
	outputExt   string
	permalink   string // cached permalink
	fileModTime time.Time
	frontMatter map[string]interface{}
}

func (f *file) String() string {
	return fmt.Sprintf("%s{Path=%v, Permalink=%v}", reflect.TypeOf(f).Name(), f.relpath, f.permalink)
}

func (f *file) OutputExt() string   { return f.outputExt }
func (f *file) Path() string        { return f.relpath }
func (f *file) Permalink() string   { return f.permalink }
func (f *file) Published() bool     { return templates.VariableMap(f.frontMatter).Bool("published", true) }
func (f *file) SiteRelPath() string { return f.relpath }

// NewFile creates a Post or StaticFile.
func NewFile(filename string, c Container, relpath string, defaults map[string]interface{}) (Document, error) {
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
func (f *file) PageVariables() map[string]interface{} {
	var (
		relpath = "/" + filepath.ToSlash(f.relpath)
		base    = path.Base(relpath)
		ext     = path.Ext(relpath)
	)

	return templates.MergeVariableMaps(f.frontMatter, map[string]interface{}{
		"path":          relpath,
		"modified_time": f.fileModTime,
		"name":          base,
		"basename":      helpers.TrimExt(base),
		"extname":       ext,
	})
}

func (f *file) categories() []string {
	if v, found := f.frontMatter["categories"]; found {
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
	return []string{templates.VariableMap(f.frontMatter).String("category", "")}
}
