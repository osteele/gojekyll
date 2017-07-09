package pipelines

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/osteele/gojekyll/frontmatter"
	"github.com/osteele/liquid"
)

// FindLayout returns a template for the named layout.
func (p *Pipeline) FindLayout(base string, fm *map[string]interface{}) (tpl liquid.Template, err error) {
	exts := []string{"", ".html"}
	for _, ext := range strings.SplitN(p.config.MarkdownExt, `,`, -1) {
		exts = append(exts, "."+ext)
	}
	var (
		filename string
		content  []byte
		found    bool
	)
	for _, ext := range exts {
		// TODO respect layout config
		filename = filepath.Join(p.LayoutsDir(), base+ext)
		content, err = ioutil.ReadFile(filename)
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
	lineNo := 1
	*fm, err = frontmatter.Read(&content, &lineNo)
	if err != nil {
		return
	}
	tpl, err = p.liquidEngine.ParseTemplate(content)
	if err != nil {
		return nil, err
	}
	tpl.SetSourceLocation(filename, lineNo)
	return
}

// LayoutsDir returns the path to the layouts directory.
func (p *Pipeline) LayoutsDir() string {
	return filepath.Join(p.SourceDir(), p.config.LayoutsDir)
}
