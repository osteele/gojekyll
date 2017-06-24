package config

import (
	"path/filepath"
	"strings"

	"github.com/osteele/gojekyll/helpers"
)

// IsMarkdown returns a boolean indicating whether the file is a Markdown file, according to the current project.
func (c *Config) IsMarkdown(name string) bool {
	ext := filepath.Ext(name)
	return c.markdownExtensions()[strings.TrimLeft(ext, ".")]
}

// IsSassPath returns a boolean indicating whether the file is a Sass (".sass" or ".scss") file.
func (c *Config) IsSassPath(name string) bool {
	return strings.HasSuffix(name, ".sass") || strings.HasSuffix(name, ".scss")
}

// MarkdownExtensions returns a set of markdown extensions, without the initial dots.
func (c *Config) markdownExtensions() map[string]bool {
	extns := strings.SplitN(c.MarkdownExt, `,`, -1)
	return helpers.StringArrayToMap(extns)
}

func (c *Config) OutputExt(pathname string) string {
	switch {
	case c.IsMarkdown(pathname):
		return ".html"
	case c.IsSassPath(pathname):
		return ".css"
	default:
		return filepath.Ext(pathname)
	}
}
