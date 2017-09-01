package utils

import (
	"regexp"
	"strings"
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
	return strings.ToLower(nonAlphanumericSequenceMatcher.ReplaceAllString(s, "-"))
}

// StringArrayToMap creates a map for use as a set.
func StringArrayToMap(a []string) map[string]bool {
	m := map[string]bool{}
	for _, s := range a {
		m[s] = true
	}
	return m
}

// StringArrayContains returns a bool indicating whether the array contains the string.
func StringArrayContains(a []string, s string) bool {
	for _, item := range a {
		if item == s {
			return true
		}
	}
	return false
}

// Titleize splits at ` `, capitalizes, and joins.
func Titleize(s string) string {
	a := strings.Split(s, "-")
	for i, s := range a {
		if len(s) > 0 {
			a[i] = strings.ToUpper(s[:1]) + s[1:]
		}
	}
	return strings.Join(a, " ")
}
