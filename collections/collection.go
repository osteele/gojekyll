package collections

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

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
	pages     []pages.Document
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

// PathPrefix is in the page.Container interface.
// PathPrefix returns the collection's directory prefix, e.g. "_posts/"
func (c *Collection) PathPrefix() string { return filepath.FromSlash("_" + c.Name + "/") }

// IsPostsCollection returns true if the collection is the special "posts" collection.
func (c *Collection) IsPostsCollection() bool { return c.Name == "posts" }

// Output returns a bool indicating whether files in this collection should be written.
func (c *Collection) Output() bool { return templates.VariableMap(c.Metadata).Bool("output", false) }

// Pages is a list of pages.
func (c *Collection) Pages() []pages.Document {
	return c.pages
}

// TemplateVariable returns an array of page objects, for use as the template variable
// value of the collection.
func (c *Collection) TemplateVariable(ctx pages.RenderingContext, includeContent bool) ([]interface{}, error) {
	d := []interface{}{}
	for _, p := range c.Pages() {
		v := p.PageVariables()
		dp, ok := p.(*pages.Page)
		if includeContent && ok {
			c, err := dp.Content(ctx)
			if err != nil {
				return nil, err
			}
			v = templates.MergeVariableMaps(v, map[string]interface{}{
				"content": string(c),
			})
		}
		d = append(d, v)
	}
	if c.IsPostsCollection() {
		generics.SortByProperty(d, "date", true)
	}
	out := make([]interface{}, len(d))
	for i, v := range d {
		out[len(d)-1-i] = v
	}
	return out, nil
}

// PermalinkPattern returns the permalink pattern for this collection.
func (c *Collection) PermalinkPattern() string {
	defaultPattern := constants.DefaultCollectionPermalinkPattern
	if c.IsPostsCollection() {
		defaultPattern = constants.DefaultPostsCollectionPermalinkPattern
	}
	return templates.VariableMap(c.Metadata).String("permalink", defaultPattern)
}

// ReadPages scans the file system for collection pages, and adds them to c.Pages.
func (c *Collection) ReadPages(sitePath string, frontMatterDefaults func(string, string) map[string]interface{}) error {
	buildTime := time.Now()
	pageDefaults := map[string]interface{}{
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
		defaultFrontMatter := templates.MergeVariableMaps(pageDefaults, frontMatterDefaults(relname, c.Name))
		if c.IsPostsCollection() {
			t, ok := DateFromFilename(relname)
			if !ok {
				return nil
			}
			if t.After(buildTime) && !c.Config().Future {
				return nil
			}
			defaultFrontMatter["date"] = t
		}
		p, err := pages.NewFile(filename, c, filepath.ToSlash(relname), defaultFrontMatter)
		switch {
		case err != nil:
			return err
		case p.Published() || c.Config().Unpublished:
			c.pages = append(c.pages, p)
		}
		return nil
	}
	return filepath.Walk(filepath.Join(sitePath, c.PathPrefix()), walkFn)
}

// DateFromFilename returns the date for a filename that uses Jekyll post convention.
// It also returns a bool indicating whether a date was found.
func DateFromFilename(s string) (time.Time, bool) {
	layout := "2006-01-02-"
	t, err := time.Parse(layout, filepath.Base(s + layout)[:len(layout)])
	return t, err == nil
}
