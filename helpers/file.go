package helpers

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
)

// VisitCreatedFile calls os.Create to create a file, and applies w to it.
func VisitCreatedFile(name string, w func(io.Writer) error) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	close := true
	defer func() {
		if close {
			_ = f.Close()
		}
	}()
	if err := w(f); err != nil {
		return err
	}
	close = false
	return f.Close()
}

// CopyFileContents copies the file contents from src to dst.
// It's not atomic and doesn't copy permissions or metadata.
func CopyFileContents(dst, src string, perm os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close() // nolint: errcheck, gas
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err = io.Copy(out, in); err != nil {
		_ = os.Remove(dst) // nolint: gas
		return err
	}
	return out.Close()
}

// ReadFileMagic returns the first four bytes of the file, with final '\r' replaced by '\n'.
func ReadFileMagic(p string) (data []byte, err error) {
	f, err := os.Open(p)
	if err != nil {
		return
	}
	defer f.Close() // nolint: errcheck
	data = make([]byte, 4)
	_, err = f.Read(data)
	if data[3] == '\r' {
		data[3] = '\n'
	}
	return
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
