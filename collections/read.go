package collections

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/osteele/gojekyll/pages"
	"github.com/osteele/gojekyll/templates"
)

// ReadPages scans the file system for collection pages, and adds them to c.Pages.
func (c *Collection) ReadPages(sitePath string, frontMatterDefaults func(string, string) map[string]interface{}) error {
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
		fm := templates.MergeVariableMaps(pageDefaults, frontMatterDefaults(relname, c.Name))
		return c.readFile(filename, relname, fm)
	}
	if c.IsPostsCollection() && c.Config().Drafts {
		if err := filepath.Walk(filepath.Join(sitePath, "_drafts"), walkFn); err != nil {
			return err
		}
	}
	return filepath.Walk(filepath.Join(sitePath, c.PathPrefix()), walkFn)
}

// readFile mutates fm.
func (c *Collection) readFile(abs string, rel string, fm map[string]interface{}) error {
	strategy := c.strategy()
	switch {
	case !strategy.collectible(rel):
		return nil
	case strategy.future(rel) && !c.Config().Future:
		return nil
	default:
		strategy.addDate(rel, fm)
	}
	p, err := pages.NewFile(abs, c, filepath.ToSlash(rel), fm)
	switch {
	case err != nil:
		return err
	case p.Published() || c.Config().Unpublished:
		c.pages = append(c.pages, p)
	}
	return nil
}
