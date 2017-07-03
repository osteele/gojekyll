package helpers

import (
	"regexp"
)

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

type replaceStringFuncError error

// SafeReplaceAllStringFunc is like regexp.ReplaceAllStringFunc but passes an
// an error back from the replacement function.
func SafeReplaceAllStringFunc(re *regexp.Regexp, src string, repl func(m string) (string, error)) (out string, err error) {
	// The ReplaceAllStringFunc callback signals errors via panic.
	// Turn them into return values.
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(replaceStringFuncError); ok {
				err = e.(error)
			} else {
				panic(r)
			}
		}
	}()
	return re.ReplaceAllStringFunc(src, func(m string) string {
		out, err := repl(m)
		if err != nil {
			panic(replaceStringFuncError(err))
		}
		return out
	}), nil
}

var nonAlphanumericSequenceMatcher = regexp.MustCompile(`[^[:alnum:]]+`)

// Slugify replaces each sequence of non-alphanumerics by a single hyphen
func Slugify(s string) string {
	return nonAlphanumericSequenceMatcher.ReplaceAllString(s, "-")
}

// StringArrayToMap creates a map for use as a set.
func StringArrayToMap(strings []string) map[string]bool {
	stringMap := map[string]bool{}
	for _, s := range strings {
		stringMap[s] = true
	}
	return stringMap
}
