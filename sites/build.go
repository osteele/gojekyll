package sites

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
	if err := s.prepareRendering(); err != nil {
		return 0, err
	}
	if err := s.Clean(options); err != nil {
		return 0, err
	}
	n, err := s.WritePages(options)
	return count + n, err
}
