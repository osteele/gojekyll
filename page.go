package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/acstech/liquid"
	"github.com/russross/blackfriday"
	yaml "gopkg.in/yaml.v2"
)

var (
	frontmatterMatcher             = regexp.MustCompile(`(?s)^---\n(.+?\n)---\n`)
	templateVariableMatcher        = regexp.MustCompile(`:\w+\b`)
	nonAlphanumericSequenceMatcher = regexp.MustCompile(`[^[:alnum:]]+`)
)

type VariableMap map[string]interface{}

var permalinkStyles = map[string]string{
	"date":    "/:categories/:year/:month/:day/:title.html",
	"pretty":  "/:categories/:year/:month/:day/:title/",
	"ordinal": "/:categories/:year/:y_day/:title.html",
	"none":    "/:categories/:title.html",
}

// A Page represents an HTML page.
type Page struct {
	Path        string // this is the relative path
	Permalink   string
	Static      bool
	Published   bool
	FrontMatter VariableMap
	Content     []byte
}

func (p *Page) String() string {
	return fmt.Sprintf("Page{Path=%v, Permalink=%v, Static=%v}",
		p.Path, p.Permalink, p.Static)
}

// PageVariables returns metadata for use in the representation of the page as a collection item
func (p *Page) PageVariables() VariableMap {
	if p.Static {
		return mergeVariableMaps(p.FrontMatter, p.staticFileData())
	}
	data := VariableMap{
		"url":  p.Permalink,
		"path": p.Source(),
		// TODO content title excerpt date id categories tags next previous
		// TODO Posts should get date, category, categories, tags
		// TODO only do the following if it's a collection document?
		"relative_path": p.Path,
		// TODO collections: output collection(name) date(of the collection)
	}
	for k, v := range p.FrontMatter {
		switch k {
		case "layout", "permalink", "published":
		default:
			data[k] = v
		}
	}
	return data
}

func (p *Page) staticFileData() VariableMap {
	var (
		path = "/" + p.Path
		base = filepath.Base(path)
		ext  = filepath.Ext(path)
	)

	return VariableMap{
		"path":          path,
		"modified_time": 0, // TODO
		"name":          base,
		"basename":      base[:len(base)-len(ext)],
		"extname":       ext,
	}
}

// Data returns the variable context for Liquid evaluation
func (p *Page) Data() VariableMap {
	return VariableMap{
		"page": p.PageVariables(),
		"site": site.Variables,
	}
}

// ReadPage reads a Page from a file, using defaults as the default front matter.
func ReadPage(path string, defaults VariableMap) (p *Page, err error) {
	var (
		frontMatter VariableMap
		static      = true
		body        []byte
	)

	// TODO don't read, parse binary files
	source, err := ioutil.ReadFile(filepath.Join(site.Source, path))
	if err != nil {
		return nil, err
	}

	if match := frontmatterMatcher.FindSubmatchIndex(source); match != nil {
		static = false
		// TODO only prepend newlines if it's markdown
		body = append(
			regexp.MustCompile(`[^\n\r]+`).ReplaceAllLiteral(source[:match[1]], []byte{}),
			source[match[1]:]...)
		frontMatter = VariableMap{}
		err = yaml.Unmarshal(source[match[2]:match[3]], &frontMatter)
		if err != nil {
			err := &os.PathError{Op: "read frontmatter", Path: path, Err: err}
			return nil, err
		}

		frontMatter = mergeVariableMaps(defaults, frontMatter)
	} else {
		frontMatter = defaults
		body = []byte{}
	}

	data := frontMatter

	permalink := "/" + path
	if val, ok := data["permalink"]; ok {
		pattern, ok := val.(string)
		if !ok {
			err := errors.New("permalink value must be a string")
			err = &os.PathError{Op: "render", Path: path, Err: err}
			return nil, err
		}
		permalink, err = expandPermalinkPattern(pattern, path, data)
		if err != nil {
			return nil, err
		}
	}

	p = &Page{
		Path:        path,
		Permalink:   permalink,
		Static:      static,
		Published:   data.Bool("published", true),
		FrontMatter: data,
		Content:     body,
	}

	return p, nil
}

// Source returns the file path of the page source.
func (p *Page) Source() string {
	return filepath.Join(site.Source, p.Path)
}

// Render applies Liquid and Markdown, as appropriate.
func (p *Page) Render(w io.Writer) error {
	if p.Static {
		source, err := ioutil.ReadFile(p.Source())
		if err != nil {
			return err
		}
		_, err = w.Write(source)
		return err
	}

	parsingLiquid := true
	defer func() {
		if parsingLiquid {
			fmt.Println("While processing", p.Source())
		}
	}()
	template, err := liquid.Parse(p.Content, nil)
	parsingLiquid = false
	if err != nil {
		err := &os.PathError{Op: "Liquid Error", Path: p.Source(), Err: err}
		return err
	}
	writer := new(bytes.Buffer)
	template.Render(writer, p.Data())
	body := writer.Bytes()

	if isMarkdown(p.Path) {
		body = blackfriday.MarkdownCommon(body)
	}

	_, err = w.Write(body)
	return err
}

func isMarkdown(path string) bool {
	ext := filepath.Ext(path)
	return site.MarkdownExtensions()[strings.TrimLeft(ext, ".")]
}

func permalinkTemplateVariables(path string, data VariableMap) map[string]string {
	var (
		collectionName string
		localPath      = path
		ext            = filepath.Ext(path)
		root           = path[:len(path)-len(ext)]
		outputExt      = ext
		name           = filepath.Base(root)
		title          = data.String("title", name)
	)

	if isMarkdown(path) {
		outputExt = ".html"
	}

	if val, found := data["collection"]; found {
		collectionName = val.(string)
		prefix := "_" + collectionName + "/"
		localPath = localPath[len(prefix):]
	}

	// replace each sequence of non-alphanumerics by a single hypen
	hyphenize := func(s string) string {
		return nonAlphanumericSequenceMatcher.ReplaceAllString(s, "-")
	}

	return map[string]string{
		"collection": collectionName,
		"ext":        strings.TrimLeft(ext, "."),
		"name":       hyphenize(name),
		"output_ext": strings.TrimLeft(outputExt, "."),
		"path":       localPath,
		"title":      hyphenize(title),
		// TODO year month imonth day i_day short_year hour minute second slug categories
	}
}

func expandPermalinkPattern(pattern string, path string, data VariableMap) (s string, err error) {
	if p, found := permalinkStyles[pattern]; found {
		pattern = p
	}
	templateVariables := permalinkTemplateVariables(path, data)
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
	return
}
