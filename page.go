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
	FrontMatter map[interface{}]interface{}
	Content     []byte
}

func (p Page) String() string {
	return fmt.Sprintf("Page{Path=%v, Permalink=%v, Static=%v}",
		p.Path, p.Permalink, p.Static)
}

// PageData returns metadata for use in the representation of the page as a collection item
func (p Page) PageData() map[interface{}]interface{} {
	// should have title, parts, url, description, due_date
	data := map[interface{}]interface{}{
		"url": p.Permalink,
		// TODO Posts should get date, category, categories, tags
		// TODO only do the following if it's a collection document?
		"path":          p.Source(),
		"relative_path": p.Path,
		// TODO collections: content output collection(name) date(of the collection)
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

// Data returns the variable context for Liquid evaluation
func (p Page) Data() map[interface{}]interface{} {
	return map[interface{}]interface{}{
		"page": p.PageData(),
		"site": site.Data,
	}
}

// ReadPage reads a Page from a file, using defaults as the default front matter.
func ReadPage(path string, defaults map[interface{}]interface{}) (p *Page, err error) {
	var (
		frontMatter map[interface{}]interface{}
		static      = true
		body        []byte
	)

	// TODO don't read, parse binary files
	source, err := ioutil.ReadFile(filepath.Join(site.Config.SourceDir, path))
	if err != nil {
		return nil, err
	}

	if match := frontmatterMatcher.FindSubmatchIndex(source); match != nil {
		static = false
		// TODO only prepend newlines if it's markdown
		body = append(
			regexp.MustCompile(`[^\n\r]+`).ReplaceAllLiteral(source[:match[1]], []byte{}),
			source[match[1]:]...)
		frontMatter = map[interface{}]interface{}{}
		err = yaml.Unmarshal(source[match[2]:match[3]], &frontMatter)
		if err != nil {
			err := &os.PathError{Op: "read frontmatter", Path: path, Err: err}
			return nil, err
		}

		frontMatter = mergeMaps(defaults, frontMatter)
	} else {
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
		permalink = expandPermalinkPattern(pattern, data, path)
	}

	p = &Page{
		Path:        path,
		Permalink:   permalink,
		Static:      static,
		Published:   getBool(data, "published", true),
		FrontMatter: data,
		Content:     body,
	}

	return p, nil
}

// Source returns the file path of the page source.
func (p *Page) Source() string {
	return filepath.Join(site.Config.SourceDir, p.Path)
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
	template.Render(writer, stringMap(p.Data()))
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

func expandPermalinkPattern(pattern string, data map[interface{}]interface{}, path string) string {
	if p, found := permalinkStyles[pattern]; found {
		pattern = p
	}

	var (
		collectionName string
		localPath      = path
		ext            = filepath.Ext(path)
		root           = path[:len(path)-len(ext)]
		outputExt      = ext
		name           = filepath.Base(root)
		title          = getString(data, "title", name)
	)

	if isMarkdown(path) {
		outputExt = ".html"
	}

	if val, found := data["collection"]; found {
		collectionName = val.(string)
		prefix := "_" + collectionName + "/"
		localPath = localPath[len(prefix):]
	}

	replaceNonalphumericsByHyphens := func(s string) string {
		return nonAlphanumericSequenceMatcher.ReplaceAllString(s, "-")
	}

	templateVariables := map[string]string{
		"collection": collectionName,
		"ext":        strings.TrimLeft(ext, "."),
		"name":       replaceNonalphumericsByHyphens(name),
		"output_ext": strings.TrimLeft(outputExt, "."),
		"path":       localPath,
		"title":      replaceNonalphumericsByHyphens(title),
		// TODO year month imonth day i_day short_year hour minute second slug categories
	}

	return templateVariableMatcher.ReplaceAllStringFunc(pattern, func(m string) string {
		varname := m[1:]
		value := templateVariables[varname]
		if value == "" {
			fmt.Printf("unknown variable %s in permalink template\n", varname)
		}
		return value
	})
}
