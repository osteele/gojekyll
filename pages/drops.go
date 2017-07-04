package pages

import (
	"path"
	"path/filepath"

	"github.com/osteele/gojekyll/helpers"
	"github.com/osteele/gojekyll/templates"
)

// ToLiquid returns the attributes of the template page object.
// See https://jekyllrb.com/docs/variables/#page-variables
func (f *file) ToLiquid() interface{} {
	var (
		relpath = "/" + filepath.ToSlash(f.relpath)
		base    = path.Base(relpath)
		ext     = path.Ext(relpath)
	)

	return templates.MergeVariableMaps(f.frontMatter, map[string]interface{}{
		"path":          relpath,
		"modified_time": f.fileModTime,
		"name":          base,
		"basename":      helpers.TrimExt(base),
		"extname":       ext,
	})
}

// ToLiquid is in the liquid.Drop interface.
func (p *page) ToLiquid() interface{} {
	var (
		relpath = p.relpath
		ext     = filepath.Ext(relpath)
		root    = helpers.TrimExt(p.relpath)
		base    = filepath.Base(root)
	)

	data := map[string]interface{}{
		"path": relpath,
		"url":  p.Permalink(),
		// TODO output

		// not documented, but present in both collection and non-collection pages
		"permalink": p.Permalink(),

		// TODO only in non-collection pages:
		// TODO dir
		// TODO name
		// TODO next previous

		// TODO Documented as present in all pages, but de facto only defined for collection pages
		"id":    base,
		"title": base, // TODO capitalize
		// TODO excerpt category? categories tags
		// TODO slug
		"categories": p.Categories(),
		"tags":       p.Tags(),

		// TODO Only present in collection pages https://jekyllrb.com/docs/collections/#documents
		"relative_path": p.Path(),
		// TODO collection(name)

		// TODO undocumented; only present in collection pages:
		"ext": ext,
	}
	for k, v := range p.frontMatter {
		switch k {
		// doc implies these aren't present, but they appear to be present in a collection page:
		// case "layout", "published":
		case "permalink":
		// omit this, in order to use the value above
		default:
			data[k] = v
		}
	}
	if p.content != nil {
		data["content"] = string(*p.content)
		// TODO excerpt
	}
	return data
}

// MarshalYAML is part of the yaml.Marshaler interface
// The variables subcommand uses this.
func (f *file) MarshalYAML() (interface{}, error) {
	return f.ToLiquid(), nil
}

// MarshalYAML is part of the yaml.Marshaler interface
// The variables subcommand uses this.
func (p *page) MarshalYAML() (interface{}, error) {
	return p.ToLiquid(), nil
}
