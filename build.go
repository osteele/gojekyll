package main

import (
	"os"
	"path/filepath"
	"strings"
)

func cleanDirectory() error {
	removeFiles := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
		if info.IsDir() {
			return nil
		}
		// TODO check for inclusion in KeepFiles
		err = os.Remove(path)
		return err
	}
	if err := filepath.Walk(siteConfig.DestinationDir, removeFiles); err != nil {
		return err
	}
	return removeEmptyDirectories(siteConfig.DestinationDir)
}

func build() error {
	if err := cleanDirectory(); err != nil {
		return err
	}
	for path, page := range siteMap {
		// TODO don't do this for js, css, etc. pages
		if !page.Static && !strings.HasSuffix(path, ".html") {
			path = filepath.Join(path, "/index.html")
		}
		src := filepath.Join(siteConfig.SourceDir, page.Path)
		dst := filepath.Join(siteConfig.DestinationDir, path)
		if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
			return err
		}
		switch {
		case page.Static && options.useHardLinks:
			if err := os.Link(src, dst); err != nil {
				return err
			}
		case page.Static:
			if err := copyFile(dst, src, 0644); err != nil {
				return err
			}
		default:
			// fmt.Println("render", filepath.Join(siteConfig.SourceDir, page.Path), "->", dst)
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
