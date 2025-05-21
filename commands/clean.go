package commands

import "github.com/osteele/gojekyll/site"

var clean = app.Command("clean", "Clean the site (removes site output) without building.")

func cleanCommand(site *site.Site) error {
	logger.label("Cleaner:", "Removing %s...", site.DestDir())
	return site.Clean()
}
