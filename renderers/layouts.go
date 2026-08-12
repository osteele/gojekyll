package renderers

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"strings"

	"github.com/osteele/gojekyll/frontmatter"
	"github.com/osteele/gojekyll/templates"
	"github.com/osteele/gojekyll/utils"
	"github.com/osteele/liquid"
)

type layoutCacheEntry struct {
	template    *liquid.Template
	frontMatter map[string]interface{}
	err         error
}

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
	if cached, ok := p.layoutCache.Load(base); ok {
		entry := cached.(layoutCacheEntry)
		if fmp != nil {
			*fmp = maps.Clone(entry.frontMatter)
		}
		return entry.template, entry.err
	}
	tpl, fm, err := p.findLayout(base)
	entry := layoutCacheEntry{template: tpl, frontMatter: fm, err: err}
	p.layoutCache.Store(base, entry)
	if fmp != nil {
		*fmp = maps.Clone(fm)
	}
	return tpl, err
}

func (p *Manager) findLayout(base string) (tpl *liquid.Template, fm map[string]interface{}, err error) {
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
			filename, err = utils.JoinWithin(dir, base+ext)
			if err != nil {
				return nil, nil, err
			}
			content, err = os.ReadFile(filename)
			if err == nil {
				found = true
				break loop
			}
			if !os.IsNotExist(err) {
				return nil, nil, err
			}
		}
	}
	if !found {
		return nil, nil, fmt.Errorf("no template for %s", base)
	}
	lineNo := 1
	fm, err = frontmatter.Read(&content, &lineNo)
	if err != nil {
		return
	}
	tpl, err = p.liquidEngine.ParseTemplateLocation(content, filename, lineNo)
	if err != nil {
		return nil, nil, err
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
