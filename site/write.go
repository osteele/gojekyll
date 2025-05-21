package site

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/osteele/gojekyll/plugins"
	"github.com/osteele/gojekyll/utils"
)

// Write cleans the destination and writes files into it.
// It sets TZ from the site config.
func (s *Site) Write() (int, error) {
	if err := s.setTimeZone(); err != nil {
		return 0, err
	}
	if err := s.ensureRendered(); err != nil {
		return 0, err
	}
	if err := s.Clean(); err != nil {
		return 0, err
	}
	return s.WriteFiles()
}

// WriteFiles writes output files.
func (s *Site) WriteFiles() (count int, err error) {
	errs := make(chan error)
	// without this, large sites run out of file descriptors
	sem := make(chan bool, 20)
	for i, n := 0, cap(sem); i < n; i++ {
		sem <- true
	}
	for _, d := range s.OutputDocs() {
		count++
		go func(d Document) {
			<-sem
			errs <- s.WriteDoc(d)
			sem <- true
		}(d)
	}
	var errList []error
	for i := 0; i < count; i++ {
		if e := <-errs; e != nil {
			errList = append(errList, e)
		}
	}
	return count, combineErrors(errList)
}

// WriteDoc writes a document to the destination directory.
func (s *Site) WriteDoc(d Document) error {
	from := d.Source()
	rel := d.URL()
	if !d.IsStatic() && filepath.Ext(rel) == "" {
		rel = filepath.Join(rel, "index.html")
	}
	to := filepath.Join(s.DestDir(), rel)
	if s.cfg.Verbose {
		fmt.Println("create", to, "from", d.Source())
	}
	if s.cfg.DryRun {
		// FIXME render the page, just don't write it
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(to), 0755); err != nil {
		return err
	}
	switch {
	case d.IsStatic():
		return utils.CopyFileContents(to, from, 0644)
	default:
		return utils.VisitCreatedFile(to, func(w io.Writer) error {
			return s.WriteDocument(w, d)
		})
	}
}

// WriteDocument writes the rendered document.
func (s *Site) WriteDocument(w io.Writer, d Document) error {
	switch p := d.(type) {
	case Page:
		return s.WritePage(w, p)
	default:
		return d.Write(w)
	}
}

// WritePage writes the rendered page. It is called as part of site.Write,
// but also, in an incremental build, to write a single page – therefore it
// also ensures that all pages have been rendered before writing this one.
func (s *Site) WritePage(w io.Writer, p Page) error {
	if err := s.ensureRendered(); err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err := p.Write(buf); err != nil {
		return err
	}
	b := buf.Bytes()
	err := s.runHooks(func(p plugins.Plugin) (err error) {
		b, err = p.PostRender(b)
		return
	})
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}
