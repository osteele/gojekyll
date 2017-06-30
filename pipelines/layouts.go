package pipelines

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/osteele/gojekyll/templates"
	"github.com/osteele/liquid"
)

// FindLayout returns a template for the named layout.
func (p *Pipeline) FindLayout(base string, fm *templates.VariableMap) (t liquid.Template, err error) {
	exts := []string{"", ".html"}
	for _, ext := range strings.SplitN(p.config.MarkdownExt, `,`, -1) {
		exts = append(exts, "."+ext)
	}
	var (
		name    string
		content []byte
		found   bool
	)
	for _, ext := range exts {
		// TODO respect layout config
		name = filepath.Join(p.LayoutsDir(), base+ext)
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
		return nil, fmt.Errorf("no template for %s", base)
	}
	*fm, err = templates.ReadFrontMatter(&content)
	if err != nil {
		return
	}
	return p.liquidEngine.ParseTemplate(content)
}

// LayoutsDir returns the path to the layouts directory.
func (p *Pipeline) LayoutsDir() string {
	return filepath.Join(p.SourceDir, p.config.LayoutsDir)
}
