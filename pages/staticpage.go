package pages

import (
	"io"
	"io/ioutil"

	"github.com/osteele/gojekyll/templates"
)

// StaticPage is a static page.
type StaticPage struct {
	pageFields
}

// Static returns a bool indicating that the page is a static page.
func (p *StaticPage) Static() bool { return true }

// Variables returns metadata for use in the representation of the page as a collection item
func (p *StaticPage) Variables() templates.VariableMap {
	return templates.MergeVariableMaps(p.frontMatter, p.pageFields.Variables())
}

func (p *StaticPage) Write(_ Context, w io.Writer) error {
	b, err := ioutil.ReadFile(p.filename)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}
