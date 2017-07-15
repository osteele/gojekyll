package cmd

import (
	"github.com/osteele/gojekyll/config"
	"github.com/osteele/gojekyll/site"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

// Command-line options
var (
	buildOptions site.BuildOptions
	watch        bool
	configFlags  = config.Flags{Watch: &watch}
	profile      = false
	quiet        = false
)

var (
	app         = kingpin.New("gojekyll", "a (somewhat) Jekyll-compatible blog generator")
	source      = app.Flag("source", "Source directory").Short('s').Default(".").ExistingDir()
	_           = app.Flag("destination", "Destination directory").Short('d').Action(stringVar("destination", &configFlags.Destination)).String()
	_           = app.Flag("drafts", "Render posts in the _drafts folder").Short('D').Action(boolVar("drafts", &configFlags.Drafts)).Bool()
	_           = app.Flag("future", "Publishes posts with a future date").Action(boolVar("future", &configFlags.Future)).Bool()
	_           = app.Flag("unpublished", "Render posts that were marked as unpublished").Action(boolVar("unpublished", &configFlags.Unpublished)).Bool()
	versionFlag = app.Flag("version", "Print the name and version").Short('v').Bool()

	build = app.Command("build", "Build your site").Alias("b")
	clean = app.Command("clean", "Clean the site (removes site output) without building.")

	benchmark = app.Command("benchmark", "Repeat build for ten seconds. Implies --profile.")

	render     = app.Command("render", "Render a file or URL path to standard output")
	renderPath = render.Arg("PATH", "Path or URL").String()

	routes        = app.Command("routes", "Display site permalinks and associated files")
	dynamicRoutes = routes.Flag("dynamic", "Only show routes to non-static files").Bool()

	serve = app.Command("serve", "Serve your site locally").Alias("server").Alias("s")
	open  = serve.Flag("open-url", "Launch your site in a browser").Short('o').Bool()
	_     = serve.Flag("host", "Host to bind to").Short('H').Action(stringVar("host", &configFlags.Host)).String()
	_     = serve.Flag("port", "Port to listen on").Short('P').Action(intVar("port", &configFlags.Port)).Int()

	variables    = app.Command("variables", "Display a file or URL path's variables").Alias("v").Alias("var").Alias("vars")
	variablePath = variables.Arg("PATH", "Path, URL, site, or site...").String()

	versionCmd = app.Command("version", "Print the name and version")
)

func init() {
	app.HelpFlag.Short('h')
	app.Flag("profile", "Create a Go pprof CPU profile").BoolVar(&profile)
	app.Flag("quiet", "Silence (some) output.").Short('q').BoolVar(&quiet)
	build.Flag("dry-run", "Dry run").Short('n').BoolVar(&buildOptions.DryRun)

	watchText := "Watch for changes and rebuild"
	build.Flag("watch", watchText).Short('w').BoolVar(&watch)
	serve.Flag("watch", watchText).Short('w').Default("true").BoolVar(&watch)
}
