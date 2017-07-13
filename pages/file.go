package pages

import (
	"fmt"
	"os"
	"time"

	"github.com/osteele/gojekyll/frontmatter"
	"github.com/osteele/gojekyll/templates"
	"github.com/osteele/gojekyll/utils"
)

// file is embedded in StaticFile and page
type file struct {
	site        Site
	filename    string // target os filepath
	relpath     string // slash-separated path relative to site or container source
	outputExt   string
	permalink   string // cached permalink
	fileModTime time.Time
	frontMatter map[string]interface{}
}

func (f *file) String() string {
	return fmt.Sprintf("%T{Path=%v, Permalink=%v}", f, f.relpath, f.permalink)
}

func (f *file) OutputExt() string  { return f.outputExt }
func (f *file) Path() string       { return utils.MustRel(f.site.Config().Source, f.filename) }
func (f *file) Permalink() string  { return f.permalink }
func (f *file) Published() bool    { return templates.VariableMap(f.frontMatter).Bool("published", true) }
func (f *file) SourcePath() string { return f.filename }

// NewFile creates a Page or StaticFile.
//
// filename is the absolute filename. relpath is the path relative to the site or collection directory.
func NewFile(s Site, filename string, relpath string, defaults map[string]interface{}) (Document, error) {
	fm, err := frontmatter.FileHasFrontMatter(filename)
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}

	fields := file{
		site:        s,
		filename:    filename,
		frontMatter: defaults,
		fileModTime: info.ModTime(),
		relpath:     relpath,
		outputExt:   s.OutputExt(relpath),
	}
	if fm {
		return makePage(filename, fields)
	}
	fields.permalink = "/" + relpath
	p := &StaticFile{fields}
	return p, nil
}
