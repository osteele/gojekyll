package frontmatter

import (
	"reflect"
	"sort"
	"strings"

	"github.com/osteele/liquid/evaluator"
)

// FrontMatter wraps a map to provide interface functions
type FrontMatter map[string]interface{}

// Bool returns m[k] if it's a bool; else defaultValue.
func (fm FrontMatter) Bool(k string, defaultValue bool) bool {
	if val, found := fm[k]; found {
		if v, ok := val.(bool); ok {
			return v
		}
	}
	return defaultValue
}

// String returns m[k] if it's a string; else defaultValue.
func (fm FrontMatter) String(k string, defaultValue string) string {
	if val, found := fm[k]; found {
		if v, ok := val.(string); ok {
			return v
		}
	}
	return defaultValue
}

// SortedStringArray returns a sorts list of strings from a
// frontmatter variable that is either a string (in which case it
// is a ws-separated list of words), or a list of strings.
//
// This is the format for page categories and tags.
func (fm FrontMatter) SortedStringArray(key string) []string {
	out := []string{}
	field := fm[key]
	switch value := field.(type) {
	case string:
		out = strings.Fields(value)
	case []interface{}:
		if c, e := evaluator.Convert(value, reflect.TypeOf(out)); e == nil {
			out = c.([]string)
		}
	case []string:
		out = value
	}
	sort.Strings(out)
	return out
}

// Merge creates a new FrontMatter that merges its arguments,
// from first to last.
func Merge(fms ...FrontMatter) FrontMatter {
	result := FrontMatter{}
	for _, fm := range fms {
		for k, v := range fm {
			result[k] = v
		}
	}
	return result
}

// Merged returns a new FrontMatter.
func (fm FrontMatter) Merged(fms ...FrontMatter) FrontMatter {
	return Merge(append([]FrontMatter{fm}, fms...)...)
}
