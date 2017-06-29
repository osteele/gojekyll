package sites

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/osteele/gojekyll/helpers"
	"github.com/osteele/gojekyll/pages"
)

func (s *Site) prepareRendering() error {
	if !s.preparedToRender {
		if err := s.initializeRenderingPipeline(); err != nil {
			return err
		}
		if err := s.initializeSiteVariables(); err != nil {
			return err
		}
		if err := s.setCollectionContent(); err != nil {
			return err
		}
		s.preparedToRender = true
	}
	return nil
}

// WriteDocument writes the document to w.
func (s *Site) WriteDocument(p pages.Document, w io.Writer) error {
	if err := s.prepareRendering(); err != nil {
		return err
	}
	return p.Write(s, w)
}

// WritePages writes output files.
// It attends to options.dry_run.
func (s *Site) WritePages(options BuildOptions) (count int, err error) {
	for _, p := range s.OutputPages() {
		count++
		if err = s.WritePage(p, options); err != nil {
			return
		}
	}
	return
}

// WritePage writes a page to the destination directory.
// It attends to options.dry_run.
func (s *Site) WritePage(p pages.Document, options BuildOptions) error {
	from := filepath.Join(s.Source, filepath.ToSlash(p.SiteRelPath()))
	to := filepath.Join(s.Destination, p.Permalink())
	if !p.Static() && filepath.Ext(to) == "" {
		to = filepath.Join(to, "index.html")
	}
	if options.Verbose {
		fmt.Println("create", to, "from", p.SiteRelPath())
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
		return helpers.CopyFileContents(to, from, 0644)
	default:
		return helpers.VisitCreatedFile(to, func(w io.Writer) error { return p.Write(s, w) })
	}
}
