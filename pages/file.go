package pages

import (
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/osteele/gojekyll/frontmatter"
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

func (f *file) OutputExt() string  { return f.outputExt }
func (f *file) Path() string       { return f.relpath }
func (f *file) Permalink() string  { return f.permalink }
func (f *file) Published() bool    { return templates.VariableMap(f.frontMatter).Bool("published", true) }
func (f *file) SourcePath() string { return f.relpath }

// NewFile creates a Post or StaticFile.
func NewFile(filename string, c Container, relpath string, defaults map[string]interface{}) (Document, error) {
	fm, err := frontmatter.FileHasFrontMatter(filename)
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
	if fm {
		return makePage(filename, fields)
	}
	p := &StaticFile{fields}
	if err = p.setPermalink(); err != nil {
		return nil, err
	}
	return p, nil
}

// Categories is in the File interface
func (f *file) Categories() []string {
	return frontmatter.FrontMatter(f.frontMatter).SortedStringArray("categories")
}

// Tags is in the File interface
func (f *file) Tags() []string {
	return frontmatter.FrontMatter(f.frontMatter).SortedStringArray("tags")
}
