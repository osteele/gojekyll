package site

import (
	"sort"

	"github.com/osteele/gojekyll/collection"
)

// Render renders the site's pages.
func (s *Site) Render() error {
	cols := make([]*collection.Collection, 0, len(s.Collections))
	copy(cols, s.Collections)
	sort.Sort(postsCollectionLast(cols))
	for _, c := range cols {
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
		err = s.Render()
		if err != nil {
			return
		}
	})
	return
}

type postsCollectionLast []*collection.Collection

func (d postsCollectionLast) Len() int {
	return len([]*collection.Collection(d))
}

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

func (d postsCollectionLast) Swap(i, j int) {
	array := []*collection.Collection(d)
	a, b := array[i], array[j]
	array[i], array[j] = b, a
}
