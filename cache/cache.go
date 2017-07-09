package cache

import (
	"crypto/md5" // nolint: gas
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

var disableCache = false

func init() {
	s := os.Getenv("GOJEKYLL_DISABLE_CACHE")
	if s != "" && s != "0" && s != "false" {
		disableCache = true
	}
}

// WithFile looks (header, content) up in a user-specific file cache.
// If found, it writes the file contents. Else it calls fn to write to
// both the writer and the file system.
//
// header and content are distinct parameters to relieve the caller from
// having to concatenate them.
func WithFile(header string, content string, fn func() (string, error)) (string, error) {
	h := md5.New()             // nolint: gas, noncrypto
	io.WriteString(h, content) // nolint: errcheck, gas
	io.WriteString(h, "\n")    // nolint: errcheck, gas
	io.WriteString(h, header)  // nolint: errcheck, gas
	sum := h.Sum(nil)

	// don't use ioutil.TempDir, because we want this to last across invocations
	// cachedir := filepath.Join("/tmp", os.ExpandEnv("gojekyll-$USER"))
	cachedir := filepath.Join(os.TempDir(), os.ExpandEnv("gojekyll-$USER"))
	cachefile := filepath.Join(cachedir, fmt.Sprintf("%x%c%x", sum[:1], filepath.Separator, sum[1:]))

	// ignore errors; if there's a missing file we don't care, and if it's
	// another error we'll pick it up during write
	// WriteFile truncates the file before writing it, so ignore empty files.
	// If the writer actually wrote an empty file, we'll end up gratuitously
	// re-running it, which is okay.
	if b, err := ioutil.ReadFile(cachefile); err == nil && len(b) > 0 && !disableCache {
		return string(b), err
	}
	s, err := fn()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(cachefile), 0700); err != nil {
		return "", err
	}
	if err := ioutil.WriteFile(cachefile, []byte(s), 0600); err != nil {
		return "", err
	}
	return s, nil
}
