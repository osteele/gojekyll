package commands

import (
	"fmt"
	"os"
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

	logger.path("Destination:", site.DestDir())
	logger.label("Generating...", "")
	count, err := site.Build()
	switch {
	case err == nil:
		elapsed := time.Since(commandStartTime)
		logger.label("", "wrote %d files in %.2fs.", count, elapsed.Seconds())
	case watch:
		fmt.Fprintln(os.Stderr, err)
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
		logger.label("Auto-regeneration:", "enabled for %q", site.SourceDir())
		for event := range events {
			fmt.Print(event)
		}
	} else {
		logger.label("Auto-regeneration:", "disabled. Use --watch to enable.")
	}
	return nil
}
