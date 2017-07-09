package pipelines

import (
	"bytes"
	"crypto/md5" // nolint: gas
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/dchest/cssmin"
	"github.com/osteele/gojekyll/cache"
	"github.com/osteele/gojekyll/utils"

	libsass "github.com/wellington/go-libsass"
)

// CopySassFileIncludes copies sass partials into a temporary directory,
// removing initial underscores.
// TODO delete the temp directory when done
func (p *Pipeline) CopySassFileIncludes() error {
	// TODO use libsass.ImportsOption instead?
	if p.sassTempDir == "" {
		dir, err := ioutil.TempDir(os.TempDir(), "_sass")
		if err != nil {
			return err
		}
		p.sassTempDir = dir
	}

	h := md5.New() // nolint: gas, noncrypto
	src := filepath.Join(p.SourceDir(), "_sass")
	dst := p.sassTempDir
	err := filepath.Walk(src, func(from string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		rel := utils.MustRel(src, from)
		to := filepath.Join(dst, strings.TrimPrefix(rel, "_"))
		in, err := os.Open(from)
		if err != nil {
			return err
		}
		defer in.Close() // nolint: errcheck
		_, err = io.Copy(h, in)
		if err != nil {
			return err
		}
		return utils.CopyFileContents(to, from, 0644)
	})
	if os.IsNotExist(err) {
		return nil
	}
	p.sassHash = fmt.Sprintf("%x", h.Sum(nil))
	return err
}

// SassIncludePaths returns an array of sass include directories.
func (p *Pipeline) SassIncludePaths() []string {
	return []string{p.sassTempDir}
}

// WriteSass converts a SASS file and writes it to w.
func (p *Pipeline) WriteSass(w io.Writer, b []byte) error {
	s, err := cache.WithFile(fmt.Sprintf("sass: %s", p.sassHash), string(b), func() (s string, err error) {
		buf := new(bytes.Buffer)
		comp, err := libsass.New(buf, bytes.NewBuffer(b))
		if err != nil {
			return "", err
		}
		err = comp.Option(libsass.IncludePaths(p.SassIncludePaths()))
		if err != nil {
			return "", err
		}
		if err = comp.Run(); err != nil {
			return "", err
		}
		b = cssmin.Minify(buf.Bytes())
		return string(b), nil
	})
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(s))
	return err
}
