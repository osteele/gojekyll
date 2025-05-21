package cache

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

var enabled = true
var cacheMx sync.Mutex

func init() {
	s := os.Getenv("GOJEKYLL_DISABLE_CACHE")
	if s != "" && s != "0" && s != "false" {
		enabled = false
	}
}

func cacheDir() string {
	return filepath.Join(os.TempDir(), os.ExpandEnv("gojekyll-$USER"))
}

// Clear clears the cache. It's used for testing.
func Clear() error {
	return os.RemoveAll(cacheDir())
}

// Enable enables the cache; for testing.
func Enable() {
	enabled = true
}

// Disable disables the cache; for testing.
func Disable() {
	enabled = false
}

// WithFile looks (header, content) up in a user-specific file cache.
// If found, it writes the file contents. Else it calls fn to write to
// both the writer and the file system.
//
// header and content are distinct parameters to relieve the caller from
// having to concatenate them.
func WithFile(header string, content string, fn func() (string, error)) (string, error) {
	h := md5.New()
	io.WriteString(h, content) // nolint: errcheck
	io.WriteString(h, "\n")    // nolint: errcheck
	io.WriteString(h, header)  // nolint: errcheck
	sum := h.Sum(nil)

	// don't use ioutil.TempDir, because we want this to last across invocations
	cachedir := cacheDir()
	cachefile := filepath.Join(cachedir, fmt.Sprintf("%x%c%x", sum[:1], filepath.Separator, sum[1:]))

	// ignore errors; if there's a missing file we don't care, and if it's
	// another error we'll pick it up during write.
	//
	// WriteFile truncates the file before writing it, so ignore empty files.
	// If the writer actually wrote an empty file, we'll end up gratuitously
	// re-running it, which is okay.
	//
	// Do as much work as possible before checking if the cache is enabled, to
	// minimize code paths and timing differences.
	if b, err := os.ReadFile(cachefile); err == nil && len(b) > 0 && enabled {
		return string(b), err
	}
	s, err := fn()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(cachefile), 0700); err != nil {
		return "", err
	}
	defer cacheMx.Unlock()
	cacheMx.Lock()
	if err := os.WriteFile(cachefile, []byte(s), 0600); err != nil {
		return "", err
	}
	return s, nil
}
