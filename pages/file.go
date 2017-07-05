package pages

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/osteele/gojekyll/helpers"
	"github.com/osteele/gojekyll/templates"
	"github.com/osteele/liquid/evaluator"
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
	if err = p.setPermalink(); err != nil {
		return nil, err
	}
	return p, nil
}

// Categories is in the File interface
func (f *file) Categories() []string {
	return sortedStringValue(f.frontMatter["categories"])
}

// Categories is in the File interface
func (f *file) Tags() []string {
	return sortedStringValue(f.frontMatter["tags"])
}

func sortedStringValue(field interface{}) []string {
	out := []string{}
	switch value := field.(type) {
	case string:
		out = strings.Fields(value)
	case []interface{}:
		if c, e := evaluator.Convert(value, reflect.TypeOf(out)); e == nil {
			out = c.([]string)
		}
	case []string:
		out = value
	}
	sort.Strings(out)
	return out
}
