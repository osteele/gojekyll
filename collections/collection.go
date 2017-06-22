package collections

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/osteele/gojekyll/pages"
	"github.com/osteele/gojekyll/templates"
)

// Collection is a Jekyll collection https://jekyllrb.com/docs/collections/.
type Collection struct {
	Name     string
	Metadata templates.VariableMap
	// context pages.Context
	pages []pages.Page
}

// NewCollection creates a new Collection with metadata m
func NewCollection(ctx pages.Context, name string, m templates.VariableMap) *Collection {
	return &Collection{
		Name:     name,
		Metadata: m,
	}
}

// DefaultPermalink returns the default Permalink for pages in a collection
// that doesn't specify a permalink in the site config.
func (c *Collection) DefaultPermalink() string { return "/:categories/:year/:month/:day/:title.html" }

// IsPosts returns true if the collection is the special "posts" collection.
func (c *Collection) IsPosts() bool { return c.Name == "posts" }

// Output returns a bool indicating whether files in this collection should be written.
func (c *Collection) Output() bool { return c.Metadata.Bool("output", false) }

// PathPrefix returns the collection's directory prefix, e.g. "_posts/"
func (c *Collection) PathPrefix() string { return filepath.FromSlash("_" + c.Name + "/") }

// Pages is a list of pages.
func (c *Collection) Pages() []pages.Page {
	return c.pages
}

// TemplateVariable returns an array of page objects, for use as the template variable
// value of the collection.
func (c *Collection) TemplateVariable() []templates.VariableMap {
	d := []templates.VariableMap{}
	for _, p := range c.Pages() {
		d = append(d, p.PageVariables())
	}
	return d
}

// ReadPages scans the file system for collection pages, and adds them to c.Pages.
func (c *Collection) ReadPages(ctx pages.Context, sitePath string, frontMatterDefaults func(string, string) templates.VariableMap) error {
	pageDefaults := templates.VariableMap{
		"collection": c.Name,
	}

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
		relname, err := filepath.Rel(sitePath, filename)
		switch {
		case strings.HasPrefix(filepath.Base(relname), "."):
			return nil
		case err != nil:
			return err
		case info.IsDir():
			return nil
		}
		defaults := templates.MergeVariableMaps(pageDefaults, frontMatterDefaults(relname, ""))
		p, err := pages.NewPageFromFile(ctx, c, filename, relname, defaults)
		switch {
		case err != nil:
			return err
		case p.Published():
			c.pages = append(c.pages, p)
		}
		return nil
	}
	return filepath.Walk(filepath.Join(sitePath, c.PathPrefix()), walkFn)
}
