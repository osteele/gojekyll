package sites

import (
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/osteele/gojekyll/config"
	"github.com/osteele/gojekyll/helpers"
	"github.com/osteele/gojekyll/liquid"
	"github.com/osteele/gojekyll/pages"
	"github.com/osteele/gojekyll/templates"
	"github.com/russross/blackfriday"
)

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
	Pages() []pages.Page
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

// OutputExt returns the output extension.
func (s *Site) OutputExt(pathname string) string {
	return s.config.OutputExt(pathname)
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

// ApplyLayout applies the named layout to bytes.
func (p *Pipeline) ApplyLayout(name string, b []byte, e templates.VariableMap) ([]byte, error) {
	for name != "" {
		var lfm templates.VariableMap
		t, err := p.FindLayout(name, &lfm)
		if err != nil {
			return nil, err
		}
		le := templates.MergeVariableMaps(e, templates.VariableMap{
			"content": string(b),
			"layout":  lfm,
		})
		b, err = t.Render(le)
		if err != nil {
			return nil, err
		}
		name = lfm.String("layout", "")
	}
	return b, nil
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
	urls := map[string]string{}
	for _, page := range p.pageSupplier.Pages() {
		urls[page.SiteRelPath()] = page.Permalink()
	}
	err = engine.FileURLMap(urls)
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
