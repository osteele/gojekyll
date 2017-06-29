package pipelines

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/osteele/gojekyll/helpers"

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

	src := filepath.Join(p.SourceDir, "_sass")
	dst := p.sassTempDir
	err := filepath.Walk(src, func(from string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		rel, err := filepath.Rel(src, from)
		if err != nil {
			return err
		}
		to := filepath.Join(dst, strings.TrimPrefix(rel, "_"))
		return helpers.CopyFileContents(to, from, 0644)
	})
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

// SassIncludePaths returns an array of sass include directories.
func (p *Pipeline) SassIncludePaths() []string {
	return []string{p.sassTempDir}
}

// WriteSass converts a SASS file and writes it to w.
func (p *Pipeline) WriteSass(w io.Writer, b []byte) error {
	comp, err := libsass.New(w, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	err = comp.Option(libsass.IncludePaths(p.SassIncludePaths()))
	if err != nil {
		log.Fatal(err)
	}
	return comp.Run()
}
