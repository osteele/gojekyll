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

func (s *Site) prepareRendering() error {
	if !s.preparedToRender {
		if err := s.initializeRenderingPipeline(); err != nil {
			return err
		}
		if err := s.setPageContent(); err != nil {
			return err
		}
		s.preparedToRender = true
	}
	return nil
}

// WriteDocument writes the document to a writer.
func (s *Site) WriteDocument(w io.Writer, p pages.Document) error {
	if err := s.prepareRendering(); err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err := p.Write(buf); err != nil {
		return err
	}
	c := buf.Bytes()
	err := s.runHooks(func(p plugins.Plugin) error {
		c = p.PostRender(c)
		return nil
	})
	if err != nil {
		return err
	}
	_, err = w.Write(c)
	return err
}

// WritePages writes output files.
// It attends to options.dry_run.
func (s *Site) WritePages(options BuildOptions) (count int, err error) {
	errs := make(chan error)
	for _, p := range s.OutputPages() {
		count++
		go func(p pages.Document) {
			errs <- s.SavePage(p, options)
		}(p)
	}
	for i := 0; i < count; i++ {
		// might as well report the last error as the first
		// TODO return an aggregate
		if e := <-errs; e != nil {
			err = e
		}
	}
	return count, err
}

// SavePage writes a document to the destination directory.
// It attends to options.dry_run.
func (s *Site) SavePage(p pages.Document, options BuildOptions) error {
	from := filepath.Join(s.SourceDir(), filepath.ToSlash(p.SourcePath()))
	to := filepath.Join(s.DestDir(), p.Permalink())
	if !p.Static() && filepath.Ext(to) == "" {
		to = filepath.Join(to, "index.html")
	}
	if options.Verbose {
		fmt.Println("create", to, "from", p.SourcePath())
	}
	if options.DryRun {
		// FIXME render the page, just don't write it
		return nil
	}
	// nolint: gas
	if err := os.MkdirAll(filepath.Dir(to), 0755); err != nil {
		return err
	}
	switch {
	case p.Static() && options.UseHardLinks:
		return os.Link(from, to)
	case p.Static():
		return utils.CopyFileContents(to, from, 0644)
	default:
		return s.SaveDocumentToFile(p, to)
	}
}

// SaveDocumentToFile writes a page to filename.
func (s *Site) SaveDocumentToFile(d pages.Document, filename string) error {
	return utils.VisitCreatedFile(filename, func(w io.Writer) error {
		return s.WriteDocument(w, d)
	})
}
