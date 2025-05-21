package renderers

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/osteele/gojekyll/cache"
	"github.com/osteele/gojekyll/utils"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"

	sass "github.com/bep/godartsass/v2"
)

const sassMIMEType = "text/css"
const sassDirName = "_sass"

// copySASSFileIncludes copies sass partials into a temporary directory,
// removing initial underscores.
func (p *Manager) copySASSFileIncludes() error {
	// TODO delete the temp directory when done?
	// TODO use libsass.ImportsOption instead?
	// FIXME this doesn't delete stale css files
	if err := p.makeSASSTempDir(); err != nil {
		return err
	}
	h := md5.New()
	if p.ThemeDir != "" {
		if err := p.copySASSFiles(filepath.Join(p.ThemeDir, sassDirName), p.sassTempDir, h); err != nil {
			return err
		}
	}
	if err := p.copySASSFiles(filepath.Join(p.sourceDir(), p.cfg.Sass.Dir), p.sassTempDir, h); err != nil {
		return err
	}
	p.sassHash = fmt.Sprintf("%x", h.Sum(nil))
	return nil
}

func (p *Manager) makeSASSTempDir() error {
	if p.sassTempDir == "" {
		dir, err := os.MkdirTemp(os.TempDir(), "_sass")
		if p.cfg.Verbose {
			fmt.Println("create", dir)
		}
		if err != nil {
			return err
		}
		p.sassTempDir = dir
	}
	return nil
}

func (p *Manager) copySASSFiles(src, dst string, h io.Writer) error {
	if p.cfg.Verbose {
		fmt.Printf("copy sass directory %s to %s\n", src, dst)
	}
	err := filepath.Walk(src, func(from string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		rel := utils.MustRel(src, from)
		to := filepath.Join(dst, strings.TrimPrefix(rel, "_"))
		if p.cfg.Verbose {
			fmt.Printf("copy sass file %s to %s\n", src, to)
		}
		in, err := os.Open(from)
		if err != nil {
			return err
		}
		defer in.Close() // nolint: errcheck
		if _, err = fmt.Fprintf(h, "--- sass file: %s ---\n", rel); err != nil {
			return err
		}
		if _, err = io.Copy(h, in); err != nil {
			return err
		}
		return utils.CopyFileContents(to, from, 0644)
	})
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

// SassIncludePaths returns an array of sass include directories.
func (p *Manager) SassIncludePaths() []string {
	return []string{p.sassTempDir}
}

// string filters
var comp, compErr = sass.Start(sass.Options{})

// WriteSass converts a SASS file and writes it to w.
func (p *Manager) WriteSass(w io.Writer, b []byte) error {
	s, err := cache.WithFile(fmt.Sprintf("sass: %s", p.sassHash), string(b), func() (s string, err error) {
		if compErr != nil {
			return "", compErr
		}
		res, err := comp.Execute(sass.Args{
			Source:       string(b),
			IncludePaths: p.SassIncludePaths(),
		})
		if err != nil {
			return "", err
		}
		m := minify.New()
		m.AddFunc(sassMIMEType, css.Minify)
		min := bytes.NewBuffer(make([]byte, 0, len(res.CSS)))
		if err := m.Minify(sassMIMEType, min, bytes.NewBuffer([]byte(res.CSS))); err != nil {
			return "", err
		}
		return min.String(), nil
	})
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, s)
	return err
}
