package utils

import yaml "gopkg.in/yaml.v2"

// UnmarshalYAMLInterface is a wrapper for yaml.Unmarshall that
// knows how to unmarshal maps and lists.
func UnmarshalYAMLInterface(b []byte, i *interface{}) error {
	var m map[interface{}]interface{}
	err := yaml.Unmarshal(b, &m)
	switch err.(type) {
	case *yaml.TypeError:
		// Work around https://github.com/go-yaml/yaml/issues/20
		var s []interface{}
		err = yaml.Unmarshal(b, &s)
		if err != nil {
			return err
		}
		*i = s
	default:
		if err != nil {
			return err
		}
		*i = m
	}
	return nil
}
