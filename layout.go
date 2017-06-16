package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	. "github.com/osteele/gojekyll/helpers"

	"github.com/acstech/liquid"
)

// FindLayout returns a template for the named layout.
func (s *Site) FindLayout(name string, fm *VariableMap) (t *liquid.Template, err error) {
	exts := []string{"", ".html"}
	for _, ext := range strings.SplitN(s.config.MarkdownExt, `,`, -1) {
		exts = append(exts, "."+ext)
	}
	var (
		path    string
		content []byte
		found   bool
	)
	for _, ext := range exts {
		// TODO respect layout config
		path = filepath.Join(s.Source, "_layouts", name+ext)
		content, err = ioutil.ReadFile(path)
		if err == nil {
			found = true
			break
		}
		if !os.IsNotExist(err) {
			return nil, err
		}
	}
	if !found {
		panic(fmt.Errorf("no template for %s", name))
	}
	*fm, err = readFrontMatter(&content)
	if err != nil {
		return
	}
	return liquid.Parse(content, nil)
}

func (p *DynamicPage) applyLayout(frontMatter VariableMap, body []byte) ([]byte, error) {
	for {
		layoutName := frontMatter.String("layout", "")
		if layoutName == "" {
			break
		}
		template, err := p.site.FindLayout(layoutName, &frontMatter)
		if err != nil {
			return nil, err
		}
		vars := MergeVariableMaps(p.TemplateVariables(), VariableMap{
			"content": body,
			"layout":  frontMatter,
		})
		body, err = RenderTemplate(template, vars)
		if err != nil {
			return nil, err
		}
	}
	return body, nil
}
