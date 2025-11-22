package plugins

import (
	"path"
	"path/filepath"
	"strings"
)

type jekyllReadmeIndexPlugin struct{ plugin }

func init() {
	register("jekyll-readme-index", jekyllReadmeIndexPlugin{})
}

func (p jekyllReadmeIndexPlugin) PostInitPage(s Site, page Page) error {
	if isReadmePage(page) {
		// Calculate the new URL for this README page
		oldURL := page.URL()
		newURL := calculateIndexURL(oldURL)

		// Set the permalink in frontmatter to change the URL
		if newURL != oldURL {
			fm := page.FrontMatter()
			fm["permalink"] = newURL
		}
	}
	return nil
}

// isReadmePage checks if a page is a README file
func isReadmePage(page Page) bool {
	source := page.Source()
	if source == "" {
		return false
	}
	basename := filepath.Base(source)
	// Check for README.md, README.markdown, README.mdown, etc.
	nameWithoutExt := strings.TrimSuffix(basename, filepath.Ext(basename))
	return strings.EqualFold(nameWithoutExt, "README")
}

// calculateIndexURL converts a README URL to its index URL
// e.g., "/README.html" -> "/"
//
//	"/foo/README.html" -> "/foo/"
//	"/foo/bar/README.html" -> "/foo/bar/"
func calculateIndexURL(url string) string {
	dir := path.Dir(url)
	if dir == "." || dir == "" {
		return "/"
	}
	// Ensure the directory path ends with a slash
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}
	return dir
}
