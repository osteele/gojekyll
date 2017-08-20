package pages

import (
	"path"
	"path/filepath"

	"github.com/osteele/gojekyll/frontmatter"
	"github.com/osteele/gojekyll/utils"
	"github.com/osteele/liquid"
)

// ToLiquid is part of the liquid.Drop interface.
func (d *StaticFile) ToLiquid() interface{} {
	return liquid.IterationKeyedMap(map[string]interface{}{
		"name":          path.Base(d.relpath),
		"basename":      utils.TrimExt(path.Base(d.relpath)),
		"path":          d.Permalink(),
		"modified_time": d.fileModTime,
		"extname":       d.OutputExt(),
		// de facto:
		"collection": nil,
	})
}

func (f *file) ToLiquid() interface{} {
	var (
		relpath = "/" + filepath.ToSlash(f.relpath)
		base    = path.Base(relpath)
		ext     = path.Ext(relpath)
	)
	return liquid.IterationKeyedMap(f.frontMatter.Merged(frontmatter.FrontMatter{
		"path":          relpath,
		"modified_time": f.fileModTime,
		"name":          base,
		"basename":      utils.TrimExt(base),
		"extname":       ext,
	}))
}

// ToLiquid is in the liquid.Drop interface.
func (p *page) ToLiquid() interface{} {
	var (
		fm      = p.frontMatter
		relpath = p.relpath
		ext     = filepath.Ext(relpath)
	)
	data := map[string]interface{}{
		"content": p.maybeContent(),
		"excerpt": p.Excerpt(),
		"path":    relpath,
		"url":     p.Permalink(),
		"slug":    fm.String("slug", utils.Slugify(utils.TrimExt(filepath.Base(p.relpath)))),
		// "output": // TODO; includes layouts

		// TODO documented as present in all pages, but de facto only defined for collection pages
		"id":         utils.TrimExt(p.Permalink()),
		"categories": p.Categories(),
		"tags":       p.Tags(),
		// "title": base,

		// TODO Only present in collection pages https://jekyllrb.com/docs/collections/#documents
		"relative_path": filepath.ToSlash(p.site.RelativePath(p.filename)),
		// TODO collection(name)

		// de facto
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
	return liquid.IterationKeyedMap(data)
}

func (p *page) maybeContent() interface{} {
	p.RLock()
	defer p.RUnlock()
	if p.rendered {
		return p.content
	}
	return p.raw
}
