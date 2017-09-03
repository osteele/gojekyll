package pages

import (
	"path"
	"path/filepath"

	"github.com/osteele/gojekyll/utils"
	"github.com/osteele/liquid"
)

// ToLiquid is part of the liquid.Drop interface.
func (d *StaticFile) ToLiquid() interface{} {
	return liquid.IterationKeyedMap(map[string]interface{}{
		"name":          path.Base(d.relPath),
		"basename":      utils.TrimExt(path.Base(d.relPath)),
		"path":          d.URL(),
		"modified_time": d.modTime,
		"extname":       d.OutputExt(),
		// de facto:
		"collection": nil,
	})
}

func (f *file) ToLiquid() interface{} {
	var (
		relpath = "/" + filepath.ToSlash(f.relPath)
		base    = path.Base(relpath)
		ext     = path.Ext(relpath)
	)
	return liquid.IterationKeyedMap(f.fm.Merged(FrontMatter{
		"path":          relpath,
		"modified_time": f.modTime,
		"name":          base,
		"basename":      utils.TrimExt(base),
		"extname":       ext,
	}))
}

// ToLiquid is in the liquid.Drop interface.
func (p *page) ToLiquid() interface{} {
	var (
		fm          = p.fm
		relpath     = p.relPath
		siteRelPath = filepath.ToSlash(p.site.RelativePath(p.filename))
		ext         = filepath.Ext(relpath)
	)
	data := map[string]interface{}{
		"categories":    p.Categories(),
		"content":       p.maybeContent(),
		"date":          fm.Get("date", p.modTime),
		"excerpt":       p.Excerpt(),
		"id":            utils.TrimExt(p.URL()),
		"path":          siteRelPath,
		"relative_path": siteRelPath,
		"slug":          fm.String("slug", utils.Slugify(utils.TrimExt(filepath.Base(p.relPath)))),
		"tags":          p.Tags(),
		"url":           p.URL(),

		// de facto
		"ext": ext,
	}
	for k, v := range p.fm {
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
