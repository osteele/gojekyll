package pipelines

import (
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/osteele/gojekyll/config"
	"github.com/osteele/gojekyll/helpers"
	"github.com/osteele/gojekyll/liquid"
	"github.com/osteele/gojekyll/templates"
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
	config       config.Config
	liquidEngine liquid.Engine
	pageSupplier PageSupplier
	sassTempDir  string
	sourceDir    string
}

// PipelineOptions configures a pipeline.
type PipelineOptions struct {
	UseRemoteLiquidEngine bool
}

// PageSupplier tells a pipeline how to resolve link tags.
type PageSupplier interface {
	FilenameURLs() map[string]string
	RelativeFilenameToURL(string) (string, bool)
}

// NewPipeline makes a rendering pipeline.
func NewPipeline(sourceDir string, c config.Config, pageSupplier PageSupplier, o PipelineOptions) (*Pipeline, error) {
	p := Pipeline{config: c, pageSupplier: pageSupplier, sourceDir: sourceDir}
	engine, err := p.makeLiquidEngine(o)
	if err != nil {
		return nil, err
	}
	p.liquidEngine = engine
	if err := p.CopySassFileIncludes(); err != nil {
		return nil, err
	}
	return &p, nil
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
		switch err := err.(type) {
		case *liquid.RenderError:
			if err.Filename == "" {
				err.Filename = filename
			}
			return nil, err
		default:
			return nil, helpers.PathError(err, "Liquid Error", filename)
		}
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

func (p *Pipeline) makeLocalLiquidEngine() liquid.Engine {
	engine := liquid.NewLocalWrapperEngine()
	engine.LinkTagHandler(p.pageSupplier.RelativeFilenameToURL)
	includeHandler := func(name string, w io.Writer, scope map[string]interface{}) error {
		filename := filepath.Join(p.sourceDir, p.config.IncludesDir, name)
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
	engine.IncludeHandler(includeHandler)
	return engine
}

func (p *Pipeline) makeLiquidClient() (engine liquid.RemoteEngine, err error) {
	engine, err = liquid.NewRPCClientEngine(liquid.DefaultServer)
	if err != nil {
		return
	}
	err = engine.FileURLMap(p.pageSupplier.FilenameURLs())
	if err != nil {
		return
	}
	err = engine.IncludeDirs([]string{filepath.Join(p.sourceDir, p.config.IncludesDir)})
	return
}

func (p *Pipeline) makeLiquidEngine(o PipelineOptions) (liquid.Engine, error) {
	if o.UseRemoteLiquidEngine {
		return p.makeLiquidClient()
	}
	return p.makeLocalLiquidEngine(), nil
}
