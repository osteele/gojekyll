package constants

// DefaultPermalinkPattern is the default permalink pattern for pages that aren't in a collection
const DefaultPermalinkPattern = "/:path:output_ext"

// DefaultCollectionPermalinkPattern is the default permalink pattern for pages in the posts collection
const DefaultCollectionPermalinkPattern = "/:collection/:path:output_ext"

// DefaultPostPermalinkPattern is the default collection permalink pattern
const DefaultPostsCollectionPermalinkPattern = "/:categories/:year/:month/:day/:title.html"
