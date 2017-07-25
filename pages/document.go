package pages

import (
	"io"

	"github.com/osteele/gojekyll/config"
	"github.com/osteele/gojekyll/pipelines"
)

// A Document is a Jekyll post, page, or file.
type Document interface {
	// Paths
	Permalink() string // relative URL path
	SourcePath() string
	OutputExt() string

	// Output
	Published() bool
	Static() bool
	Write(io.Writer) error

	Reload() error
}

// Site is the interface that the site provides to a page.
type Site interface {
	Config() *config.Config
	RelativePath(string) string
	RenderingPipeline() pipelines.PipelineInterface
	OutputExt(pathname string) string
}
