package collections

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/osteele/gojekyll/helpers"
	"github.com/osteele/gojekyll/pages"
	"github.com/osteele/gojekyll/templates"
)

const draftsPath = "_drafts"

// ReadPages scans the file system for collection pages, and adds them to c.Pages.
func (c *Collection) ReadPages() error {
	sitePath := c.Config().Source
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
		relname := helpers.MustRel(sitePath, filename)
		switch {
		case strings.HasPrefix(filepath.Base(relname), "."):
			return nil
		case err != nil:
			return err
		case info.IsDir():
			return nil
		}
		fm := templates.MergeVariableMaps(pageDefaults, c.Config().GetFrontMatterDefaults(c.Name, relname))
		return c.readFile(filename, relname, fm)
	}
	if c.IsPostsCollection() && c.Config().Drafts {
		if err := filepath.Walk(filepath.Join(sitePath, draftsPath), walkFn); err != nil {
			return err
		}
		sort.Sort(pagesByDate{c.pages})
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
	f, err := pages.NewFile(abs, c, filepath.ToSlash(rel), fm)
	switch {
	case err != nil:
		return err
	case f.Static():
		return nil
	case f.Published() || c.Config().Unpublished:
		p := f.(pages.Page)
		c.pages = append(c.pages, p)
	}
	return nil
}
