package pages

import (
	"io"

	"github.com/osteele/gojekyll/templates"
)

// Page is a Jekyll page.
type Page interface {
	// Paths
	SiteRelPath() string // relative to the site source directory
	Permalink() string   // relative URL path
	OutputExt() string

	// Output
	Published() bool
	Static() bool
	Write(RenderingContext, io.Writer) error

	// Variables
	PageVariables() templates.VariableMap

	// internal
	initPermalink() error
}

// RenderingContext provides context information to a Page.
type RenderingContext interface {
	ApplyLayout(string, []byte, templates.VariableMap) ([]byte, error)
	OutputExt(pathname string) string
	Render(io.Writer, []byte, string, templates.VariableMap) ([]byte, error)
	SiteVariables() templates.VariableMap // value of the "site" template variable
}

// Container is the Page container
type Container interface {
	PathPrefix() string // PathPrefix is the relative prefix, "" for the site and "_coll/" for a collection
}
