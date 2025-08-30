package pages

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/osteele/gojekyll/utils"
)

// DefaultPermalinkPattern is the default permalink pattern for pages that aren't in a collection
const DefaultPermalinkPattern = "/:path:output_ext"

// PermalinkStyles defines built-in styles from https://jekyllrb.com/docs/permalinks/#builtinpermalinkstyles
var PermalinkStyles = map[string]string{
	"date":    "/:categories/:year/:month/:day/:title.html",
	"pretty":  "/:categories/:year/:month/:day/:title/",
	"ordinal": "/:categories/:year/:y_day/:title.html",
	"none":    "/:categories/:title.html",
}

// permalinkDateVariables maps Jekyll permalink template variable names
// to time.Format layout strings
var permalinkDateVariables = map[string]string{
	"month":      "01",
	"imonth":     "1",
	"day":        "02",
	"i_day":      "2",
	"hour":       "15",
	"minute":     "04",
	"second":     "05",
	"year":       "2006",
	"short_year": "06",
}

var templateVariableMatcher = regexp.MustCompile(`:\w+\b`)

// See https://jekyllrb.com/docs/permalinks/#template-variables
func (p *page) permalinkVariables() map[string]string {
	var (
		relpath = p.relPath
		root    = utils.TrimExt(relpath)
		name    = filepath.Base(root)
		slug    = p.fm.String("slug", utils.Slugify(name))
		// date      = p.fileModTime
		date = p.PostDate().In(time.Local)
	)
	vars := map[string]string{
		"categories": strings.Join(p.Categories(), "/"),
		"collection": p.fm.String("collection", ""),
		"name":       utils.Slugify(name),
		"path":       "/" + root, // TODO are we removing and then adding this?
		"slug":       slug,
		"title":      utils.Slugify(p.fm.String("title", name)),
		"y_day":      strconv.Itoa(p.modTime.YearDay()),
		// Undocumented but evident:
		"output_ext": p.OutputExt(),
	}
	for k, v := range permalinkDateVariables {
		vars[k] = date.Format(v)
	}
	return vars
}

func (p *page) computePermalink(vars map[string]string) (src string, err error) {
	// First check for permalink in front matter
	var pattern string
	if permalink, hasFrontMatterPermalink := p.fm["permalink"]; hasFrontMatterPermalink {
		pattern = fmt.Sprintf("%v", permalink)
	} else {
		// If no front matter permalink, check global config
		if globalPermalink := p.site.Config().Permalink; globalPermalink != "" {
			pattern = globalPermalink
		} else {
			pattern = DefaultPermalinkPattern
		}
	}

	// Apply built-in permalink styles
	if pat, found := PermalinkStyles[pattern]; found {
		pattern = pat
	}

	// Jekyll Compatibility: Remove date/category placeholders for non-posts
	// Posts use the full permalink pattern, while pages and other collections
	// ignore date and category placeholders. This distinction is required for
	// Jekyll compatibility. See docs/PERMALINKS.md for detailed explanation.
	if !p.IsPost() {
		pattern = removePostOnlyPlaceholders(pattern)
	}

	templateVariables := p.permalinkVariables()
	s, err := utils.SafeReplaceAllStringFunc(templateVariableMatcher, pattern, func(m string) (string, error) {
		varname := m[1:]
		value, found := templateVariables[varname]
		if !found {
			return "", fmt.Errorf("unknown variable %q in permalink template %q", varname, pattern)
		}
		return value, nil
	})
	if err != nil {
		return "", err
	}
	return utils.URLPathClean("/" + s), nil
}

// removePostOnlyPlaceholders removes date and category placeholders from permalink patterns
// for non-post documents (pages and non-post collections).
// This matches Jekyll's behavior where these placeholders are ignored for non-posts.
func removePostOnlyPlaceholders(pattern string) string {
	// Remove category placeholder
	pattern = strings.ReplaceAll(pattern, ":categories/", "")
	pattern = strings.ReplaceAll(pattern, "/:categories", "")

	// Remove date-related placeholders
	datePatterns := []string{
		":year/", ":month/", ":day/", ":hour/", ":minute/", ":second/",
		":i_month/", ":i_day/", ":short_year/", ":y_day/",
		"/:year", "/:month", "/:day", "/:hour", "/:minute", "/:second",
		"/:i_month", "/:i_day", "/:short_year", "/:y_day",
	}
	for _, dp := range datePatterns {
		pattern = strings.ReplaceAll(pattern, dp, "")
	}

	// Clean up any double slashes that might result
	for strings.Contains(pattern, "//") {
		pattern = strings.ReplaceAll(pattern, "//", "/")
	}

	// Ensure pattern starts with /
	if !strings.HasPrefix(pattern, "/") {
		pattern = "/" + pattern
	}

	return pattern
}

func (p *page) setPermalink() (err error) {
	p.permalink, err = p.computePermalink(p.permalinkVariables())
	return
}
