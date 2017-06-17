package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/osteele/gojekyll/helpers"

	"github.com/acstech/liquid"
	"github.com/acstech/liquid/core"
)

// FindLayout returns a template for the named layout.
func (s *Site) FindLayout(base string, fm *VariableMap) (t *liquid.Template, err error) {
	exts := []string{"", ".html"}
	for _, ext := range strings.SplitN(s.config.MarkdownExt, `,`, -1) {
		exts = append(exts, "."+ext)
	}
	var (
		name    string
		content []byte
		found   bool
	)
	for _, ext := range exts {
		// TODO respect layout config
		name = filepath.Join(s.LayoutsDir(), base+ext)
		content, err = ioutil.ReadFile(name)
		if err == nil {
			found = true
			break
		}
		if !os.IsNotExist(err) {
			return nil, err
		}
	}
	if !found {
		panic(fmt.Errorf("no template for %s", base))
	}
	*fm, err = readFrontMatter(&content)
	if err != nil {
		return
	}
	return liquid.Parse(content, s.LiquidConfiguration())
}

func (p *DynamicPage) applyLayout(frontMatter VariableMap, body []byte, config *core.Configuration) ([]byte, error) {
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
		body, err = helpers.RenderTemplate(template, vars)
		if err != nil {
			return nil, err
		}
	}
	return body, nil
}
