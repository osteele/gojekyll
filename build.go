package gojekyll

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/osteele/gojekyll/helpers"
)

// BuildOptions holds options for Build and Clean
type BuildOptions struct {
	DryRun       bool
	UseHardLinks bool
	Verbose      bool
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
func (site *Site) Build(options BuildOptions) (count int, err error) {
	if err != nil {
		return
	}
	if err = site.Clean(options); err != nil {
		return
	}
	err = site.CopySassFileIncludes()
	if err != nil {
		return
	}
	for _, page := range site.Paths {
		count++
		if err = site.WritePage(page, options); err != nil {
			return
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
