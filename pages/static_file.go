package pages

import (
	"io"
	"os"
)

// StaticFile is a static file.
type StaticFile struct {
	file
}

// Static is in the File interface.
func (p *StaticFile) Static() bool { return true }

// Write returns a bool indicating that the page is a static page.
func (p *StaticFile) Write(w io.Writer, _ RenderingContext) error {
	in, err := os.Open(p.filename)
	if err != nil {
		return err
	}
	defer in.Close() // nolint: errcheck, gas
	_, err = io.Copy(w, in)
	return err
}
