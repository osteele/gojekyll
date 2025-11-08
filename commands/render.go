package commands

import (
	"os"
	"path/filepath"

	"github.com/osteele/gojekyll/site"
)

var render = app.Command("render", "Render a file or URL path to standard output")
var renderPath = render.Arg("PATH", "Path or URL").String()

func renderCommand(site *site.Site) error {
	p, err := pageFromPathOrRoute(site, *renderPath)
	if err != nil {
		return err
	}
	bannerLog.path("Render:", filepath.Join(site.SourceDir(), p.Source()))
	//nolint:govet
	bannerLog.label("URL:", p.URL())
	bannerLog.label("Content:", "")
	return site.WriteDocument(os.Stdout, p)
}
