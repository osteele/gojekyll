package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func cleanDirectory() error {
	removeFiles := func(path string, info os.FileInfo, err error) error {
		stat, err := os.Stat(path)
		if os.IsNotExist(err) {
			return nil
		}
		if err != nil {
			return err
		}

		if stat.IsDir() {
			return nil
		}
		// TODO check for inclusion in KeepFiles
		err = os.Remove(path)
		return err
	}
	err := filepath.Walk(siteConfig.DestinationDir, removeFiles)
	if err == nil {
		err = removeEmptyDirectories(siteConfig.DestinationDir)
	}
	return err
}

func build() error {
	err := cleanDirectory()
	if err != nil {
		return err
	}
	for _, page := range siteMap {
		if !page.Static {
			page, err = readFile(page.Path, true)
		}
		if err != nil {
			return err
		}
		path := page.Permalink
		// TODO only do this for MIME pages
		if !page.Static && !strings.HasSuffix(path, ".html") {
			path += "/index.html"
		}
		destPath := filepath.Join(siteConfig.DestinationDir, path)
		os.MkdirAll(filepath.Dir(destPath), 0777)
		if page.Static {
			os.Link(filepath.Join(siteConfig.SourceDir, page.Path), destPath)
		} else {
			fmt.Println("render", filepath.Join(siteConfig.SourceDir, page.Path), "->", destPath)
			ioutil.WriteFile(destPath, page.Body, 0644)
		}
	}
	return nil
}
