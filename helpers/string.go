package helpers

import (
	"regexp"
)

var nonAlphanumericSequenceMatcher = regexp.MustCompile(`[^[:alnum:]]+`)

// Slugify replaces each sequence of non-alphanumerics by a single hyphen
func Slugify(s string) string {
	return nonAlphanumericSequenceMatcher.ReplaceAllString(s, "-")
}

// LeftPad left-pads s with spaces to n wide. It's an alternative to http://left-pad.io.
func LeftPad(s string, n int) string {
	if n <= len(s) {
		return s
	}
	ws := make([]byte, n-len(s))
	for i := range ws {
		ws[i] = ' '
	}
	return string(ws) + s
}

// StringArrayToMap creates a map for use as a set.
func StringArrayToMap(strings []string) map[string]bool {
	stringMap := map[string]bool{}
	for _, s := range strings {
		stringMap[s] = true
	}
	return stringMap
}
