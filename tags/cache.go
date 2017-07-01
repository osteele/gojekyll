package tags

import (
	"bytes"
	"crypto/md5" // nolint: gas
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// withFileCache looks (header, content) up in a user-specific file cache.
// If found, it writes the file contents. Else it calls fn to write to
// both the writer and the file system.
//
// header and content are distinct parameters to relieve the caller from
// having to concatenate them.
func withFileCache(w io.Writer, header string, content string, fn func(w io.Writer) error) error {
	h := md5.New()             // nolint: gas
	io.WriteString(h, content) // nolint: errcheck, gas
	io.WriteString(h, "\n")    // nolint: errcheck, gas
	io.WriteString(h, header)  // nolint: errcheck, gas
	sum := h.Sum(nil)

	// don't use ioutil.TempDir, because we want this to last across invocations
	dirname := filepath.Join(os.TempDir(), os.ExpandEnv("gojekyll-$USER"))
	cachepath := filepath.Join(dirname, fmt.Sprintf("%x%c%x", sum[:1], filepath.Separator, sum[1:]))

	// ignore errors; if there's a missing file we don't care, and if it's
	// another error we'll pick it up during write
	// WriteFile truncates the file before writing it, so ignore empty files.
	// If the writer actually wrote an empty file, we'll end up gratuitously
	// re-running it, which is okay.
	if b, err := ioutil.ReadFile(cachepath); err == nil && len(b) > 0 {
		_, err = w.Write(b)
		return err
	}

	buf := new(bytes.Buffer)
	if err := fn(buf); err != nil {
		return err
	}
	out := buf.Bytes()

	if err := os.MkdirAll(filepath.Dir(cachepath), 0700); err != nil {
		return err
	}
	if err := ioutil.WriteFile(cachepath, out, 0600); err != nil {
		return err
	}
	_, err := w.Write(out)
	return err
}
