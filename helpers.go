package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func getBool(m map[interface{}]interface{}, k string, defaultValue bool) bool {
	if val, found := m[k]; found {
		if v, ok := val.(bool); ok {
			return v
		}
	}
	return defaultValue
}

func getString(m map[interface{}]interface{}, k string, defaultValue string) string {
	if val, found := m[k]; found {
		if v, ok := val.(string); ok {
			return v
		}
	}
	return defaultValue
}

// alternative to http://left-pad.io
func leftPad(s string, n int) string {
	if n <= len(s) {
		return s
	}
	ws := make([]byte, n-len(s))
	for i := range ws {
		ws[i] = ' '
	}
	return string(ws) + s
}

func mergeMaps(a map[interface{}]interface{}, b map[interface{}]interface{}) map[interface{}]interface{} {
	result := map[interface{}]interface{}{}
	for k, v := range a {
		result[k] = v
	}
	for k, v := range b {
		result[k] = v
	}
	return result
}

// stringMap returns a string-indexed map with the same values as its argument.
// Non-strings keys are converted to strings.
func stringMap(m map[interface{}]interface{}) map[string]interface{} {
	result := map[string]interface{}{}
	for k, v := range m {
		stringer, ok := k.(fmt.Stringer)
		if ok {
			result[stringer.String()] = v
		} else {
			result[fmt.Sprintf("%v", k)] = v
		}
	}
	return result
}

func postfixWalk(path string, walkFn filepath.WalkFunc) error {
	if files, err := ioutil.ReadDir(path); err == nil {
		for _, stat := range files {
			if stat.IsDir() {
				if err = postfixWalk(filepath.Join(path, stat.Name()), walkFn); err != nil {
					return err
				}
			}
		}
	}

	info, err := os.Stat(path)
	return walkFn(path, info, err)
}

func removeEmptyDirectories(path string) error {
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}

		if info.IsDir() {
			if err := os.Remove(path); err != nil {
				if os.IsNotExist(err) {
					return nil
				}
				// TODO swallow the error if it's because the directory isn't
				// empty. This can happen if there's an entry in _config.keepfiles
				return err
			}
		}
		return nil
	}
	return postfixWalk(path, walkFn)
}

func stringArrayToMap(strings []string) map[string]bool {
	stringMap := map[string]bool{}
	for _, s := range strings {
		stringMap[s] = true
	}
	return stringMap
}
