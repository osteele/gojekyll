package collection

import (
	"time"

	"github.com/osteele/gojekyll/config"
	"github.com/osteele/gojekyll/utils"
)

// A collectionStrategy encapsulates behavior differences between the `_post`
// collection and other collections.
type collectionStrategy interface {
	defaultPermalinkPattern(*config.Config) string
	isCollectible(filename string) bool
	isFuture(filename string) bool
	parseFilename(string, map[string]interface{})
}

func (c *Collection) strategy() collectionStrategy {
	if c.IsPostsCollection() {
		return postsStrategy{}
	}
	return defaultStrategy{}
}

type defaultStrategy struct{}

func (s defaultStrategy) parseFilename(_ string, fm map[string]interface{}) {
	// de facto
	fm["draft"] = false
}

func (s defaultStrategy) isCollectible(string) bool { return true }
func (s defaultStrategy) isFuture(string) bool      { return false }

type postsStrategy struct{}

func (s postsStrategy) parseFilename(filename string, fm map[string]interface{}) {
	if t, title, found := utils.ParseFilenameDateTitle(filename); found {
		fm["date"] = t
		fm["title"] = title
		fm["slug"] = utils.Slugify(title)
	}
}

func (s postsStrategy) isCollectible(filename string) bool {
	_, _, ok := utils.ParseFilenameDateTitle(filename)
	return ok
}

func (s postsStrategy) isFuture(filename string) bool {
	t, _, ok := utils.ParseFilenameDateTitle(filename)
	return ok && t.After(time.Now())
}

// DefaultCollectionPermalinkPattern is the default permalink pattern for pages in the posts collection
const DefaultCollectionPermalinkPattern = "/:collection/:path:output_ext"

// DefaultPostsCollectionPermalinkPattern is the default collection permalink pattern
const DefaultPostsCollectionPermalinkPattern = "/:categories/:year/:month/:day/:title.html"

func (s defaultStrategy) defaultPermalinkPattern(*config.Config) string {
	return DefaultCollectionPermalinkPattern
}

func (s postsStrategy) defaultPermalinkPattern(cfg *config.Config) string {
	if s, ok := cfg.String("permalink"); ok {
		return s
	}
	return DefaultPostsCollectionPermalinkPattern
}
