package pages

import (
	"io"
	"io/ioutil"
)

// StaticPage is a static page.
type StaticPage struct {
	pageFields
}

// Static returns a bool indicating that the page is a static page.
func (p *StaticPage) Static() bool { return true }

func (p *StaticPage) Write(_ Context, w io.Writer) error {
	b, err := ioutil.ReadFile(p.filename)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}
