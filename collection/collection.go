package collection

import (
	"path/filepath"

	"github.com/osteele/gojekyll/config"
	"github.com/osteele/gojekyll/pages"
	"github.com/osteele/gojekyll/templates"
)

// Collection is a Jekyll collection https://jekyllrb.com/docs/collections/.
type Collection struct {
	Name     string
	Metadata map[string]interface{}

	config *config.Config
	pages  []pages.Page
	site   Site
}

// Site is the interface a site provides to collections it contains.
type Site interface {
	Config() *config.Config
	OutputExt(pathname string) string
}

// New creates a new Collection
func New(s Site, name string, metadata map[string]interface{}) *Collection {
	return &Collection{
		Name:     name,
		Metadata: metadata,
		config:   s.Config(),
		site:     s,
	}
}

// AbsDir is in the page.Container interface.
func (c *Collection) AbsDir() string {
	return filepath.Join(c.config.SourceDir(), c.PathPrefix())
}

// Config is in the page.Container interface.
func (c *Collection) Config() *config.Config {
	return c.config
}

// OutputExt is in the page.Container interface.
func (c *Collection) OutputExt(pathname string) string {
	return c.site.OutputExt(pathname)
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

// SetPageContent sets up the collection's pages' "content".
func (c *Collection) SetPageContent(ctx pages.RenderingContext) error {
	for _, p := range c.Pages() {
		_, err := p.Content(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

// ToLiquid returns the value of the collection in the template
// "collections" array.
func (c *Collection) ToLiquid() interface{} {
	return templates.MergeVariableMaps(
		c.Metadata,
		map[string]interface{}{
			"label":              c.Name,
			"docs":               c.pages,
			"files":              []string{},
			"relative_directory": c.PathPrefix(),
			"directory":          c.AbsDir(),
		})
}

// PermalinkPattern returns the default permalink pattern for this collection.
func (c *Collection) PermalinkPattern() string {
	defaultPattern := c.strategy().defaultPermalinkPattern()
	return templates.VariableMap(c.Metadata).String("permalink", defaultPattern)
}
