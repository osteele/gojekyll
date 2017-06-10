package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"

	"github.com/acstech/liquid"
	"github.com/russross/blackfriday"
	yaml "gopkg.in/yaml.v2"
)

var frontmatterMatcher = regexp.MustCompile(`(?s)^---\n(.+?\n)---\n`)
var templateVariableMatcher = regexp.MustCompile(`:(?:collection|path|name|title)\b`)
var nonAlphanumericSequenceMatcher = regexp.MustCompile(`[^[:alnum:]]+`)

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
		body = source[match[1]:]
		fm := map[interface{}]interface{}{}
		err = yaml.Unmarshal(source[match[2]:match[3]], &fm)
		if err != nil {
			return nil, err
		}

		data = mergeMaps(data, fm)
	}

	ext := filepath.Ext(path)

	// var title string
	// if val, ok := data["permalink"]; ok {
	// 	title = fmt.Sprintf("%v", val)
	// } else {
	// 	title = filepath.Base(path)
	// 	title = title[:len(title)-len(ext)]
	// }

	permalink := path
	if val, ok := data["permalink"]; ok {
		permalink, ok = val.(string)
		if !ok {
			return nil, errors.New("Required string value for permalink")
		}
	}
	templateVariables := map[string]string{
		"output_ext": ".html",
		"path":       regexp.MustCompile(`\.md$`).ReplaceAllLiteralString(path, ""),
		"name":       nonAlphanumericSequenceMatcher.ReplaceAllString(filepath.Base(path), "-"),
	}
	if val, found := data["collection"]; found {
		collectionName := val.(string)
		collectionPath := "_" + collectionName + "/"
		templateVariables["collection"] = collectionName
		templateVariables["path"] = templateVariables["path"][len(collectionPath):]
	}
	permalink = templateVariableMatcher.ReplaceAllStringFunc(permalink, func(m string) string {
		return templateVariables[m[1:]]
	})

	if expand {
		template, err := liquid.Parse(body, nil)
		if err != nil {
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
