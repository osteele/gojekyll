package gojekyll

import (
	"path/filepath"
	"strings"

	"github.com/osteele/gojekyll/helpers"
)

// IsMarkdown returns a boolean indicating whether the file is a Markdown file, according to the current project.
func (s *Site) IsMarkdown(name string) bool {
	ext := filepath.Ext(name)
	return s.MarkdownExtensions()[strings.TrimLeft(ext, ".")]
}

// MarkdownExtensions returns a set of markdown extension, without the final dots.
func (s *Site) MarkdownExtensions() map[string]bool {
	extns := strings.SplitN(s.config.MarkdownExt, `,`, -1)
	return helpers.StringArrayToMap(extns)
}
