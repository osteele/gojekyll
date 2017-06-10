package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// alternative to http://left-pad.io
func leftPad(s string, n int) string {
	ws := make([]byte, n)
	for i := range ws {
		ws[i] = ' '
	}
	return string(ws) + s
}

func postfixWalk(path string, walkFn filepath.WalkFunc) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, stat := range files {
		if stat.IsDir() {
			postfixWalk(filepath.Join(path, stat.Name()), walkFn)
		}
	}

	info, err := os.Stat(path)
	err = walkFn(path, info, err)
	if err != nil {
		return err
	}
	return nil
}

func removeEmptyDirectories(path string) error {
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		stat, err := os.Stat(path)
		if os.IsNotExist(err) {
			return nil
		}
		if err != nil {
			return nil
		}
		if stat.IsDir() {
			err = os.Remove(path)
			// TODO swallow the error if it's because the directory isn't
			// empty. This can happen if there's an entry in _config.keepfiles
		}
		return err
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
