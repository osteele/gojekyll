package commands

import (
	"github.com/osteele/gojekyll/plugins"
)

var pluginsApp = app.Command("plugins", "List emulated plugins")

func pluginsCommand() error {
	bannerLog.label("Plugins:", "")
	for _, name := range plugins.Names() {
		log.Printf("  %s\n", name)
	}
	log.Println("\nhttps://github.com/osteele/gojekyll/blob/master/docs/plugins.md describes plugin implementation status.")
	log.Println("(This may not accurately describe your installation, if you are running an older version of gojekyll.)")
	return nil
}
