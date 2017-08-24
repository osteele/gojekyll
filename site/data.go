package site

import (
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/osteele/gojekyll/utils"
)

func (s *Site) readDataFiles() error {
	s.data = map[string]interface{}{}
	dataDir := filepath.Join(s.SourceDir(), s.cfg.DataDir)
	files, err := ioutil.ReadDir(dataDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	for _, f := range files {
		if f.IsDir() {
			break
		}
		var (
			filename  = filepath.Join(dataDir, f.Name())
			basename  = utils.TrimExt(filepath.Base(f.Name()))
			data, err = readDataFile(filename)
		)
		if err != nil {
			return utils.WrapPathError(err, filename)
		}
		if data != nil {
			s.data[basename] = data
		}
	}
	return nil
}

func readDataFile(filename string) (interface{}, error) {
	switch filepath.Ext(filename) {
	case ".csv":
		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer f.Close() // nolint: errcheck
		r := csv.NewReader(f)
		return r.ReadAll()
	case ".json":
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		var d interface{}
		err = json.Unmarshal(b, &d)
		return d, err
	case ".yaml", ".yml":
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		var d interface{}
		err = utils.UnmarshalYAMLInterface(b, &d)
		return d, err
	}
	return nil, nil
}
