package main

import (
	"os"
	"path/filepath"
	"runtime/pprof"

	"github.com/osteele/gojekyll/config"
	"github.com/osteele/gojekyll/sites"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Command-line options
var (
	buildOptions sites.BuildOptions
	configFlags  = config.Flags{}
	profile      = false
	quiet        = false
)

var (
	app    = kingpin.New("gojekyll", "a (somewhat) Jekyll-compatible blog generator")
	source = app.Flag("source", "Source directory").Short('s').Default(".").String()
	_      = app.Flag("destination", "Destination directory").Short('d').Action(stringVar("destination", &configFlags.Destination)).String()
	_      = app.Flag("drafts", "Render posts in the _drafts folder").Short('D').Action(boolVar("drafts", &configFlags.Drafts)).Bool()
	_      = app.Flag("future", "Publishes posts with a future date").Action(boolVar("future", &configFlags.Future)).Bool()
	_      = app.Flag("unpublished", "Render posts that were marked as unpublished").Action(boolVar("unpublished", &configFlags.Unpublished)).Bool()

	build = app.Command("build", "Build your site").Alias("b")
	clean = app.Command("clean", "Clean the site (removes site output) without building.")

	benchmark = app.Command("profile", "Repeat build for ten seconds. Implies --profile.")

	render     = app.Command("render", "Render a file or URL path to standard output")
	renderPath = render.Arg("PATH", "Path or URL").String()

	routes        = app.Command("routes", "Display site permalinks and associated files")
	dynamicRoutes = routes.Flag("dynamic", "Only show routes to non-static files").Bool()

	serve = app.Command("serve", "Serve your site locally").Alias("server").Alias("s")
	open  = serve.Flag("open-url", "Launch your site in a browser").Short('o').Bool()

	variables    = app.Command("variables", "Display a file or URL path's variables").Alias("v").Alias("var").Alias("vars")
	variablePath = variables.Arg("PATH", "Path, URL, site, or site...").String()
)

func init() {
	app.Flag("profile", "Create a Go pprof CPU profile").BoolVar(&profile)
	app.Flag("quiet", "Silence (some) output.").Short('q').BoolVar(&quiet)
	build.Flag("dry-run", "Dry run").Short('n').BoolVar(&buildOptions.DryRun)
}

func main() {
	parseAndRun(os.Args[1:])
}

func parseAndRun(args []string) {
	app.HelpFlag.Short('h')
	cmd := kingpin.MustParse(app.Parse(args))
	if configFlags.Destination != nil {
		dest, err := filepath.Abs(*configFlags.Destination)
		app.FatalIfError(err, "")
		configFlags.Destination = &dest
	}
	if buildOptions.DryRun {
		buildOptions.Verbose = true
	}
	if cmd == benchmark.FullCommand() {
		profile = true
	}
	app.FatalIfError(run(cmd), "")
}

func run(cmd string) error { // nolint: gocyclo
	if profile {
		setupProfiling()
	}
	site, err := loadSite(*source, configFlags)
	if err != nil {
		return err
	}

	switch cmd {
	case benchmark.FullCommand():
		return benchmarkCommand(site)
	case build.FullCommand():
		return buildCommand(site)
	case clean.FullCommand():
		return cleanCommand(site)
	case render.FullCommand():
		return renderCommand(site)
	case routes.FullCommand():
		return routesCommand(site)
	case serve.FullCommand():
		return serveCommand(site)
	case variables.FullCommand():
		return varsCommand(site)
	default:
		// kingpin should have provided help and exited before here
		panic("unknown command")
	}
}

// Load the site, and print the common banner settings.
func loadSite(source string, flags config.Flags) (*sites.Site, error) {
	site, err := sites.NewSiteFromDirectory(source, flags)
	if err != nil {
		return nil, err
	}
	const configurationFileLabel = "Configuration file:"
	if site.ConfigFile != nil {
		logger.path(configurationFileLabel, *site.ConfigFile)
	} else {
		logger.label(configurationFileLabel, "none")
	}
	logger.label("Source:", site.SourceDir())
	err = site.Load()
	return site, err
}

func setupProfiling() {
	profilePath := "gojekyll.prof"
	logger.label("Profiling...", "")
	f, err := os.Create(profilePath)
	app.FatalIfError(err, "")
	err = pprof.StartCPUProfile(f)
	app.FatalIfError(err, "")
	defer func() {
		pprof.StopCPUProfile()
		err = f.Close()
		app.FatalIfError(err, "")
		logger.Info("Wrote", profilePath)
	}()
}
