package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// Collection is a Jekyll collection.
type Collection struct {
	Site   *Site
	Name   string
	Data   VariableMap
	Output bool
	Pages  []Page
}

func makeCollection(s *Site, name string, d VariableMap) *Collection {
	return &Collection{
		Site:   s,
		Name:   name,
		Data:   d,
		Output: d.Bool("output", false),
	}
}

// PageArrayVariableValue returns an array of a page data, for use as the template variable
// value of the collection.
func (c *Collection) PageArrayVariableValue() (d []VariableMap) {
	for _, p := range c.Pages {
		d = append(d, p.PageVariables())
	}
	return
}

// Posts returns true if the collection is the special "posts" collection.
func (c *Collection) IsPosts() bool {
	return c.Name == "posts"
}

// Source returns the source directory for pages in the collection.
func (c *Collection) Source() string {
	return filepath.Join(c.Site.Source, "_"+c.Name)
}

// ReadPages scans the file system for collection pages, and adds them to c.Pages.
func (c *Collection) ReadPages() error {
	basePath := c.Site.Source
	defaults := mergeVariableMaps(c.Data, VariableMap{
		"collection": c.Name,
	})

	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// if the issue is simply that the directory doesn't exist, warn instead of error
			if os.IsNotExist(err) {
				if !c.IsPosts() {
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
		p, err := ReadPage(rel, defaults)
		switch {
		case err != nil:
			return err
		case p.Static():
			fmt.Printf("skipping static file inside collection: %s\n", path)
		case p.Published():
			c.Site.Paths[p.Permalink()] = p
			c.Pages = append(c.Pages, p)
		}
		return nil
	}
	return filepath.Walk(c.Source(), walkFn)
}
