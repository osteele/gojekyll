package collections

import (
	"path/filepath"

	"github.com/osteele/gojekyll/config"
	"github.com/osteele/gojekyll/constants"
	"github.com/osteele/gojekyll/pages"
	"github.com/osteele/gojekyll/templates"
	"github.com/osteele/liquid/generics"
)

// Collection is a Jekyll collection https://jekyllrb.com/docs/collections/.
type Collection struct {
	Name      string
	Metadata  map[string]interface{}
	container pages.Container
	pages     []pages.Page
}

// NewCollection creates a new Collection
func NewCollection(name string, metadata map[string]interface{}, c pages.Container) *Collection {
	return &Collection{
		Name:      name,
		Metadata:  metadata,
		container: c,
	}
}

// Config is in the page.Container interface.
func (c *Collection) Config() config.Config {
	return c.container.Config()
}

// OutputExt is in the page.Container interface.
func (c *Collection) OutputExt(pathname string) string {
	return c.container.OutputExt(pathname)
}

// AbsDir is in the page.Container interface.
func (c *Collection) AbsDir() string {
	return filepath.Join(c.container.AbsDir(), c.PathPrefix())
}

// PathPrefix is in the page.Container interface.
// PathPrefix returns the collection's directory prefix, e.g. "_posts/"
func (c *Collection) PathPrefix() string { return filepath.FromSlash("_" + c.Name + "/") }

// IsPostsCollection returns true if the collection is the special "posts" collection.
func (c *Collection) IsPostsCollection() bool { return c.Name == "posts" }

// Output returns a bool indicating whether files in this collection should be written.
func (c *Collection) Output() bool { return templates.VariableMap(c.Metadata).Bool("output", false) }

// Pages is a list of pages.
func (c *Collection) Pages() []pages.Page {
	return c.pages
}

// TemplateVariable returns an array of page objects, for use as the template variable
// value of the collection.
func (c *Collection) TemplateVariable(ctx pages.RenderingContext, includeContent bool) ([]interface{}, error) {
	pages := []interface{}{}
	for _, p := range c.Pages() {
		v := p.PageVariables()
		if includeContent {
			c, err := p.Content(ctx)
			if err != nil {
				return nil, err
			}
			v = templates.MergeVariableMaps(v, map[string]interface{}{
				"content": string(c),
			})
		}
		pages = append(pages, v)
	}
	if c.IsPostsCollection() {
		generics.SortByProperty(pages, "date", true)
		reversed := make([]interface{}, len(pages))
		for i, v := range pages {
			reversed[len(pages)-1-i] = v
		}
		pages = reversed
	}
	return pages, nil
}

// TemplateObject returns the value of the collection in the template
// "collections" array.
func (c *Collection) TemplateObject(pages interface{}) interface{} {
	return templates.MergeVariableMaps(
		c.Metadata,
		map[string]interface{}{
			"label":              c.Name,
			"docs":               pages,
			"files":              []string{},
			"relative_directory": c.PathPrefix(),
			"directory":          c.AbsDir(),
		})
}

// PermalinkPattern returns the permalink pattern for this collection.
func (c *Collection) PermalinkPattern() string {
	defaultPattern := constants.DefaultCollectionPermalinkPattern
	if c.IsPostsCollection() {
		defaultPattern = constants.DefaultPostsCollectionPermalinkPattern
	}
	return templates.VariableMap(c.Metadata).String("permalink", defaultPattern)
}
