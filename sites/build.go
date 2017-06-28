package sites

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/osteele/gojekyll/helpers"
	"github.com/osteele/gojekyll/pages"
)

// BuildOptions holds options for Build and Clean
type BuildOptions struct {
	DryRun       bool
	UseHardLinks bool
	Verbose      bool
}

// Clean the destination. Remove files that aren't in keep_files, and resulting empty diretories.
func (s *Site) Clean(options BuildOptions) error {
	// TODO PERF when called as part of build, keep files that will be re-generated
	removeFiles := func(name string, info os.FileInfo, err error) error {
		if options.Verbose {
			fmt.Println("rm", name)
		}
		switch {
		case err != nil && os.IsNotExist(err):
			return nil
		case err != nil:
			return err
		case info.IsDir():
			return nil
		case s.KeepFile(name):
			return nil
		case options.DryRun:
			return nil
		default:
			return os.Remove(name)
		}
	}
	if err := filepath.Walk(s.Destination, removeFiles); err != nil {
		return err
	}
	return helpers.RemoveEmptyDirectories(s.Destination)
}

// Build cleans the destination and create files in it.
// It attends to the global options.dry_run.
func (s *Site) Build(options BuildOptions) (int, error) {
	count := 0
	if err := s.InitializeRenderingPipeline(); err != nil {
		return 0, err
	}
	if err := s.Clean(options); err != nil {
		return 0, err
	}
	if err := s.SetPageContentTemplateValues(); err != nil {
		return 0, err
	}
	n, err := s.WritePages(options)
	return count + n, err
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
func (s *Site) WritePage(p pages.Page, options BuildOptions) error {
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
