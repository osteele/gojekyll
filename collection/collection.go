package collection

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/osteele/gojekyll/config"
	"github.com/osteele/gojekyll/pages"
	"github.com/osteele/gojekyll/renderers"
	"github.com/osteele/gojekyll/templates"
	"github.com/osteele/gojekyll/utils"
	"github.com/osteele/liquid"
)

// Collection is a Jekyll collection https://jekyllrb.com/docs/collections/.
type Collection struct {
	Name     string
	Metadata map[string]interface{}

	cfg   *config.Config
	pages []Page
	site  Site
}

// Site is the interface a site provides to its collections.
type Site interface {
	Config() *config.Config
	Exclude(string) bool
	RelativePath(string) string
	RendererManager() renderers.Renderers
}

// Page is in package pages.
type Page = pages.Page

const draftsPath = "_drafts"
const postsName = "posts"

// New creates a new Collection
func New(s Site, name string, metadata map[string]interface{}) *Collection {
	return &Collection{
		Name:     name,
		Metadata: metadata,
		cfg:      s.Config(),
		site:     s,
	}
}

func (c *Collection) String() string {
	return fmt.Sprintf("%T{Name=%q}", c, c.Name)
}

// AbsDir returns the absolute path to the collection directory.
func (c *Collection) AbsDir() string {
	return filepath.Join(c.cfg.SourceDir(), c.PathPrefix())
}

// PathPrefix returns the collection's directory prefix, e.g. "_posts/"
func (c *Collection) PathPrefix() string { return filepath.FromSlash("_" + c.Name + "/") }

// IsPostsCollection returns true if the collection is the special "posts" collection.
func (c *Collection) IsPostsCollection() bool { return c.Name == postsName }

// Output returns a bool indicating whether files in this collection should be written.
func (c *Collection) Output() bool { return templates.VariableMap(c.Metadata).Bool("output", false) }

// Pages in the Post collection are ordered by date.
func (c *Collection) Pages() []Page {
	return c.pages
}

// Render renders the collection's pages.
func (c *Collection) Render() error {
	var errs []error
	for _, p := range c.Pages() {
		err := p.Render()
		if err != nil {
			errs = append(errs, err)
		}
	}
	return utils.CombineErrors(errs)
}

// ToLiquid returns the value of the collection in the template
// "collections" array.
func (c *Collection) ToLiquid() interface{} {
	return liquid.IterationKeyedMap(templates.MergeVariableMaps(
		c.Metadata,
		map[string]interface{}{
			"label":              c.Name,
			"docs":               c.pages,
			"files":              []string{},
			"relative_directory": strings.TrimSuffix(c.PathPrefix(), "/"),
			"directory":          c.AbsDir(),
		}))
}

// PermalinkPattern returns the default permalink pattern for this collection.
func (c *Collection) PermalinkPattern() string {
	pattern := c.strategy().defaultPermalinkPattern(c.cfg)
	return templates.VariableMap(c.Metadata).String("permalink", pattern)
}
