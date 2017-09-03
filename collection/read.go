package collection

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/osteele/gojekyll/pages"
	"github.com/osteele/gojekyll/utils"
)

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
		addPrevNext(c.pages)
	}
	return nil
}

func addPrevNext(ps []Page) {
	const prevPageField = "previous"
	const nextPageField = "next"
	var prev Page
	for _, p := range ps {
		p.FrontMatter()[prevPageField] = prev
		if prev != nil {
			prev.FrontMatter()[nextPageField] = p
		}
		prev = p
	}
	if prev != nil {
		prev.FrontMatter()[nextPageField] = nil
	}
}

// scanDirectory scans the file system for collection pages, and adds them to c.Pages.
//
// This function is distinct from ReadPages so that the posts collection can call it twice.
func (c *Collection) scanDirectory(dirname string) error {
	sitePath := c.cfg.Source
	dir := filepath.Join(sitePath, dirname)
	return filepath.Walk(dir, func(filename string, info os.FileInfo, err error) error {
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
	})
}

func (c *Collection) readPost(path string, rel string) error {
	siteRel := utils.MustRel(c.cfg.Source, path)
	strategy := c.strategy()
	switch {
	case !strategy.isCollectible(rel):
		return nil
	case strategy.isFuture(rel) && !c.cfg.Future:
		return nil
	}
	fm := pages.FrontMatter{
		"collection": c.Name,
		"permalink":  c.PermalinkPattern(),
	}.Merged(c.cfg.GetFrontMatterDefaults(c.Name, siteRel))
	strategy.parseFilename(rel, fm)
	f, err := pages.NewFile(c.site, path, filepath.ToSlash(rel), fm)
	switch {
	case err != nil:
		return err
	case f.IsStatic():
		return nil
	case f.Published() || c.cfg.Unpublished:
		p := f.(Page) // f.Static() guarantees this
		c.pages = append(c.pages, p)
	}
	return nil
}
