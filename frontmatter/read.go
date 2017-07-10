package frontmatter

import (
	"bytes"
	"regexp"

	"github.com/osteele/gojekyll/templates"

	yaml "gopkg.in/yaml.v2"
)

var (
	frontMatterMatcher     = regexp.MustCompile(`(?s)^---\n(.+?\n)---\n+`)
	emptyFontMatterMatcher = regexp.MustCompile(`(?s)^---\n+---\n+`)
)

// Read reads the frontmatter from a document. It modifies srcPtr to point to the
// content after the frontmatter, and sets firstLine to its 1-indexed line number.
func Read(sourcePtr *[]byte, firstLine *int) (frontMatter templates.VariableMap, err error) {
	var (
		source = *sourcePtr
		start  = 0
	)
	// Replace Windows line feeds. This allows the following regular expressions to work.
	source = bytes.Replace(source, []byte("\r\n"), []byte("\n"), -1)
	if match := frontMatterMatcher.FindSubmatchIndex(source); match != nil {
		start = match[1]
		if err = yaml.Unmarshal(source[match[2]:match[3]], &frontMatter); err != nil {
			return
		}
	} else if match := emptyFontMatterMatcher.FindSubmatchIndex(source); match != nil {
		start = match[1]
	}
	if firstLine != nil {
		*firstLine = 1 + bytes.Count(source[:start], []byte("\n"))
	}
	*sourcePtr = source[start:]
	return
}
