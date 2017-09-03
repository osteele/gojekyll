package frontmatter

import (
	"fmt"
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

// Get returns m[k] if present; else defaultValue.
func (fm FrontMatter) Get(k string, defaultValue interface{}) interface{} {
	if val, found := fm[k]; found {
		return val
	}
	return defaultValue
}

// String returns m[k] if it's a string or can be stringified; else defaultValue.
func (fm FrontMatter) String(k string, defaultValue string) string {
	if val, found := fm[k]; found {
		switch v := val.(type) {
		case string:
			return v
		case fmt.Stringer:
			return v.String()
		}
	}
	return defaultValue
}

// StringArray returns m[k] if it's a []string or string array
func (fm FrontMatter) StringArray(k string) []string {
	if value, ok := fm[k]; ok {
		switch value := value.(type) {
		case []string:
			return value
		case []interface{}:
			a := make([]string, len(value))
			for i, item := range value {
				a[i] = fmt.Sprintf("%s", item)
			}
			return a
		case string:
			return []string{value}
		}
	}
	return nil
}

// SortedStringArray returns a sorts list of strings from a
// frontmatter variable that is either a string (in which case it
// is a ws-separated list of words), or a list of strings.
//
// This is the format for page categories and tags.
func (fm FrontMatter) SortedStringArray(key string) []string {
	var result []string
	switch v := fm[key].(type) {
	case string:
		result = strings.Fields(v)
	case []interface{}:
		if c, e := evaluator.Convert(v, reflect.TypeOf(result)); e == nil {
			result = c.([]string)
		}
	case []string:
		result = v
	}
	sort.Strings(result)
	return result
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
