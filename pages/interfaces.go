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
	Write(io.Writer) error

	Categories() []string
	Tags() []string
}

// Site is the interface that the site provides to a page.
type Site interface {
	Config() *config.Config
	RenderingPipeline() pipelines.PipelineInterface
	OutputExt(pathname string) string
}
