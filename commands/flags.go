package commands

import (
	kingpin "github.com/alecthomas/kingpin/v2"
	"github.com/osteele/gojekyll/config"
)

// Command-line options
var (
	options config.Flags
	profile = false
	quiet   = false
)

var (
	app         = kingpin.New("gojekyll", "a (somewhat) Jekyll-compatible blog generator")
	source      = app.Flag("source", "Source directory").Short('s').Default(".").ExistingDir()
	_           = app.Flag("destination", "Destination directory").Short('d').Action(stringVar("destination", &options.Destination)).String()
	_           = app.Flag("drafts", "Render posts in the _drafts folder").Short('D').Action(boolVar("drafts", &options.Drafts)).Bool()
	_           = app.Flag("future", "Publishes posts with a future date").Action(boolVar("future", &options.Future)).Bool()
	_           = app.Flag("unpublished", "Render posts that were marked as unpublished").Action(boolVar("unpublished", &options.Unpublished)).Bool()
	_           = app.Flag("baseurl", "Serve the website from the given base URL").Action(stringVar("baseurl", &options.BaseURL)).String()
	versionFlag = app.Flag("version", "Print the name and version").Short('v').Bool()
)

func init() {
	app.Flag("config", "Custom configuration file").StringVar(&options.ConfigFile)
}

func init() {
	app.HelpFlag.Short('h')
	app.Flag("profile", "Create a Go pprof CPU profile").BoolVar(&profile)
	app.Flag("quiet", "Silence (some) output.").Short('q').BoolVar(&quiet)
	_ = app.Flag("verbose", "Print verbose output.").Short('V').Action(boolVar("verbose", &options.Verbose)).Bool()

	// these flags are just present on build and serve, but I don't see a DRY way to say this
	app.Flag("incremental", "Enable incremental rebuild.").Short('I').Action(boolVar("incremental", &options.Incremental)).Bool()
	app.Flag("force_polling", "Force watch to use polling").BoolVar(&options.ForcePolling)

	// --watch has different defaults for build and serve
	watchText := "Watch for changes and rebuild"
	build.Flag("watch", watchText).Short('w').BoolVar(&options.Watch)
	serve.Flag("watch", watchText).Short('w').Default("true").BoolVar(&options.Watch)
}
