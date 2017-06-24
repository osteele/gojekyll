package templates

import (
	"bytes"
	"regexp"

	yaml "gopkg.in/yaml.v2"
)

var (
	frontMatterMatcher     = regexp.MustCompile(`(?s)^---\n(.+?\n)---\n`)
	emptyFontMatterMatcher = regexp.MustCompile(`(?s)^---\n+---\n`)
)

// ReadFrontMatter reads the front matter from a document.
func ReadFrontMatter(sourcePtr *[]byte) (frontMatter VariableMap, err error) {
	var (
		source = *sourcePtr
		start  = 0
	)
	// Replace Windows linefeeds. This allows the following regular expressions to work.
	source = bytes.Replace(source, []byte("\r\n"), []byte("\n"), -1)
	if match := frontMatterMatcher.FindSubmatchIndex(source); match != nil {
		start = match[1]
		if err = yaml.Unmarshal(source[match[2]:match[3]], &frontMatter); err != nil {
			return
		}
	} else if match := emptyFontMatterMatcher.FindSubmatchIndex(source); match != nil {
		start = match[1]
	}
	// This fixes the line numbers, so that template errors show with the correct line.
	// TODO find a less hack-ey solution
	*sourcePtr = append(
		regexp.MustCompile(`[^\n\r]+`).ReplaceAllLiteral(source[:start], []byte{}),
		source[start:]...)
	return
}
