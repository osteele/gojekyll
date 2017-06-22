package gojekyll

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/osteele/gojekyll/helpers"
	"github.com/osteele/gojekyll/pages"
	"github.com/osteele/gojekyll/templates"
)

// SiteVariables returns the site variable for template evaluation.
func (s *Site) SiteVariables() templates.VariableMap {
	return s.Variables
}

func (s *Site) initSiteVariables() error {
	data, err := s.readDataFiles()
	if err != nil {
		return err
	}
	s.Variables = templates.MergeVariableMaps(s.Variables, templates.VariableMap{
		"data": data,
		// TODO read time from _config, if it's available
		"time": time.Now(),
		// TODO pages, posts, related_posts, static_files, html_pages, html_files, collections, data, documents, categories.CATEGORY, tags.TAG
	})
	s.updateCollectionVariables()
	return nil
}

func (s *Site) SetPageContentTemplateValues() error {
	for _, c := range s.Collections {
		for _, p := range c.Pages() {
			switch p := p.(type) {
			case *pages.DynamicPage:
				if err := p.ComputeContent(c.Site); err != nil {
					return err
				}
			}
		}
	}
	s.updateCollectionVariables()
	return nil
}

func (s *Site) updateCollectionVariables() {
	for _, c := range s.Collections {
		s.Variables[c.Name] = c.TemplateVariable()
	}
}

func (s *Site) readDataFiles() (templates.VariableMap, error) {
	data := templates.VariableMap{}
	dataDir := filepath.Join(s.Source, s.config.DataDir)
	files, err := ioutil.ReadDir(dataDir)
	if err != nil {
		if os.IsNotExist(err) {
			return templates.VariableMap{}, nil
		}
		return nil, err
	}
	for _, f := range files {
		if f.IsDir() {
			break
		}
		filename := filepath.Join(dataDir, f.Name())
		switch filepath.Ext(f.Name()) {
		case ".yaml", ".yml":
			b, err := ioutil.ReadFile(filename)
			if err != nil {
				return nil, err
			}
			var d interface{} // map or slice
			err = helpers.UnmarshalYAMLInterface(b, &d)
			if err != nil {
				return nil, helpers.PathError(err, "read YAML", filename)
			}
			basename := helpers.TrimExt(filepath.Base(f.Name()))
			data[basename] = d
		}
	}
	return data, nil
}
