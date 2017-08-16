package site

import (
	"sort"

	"github.com/osteele/gojekyll/collection"
)

// render renders the site's pages.
func (s *Site) render() error {
	for _, c := range s.sortedCollections() {
		if err := c.Render(); err != nil {
			return err
		}
	}
	for _, c := range s.nonCollectionPages {
		if err := c.Render(); err != nil {
			return err
		}
	}
	return nil
}

func (s *Site) ensureRendered() (err error) {
	s.renderOnce.Do(func() {
		err = s.initializeRenderingPipeline()
		if err != nil {
			return
		}
		err = s.render()
		if err != nil {
			return
		}
	})
	return
}

// returns a slice of collections, sorted by name but with _posts last.
func (s *Site) sortedCollections() []*collection.Collection {
	cols := make([]*collection.Collection, len(s.Collections))
	copy(cols, s.Collections)
	sort.Slice(cols, postsCollectionLast(cols).Less)
	return cols
}

type postsCollectionLast []*collection.Collection

func (d postsCollectionLast) Less(i, j int) bool {
	array := []*collection.Collection(d)
	a, b := array[i], array[j]
	switch {
	case a.IsPostsCollection():
		return false
	case b.IsPostsCollection():
		return true
	default:
		return a.Name < b.Name
	}
}
