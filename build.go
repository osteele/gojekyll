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
func (s *Site) Build() (count int, err error) {
	if err = s.Clean(); err != nil {
		return
	}
	for _, page := range s.Paths {
		count++
		if err = s.WritePage(page); err != nil {
			return
		}
	}
	return
}

// WritePage writes a page to the destination directory.
func (s *Site) WritePage(page Page) error {
	src := filepath.Join(s.Source, page.Path())
	dst := filepath.Join(s.Dest, page.Path())
	if !page.Static() && filepath.Ext(dst) == "" {
		dst = filepath.Join(dst, "/index.html")
	}
	// nolint: gas
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	switch {
	case options.dryRun:
		fmt.Println("create", dst, "from", page.Source())
		return nil
	case page.Static() && options.useHardLinks:
		return os.Link(src, dst)
	case page.Static():
		return copyFile(dst, src, 0644)
	default:
		f, err := os.Create(dst)
		if err != nil {
			return err
		}
		if err := page.Write(f); err != nil {
			_ = f.Close()
			return err
		}
		return f.Close()
	}
}
