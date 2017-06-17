package helpers

import "path/filepath"

// PathWithoutExtension returns a path without its extension, if any
func PathWithoutExtension(name string) string {
	return name[:len(name)-len(filepath.Ext(name))]
}
