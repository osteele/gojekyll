package site

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/osteele/gojekyll/pages"
	"github.com/osteele/gojekyll/plugins"
	"github.com/osteele/gojekyll/utils"
)

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
	var errList []error
	for i := 0; i < count; i++ {
		if e := <-errs; e != nil {
			errList = append(errList, e)
		}
	}
	return count, combineErrors(errList)
}

// SavePage writes a document to the destination directory.
// It attends to options.dry_run.
func (s *Site) SavePage(p pages.Document, options BuildOptions) error {
	from := p.SourcePath()
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

// WriteDocument writes the document to a writer.
func (s *Site) WriteDocument(w io.Writer, d pages.Document) error {
	if err := s.prepareRendering(); err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err := d.Write(buf); err != nil {
		return err
	}
	b := buf.Bytes()
	err := s.runHooks(func(p plugins.Plugin) error {
		b = p.PostRender(b)
		return nil
	})
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}

func combineErrors(errs []error) error {
	switch len(errs) {
	case 0:
		return nil
	case 1:
		return errs[0]
	default:
		messages := make([]string, len(errs))
		for i, e := range errs {
			messages[i] = e.Error()
		}
		return fmt.Errorf(strings.Join(messages, "\n"))
	}
}

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
