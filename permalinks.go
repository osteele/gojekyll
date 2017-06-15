package main

import (
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"time"
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
func permalinkTemplateVariables(path string, frontMatter VariableMap) map[string]string {
	var (
		collectionName string
		localPath      = path
		ext            = filepath.Ext(path)
		outputExt      = ext
		root           = path[:len(path)-len(ext)]
		name           = filepath.Base(root)
		title          = frontMatter.String("title", name)
	)
	if isMarkdown(path) {
		outputExt = ".html"
	}
	if val, found := frontMatter["collection"]; found {
		collectionName = val.(string)
		prefix := "_" + collectionName + "/"
		localPath = localPath[len(prefix):]
	}
	vs := map[string]string{
		"collection": collectionName,
		"name":       name,
		"path":       "/" + localPath,
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

func expandPermalinkPattern(pattern string, rel string, frontMatter VariableMap) (s string, err error) {
	if p, found := PermalinkStyles[pattern]; found {
		pattern = p
	}
	templateVariables := permalinkTemplateVariables(rel, frontMatter)
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
