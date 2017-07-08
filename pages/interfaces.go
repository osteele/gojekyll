package pages

import (
	"io"

	"github.com/osteele/gojekyll/config"
	"github.com/osteele/gojekyll/pipelines"
)

// A Document is a Jekyll post, page, or file.
type Document interface {
	// Paths
	Permalink() string  // relative URL path
	SourcePath() string // relative to the site source directory
	OutputExt() string

	// Output
	Published() bool
	Static() bool
	Write(io.Writer, RenderingContext) error

	Categories() []string
	Tags() []string
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
	Config() *config.Config
	OutputExt(pathname string) string
	PathPrefix() string // PathPrefix is the relative prefix, "" for the site and "_coll/" for a collection
}
