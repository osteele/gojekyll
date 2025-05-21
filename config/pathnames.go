package config

import (
	"path/filepath"
	"strings"

	"github.com/osteele/gojekyll/utils"
)

// IsMarkdown returns a boolean indicating whether the file is a Markdown file, according to the current project.
func (c *Config) IsMarkdown(name string) bool {
	ext := filepath.Ext(name)
	return c.markdownExtensions()[strings.TrimLeft(ext, ".")]
}

// IsSASSPath returns a boolean indicating whether the file is a Sass (".sass" or ".scss") file.
func (c *Config) IsSASSPath(name string) bool {
	return strings.HasSuffix(name, ".sass") || strings.HasSuffix(name, ".scss")
}

// markdownExtensions returns a set of markdown extensions, without the initial dots.
func (c *Config) markdownExtensions() map[string]bool {
	exts := strings.Split(c.MarkdownExt, `,`)
	return utils.StringArrayToMap(exts)
}

// MarkdownExtensions returns a list of markdown extensions, with dotsa.
func (c *Config) MarkdownExtensions() []string {
	exts := strings.Split(c.MarkdownExt, `,`)
	for i, k := range exts {
		exts[i] = "." + k
	}
	return exts
}

// OutputExt returns the pathname's output extension.
// This is generally the pathname extension;
// exceptions are *.md -> *.html, and *.{sass,scss} -> *.css.
func (c *Config) OutputExt(pathname string) string {
	switch {
	case c.IsMarkdown(pathname):
		return ".html"
	case c.IsSASSPath(pathname):
		return ".css"
	default:
		return filepath.Ext(pathname)
	}
}
