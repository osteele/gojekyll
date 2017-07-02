package pages

import (
	"io"

	"github.com/osteele/gojekyll/config"
	"github.com/osteele/gojekyll/pipelines"
)

// Document is a Jekyll page or file.
type Document interface {
	// Paths
	SiteRelPath() string // relative to the site source directory
	Permalink() string   // relative URL path
	OutputExt() string

	// Output
	Published() bool
	Static() bool
	Write(RenderingContext, io.Writer) error

	// Variables
	PageVariables() map[string]interface{}

	// Document initialization uses this.
	initPermalink() error
}

// RenderingContext provides context information for rendering.
type RenderingContext interface {
	RenderingPipeline() pipelines.PipelineInterface
	// SiteVariables is the value of the "site" template variable.
	SiteVariables() map[string]interface{} // value of the "site" template variable
}

// Container is the document container.
// It's either the Site or Collection that immediately contains the document.
type Container interface {
	AbsDir() string
	Config() config.Config
	OutputExt(pathname string) string
	PathPrefix() string // PathPrefix is the relative prefix, "" for the site and "_coll/" for a collection
}
