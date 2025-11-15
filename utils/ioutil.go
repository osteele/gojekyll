package utils

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"syscall"
)

// CopyFileContents copies the file contents from src to dst.
// It's not atomic and doesn't copy permissions or metadata.
func CopyFileContents(dst, src string, perm os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close() // nolint: errcheck
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err = io.Copy(out, in); err != nil {
		_ = os.Remove(dst)
		return err
	}
	return out.Close()
}

// ReadFileMagic returns the first four bytes of the file, with final '\r' replaced by '\n'.
func ReadFileMagic(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close() // nolint: errcheck
	b := make([]byte, 4)
	_, err = f.Read(b)
	if err != nil && err != io.EOF {
		return nil, err
	}
	// Normalize windows linefeeds. This function is used to
	// recognize frontmatter, so we only need to look at the fourth byte.
	if b[3] == '\r' {
		b[3] = '\n'
	}
	return b, nil
}

// PostfixWalk is like filepath.Walk, but visits each directory after visiting its children instead of before.
// It does not implement SkipDir.
func PostfixWalk(root string, walkFn filepath.WalkFunc) error {
	if files, err := os.ReadDir(root); err == nil {
		for _, stat := range files {
			if stat.IsDir() {
				if err = PostfixWalk(filepath.Join(root, stat.Name()), walkFn); err != nil {
					return err
				}
			}
		}
	}
	info, err := os.Stat(root)
	return walkFn(root, info, err)
}

// IsNotEmpty returns a boolean indicating whether the error is known to report that a directory is not empty.
func IsNotEmpty(err error) bool {
	if pathErr, ok := err.(*os.PathError); ok {
		if errno, ok := pathErr.Err.(syscall.Errno); ok {
			// Check for platform-specific "directory not empty" errors
			// Unix/Linux/macOS: ENOTEMPTY
			// Windows: ERROR_DIR_NOT_EMPTY (checked via platform-specific helper)
			return errno == syscall.ENOTEMPTY || isWindowsDirNotEmpty(errno)
		}
	}
	return false
}

// NewPathError returns an os.PathError that formats as the given text.
func NewPathError(op, name, text string) *os.PathError {
	return &os.PathError{Op: op, Path: name, Err: errors.New(text)}
}

// RemoveEmptyDirectories recursively removes empty directories.
// It will not remove the root directory itself, only subdirectories within it.
func RemoveEmptyDirectories(root string) error {
	// Get the absolute path to ensure consistent comparison
	rootAbs, absErr := filepath.Abs(root)
	if absErr != nil {
		return absErr
	}

	// Normalize the root path
	root = filepath.Clean(root)
	rootAbs = filepath.Clean(rootAbs)

	walkFn := func(name string, info os.FileInfo, err error) error {
		cleanName := filepath.Clean(name)
		cleanNameAbs, _ := filepath.Abs(name)
		cleanNameAbs = filepath.Clean(cleanNameAbs)

		switch {
		case err != nil && os.IsNotExist(err):
			// It's okay to call this on a directory that doesn't exist.
			// It's also okay if another process removed a file during traversal.
			return nil
		case err != nil:
			return err
		case info.IsDir():
			// Don't remove the root directory itself, even if it's empty
			// This preserves symlinks and prevents the destination directory from being removed
			// We compare using absolute paths to handle both relative and absolute path cases
			if cleanName == root || cleanNameAbs == rootAbs {
				return nil
			}
			err := os.Remove(name)
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
	return PostfixWalk(root, walkFn)
}

// VisitCreatedFile calls os.Create to create a file, and applies w to it.
func VisitCreatedFile(filename string, w func(io.Writer) error) (err error) {
	f, err := os.Create(filename)
	if err != nil {
		return
	}
	defer func() {
		if e := f.Close(); e != nil && err == nil {
			err = e
		}
	}()
	err = w(f)
	return
}
