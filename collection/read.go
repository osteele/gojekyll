package collection

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/osteele/gojekyll/pages"
	"github.com/osteele/gojekyll/templates"
	"github.com/osteele/gojekyll/utils"
)

const draftsPath = "_drafts"

// ScanDirectory scans the file system for collection pages, and adds them to c.Pages.
func (c *Collection) ScanDirectory(dirname string) error {
	sitePath := c.config.Source
	pageDefaults := map[string]interface{}{
		"collection": c.Name,
		"permalink":  c.PermalinkPattern(),
	}
	walkFn := func(filename string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
		relname := utils.MustRel(sitePath, filename)
		switch {
		case strings.HasPrefix(filepath.Base(relname), "."):
			return nil
		case err != nil:
			return err
		case info.IsDir():
			return nil
		}
		fm := templates.MergeVariableMaps(pageDefaults, c.config.GetFrontMatterDefaults(c.Name, relname))
		return c.readFile(filename, relname, fm)
	}
	return filepath.Walk(filepath.Join(sitePath, dirname), walkFn)
}

// ReadPages scans the file system for collection pages, and adds them to c.Pages.
func (c *Collection) ReadPages() error {
	if c.IsPostsCollection() && c.config.Drafts {
		if err := c.ScanDirectory(draftsPath); err != nil {
			return err
		}
	}
	if err := c.ScanDirectory(c.PathPrefix()); err != nil {
		return err
	}
	if c.IsPostsCollection() {
		sort.Sort(pagesByDate{c.pages})
	}
	return nil
}

// readFile mutates fm.
func (c *Collection) readFile(abs string, rel string, fm map[string]interface{}) error {
	strategy := c.strategy()
	switch {
	case !strategy.collectible(rel):
		return nil
	case strategy.future(rel) && !c.config.Future:
		return nil
	default:
		strategy.addDate(rel, fm)
	}
	f, err := pages.NewFile(c.site, abs, filepath.ToSlash(rel), fm)
	switch {
	case err != nil:
		return err
	case f.Static():
		return nil
	case f.Published() || c.config.Unpublished:
		p := f.(pages.Page)
		c.pages = append(c.pages, p)
	}
	return nil
}
