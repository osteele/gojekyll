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

type LocalEngine interface {
	Engine
	IncludeHandler(func(string, io.Writer, map[string]interface{}))
	LinkHandler(LinkHandler)
}

type RemoteEngine interface {
	Engine
	FileUrlMap(map[string]string)
	IncludeDirs([]string)
}
