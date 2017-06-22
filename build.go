package gojekyll

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/osteele/gojekyll/helpers"
)

// PageContainer has a slice of pages
type PageContainer interface {
	Pages() []Page
}

// BuildOptions holds options for Build and Clean
type BuildOptions struct {
	DryRun       bool
	UseHardLinks bool
	Verbose      bool
}

// Pages is a list of pages.
func (s *Site) Pages() []Page {
	pages := make([]Page, len(s.Paths))
	i := 0
	for _, p := range s.Paths {
		pages[i] = p
		i++
	}
	return pages
}

// Pages is a list of pages.
func (c *Collection) Pages() []Page {
	return c.pages
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
	if err := s.Clean(options); err != nil {
		return count, err
	}
	if err := s.CopySassFileIncludes(); err != nil {
		return count, err
	}
	for _, coll := range s.Collections {
		n, err := s.WritePages(coll, options)
		if err != nil {
			return count, err
		}
		count += n
	}
	s.updateCollectionVariables()
	n, err := s.WritePages(s, options)
	return count + n, err
}

// WritePages cleans the destination and create files in it.
// It attends to the global options.dry_run.
func (s *Site) WritePages(container PageContainer, options BuildOptions) (count int, err error) {
	for _, page := range container.Pages() {
		if page.Output() {
			count++
			if err = s.WritePage(page, options); err != nil {
				return
			}
		}
	}
	return
}

// WritePage writes a page to the destination directory.
func (s *Site) WritePage(page Page, options BuildOptions) error {
	from := filepath.Join(s.Source, page.SiteRelPath())
	to := filepath.Join(s.Destination, page.Permalink())
	if !page.Static() && filepath.Ext(to) == "" {
		to = filepath.Join(to, "/index.html")
	}
	if options.Verbose {
		fmt.Println("create", to, "from", page.SiteRelPath())
	}
	if options.DryRun {
		// FIXME render the page in dry run mode, just don't write it
		return nil
	}
	// nolint: gas
	if err := os.MkdirAll(filepath.Dir(to), 0755); err != nil {
		return err
	}
	switch {
	case page.Static() && options.UseHardLinks:
		return os.Link(from, to)
	case page.Static():
		return helpers.CopyFileContents(to, from, 0644)
	default:
		return helpers.VisitCreatedFile(to, func(w io.Writer) error { return page.Write(s, w) })
	}
}
