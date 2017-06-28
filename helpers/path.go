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
func URLPathClean(url string) string {
	finalSlash := false
	if strings.HasSuffix(url, "/") && len(url) > 1 {
		finalSlash = true
	}
	cleaned := path.Clean(url)
	if finalSlash && !strings.HasSuffix(cleaned, "/") {
		cleaned += "/"
	}
	return cleaned
}
