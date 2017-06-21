package gojekyll

import (
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/osteele/gojekyll/helpers"
)

func (site *Site) initSiteVariables() error {
	data, err := site.readDataFiles()
	if err != nil {
		return err
	}
	site.Variables = MergeVariableMaps(site.Variables, VariableMap{
		"data": data,
		// TODO read time from _config, if it's available
		"time": time.Now(),
		// TODO pages, posts, related_posts, static_files, html_pages, html_files, collections, data, documents, categories.CATEGORY, tags.TAG
	})
	site.updateCollectionVariables()
	return nil
}

func (site *Site) updateCollectionVariables() {
	for _, c := range site.Collections {
		site.Variables[c.Name] = c.CollectionValue()
	}
}

func (site *Site) readDataFiles() (VariableMap, error) {
	data := VariableMap{}
	dataDir := filepath.Join(site.Source, site.config.DataDir)
	files, err := ioutil.ReadDir(dataDir)
	if err != nil {
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
