package renderers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/osteele/gojekyll/frontmatter"
	"github.com/osteele/gojekyll/templates"
	"github.com/osteele/gojekyll/utils"
	"github.com/osteele/liquid"
)

// ApplyLayout applies the named layout to the content.
func (p *Manager) ApplyLayout(name string, content []byte, vars liquid.Bindings) ([]byte, error) {
	for name != "" {
		var lfm map[string]interface{}
		tpl, err := p.FindLayout(name, &lfm)
		if err != nil {
			return nil, err
		}
		b := utils.MergeStringMaps(vars, map[string]interface{}{
			"content": string(content),
			"layout":  lfm,
		})
		content, err = tpl.Render(b)
		if err != nil {
			return nil, utils.WrapPathError(err, name)
		}
		name = templates.VariableMap(lfm).String("layout", "")
	}
	return content, nil
}

// FindLayout returns a template for the named layout.
func (p *Manager) FindLayout(base string, fmp *map[string]interface{}) (tpl *liquid.Template, err error) {
	// not cached, but the time here is negligible
	exts := []string{"", ".html"}
	for _, ext := range strings.Split(p.cfg.MarkdownExt, `,`) {
		exts = append(exts, "."+ext)
	}
	var (
		filename string
		content  []byte
		found    bool
	)
loop:
	for _, dir := range p.layoutDirs() {
		for _, ext := range exts {
			filename = filepath.Join(dir, base+ext)
			content, err = os.ReadFile(filename)
			if err == nil {
				found = true
				break loop
			}
			if !os.IsNotExist(err) {
				return nil, err
			}
		}
	}
	if !found {
		return nil, fmt.Errorf("no template for %s", base)
	}
	lineNo := 1
	fm, err := frontmatter.Read(&content, &lineNo)
	if err != nil {
		return
	}
	if fmp != nil {
		*fmp = fm
	}
	tpl, err = p.liquidEngine.ParseTemplateLocation(content, filename, lineNo)
	if err != nil {
		return nil, err
	}
	return
}

// LayoutsDir returns the path to the layouts directory.
func (p *Manager) layoutDirs() []string {
	dirs := []string{filepath.Join(p.sourceDir(), p.cfg.LayoutsDir)}
	if p.ThemeDir != "" {
		dirs = append(dirs, filepath.Join(p.ThemeDir, "_layouts"))
	}
	return dirs
}
