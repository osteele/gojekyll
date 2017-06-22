package gojekyll

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/osteele/gojekyll/pages"
	"github.com/osteele/gojekyll/templates"
)

// Collection is a Jekyll collection.
type Collection struct {
	Site  *Site
	Name  string
	Data  templates.VariableMap
	pages []pages.Page
}

// NewCollection creates a new Collection with defaults d
func NewCollection(s *Site, name string, d templates.VariableMap) *Collection {
	return &Collection{
		Site: s,
		Name: name,
		Data: d,
	}
}

// IsPosts returns true if the collection is the special "posts" collection.
func (c *Collection) IsPosts() bool { return c.Name == "posts" }

// Output returns a bool indicating whether files in this collection should be written.
func (c *Collection) Output() bool { return c.Data.Bool("output", false) }

// PathPrefix returns the collection's directory prefix, e.g. "_posts/"
func (c *Collection) PathPrefix() string { return filepath.FromSlash("_" + c.Name + "/") }

// Source returns the source directory for pages in the collection.
func (c *Collection) Source() string { return filepath.Join(c.Site.Source, "_"+c.Name) }

// TemplateVariable returns an array of page objects, for use as the template variable
// value of the collection.
func (c *Collection) TemplateVariable() (d []templates.VariableMap) {
	for _, page := range c.Pages() {
		d = append(d, page.PageVariables())
	}
	return
}

// ReadPages scans the file system for collection pages, and adds them to c.Pages.
func (c *Collection) ReadPages() error {
	collectionDefaults := templates.MergeVariableMaps(c.Data, templates.VariableMap{
		"collection": c.Name,
	})

	walkFn := func(filename string, info os.FileInfo, err error) error {
		if err != nil {
			// if the issue is simply that the directory doesn't exist, warn instead of error
			if os.IsNotExist(err) {
				if !c.IsPosts() {
					fmt.Printf("Missing collection directory: _%s\n", c.Name)
				}
				return nil
			}
			return err
		}
		relname, err := filepath.Rel(c.Site.Source, filename)
		switch {
		case strings.HasPrefix(filepath.Base(relname), "."):
			return nil
		case err != nil:
			return err
		case info.IsDir():
			return nil
		}
		defaults := templates.MergeVariableMaps(c.Site.GetFrontMatterDefaults(relname, ""), collectionDefaults)
		p, err := pages.NewPageFromFile(c.Site, c, filename, relname, defaults)
		switch {
		case err != nil:
			return err
		case p.Static():
			fmt.Printf("skipping static file inside collection: %s\n", filename)
		case p.Published():
			c.pages = append(c.pages, p)
		}
		return nil
	}
	return filepath.Walk(c.Source(), walkFn)
}
