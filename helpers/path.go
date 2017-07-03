package helpers

import (
	"path"
	"path/filepath"
	"strings"
	"time"
)

// FilenameDate returns the date for a filename that uses Jekyll post convention.
// It also returns a bool indicating whether a date was found.
func FilenameDate(s string) (time.Time, bool) {
	layout := "2006-01-02-"
	t, err := time.Parse(layout, filepath.Base(s + layout)[:len(layout)])
	return t, err == nil
}

// MustRel is like filepath.Rel, but panics if the path cannot be relativized.
func MustRel(basepath, targpath string) string {
	relpath, err := filepath.Rel(basepath, targpath)
	if err != nil {
		panic(err)
	}
	return relpath
}

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
