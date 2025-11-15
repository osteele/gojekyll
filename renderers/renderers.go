package renderers

import (
	"io"
	"path/filepath"
	"strings"
	"sync"

	"github.com/osteele/gojekyll/config"
	"github.com/osteele/gojekyll/filters"
	"github.com/osteele/gojekyll/tags"
	"github.com/osteele/gojekyll/utils"
	"github.com/osteele/liquid"
	sass "github.com/bep/godartsass/v2"
)

// Renderers applies transformations to a document.
type Renderers interface {
	ApplyLayout(string, []byte, liquid.Bindings) ([]byte, error)
	Render(io.Writer, []byte, liquid.Bindings, string, int) error
	RenderTemplate([]byte, liquid.Bindings, string, int) ([]byte, error)
}

// Manager applies a rendering transformation to a file.
type Manager struct {
	Options
	cfg            config.Config
	liquidEngine   *liquid.Engine
	sassTempDir    string
	sassHash       string
	sassTranspiler *sass.Transpiler
	sassInitOnce   sync.Once
	sassInitErr    error
}

// Options configures a rendering manager.
type Options struct {
	RelativeFilenameToURL tags.LinkTagHandler
	ThemeDir              string
}

// New makes a rendering manager.
func New(c config.Config, options Options) (*Manager, error) {
	p := Manager{Options: options, cfg: c}
	p.liquidEngine = p.makeLiquidEngine()
	if err := p.copySASSFileIncludes(); err != nil {
		return nil, err
	}
	return &p, nil
}

// sourceDir returns the site source directory. Seeing how far we can bend
// the Law of Demeter.
func (p *Manager) sourceDir() string {
	return p.cfg.Source
}

// TemplateEngine returns the Liquid engine.
func (p *Manager) TemplateEngine() *liquid.Engine {
	return p.liquidEngine
}

// Render sends content through SASS and/or Liquid -> Markdown
func (p *Manager) Render(w io.Writer, src []byte, vars liquid.Bindings, filename string, lineNo int) error {
	if p.cfg.IsSASSPath(filename) {
		return p.WriteSass(w, src)
	}
	src, err := p.RenderTemplate(src, vars, filename, lineNo)
	if err != nil {
		return err
	}
	if p.cfg.IsMarkdown(filename) {
		src, err = renderMarkdownWithOptions(src, p.getTOCOptions())
		if err != nil {
			return err
		}
	}
	_, err = w.Write(src)
	return err
}

// getTOCOptions extracts TOC configuration from kramdown settings in _config.yml
func (p *Manager) getTOCOptions() *TOCOptions {
	opts := &TOCOptions{
		MinLevel: 1,
		MaxLevel: 6,
		UseJekyllHTML: false,
	}

	// Check for kramdown configuration
	if kramdown, ok := p.cfg.Map("kramdown"); ok {
		// Parse toc_levels (e.g., "1..6" or "2..3")
		if tocLevels, ok := kramdown["toc_levels"]; ok {
			minLevel, maxLevel := parseTOCLevels(tocLevels)
			if minLevel > 0 && maxLevel > 0 {
				opts.MinLevel = minLevel
				opts.MaxLevel = maxLevel
			}
		}
	}

	return opts
}

// parseTOCLevels parses Jekyll's toc_levels format (e.g., "1..6", "2..3", [1, 2, 3])
func parseTOCLevels(value interface{}) (int, int) {
	switch v := value.(type) {
	case string:
		// Parse "1..6" format
		parts := strings.Split(v, "..")
		if len(parts) == 2 {
			minLevel := parseInt(parts[0], 1)
			maxLevel := parseInt(parts[1], 6)
			return minLevel, maxLevel
		}
	case []interface{}:
		// Parse array format [1, 2, 3, 4]
		if len(v) > 0 {
			minLevel := 6
			maxLevel := 1
			for _, item := range v {
				if level, ok := item.(int); ok {
					if level < minLevel {
						minLevel = level
					}
					if level > maxLevel {
						maxLevel = level
					}
				}
			}
			if minLevel <= maxLevel {
				return minLevel, maxLevel
			}
		}
	}
	return 1, 6
}

// parseInt parses a string to int with a default value
func parseInt(s string, defaultVal int) int {
	s = strings.TrimSpace(s)
	val := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			val = val*10 + int(c-'0')
		}
	}
	if val == 0 {
		return defaultVal
	}
	return val
}

// RenderTemplate renders a Liquid template
func (p *Manager) RenderTemplate(src []byte, vars liquid.Bindings, filename string, lineNo int) ([]byte, error) {
	tpl, err := p.liquidEngine.ParseTemplateLocation(src, filename, lineNo)
	if err != nil {
		return nil, utils.WrapPathError(err, filename)
	}
	out, err := tpl.Render(vars)
	if err != nil {
		return nil, utils.WrapPathError(err, filename)
	}
	return out, err
}

func (p *Manager) makeLiquidEngine() *liquid.Engine {
	dirs := []string{filepath.Join(p.cfg.Source, p.cfg.IncludesDir)}
	if p.ThemeDir != "" {
		dirs = append(dirs, filepath.Join(p.ThemeDir, "_includes"))
	}
	engine := liquid.NewEngine()
	filters.AddJekyllFilters(engine, &p.cfg)
	tags.AddJekyllTags(engine, &p.cfg, dirs, p.RelativeFilenameToURL)
	return engine
}

// getSassTranspiler returns the SASS transpiler, initializing it if necessary.
// This uses lazy initialization to avoid creating the transpiler at package load time,
// which can cause "connection is shut down" errors in CI/CD environments.
func (p *Manager) getSassTranspiler() (*sass.Transpiler, error) {
	p.sassInitOnce.Do(func() {
		p.sassTranspiler, p.sassInitErr = sass.Start(sass.Options{})
	})
	return p.sassTranspiler, p.sassInitErr
}
