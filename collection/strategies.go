package collection

import (
	"time"

	"github.com/osteele/gojekyll/utils"
)

// A collectionStrategy encapsulates behavior differences between the _post
// collection and other collection.
type collectionStrategy interface {
	addDate(filename string, fm map[string]interface{})
	collectible(filename string) bool
	defaultPermalinkPattern() string
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
	if t, found := utils.FilenameDate(filename); found {
		fm["date"] = t
	}
}

func (s postsStrategy) collectible(filename string) bool {
	_, ok := utils.FilenameDate(filename)
	return ok
}

func (s postsStrategy) future(filename string) bool {
	t, ok := utils.FilenameDate(filename)
	return ok && t.After(time.Now())
}

// DefaultCollectionPermalinkPattern is the default permalink pattern for pages in the posts collection
const DefaultCollectionPermalinkPattern = "/:collection/:path:output_ext"

// DefaultPostsCollectionPermalinkPattern is the default collection permalink pattern
const DefaultPostsCollectionPermalinkPattern = "/:categories/:year/:month/:day/:title.html"

func (s defaultStrategy) defaultPermalinkPattern() string {
	return DefaultCollectionPermalinkPattern
}

func (s postsStrategy) defaultPermalinkPattern() string {
	return DefaultPostsCollectionPermalinkPattern
}
