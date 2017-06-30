package pipelines

import (
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/osteele/gojekyll/config"
	"github.com/osteele/gojekyll/filters"
	"github.com/osteele/gojekyll/helpers"
	"github.com/osteele/gojekyll/tags"
	"github.com/osteele/gojekyll/templates"
	"github.com/osteele/liquid"
	"github.com/russross/blackfriday"
)

// PipelineInterface applies transformations to a document.
type PipelineInterface interface {
	ApplyLayout(string, []byte, templates.VariableMap) ([]byte, error)
	OutputExt(pathname string) string
	Render(io.Writer, []byte, string, templates.VariableMap) ([]byte, error)
}

// Pipeline applies a rendering transformation to a file.
type Pipeline struct {
	PipelineOptions
	config       config.Config
	liquidEngine liquid.Engine
	sassTempDir  string
}

// PipelineOptions configures a pipeline.
type PipelineOptions struct {
	SourceDir             string
	RelativeFilenameToURL tags.LinkTagHandler
}

// NewPipeline makes a rendering pipeline.
func NewPipeline(c config.Config, options PipelineOptions) (*Pipeline, error) {
	p := Pipeline{PipelineOptions: options, config: c}
	p.liquidEngine = p.makeLiquidEngine()
	if err := p.CopySassFileIncludes(); err != nil {
		return nil, err
	}
	return &p, nil
}

// TemplateEngine returns the Liquid engine.
func (p *Pipeline) TemplateEngine() liquid.Engine {
	return p.liquidEngine
}

// OutputExt returns the output extension.
func (p *Pipeline) OutputExt(pathname string) string {
	return p.config.OutputExt(pathname)
}

// Render returns nil iff it wrote to the writer
func (p *Pipeline) Render(w io.Writer, b []byte, filename string, e templates.VariableMap) ([]byte, error) {
	if p.config.IsSassPath(filename) {
		return nil, p.WriteSass(w, b)
	}
	b, err := p.renderTemplate(b, e, filename)
	if err != nil {
		return nil, err
	}
	if p.config.IsMarkdown(filename) {
		b = blackfriday.MarkdownCommon(b)
	}
	return b, nil
}

func (p *Pipeline) renderTemplate(b []byte, e templates.VariableMap, filename string) ([]byte, error) {
	b, err := p.liquidEngine.ParseAndRender(b, e)
	if err != nil {
		return nil, helpers.PathError(err, "Liquid Error", filename)
	}
	return b, err
}

// ApplyLayout applies the named layout to the data.
func (p *Pipeline) ApplyLayout(name string, data []byte, e templates.VariableMap) ([]byte, error) {
	for name != "" {
		var lfm templates.VariableMap
		t, err := p.FindLayout(name, &lfm)
		if err != nil {
			return nil, err
		}
		le := templates.MergeVariableMaps(e, templates.VariableMap{
			"content": string(data),
			"layout":  lfm,
		})
		data, err = t.Render(le)
		if err != nil {
			return nil, helpers.PathError(err, "render template", name)
		}
		name = lfm.String("layout", "")
	}
	return data, nil
}

func (p *Pipeline) makeLiquidEngine() liquid.Engine {
	engine := liquid.NewEngine()
	includeHandler := func(name string, w io.Writer, scope map[string]interface{}) error {
		filename := filepath.Join(p.SourceDir, p.config.IncludesDir, name)
		template, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}
		text, err := engine.ParseAndRender(template, scope)
		if err != nil {
			return err
		}
		_, err = w.Write(text)
		return err
	}
	filters.AddJekyllFilters(engine, p.config)
	tags.AddJekyllTags(engine, p.config, includeHandler, p.RelativeFilenameToURL)
	return engine
}
