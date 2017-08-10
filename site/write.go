package site

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/osteele/gojekyll/pages"
	"github.com/osteele/gojekyll/plugins"
	"github.com/osteele/gojekyll/utils"
)

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
		go func(d pages.Document) {
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
func (s *Site) WriteDoc(d pages.Document) error {
	from := d.SourcePath()
	rel := d.Permalink()
	if !d.Static() && filepath.Ext(rel) == "" {
		rel = filepath.Join(rel, "index.html")
	}
	to := filepath.Join(s.DestDir(), rel)
	if s.config.Verbose {
		fmt.Println("create", to, "from", d.SourcePath())
	}
	if s.config.DryRun {
		// FIXME render the page, just don't write it
		return nil
	}
	// nolint: gas
	if err := os.MkdirAll(filepath.Dir(to), 0755); err != nil {
		return err
	}
	switch {
	case d.Static():
		return utils.CopyFileContents(to, from, 0644)
	default:
		return utils.VisitCreatedFile(to, func(w io.Writer) error {
			return s.WriteDocument(w, d)
		})
	}
}

// WriteDocument writes the rendered document.
func (s *Site) WriteDocument(w io.Writer, d pages.Document) error {
	switch p := d.(type) {
	case pages.Page:
		return s.WritePage(w, p)
	default:
		return d.Write(w)
	}
}

// WritePage writes the rendered page.
func (s *Site) WritePage(w io.Writer, p pages.Page) error {
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
