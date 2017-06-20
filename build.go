package gojekyll

import (
	"fmt"
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
func (site *Site) Pages() []Page {
	pages := make([]Page, len(site.Paths))
	i := 0
	for _, p := range site.Paths {
		pages[i] = p
		i++
	}
	return pages
}

// Pages is a list of pages.
func (coll *Collection) Pages() []Page {
	return coll.pages
}

// Clean the destination. Remove files that aren't in keep_files, and resulting empty diretories.
func (site *Site) Clean(options BuildOptions) error {
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
		case site.KeepFile(name):
			return nil
		case options.DryRun:
			return nil
		default:
			return os.Remove(name)
		}
	}
	if err := filepath.Walk(site.Destination, removeFiles); err != nil {
		return err
	}
	return helpers.RemoveEmptyDirectories(site.Destination)
}

// Build cleans the destination and create files in it.
// It attends to the global options.dry_run.
func (site *Site) Build(options BuildOptions) (int, error) {
	count := 0
	if err := site.Clean(options); err != nil {
		return count, err
	}
	if err := site.CopySassFileIncludes(); err != nil {
		return count, err
	}
	for _, coll := range site.Collections {
		n, err := site.WritePages(coll, options)
		if err != nil {
			return count, err
		}
		count += n
	}
	site.updateCollectionVariables()
	n, err := site.WritePages(site, options)
	return count + n, err
}

// WritePages cleans the destination and create files in it.
// It attends to the global options.dry_run.
func (site *Site) WritePages(container PageContainer, options BuildOptions) (count int, err error) {
	for _, page := range container.Pages() {
		if page.Output() {
			count++
			if err = site.WritePage(page, options); err != nil {
				return
			}
		}
	}
	return
}

// WritePage writes a page to the destination directory.
func (site *Site) WritePage(page Page, options BuildOptions) error {
	from := filepath.Join(site.Source, page.Path())
	to := filepath.Join(site.Destination, page.Permalink())
	if !page.Static() && filepath.Ext(to) == "" {
		to = filepath.Join(to, "/index.html")
	}
	if options.Verbose {
		fmt.Println("create", to, "from", page.Source())
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
		return helpers.VisitCreatedFile(to, page.Write)
	}
}
