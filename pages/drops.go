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
	// In Jekyll, page.date is only defined for posts and collection documents.
	// For regular pages, it's only present if explicitly set in frontmatter.
	if _, hasDate := fm["date"]; hasDate {
		data["date"] = fm["date"]
	} else if p.IsPost() {
		data["date"] = p.modTime
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
	p.m.RLock()
	defer p.m.RUnlock()
	if p.rendered {
		return p.content
	}
	return p.raw
}
