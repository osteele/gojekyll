package frontmatter

import (
	"bytes"
	"regexp"

	"github.com/osteele/gojekyll/utils"

	yaml "gopkg.in/yaml.v2"
)

// The first four bytes of a file with front matter.
const fmMagic = "---\n"

var frontMatterMatcher = regexp.MustCompile(`(?s)^---\n(.+?\n)---\n+`)
var emptyFontMatterMatcher = regexp.MustCompile(`(?s)^---\n+---\n+`)

// FileHasFrontMatter returns a bool indicating whether the
// file looks like it has frontmatter.
func FileHasFrontMatter(filename string) (bool, error) {
	magic, err := utils.ReadFileMagic(filename)
	if err != nil {
		return false, err
	}
	return string(magic) == fmMagic, nil
}

// Read reads the frontmatter from a document. It modifies srcPtr to point to the
// content after the frontmatter, and sets firstLine to its 1-indexed line number.
func Read(sourcePtr *[]byte, firstLine *int) (fm FrontMatter, err error) {
	var (
		source = *sourcePtr
		start  = 0
	)
	// Replace Windows line feeds. This allows the following regular expressions to work.
	source = bytes.ReplaceAll(source, []byte("\r\n"), []byte("\n"))
	if match := frontMatterMatcher.FindSubmatchIndex(source); match != nil {
		start = match[1]
		if err = yaml.Unmarshal(source[match[2]:match[3]], &fm); err != nil {
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
