package gojekyll

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Collection is a Jekyll collection.
type Collection struct {
	Site   *Site
	Name   string
	Data   VariableMap
	Output bool
	pages  []Page
}

// NewCollection creates a new Collection with defaults d
func NewCollection(s *Site, name string, d VariableMap) *Collection {
	return &Collection{
		Site:   s,
		Name:   name,
		Data:   d,
		Output: d.Bool("output", false),
	}
}

// CollectionValue returns an array of page objects, for use as the template variable
// value of the collection.
func (coll *Collection) CollectionValue() (d []VariableMap) {
	for _, page := range coll.Pages() {
		d = append(d, page.TemplateObject())
	}
	return
}

// IsPosts returns true if the collection is the special "posts" collection.
func (coll *Collection) IsPosts() bool {
	return coll.Name == "posts"
}

// PathPrefix returns the collection's directory prefix, e.g. "_posts/"
func (coll *Collection) PathPrefix() string {
	return filepath.FromSlash("_" + coll.Name + "/")
}

// Source returns the source directory for pages in the collection.
func (coll *Collection) Source() string {
	return filepath.Join(coll.Site.Source, "_"+coll.Name)
}

// ReadPages scans the file system for collection pages, and adds them to coll.Pages.
func (coll *Collection) ReadPages() error {
	basePath := coll.Site.Source
	collectionDefaults := MergeVariableMaps(coll.Data, VariableMap{
		"collection": coll.Name,
	})

	walkFn := func(name string, info os.FileInfo, err error) error {
		if err != nil {
			// if the issue is simply that the directory doesn't exist, warn instead of error
			if os.IsNotExist(err) {
				if !coll.IsPosts() {
					fmt.Printf("Missing collection directory: _%s\n", coll.Name)
				}
				return nil
			}
			return err
		}
		relname, err := filepath.Rel(basePath, name)
		switch {
		case strings.HasPrefix(filepath.Base(relname), "."):
			return nil
		case err != nil:
			return err
		case info.IsDir():
			return nil
		}
		defaults := MergeVariableMaps(coll.Site.GetFrontMatterDefaults(relname, ""), collectionDefaults)
		p, err := ReadPage(coll.Site, coll, relname, defaults)
		switch {
		case err != nil:
			return err
		case p.Static():
			fmt.Printf("skipping static file inside collection: %s\n", name)
		case p.Published():
			coll.Site.Paths[p.Permalink()] = p
			coll.pages = append(coll.pages, p)
		}
		return nil
	}
	return filepath.Walk(coll.Source(), walkFn)
}
