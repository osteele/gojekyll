package commands

import (
	"time"

	"github.com/osteele/gojekyll/site"
)

// main sets this
var commandStartTime = time.Now()

var build = app.Command("build", "Build your site").Alias("b")

func init() {
	build.Flag("dry-run", "Dry run").Short('n').BoolVar(&options.DryRun)
}

func buildCommand(site *site.Site) error {
	watch := site.Config().Watch

	bannerLog.path("Destination:", site.DestDir())
	bannerLog.label("Generating...", "")
	count, err := site.Write()
	switch {
	case err == nil:
		elapsed := time.Since(commandStartTime)
		bannerLog.label("", "wrote %d files in %.2fs.", count, elapsed.Seconds())
	case watch:
		log.Error(err.Error())
	default:
		return err
	}

	// FIXME the watch will miss files that changed during the first build

	// server watch is implemented inside Server.Run, in contrast to this command
	if watch {
		events, err := site.WatchRebuild()
		if err != nil {
			return err
		}
		bannerLog.label("Auto-regeneration:", "enabled for %q", site.SourceDir())
		for event := range events {
			log.Printf("%s", event)
		}
	} else {
		bannerLog.label("Auto-regeneration:", "disabled. Use --watch to enable.")
	}
	return nil
}
