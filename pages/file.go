package pages

import (
	"fmt"
	"os"
	"time"

	"github.com/osteele/gojekyll/frontmatter"
)

// file is embedded in StaticFile and page
type file struct {
	site      Site
	filename  string // target filepath
	relPath   string // slash-separated path relative to site or container source
	outputExt string
	permalink string // cached permalink
	modTime   time.Time
	dfm       FrontMatter // default frontMatter
	fm        FrontMatter
}

// NewFile creates a Page or StaticFile.
//
// filename is the absolute filename. relpath is the path relative to the site or collection directory.
func NewFile(s Site, filename string, relpath string, fm FrontMatter) (Document, error) {
	hasFM, err := frontmatter.FileHasFrontMatter(filename)
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}
	fields := file{
		site:      s,
		filename:  filename,
		dfm:       fm,
		fm:        fm,
		modTime:   info.ModTime(),
		relPath:   relpath,
		outputExt: s.Config().OutputExt(relpath),
	}
	if hasFM || !s.Config().RequiresFrontMatter(relpath) {
		return makePage(filename, fields)
	}
	fields.permalink = "/" + relpath
	p := &StaticFile{fields}
	return p, nil
}

func (f *file) String() string {
	return fmt.Sprintf("%T{Path=%v, Permalink=%v}", f, f.relPath, f.permalink)
}

func (f *file) OutputExt() string { return f.outputExt }
func (f *file) URL() string       { return f.permalink }
func (f *file) Published() bool   { return f.fm.Bool("published", true) }
func (f *file) Source() string    { return f.filename }

// const requiresReloadError = error.Error("requires reload")

func (f *file) Reload() error {
	info, err := os.Stat(f.filename)
	if err != nil {
		return err
	}
	f.modTime = info.ModTime()
	return nil
}
