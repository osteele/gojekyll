package liquid

import "io"

// Engine is a configured liquid engine.
type Engine interface {
	Parse([]byte) (Template, error)
	ParseAndRender([]byte, map[string]interface{}) ([]byte, error)
}

// Template is a liquid template.
type Template interface {
	Render(map[string]interface{}) ([]byte, error)
}

// LocalEngine runs in the same process, and can resolve tag arguments via local handlers.
type LocalEngine interface {
	Engine
	IncludeHandler(IncludeTagHandler)
	LinkTagHandler(LinkTagHandler)
}

// RemoteEngine runs out-of-process. It needs static tables for tag argument resolution.
type RemoteEngine interface {
	Engine
	FileURLMap(map[string]string) error
	IncludeDirs([]string) error
}

// IncludeTagHandler resolves the filename in a Liquid include tag into the expanded content
// of the included file.
type IncludeTagHandler func(string, io.Writer, map[string]interface{}) error
