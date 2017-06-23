package helpers

import (
	"path"
	"strings"
)

// TrimExt returns a path without its extension, if any
func TrimExt(name string) string {
	return name[:len(name)-len(path.Ext(name))]
}

// URLPathClean removes internal // etc. Unlike path.Clean, it
// leaves the final "/" intact.
func URLPathClean(filepath string) string {
	finalSlash := ""
	if strings.HasSuffix(filepath, "/") && len(filepath) > 1 {
		finalSlash = "/"
	}
	return path.Clean(filepath) + finalSlash
}
