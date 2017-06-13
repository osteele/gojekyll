package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// Collection is a Jekyll collection.
type Collection struct {
	Name   string
	Data   map[interface{}]interface{}
	Output bool
	Pages  []*Page
}

func makeCollection(name string, data map[interface{}]interface{}) *Collection {
	return &Collection{
		Name:   name,
		Data:   data,
		Output: getBool(data, "output", false),
	}
}

// PageData returns an array of a page data, for use as the template variable
// value of the collection.
func (c *Collection) PageData() (d []interface{}) {
	for _, p := range c.Pages {
		d = append(d, p.PageData())
	}
	return
}

// Posts returns true if the collection is the special "posts" collection.
func (c *Collection) Posts() bool {
	return c.Name == "posts"
}

// SourceDir returns the source directory for pages in the collection.
func (c *Collection) SourceDir() string {
	return filepath.Join(siteConfig.SourceDir, "_"+c.Name)
}

// ReadPages scans the file system for collection pages, and adds them to c.Pages.
func (c *Collection) ReadPages(fileMap map[string]*Page) error {
	basePath := siteConfig.SourceDir
	d := map[interface{}]interface{}{
		"site":       siteData,
		"collection": c.Name,
	}
	d = mergeMaps(c.Data, d)

	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// if the issue is simply that the directory doesn't exist, ignore the error
			if os.IsNotExist(err) {
				if !c.Posts() {
					fmt.Println("Missing directory for collection", c.Name)
				}
				return nil
			}
			return err
		}
		rel, err := filepath.Rel(basePath, path)
		switch {
		case err != nil:
			return err
		case info.IsDir():
			return nil
		}
		p, err := ReadPage(rel, d)
		switch {
		case err != nil:
			return err
		case p.Static:
			fmt.Printf("skipping static file inside collection: %s\n", path)
		case p.Published:
			fileMap[p.Permalink] = p
			c.Pages = append(c.Pages, p)
		}
		return nil
	}
	return filepath.Walk(c.SourceDir(), walkFn)
}
