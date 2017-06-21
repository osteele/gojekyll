package helpers

import (
	"path/filepath"
)

// TrimExt returns a path without its extension, if any
func TrimExt(name string) string {
	return name[:len(name)-len(filepath.Ext(name))]
}
