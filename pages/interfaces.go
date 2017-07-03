package pages

import (
	"io"
	"time"

	"github.com/osteele/gojekyll/pipelines"
	"github.com/osteele/liquid"
	"gopkg.in/yaml.v2"
)

// Document is a Jekyll page or file.
type Document interface {
	liquid.Drop
	yaml.Marshaler

	// Paths
	SiteRelPath() string // relative to the site source directory
	Permalink() string   // relative URL path
	OutputExt() string

	// Output
	Published() bool
	Static() bool
	Write(RenderingContext, io.Writer) error

	Categories() []string
	Tags() []string

	// Document initialization
	setPermalink() error
}

// Page is a document with frontmatter.
type Page interface {
	Document
	// Content asks a page to compute its content.
	// This has the side effect of causing the content to subsequently appear in the drop.
	Content(rc RenderingContext) ([]byte, error)
	// PostDate returns the date computed from the filename or frontmatter.
	// It is an uncaught error to call this on a page that is not a Post.
	// TODO Should posts have their own interface?
	PostDate() time.Time
}

// RenderingContext provides context information for rendering.
type RenderingContext interface {
	RenderingPipeline() pipelines.PipelineInterface
	// Site is the value of the "site" template variable.
	Site() interface{} // used as a drop in the rendering context
}

// Container is the document container.
// It's either the Site or Collection that immediately contains the document.
type Container interface {
	OutputExt(pathname string) string
	PathPrefix() string // PathPrefix is the relative prefix, "" for the site and "_coll/" for a collection
}
