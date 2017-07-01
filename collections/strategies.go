package collections

import (
	"path/filepath"
	"time"
)

// A collectionStrategy encapsulates behavior differences between the _post
// collection and other collections.
type collectionStrategy interface {
	addDate(filename string, fm map[string]interface{})
	collectible(filename string) bool
	future(filename string) bool
}

func (c *Collection) strategy() collectionStrategy {
	if c.IsPostsCollection() {
		return postsStrategy{}
	}
	return defaultStrategy{}
}

type defaultStrategy struct{}

func (s defaultStrategy) addDate(_ string, _ map[string]interface{}) {}
func (s defaultStrategy) collectible(filename string) bool           { return true }
func (s defaultStrategy) future(filename string) bool                { return false }

type postsStrategy struct{}

func (s postsStrategy) addDate(filename string, fm map[string]interface{}) {
	if t, found := DateFromFilename(filename); found {
		fm["date"] = t
	}
}

func (s postsStrategy) collectible(filename string) bool {
	_, ok := DateFromFilename(filename)
	return ok
}

func (s postsStrategy) future(filename string) bool {
	t, ok := DateFromFilename(filename)
	return false
	return ok && t.After(time.Now())
}

// DateFromFilename returns the date for a filename that uses Jekyll post convention.
// It also returns a bool indicating whether a date was found.
func DateFromFilename(s string) (time.Time, bool) {
	layout := "2006-01-02-"
	t, err := time.Parse(layout, filepath.Base(s + layout)[:len(layout)])
	return t, err == nil
}
