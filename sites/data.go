package sites

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/osteele/gojekyll/helpers"
)

func (s *Site) readDataFiles() error {
	s.data = map[string]interface{}{}
	dataDir := filepath.Join(s.Source, s.config.DataDir)
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
		filename := filepath.Join(dataDir, f.Name())
		switch filepath.Ext(f.Name()) {
		case ".csv", ".json":
			return fmt.Errorf("unimplemented reading %s", filepath.Ext(f.Name()))
		case ".yaml", ".yml":
			b, err := ioutil.ReadFile(filename)
			if err != nil {
				return err
			}
			var d interface{} // map or slice
			err = helpers.UnmarshalYAMLInterface(b, &d)
			if err != nil {
				return helpers.PathError(err, "read YAML", filename)
			}
			basename := helpers.TrimExt(filepath.Base(f.Name()))
			s.data[basename] = d
		}
	}
	return nil
}
