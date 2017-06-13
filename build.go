package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// Clean the destination. Remove files that aren't in keep_files, and resulting empty diretories.
// It attends to the global options.dry_run.
func (s *Site) Clean() error {
	removeFiles := func(path string, info os.FileInfo, err error) error {
		switch {
		case err != nil && os.IsNotExist(err):
			return nil
		case err != nil:
			return err
		case info.IsDir():
			return nil
		case site.KeepFile(path):
			return nil
		case options.dryRun:
			fmt.Println("rm", path)
		default:
			return os.Remove(path)
		}
		return nil
	}
	if err := filepath.Walk(s.Dest, removeFiles); err != nil {
		return err
	}
	return RemoveEmptyDirectories(s.Dest)
}

// Build cleans the destination and create files in it.
// It attends to the global options.dry_run.
func (s *Site) Build() error {
	if err := s.Clean(); err != nil {
		return err
	}
	for path, page := range s.Paths {
		if !page.Static && filepath.Ext(path) == "" {
			path = filepath.Join(path, "/index.html")
		}
		src := filepath.Join(s.Source, page.Path)
		dst := filepath.Join(s.Dest, path)
		if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
			return err
		}
		switch {
		case options.dryRun:
			fmt.Println("create", dst, "from", page.Source())
		case page.Static && options.useHardLinks:
			if err := os.Link(src, dst); err != nil {
				return err
			}
		case page.Static:
			if err := copyFile(dst, src, 0644); err != nil {
				return err
			}
		default:
			f, err := os.Create(dst)
			if err != nil {
				return err
			}
			defer func() { _ = f.Close() }()
			if err := page.Render(f); err != nil {
				return err
			}
		}
	}
	return nil
}
