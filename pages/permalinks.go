package pages

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/osteele/gojekyll/logger"
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
		// date = p.PostDate().In(time.Local)
	)
	loc := time.Local
	// Check standard Jekyll timezone config first for compatibility
	tzName := p.site.Config().Timezone
	// Fall back to permalink_timezone if standard timezone is not set
	if tzName == "" {
		tzName = p.site.Config().PermalinkTimezone
	}
	if tzName != "" {
		l, err := time.LoadLocation(tzName)
		if err != nil {
			log := logger.Default()
			log.Warn("Could not load timezone %q for permalink: %s. Using local time zone instead.", tzName, err)
		} else {
			loc = l
		}
	}
	date := p.PostDate().In(loc)
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
			// For non-posts (pages and collections), only apply built-in permalink styles
			// Custom patterns should only affect posts
			if !p.IsPost() {
				if _, isBuiltInStyle := PermalinkStyles[globalPermalink]; !isBuiltInStyle {
					// Not a built-in style, use default pattern for pages
					pattern = DefaultPermalinkPattern
				} else {
					// Built-in style, apply it
					pattern = globalPermalink
				}
			} else {
				// Posts use the global permalink regardless of whether it's built-in or custom
				pattern = globalPermalink
			}
		} else {
			pattern = DefaultPermalinkPattern
		}
	}

	// Check if pattern is a built-in style
	isBuiltInStyle := false
	if pat, found := PermalinkStyles[pattern]; found {
		pattern = pat
		isBuiltInStyle = true
	}

	// Jekyll Compatibility: Custom patterns (non-built-in styles) should only
	// apply to posts when set globally. Pages and other collections should use
	// the default pattern when a custom pattern is configured globally.
	//
	// However, custom patterns explicitly set in a page's front matter should
	// still be honored (with date/category placeholders removed).
	//
	// Built-in styles (pretty, date, ordinal, none) apply to all document types,
	// but date/category placeholders are removed for non-posts.
	if !p.IsPost() {
		_, hasFrontMatterPermalink := p.fm["permalink"]

		if !isBuiltInStyle && !hasFrontMatterPermalink {
			// Custom global patterns don't apply to pages - use default pattern
			pattern = DefaultPermalinkPattern
		} else {
			// Built-in styles or explicit front matter permalinks: remove date/category placeholders
			pattern = removePostOnlyPlaceholders(pattern)
		}
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
	originalPattern := pattern

	// Use regex to remove category placeholders more comprehensively
	// This handles :categories in various positions and contexts
	categoryRegex := regexp.MustCompile(`:categories\b/?`)
	pattern = categoryRegex.ReplaceAllString(pattern, "")

	// Remove date-related placeholders using regex for better coverage
	// This handles all date placeholders regardless of position
	dateRegex := regexp.MustCompile(`:(?:year|month|i_month|day|i_day|hour|minute|second|short_year|y_day)\b/?`)
	pattern = dateRegex.ReplaceAllString(pattern, "")

	// Clean up any double slashes that might result
	pattern = regexp.MustCompile(`/+`).ReplaceAllString(pattern, "/")

	// Remove trailing slash temporarily for processing
	pattern = strings.TrimSuffix(pattern, "/")

	// Special case: if pattern becomes empty or just "/", use "/:title"
	if pattern == "" || pattern == "/" {
		pattern = "/:title"
	}

	// Ensure pattern starts with /
	if !strings.HasPrefix(pattern, "/") {
		pattern = "/" + pattern
	}

	// Preserve trailing slash if the original pattern ended with "/:title/" or similar
	// non-date/category placeholder patterns, or if it was a plain "/"
	if strings.HasSuffix(originalPattern, "/") {
		// Check what comes before the trailing slash
		beforeSlash := originalPattern[:len(originalPattern)-1]
		// If it ends with :title, :slug, :name, :path, :output_ext, :collection, etc.
		// (basically any placeholder that's not a date or category), keep the slash
		if strings.HasSuffix(beforeSlash, ":title") ||
			strings.HasSuffix(beforeSlash, ":slug") ||
			strings.HasSuffix(beforeSlash, ":name") ||
			strings.HasSuffix(beforeSlash, ":path") ||
			strings.HasSuffix(beforeSlash, ":collection") ||
			!regexp.MustCompile(`:\w+$`).MatchString(beforeSlash) {
			pattern += "/"
		}
	}

	return pattern
}

func (p *page) setPermalink() (err error) {
	p.permalink, err = p.computePermalink(p.permalinkVariables())
	return
}
