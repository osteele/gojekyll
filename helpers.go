package main

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
)

// VariableMap is a map of strings to interface values, for use in template processing.
type VariableMap map[string]interface{}

// copyFile copies from file src to dst. It's not atomic and doesn't copy permissions or metadata.
// This is sufficient for its use within this package.
func copyFile(dst, src string, perm os.FileMode) error {
	inf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer inf.Close() // nolint: errcheck, gas
	outf, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err = io.Copy(outf, inf); err != nil {
		_ = os.Remove(dst) // nolint: gas
		return err
	}
	return outf.Close()
}

// Bool returns m[k] if it's a bool; else defaultValue.
func (m VariableMap) Bool(k string, defaultValue bool) bool {
	if val, found := m[k]; found {
		if v, ok := val.(bool); ok {
			return v
		}
	}
	return defaultValue
}

// String returns m[k] if it's a string; else defaultValue.
func (m VariableMap) String(k string, defaultValue string) string {
	if val, found := m[k]; found {
		if v, ok := val.(string); ok {
			return v
		}
	}
	return defaultValue
}

// LeftPad left-pads s with spaces to n wide. It's an alternative to http://left-pad.io.
func LeftPad(s string, n int) string {
	if n <= len(s) {
		return s
	}
	ws := make([]byte, n-len(s))
	for i := range ws {
		ws[i] = ' '
	}
	return string(ws) + s
}

func mergeVariableMaps(ms ...VariableMap) VariableMap {
	result := VariableMap{}
	for _, m := range ms {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

// PostfixWalk is like filepath.Walk, but visits each directory after visiting its children instead of before.
// It does not implement SkipDir.
func PostfixWalk(path string, walkFn filepath.WalkFunc) error {
	if files, err := ioutil.ReadDir(path); err == nil {
		for _, stat := range files {
			if stat.IsDir() {
				if err = PostfixWalk(filepath.Join(path, stat.Name()), walkFn); err != nil {
					return err
				}
			}
		}
	}

	info, err := os.Stat(path)
	return walkFn(path, info, err)
}

// IsNotEmpty returns a boolean indicating whether the error is known to report that a directory is not empty.
func IsNotEmpty(err error) bool {
	if err, ok := err.(*os.PathError); ok {
		return err.Err.(syscall.Errno) == syscall.ENOTEMPTY
	}
	return false
}

// RemoveEmptyDirectories recursively removes empty directories.
func RemoveEmptyDirectories(path string) error {
	walkFn := func(path string, info os.FileInfo, err error) error {
		switch {
		case err != nil && os.IsNotExist(err):
			// It's okay to call this on a directory that doesn't exist.
			// It's also okay if another process removed a file during traversal.
			return nil
		case err != nil:
			return err
		case info.IsDir():
			err := os.Remove(path)
			switch {
			case err == nil:
				return nil
			case os.IsNotExist(err):
				return nil
			case IsNotEmpty(err):
				return nil
			default:
				return err
			}
		}
		return nil
	}
	return PostfixWalk(path, walkFn)
}

func stringArrayToMap(strings []string) map[string]bool {
	stringMap := map[string]bool{}
	for _, s := range strings {
		stringMap[s] = true
	}
	return stringMap
}
