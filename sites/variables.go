package sites

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/osteele/gojekyll/helpers"
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
	s.Variables = templates.MergeVariableMaps(s.config.Variables, templates.VariableMap{
		"data": data,
		// TODO read time from _config, if it's available
		"time": time.Now(),
		// TODO pages, posts, related_posts, static_files, html_pages, html_files, collections, data, documents, categories.CATEGORY, tags.TAG
	})
	return s.updateCollectionVariables(false)
}

// SetPageContentTemplateValues sets the site[collection][i].content
// template variables
func (s *Site) SetPageContentTemplateValues() error {
	return s.updateCollectionVariables(true)
}

func (s *Site) updateCollectionVariables(includeContent bool) error {
	for _, c := range s.Collections {
		v, err := c.TemplateVariable(s, includeContent)
		if err != nil {
			return err
		}
		s.Variables[c.Name] = v
	}
	return nil
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
