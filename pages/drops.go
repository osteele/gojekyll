package pages

import (
	"path"
	"path/filepath"

	"github.com/osteele/gojekyll/utils"
	"github.com/osteele/liquid"
)

type pageDropData struct {
	categories  []string
	tags        []string
	slug        string
	siteRelPath string
	id          string
}

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
		fm     = p.fm
		stable = p.stableDropData()
		ext    = filepath.Ext(p.relPath)
	)
	data := map[string]interface{}{
		"categories":    stable.categories,
		"content":       p.maybeContent(),
		"excerpt":       p.Excerpt(),
		"id":            stable.id,
		"path":          stable.siteRelPath,
		"relative_path": stable.siteRelPath,
		"slug":          stable.slug,
		"tags":          stable.tags,
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

func (p *page) stableDropData() pageDropData {
	p.dropOnce.Do(func() {
		p.dropData = pageDropData{
			categories:  p.Categories(),
			tags:        p.Tags(),
			slug:        p.fm.String("slug", utils.Slugify(utils.TrimExt(filepath.Base(p.relPath)))),
			siteRelPath: filepath.ToSlash(p.site.RelativePath(p.filename)),
			id:          utils.TrimExt(p.URL()),
		}
	})
	return p.dropData
}

func (p *page) maybeContent() interface{} {
	p.m.RLock()
	defer p.m.RUnlock()
	if p.rendered {
		return p.content
	}
	return p.raw
}
