package utils

import (
	"path"
	"path/filepath"
	"strings"
	"time"
)

// ParseFilenameDateTitle returns the date for a filename that uses Jekyll post convention.
// It also returns a bool indicating whether a date was found.
func ParseFilenameDateTitle(s string) (t time.Time, title string, found bool) {
	var (
		base   = TrimExt(filepath.Base(s))
		layout = "2006-01-02-"
	)
	if len(base) < len(layout) {
		found = false
		return
	}
	t, err := time.ParseInLocation(layout, base[:len(layout)], time.Local)
	if err != nil {
		return
	}
	title = Titleize(base[len(layout):])
	found = true
	return
}

// MatchList implement Jekyll include: and exclude: configurations.
// It reports whether any glob pattern matches the path.
// It panics with ErrBadPattern if any pattern is malformed.
// To match Jekyll, a string "dir/" matches that begins with this prefix.
func MatchList(patterns []string, name string) bool {
	for _, p := range patterns {
		match, err := filepath.Match(p, name)
		if err != nil {
			panic(err)
		}
		if match {
			return true
		}
		if strings.HasSuffix(p, "/") && strings.HasPrefix(name, p) {
			return true
		}
	}
	return false
}

// MustAbs is like filepath.Abs, but panics instead of returning an error.
func MustAbs(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	return abs
}

// MustRel is like filepath.Rel, but panics if the path cannot be relativized.
func MustRel(basepath, targpath string) string {
	rel, err := filepath.Rel(basepath, targpath)
	if err != nil {
		panic(err)
	}
	return rel
}

// TrimExt returns a path without its extension, if any
func TrimExt(name string) string {
	return name[:len(name)-len(path.Ext(name))]
}

// URLPathClean removes internal // etc. Unlike path.Clean, it
// leaves the final "/" intact.
func URLPathClean(url string) string {
	finalSlash := strings.HasSuffix(url, "/") && len(url) > 1
	cleaned := path.Clean(url)
	if finalSlash && !strings.HasSuffix(cleaned, "/") {
		cleaned += "/"
	}
	return cleaned
}
