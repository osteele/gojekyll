package pipelines

import (
	"bytes"
	"io"
	"path/filepath"

	"github.com/osteele/gojekyll/config"
	"github.com/osteele/gojekyll/filters"
	"github.com/osteele/gojekyll/tags"
	"github.com/osteele/gojekyll/utils"
	"github.com/osteele/liquid"
)

// PipelineInterface applies transformations to a document.
type PipelineInterface interface {
	ApplyLayout(string, []byte, map[string]interface{}) ([]byte, error)
	OutputExt(pathname string) string
	Render(io.Writer, []byte, string, int, map[string]interface{}) ([]byte, error)
}

// Pipeline applies a rendering transformation to a file.
type Pipeline struct {
	PipelineOptions
	cfg          config.Config
	liquidEngine *liquid.Engine
	sassTempDir  string
	sassHash     string
}

// PipelineOptions configures a pipeline.
type PipelineOptions struct {
	RelativeFilenameToURL tags.LinkTagHandler
	ThemeDir              string
}

// NewPipeline makes a rendering pipeline.
func NewPipeline(c config.Config, options PipelineOptions) (*Pipeline, error) {
	p := Pipeline{PipelineOptions: options, cfg: c}
	p.liquidEngine = p.makeLiquidEngine()
	if err := p.CopySassFileIncludes(); err != nil {
		return nil, err
	}
	return &p, nil
}

// SourceDir returns the site source directory. Seeing how far we can bend
// the Law of Demeter.
func (p *Pipeline) SourceDir() string {
	return p.cfg.Source
}

// TemplateEngine returns the Liquid engine.
func (p *Pipeline) TemplateEngine() *liquid.Engine {
	return p.liquidEngine
}

// OutputExt returns the output extension.
func (p *Pipeline) OutputExt(pathname string) string {
	return p.cfg.OutputExt(pathname)
}

// Render returns nil iff it wrote to the writer
func (p *Pipeline) Render(w io.Writer, b []byte, filename string, lineNo int, e map[string]interface{}) ([]byte, error) {
	if p.cfg.IsSASSPath(filename) {
		buf := new(bytes.Buffer)
		if err := p.WriteSass(buf, b); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}
	b, err := p.renderTemplate(b, e, filename, lineNo)
	if err != nil {
		return nil, err
	}
	if p.cfg.IsMarkdown(filename) {
		b = markdownRenderer(b)
	}
	return b, nil
}

func (p *Pipeline) renderTemplate(src []byte, b map[string]interface{}, filename string, lineNo int) ([]byte, error) {
	tpl, err := p.liquidEngine.ParseTemplateLocation(src, filename, lineNo)
	if err != nil {
		return nil, utils.WrapPathError(err, filename)
	}
	out, err := tpl.Render(b)
	if err != nil {
		return nil, utils.WrapPathError(err, filename)
	}
	return out, err
}

func (p *Pipeline) makeLiquidEngine() *liquid.Engine {
	dirs := []string{filepath.Join(p.cfg.Source, p.cfg.IncludesDir)}
	if p.ThemeDir != "" {
		dirs = append(dirs, filepath.Join(p.ThemeDir, "_includes"))
	}
	engine := liquid.NewEngine()
	filters.AddJekyllFilters(engine, &p.cfg)
	tags.AddJekyllTags(engine, &p.cfg, dirs, p.RelativeFilenameToURL)
	return engine
}
