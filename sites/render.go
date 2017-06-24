package sites

import (
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/osteele/gojekyll/helpers"
	"github.com/osteele/gojekyll/liquid"
	"github.com/osteele/gojekyll/templates"
	"github.com/russross/blackfriday"
)

// OutputExt returns the output extension.
func (s *Site) OutputExt(pathname string) string {
	switch {
	case s.IsMarkdown(pathname):
		return ".html"
	case s.IsSassPath(pathname):
		return ".css"
	default:
		return filepath.Ext(pathname)
	}
}

// Render returns nil iff it wrote to the writer
func (s *Site) Render(w io.Writer, b []byte, filename string, e templates.VariableMap) ([]byte, error) {
	if s.IsSassPath(filename) {
		return nil, s.WriteSass(w, b)
	}
	b, err := s.renderTemplate(b, e, filename)
	if err != nil {
		return nil, err
	}
	if s.IsMarkdown(filename) {
		b = blackfriday.MarkdownCommon(b)
	}
	return b, nil
}

func (s *Site) renderTemplate(b []byte, e templates.VariableMap, filename string) ([]byte, error) {
	b, err := s.liquidEngine.ParseAndRender(b, e)
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
func (s *Site) ApplyLayout(name string, b []byte, e templates.VariableMap) ([]byte, error) {
	for name != "" {
		var lfm templates.VariableMap
		t, err := s.FindLayout(name, &lfm)
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

func (s *Site) makeLocalLiquidEngine() liquid.Engine {
	engine := liquid.NewLocalWrapperEngine()
	engine.LinkTagHandler(s.RelPathURL)
	includeHandler := func(name string, w io.Writer, scope map[string]interface{}) error {
		filename := filepath.Join(s.Source, s.config.IncludesDir, name)
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

func (s *Site) makeLiquidClient() (engine liquid.RemoteEngine, err error) {
	engine, err = liquid.NewRPCClientEngine(liquid.DefaultServer)
	if err != nil {
		return
	}
	urls := map[string]string{}
	for _, p := range s.Paths {
		urls[p.SiteRelPath()] = p.Permalink()
	}
	err = engine.FileURLMap(urls)
	if err != nil {
		return
	}
	err = engine.IncludeDirs([]string{filepath.Join(s.Source, s.config.IncludesDir)})
	return
}

func (s *Site) makeLiquidEngine() (liquid.Engine, error) {
	if s.UseRemoteLiquidEngine {
		return s.makeLiquidClient()
	}
	return s.makeLocalLiquidEngine(), nil
}

// TemplateEngine creates a liquid engine configured to with include paths and link tag resolution
// for this site.
func (s *Site) TemplateEngine() liquid.Engine {
	return s.liquidEngine
}
