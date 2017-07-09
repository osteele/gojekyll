package site

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/osteele/gojekyll/utils"
)

func (s *Site) readDataFiles() error {
	s.data = map[string]interface{}{}
	dataDir := filepath.Join(s.SourceDir(), s.config.DataDir)
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
			err = utils.UnmarshalYAMLInterface(b, &d)
			if err != nil {
				return utils.PathError(err, "read YAML", filename)
			}
			basename := utils.TrimExt(filepath.Base(f.Name()))
			s.data[basename] = d
		}
	}
	return nil
}
