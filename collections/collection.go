package collections

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/osteele/gojekyll/constants"
	"github.com/osteele/gojekyll/pages"
	"github.com/osteele/gojekyll/templates"
)

// Collection is a Jekyll collection https://jekyllrb.com/docs/collections/.
type Collection struct {
	Name     string
	Metadata templates.VariableMap
	pages    []pages.Page
}

// NewCollection creates a new Collection
func NewCollection(ctx pages.RenderingContext, name string, metadata templates.VariableMap) *Collection {
	return &Collection{
		Name:     name,
		Metadata: metadata,
	}
}

// IsPostsCollection returns true if the collection is the special "posts" collection.
func (c *Collection) IsPostsCollection() bool { return c.Name == "posts" }

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
func (c *Collection) TemplateVariable(ctx pages.RenderingContext, includeContent bool) ([]templates.VariableMap, error) {
	d := []templates.VariableMap{}
	for _, p := range c.Pages() {
		v := p.PageVariables()
		dp, ok := p.(*pages.DynamicPage)
		if includeContent && ok {
			c, err := dp.ComputeContent(ctx)
			if err != nil {
				return nil, err
			}
			v = templates.MergeVariableMaps(v, templates.VariableMap{
				"content": string(c),
			})
		}
		d = append(d, v)
	}
	return d, nil
}

// PermalinkPattern returns the permalink pattern for this collection.
func (c *Collection) PermalinkPattern() string {
	defaultPattern := constants.DefaultCollectionPermalinkPattern
	if c.IsPostsCollection() {
		defaultPattern = constants.DefaultPostsCollectionPermalinkPattern
	}
	return c.Metadata.String("permalink", defaultPattern)
}

// ReadPages scans the file system for collection pages, and adds them to c.Pages.
func (c *Collection) ReadPages(ctx pages.RenderingContext, sitePath string, frontMatterDefaults func(string, string) templates.VariableMap) error {
	pageDefaults := templates.VariableMap{
		"collection": c.Name,
		"permalink":  c.PermalinkPattern(),
	}

	walkFn := func(filename string, info os.FileInfo, err error) error {
		if err != nil {
			// if the issue is simply that the directory doesn't exist, warn instead of error
			if os.IsNotExist(err) {
				if !c.IsPostsCollection() {
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
		defaultFrontmatter := templates.MergeVariableMaps(pageDefaults, frontMatterDefaults(relname, c.Name))
		p, err := pages.NewPageFromFile(ctx, c, filename, filepath.ToSlash(relname), defaultFrontmatter)
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
