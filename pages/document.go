package pages

import (
	"io"

	"github.com/osteele/gojekyll/config"
	"github.com/osteele/gojekyll/renderers"
)

// A Document is a Jekyll post, page, or file.
type Document interface {
	// Paths
	URL() string // relative to site base
	Source() string
	OutputExt() string

	// Output
	Published() bool
	IsStatic() bool
	Write(io.Writer) error

	Reload() error
}

// Site is the interface that the site provides to a page.
type Site interface {
	Config() *config.Config
	RelativePath(string) string
	RendererManager() renderers.Renderers
}
