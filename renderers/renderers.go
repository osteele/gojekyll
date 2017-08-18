package renderers

import (
	"io"
	"path/filepath"

	"github.com/osteele/gojekyll/config"
	"github.com/osteele/gojekyll/filters"
	"github.com/osteele/gojekyll/tags"
	"github.com/osteele/gojekyll/utils"
	"github.com/osteele/liquid"
)

// Renderers applies transformations to a document.
type Renderers interface {
	ApplyLayout(string, []byte, liquid.Bindings) ([]byte, error)
	OutputExt(pathname string) string
	Render(io.Writer, []byte, liquid.Bindings, string, int) error
	RenderTemplate([]byte, liquid.Bindings, string, int) ([]byte, error)
}

// Manager applies a rendering transformation to a file.
type Manager struct {
	Options
	cfg          config.Config
	liquidEngine *liquid.Engine
	sassTempDir  string
	sassHash     string
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
	if err := p.CopySassFileIncludes(); err != nil {
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

// OutputExt returns the output extension.
func (p *Manager) OutputExt(pathname string) string {
	return p.cfg.OutputExt(pathname)
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
		src = renderMarkdown(src)
	}
	_, err = w.Write(src)
	return err
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
