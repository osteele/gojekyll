package pages

import (
	"io"

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

// RenderingContext provides context information for rendering.
type RenderingContext interface {
	RenderingPipeline() pipelines.PipelineInterface
	// SiteVariables is the value of the "site" template variable.
	SiteVariables() map[string]interface{} // value of the "site" template variable
}

// Container is the document container.
// It's either the Site or Collection that immediately contains the document.
type Container interface {
	OutputExt(pathname string) string
	PathPrefix() string // PathPrefix is the relative prefix, "" for the site and "_coll/" for a collection
}
