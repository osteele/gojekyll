package pages

import (
	"io"
	"os"
)

// StaticPage is a static page.
type StaticPage struct {
	pageFields
}

// Static returns a bool indicating that the page is a static page.
func (p *StaticPage) Static() bool { return true }

func (p *StaticPage) Write(_ RenderingContext, w io.Writer) error {
	in, err := os.Open(p.filename)
	if err != nil {
		return err
	}
	defer in.Close() // nolint: errcheck, gas
	_, err = io.Copy(w, in)
	return err
}
