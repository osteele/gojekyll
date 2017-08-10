package collection

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/osteele/gojekyll/pages"
	"github.com/osteele/gojekyll/templates"
	"github.com/osteele/gojekyll/utils"
)

const draftsPath = "_drafts"

// ReadPages scans the file system for collection pages, and adds them to c.Pages.
func (c *Collection) ReadPages() error {
	if c.IsPostsCollection() && c.cfg.Drafts {
		if err := c.scanDirectory(draftsPath); err != nil {
			return err
		}
	}
	if err := c.scanDirectory(c.PathPrefix()); err != nil {
		return err
	}
	if c.IsPostsCollection() {
		sort.Sort(pagesByDate{c.pages})
		var prev pages.Page
		for _, p := range c.pages {
			p.FrontMatter()["previous"] = prev
			if prev != nil {
				prev.FrontMatter()["next"] = p
			}
			prev = p
		}
		if prev != nil {
			prev.FrontMatter()["next"] = nil
		}
	}
	return nil
}

// scanDirectory scans the file system for collection pages, and adds them to c.Pages.
//
// This function is distinct from ReadPages so that the posts collection can call it twice.
func (c *Collection) scanDirectory(dirname string) error {
	var (
		sitePath = c.cfg.Source
		dir      = filepath.Join(sitePath, dirname)
	)
	walkFn := func(filename string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
		siteRel := utils.MustRel(sitePath, filename)
		switch {
		case info.IsDir():
			return nil
		case c.site.Exclude(siteRel):
			return nil
		default:
			return c.readPost(filename, utils.MustRel(dir, filename))
		}
	}
	return filepath.Walk(dir, walkFn)
}

func (c *Collection) readPost(abs string, rel string) error {
	siteRel := utils.MustRel(c.cfg.Source, abs)
	strategy := c.strategy()
	switch {
	case !strategy.collectible(rel):
		return nil
	case strategy.future(rel) && !c.cfg.Future:
		return nil
	}
	pageDefaults := map[string]interface{}{
		"collection": c.Name,
		"permalink":  c.PermalinkPattern(),
	}
	fm := templates.MergeVariableMaps(pageDefaults, c.cfg.GetFrontMatterDefaults(c.Name, siteRel))
	strategy.addDate(rel, fm)
	f, err := pages.NewFile(c.site, abs, filepath.ToSlash(rel), fm)
	switch {
	case err != nil:
		return err
	case f.Static():
		return nil
	case f.Published() || c.cfg.Unpublished:
		p := f.(pages.Page) // f.Static() guarantees this
		c.pages = append(c.pages, p)
	}
	return nil
}
