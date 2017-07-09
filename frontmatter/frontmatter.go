package frontmatter

import (
	"reflect"
	"sort"
	"strings"

	"github.com/osteele/gojekyll/utils"
	"github.com/osteele/liquid/evaluator"
)

// FrontMatter wraps a map to provide interface functions
type FrontMatter map[string]interface{}

// The first four bytes of a file with front matter.
const fmMagic = "---\n"

// FileHasFrontMatter returns a bool indicating whether the
// file looks like it has frontmatter.
func FileHasFrontMatter(filename string) (bool, error) {
	magic, err := utils.ReadFileMagic(filename)
	if err != nil {
		return false, err
	}
	return string(magic) == fmMagic, nil
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
