package site

import (
	"bytes"
	"io"

	"github.com/osteele/gojekyll/frontmatter"
	"github.com/osteele/gojekyll/pages"
	"github.com/osteele/gojekyll/plugins"
	"github.com/osteele/gojekyll/utils"
	"github.com/osteele/liquid"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/html"
)

// AddHTMLPage is in the plugins.Site interface.
func (s *Site) AddHTMLPage(url string, src string, fm frontmatter.FrontMatter) {
	tpl, err := s.TemplateEngine().ParseTemplate([]byte(src))
	if err != nil {
		panic(err)
	}
	d := &templateDoc{pages.PageEmbed{Path: url}, s, tpl}
	s.AddDocument(d, true)
}

func (s *Site) installPlugins() error {
	s.plugins = s.cfg.Plugins
	installed := utils.StringSet{}
	// Install plugins and call their ModifyPluginList methods.
	// Repeat until no plugins have been added.
	for len(s.plugins) > len(installed) {
		// Collect plugins into a list instead of map, in order to preserve order
		pending := utils.StringList(s.plugins).Reject(installed.Contains)
		if err := plugins.Install(pending, s); err != nil {
			return err
		}
		for _, name := range pending {
			p, ok := plugins.Lookup(name)
			if ok {
				s.plugins = p.ModifyPluginList(s.plugins)
			}
		}
		installed.AddStrings(pending)
	}
	return nil
}

func (s *Site) runHooks(h func(plugins.Plugin) error) error {
	for _, name := range s.plugins {
		p, ok := plugins.Lookup(name)
		if ok {
			if err := h(p); err != nil {
				return utils.WrapError(err, "running plugin")
			}
		}
	}
	return nil
}

type templateDoc struct {
	pages.PageEmbed
	site *Site
	tpl  *liquid.Template
}

func (d *templateDoc) Content() string {
	bindings := map[string]interface{}{"site": d.site}
	b, err := d.tpl.Render(bindings)
	if err != nil {
		panic(err)
	}
	m := minify.New()
	m.AddFunc("text/html", html.Minify)
	min := bytes.NewBuffer(make([]byte, 0, len(b)))
	if err := m.Minify("text/html", min, bytes.NewBuffer(b)); err != nil {
		panic(err)
	}
	return min.String()
}

func (d *templateDoc) Write(w io.Writer) error {
	_, err := io.WriteString(w, d.Content())
	return err
}
