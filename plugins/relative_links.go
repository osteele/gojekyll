package plugins

import (
	"path"
	"path/filepath"
	"strings"

	"github.com/osteele/gojekyll/utils"
)

type jekyllRelativeLinksPlugin struct {
	plugin
	site Site
}

func init() {
	register("jekyll-relative-links", &jekyllRelativeLinksPlugin{})
}

func (p *jekyllRelativeLinksPlugin) AfterInitSite(s Site) error {
	p.site = s
	return nil
}

func (p *jekyllRelativeLinksPlugin) PostRender(b []byte) ([]byte, error) {
	return utils.ProcessAnchorHrefs(b, func(href string) string {
		return p.processHref(href)
	}), nil
}

// processHref converts relative links to markdown files to their rendered URLs
func (p *jekyllRelativeLinksPlugin) processHref(href string) string {
	// Skip absolute URLs (http://, https://, //)
	if strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://") || strings.HasPrefix(href, "//") {
		return href
	}

	// Skip anchors and query strings without paths
	if strings.HasPrefix(href, "#") || strings.HasPrefix(href, "?") {
		return href
	}

	// Skip mail links
	if strings.HasPrefix(href, "mailto:") {
		return href
	}

	// Skip absolute paths that don't point to markdown files
	if strings.HasPrefix(href, "/") && !p.isMarkdownFile(href) {
		return href
	}

	// Check if this is a markdown file
	if !p.isMarkdownFile(href) {
		return href
	}

	// Split the href to handle fragments and query strings
	parts := strings.SplitN(href, "#", 2)
	fragment := ""
	if len(parts) > 1 {
		fragment = "#" + parts[1]
	}
	linkPath := parts[0]

	parts = strings.SplitN(linkPath, "?", 2)
	query := ""
	if len(parts) > 1 {
		query = "?" + parts[1]
	}
	linkPath = parts[0]

	// Clean the path
	linkPath = path.Clean(linkPath)

	// Try to find the page and get its URL
	url, found := p.site.FilenameURLPath(linkPath)
	if found {
		return url + query + fragment
	}

	// If not found as-is, try without leading slash
	if strings.HasPrefix(linkPath, "/") {
		url, found = p.site.FilenameURLPath(strings.TrimPrefix(linkPath, "/"))
		if found {
			return url + query + fragment
		}
	}

	// If the file isn't found, return the original href
	return href
}

// isMarkdownFile checks if a path points to a markdown file based on extension
func (p *jekyllRelativeLinksPlugin) isMarkdownFile(filePath string) bool {
	// Get the base path without query string or fragment
	basePath := filePath
	if idx := strings.IndexAny(filePath, "?#"); idx != -1 {
		basePath = filePath[:idx]
	}

	ext := strings.ToLower(filepath.Ext(basePath))
	if ext == "" {
		return false
	}

	// Remove the leading dot
	ext = strings.TrimPrefix(ext, ".")

	// Get configured markdown extensions
	markdownExts := p.site.Config().MarkdownExtensions()
	for _, mdExt := range markdownExts {
		// MarkdownExtensions() returns extensions with dots
		if "."+ext == strings.ToLower(mdExt) {
			return true
		}
	}

	return false
}
