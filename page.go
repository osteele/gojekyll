package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/acstech/liquid"
	"github.com/russross/blackfriday"
	yaml "gopkg.in/yaml.v2"
)

var (
	frontmatterMatcher             = regexp.MustCompile(`(?s)^---\n(.+?\n)---\n`)
	templateVariableMatcher        = regexp.MustCompile(`:(?:collection|file_ext|name|path|title)\b`)
	nonAlphanumericSequenceMatcher = regexp.MustCompile(`[^[:alnum:]]+`)
)

var permalinkStyles = map[string]string{
	"date":    "/:categories/:year/:month/:day/:title.html",
	"pretty":  "/:categories/:year/:month/:day/:title/",
	"ordinal": "/:categories/:year/:y_day/:title.html",
	"none":    "/:categories/:title.html",
}

// A Page represents an HTML page.
type Page struct {
	Path      string
	Permalink string
	Static    bool
	Expanded  bool
	Published bool
	Body      []byte
}

func (p Page) String() string {
	return fmt.Sprintf("Page{Path=%v, Permalink=%v, Static=%v}",
		p.Path, p.Permalink, p.Static)
}

func readFile(path string, defaults map[interface{}]interface{}, expand bool) (*Page, error) {
	// TODO don't read, parse binary files

	source, err := ioutil.ReadFile(filepath.Join(siteConfig.SourceDir, path))
	if err != nil {
		return nil, err
	}

	static := true
	data := defaults
	body := source

	if match := frontmatterMatcher.FindSubmatchIndex(source); match != nil {
		static = false
		if expand {
			// TODO only prepend newlines if it's markdown
			body = append(
				regexp.MustCompile(`[^\n\r]+`).ReplaceAllLiteral(source[:match[1]], []byte{}),
				source[match[1]:]...)
		}
		fm := map[interface{}]interface{}{}
		err = yaml.Unmarshal(source[match[2]:match[3]], &fm)
		if err != nil {
			err := &os.PathError{Op: "read frontmatter", Path: path, Err: err}
			return nil, err
		}

		data = mergeMaps(data, fm)
	}

	ext := filepath.Ext(path)

	permalink := path
	if val, ok := data["permalink"]; ok {
		pattern, ok := val.(string)
		if !ok {
			err := errors.New("permalink value must be a string")
			err = &os.PathError{Op: "render", Path: path, Err: err}
			return nil, err
		}
		permalink = expandPermalinkPattern(pattern, data, path)
	}

	if expand {
		template, err := liquid.Parse(body, nil)
		if err != nil {
			fmt.Println(data)
			err := &os.PathError{Op: "Liquid Error", Path: path, Err: err}
			return nil, err
		}
		writer := new(bytes.Buffer)
		template.Render(writer, stringMap(data))
		body = writer.Bytes()
		if ext == ".md" {
			body = blackfriday.MarkdownBasic(body)
		}
	} else {
		body = []byte{}
	}

	return &Page{
		Path:      path,
		Permalink: permalink,
		Expanded:  expand,
		Static:    static,
		Published: getBool(data, "published", true),
		Body:      body,
	}, nil
}

func expandPermalinkPattern(pattern string, data map[interface{}]interface{}, path string) string {
	if p, found := permalinkStyles[pattern]; found {
		pattern = p
	}

	var (
		collectionName string
		ext            = filepath.Ext(path)
		localPath      = path
		outputExt      = ext
		name           = filepath.Base(localPath)
		title          = getString(data, "title", name)
	)

	if ext == ".md" {
		outputExt = ""
		localPath = localPath[:len(localPath)-len(ext)]
	}

	if val, found := data["collection"]; found {
		collectionName = val.(string)
		collectionPath := "_" + collectionName + "/"
		localPath = localPath[len(collectionPath):]
	}

	hyphenize := func(s string) string {
		return nonAlphanumericSequenceMatcher.ReplaceAllString(s, "-")
	}

	templateVariables := map[string]string{
		"collection": collectionName,
		"name":       hyphenize(name),
		"output_ext": outputExt,
		"path":       localPath,
		"title":      hyphenize(title),
		// TODO year month imonth day i_day short_year hour minute second slug categories
	}

	return templateVariableMatcher.ReplaceAllStringFunc(pattern, func(m string) string {
		varname := m[1:]
		value := templateVariables[varname]
		if value == "" {
			fmt.Printf("%s: unknown variable %s in permalink template\n", varname)
		}
		return value
	})
}
