package gojekyll

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/osteele/gojekyll/helpers"
)

// SiteVariables returns the site variable for template evaluation.
func (s *Site) SiteVariables() VariableMap {
	return s.Variables
}

func (s *Site) initSiteVariables() error {
	data, err := s.readDataFiles()
	if err != nil {
		return err
	}
	s.Variables = MergeVariableMaps(s.Variables, VariableMap{
		"data": data,
		// TODO read time from _config, if it's available
		"time": time.Now(),
		// TODO pages, posts, related_posts, static_files, html_pages, html_files, collections, data, documents, categories.CATEGORY, tags.TAG
	})
	s.updateCollectionVariables()
	return nil
}

func (s *Site) updateCollectionVariables() {
	for _, c := range s.Collections {
		s.Variables[c.Name] = c.CollectionValue()
	}
}

func (s *Site) readDataFiles() (VariableMap, error) {
	data := VariableMap{}
	dataDir := filepath.Join(s.Source, s.config.DataDir)
	files, err := ioutil.ReadDir(dataDir)
	if err != nil {
		if os.IsNotExist(err) {
			return VariableMap{}, nil
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
