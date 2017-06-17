package main

import (
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	. "github.com/osteele/gojekyll/helpers"
)

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
func (p *pageFields) permalinkTemplateVariables() map[string]string {
	var (
		collectionName string
		path           = p.relpath
		ext            = filepath.Ext(path)
		outputExt      = ext
		root           = PathWithoutExtension(path)
		name           = filepath.Base(root)
		title          = p.frontMatter.String("title", name)
	)
	switch {
	case p.site.IsMarkdown(path):
		outputExt = ".html"
	case p.site.IsSassPath(path):
		outputExt = ".css"
	}
	if val, found := p.frontMatter["collection"]; found {
		collectionName = val.(string)
		prefix := "_" + collectionName + "/"
		if !strings.HasPrefix(path, prefix) {
			panic(fmt.Errorf("Expected %s to start with %s", path, prefix))
		}
		root = root[len(prefix):]
	}
	vs := map[string]string{
		"collection": collectionName,
		"name":       Slugify(name),
		"path":       "/" + root,
		"title":      title,
		"slug":       Slugify(name),
		// TODO categories
		// The following isn't documented, but is evident
		"output_ext": outputExt,
	}
	d := time.Now() // TODO read from frontMatter or use file modtime
	for name, f := range permalinkDateVariables {
		vs[name] = d.Format(f)
	}
	return vs
}

func (p *pageFields) expandPermalink() (s string, err error) {
	pattern := p.frontMatter.String("permalink", ":path:output_ext")
	if p, found := PermalinkStyles[pattern]; found {
		pattern = p
	}
	templateVariables := p.permalinkTemplateVariables()
	// The ReplaceAllStringFunc callback signals errors via panic.
	// Turn them into return values.
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				panic(r)
			}
		}
	}()
	s = templateVariableMatcher.ReplaceAllStringFunc(pattern, func(m string) string {
		varname := m[1:]
		value, found := templateVariables[varname]
		if !found {
			panic(fmt.Errorf("unknown variable %s in permalink template %s", varname, pattern))
		}
		return value
	})
	return path.Clean(s), nil
}

// The permalink is computed once instead of on demand, so that subsequent
// access needn't check for an error.
func (p *pageFields) initPermalink() (err error) {
	p.permalink, err = p.expandPermalink()
	return
}
