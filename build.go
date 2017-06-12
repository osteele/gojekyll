package main

import (
	"io/ioutil"
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
		if !page.Static {
			p, err := readFile(page.Path, siteData, true)
			if err != nil {
				return err
			}
			page = p
		}
		// TODO don't do this for js, css, etc. pages
		if !page.Static && !strings.HasSuffix(path, ".html") {
			path = filepath.Join(path, "/index.html")
		}
		destPath := filepath.Join(siteConfig.DestinationDir, path)
		if err := os.MkdirAll(filepath.Dir(destPath), 0777); err != nil {
			return err
		}
		if page.Static {
			if err := os.Link(filepath.Join(siteConfig.SourceDir, page.Path), destPath); err != nil {
				return err
			}
		} else {
			// fmt.Println("render", filepath.Join(siteConfig.SourceDir, page.Path), "->", destPath)
			if err := ioutil.WriteFile(destPath, page.Body, 0644); err != nil {
				return err
			}
		}
	}
	return nil
}
