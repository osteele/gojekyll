package gojekyll

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/osteele/gojekyll/liquid"
	"github.com/osteele/gojekyll/templates"
)

// FindLayout returns a template for the named layout.
func (s *Site) FindLayout(base string, fm *templates.VariableMap) (t liquid.Template, err error) {
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
	return s.TemplateEngine().Parse(content)
}
