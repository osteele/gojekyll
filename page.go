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

var frontmatterRe = regexp.MustCompile(`(?s)^---\n(.+?)\n---\n`)

// A Page represents an HTML page.
type Page struct {
	Path      string
	Permalink string
	Static    bool
	Expanded  bool
	Body      []byte
}

func (p Page) String() string {
	return fmt.Sprintf("Page{Path=%v, Permalink=%v}", p.Path, p.Permalink)
}

func readFile(path string, expand bool) (*Page, error) {
	// TODO don't read, parse binary files

	source, err := ioutil.ReadFile(filepath.Join(siteConfig.SourceDir, path))
	if err != nil {
		return nil, err
	}

	static := true
	data := map[string]interface{}{}
	body := source

	fmMatchIndex := frontmatterRe.FindSubmatchIndex(source)
	if fmMatchIndex != nil {
		static = false
		body = source[fmMatchIndex[1]:]
		fmBytes := source[fmMatchIndex[2]:fmMatchIndex[3]]
		var fmMap interface{}
		err = yaml.Unmarshal(fmBytes, &fmMap)
		if err != nil {
			return nil, err
		}
		fmStringMap, ok := fmMap.(map[interface{}]interface{})
		if !ok {
			return nil, errors.New("YAML frontmatter is not a map")
		}
		for k, v := range fmStringMap {
			stringer, ok := k.(fmt.Stringer)
			if ok {
				data[stringer.String()] = v
			} else {
				data[fmt.Sprintf("%v", k)] = v
			}
		}
	}

	ext := filepath.Ext(path)

	var title string
	if val, ok := data["permalink"]; ok {
		title = fmt.Sprintf("%v", val)
	} else {
		title = filepath.Base(path)
		title = title[:len(title)-len(ext)]
	}

	// TODO use site, collection default; expand components
	permalink := "/" + path[:len(path)-len(ext)]
	if val, ok := data["permalink"]; ok {
		permalink = val.(string) // TODO what if it's not a string?
	}

	if expand && ext == ".md" {
		template, err := liquid.Parse(body, nil)
		if err != nil {
			return nil, err
		}
		writer := new(bytes.Buffer)
		template.Render(writer, data)
		body = blackfriday.MarkdownBasic(writer.Bytes())
	}

	if !expand {
		body = []byte{}
	}

	return &Page{
		Path:      path,
		Permalink: permalink,
		Expanded:  expand,
		Static:    static,
		Body:      body,
	}, nil
}
