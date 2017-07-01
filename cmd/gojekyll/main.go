package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/osteele/gojekyll/config"
	"github.com/osteele/gojekyll/helpers"
	"github.com/osteele/gojekyll/sites"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Command-line options
var (
	buildOptions sites.BuildOptions
)

var configFlags = config.Flags{}

var (
	app    = kingpin.New("gojekyll", "a (maybe someday) Jekyll-compatible blog generator in Go")
	source = app.Flag("source", "Source directory").Short('s').Default(".").String()
	_      = app.Flag("destination", "Destination directory").Short('d').Action(stringAction("destination", &configFlags.Destination)).String()
	_      = app.Flag("future", "Publishes posts with a future date").Action(boolAction("future", &configFlags.Future)).Bool()
	_      = app.Flag("drafts", "Render posts in the _drafts folder").Short('D').Action(boolAction("drafts", &configFlags.Drafts)).Bool()
	_      = app.Flag("unpublished", "Render posts that were marked as unpublished").Action(boolAction("unpublished", &configFlags.Unpublished)).Bool()

	serve = app.Command("serve", "Serve your site locally").Alias("server").Alias("s")
	open  = serve.Flag("open-url", "Launch your site in a browser").Short('o').Bool()

	build = app.Command("build", "Build your site").Alias("b")

	profile = app.Command("profile", "Build several times, and write a profile file")

	variables    = app.Command("variables", "Display a file or URL path's variables").Alias("v").Alias("var").Alias("vars")
	dataVariable = variables.Flag("data", "Display site.data").Bool()
	siteVariable = variables.Flag("site", "Display site variables instead of page variables").Bool()
	variablePath = variables.Arg("PATH", "Path or URL").String()

	routes        = app.Command("routes", "Display site permalinks and associated files")
	dynamicRoutes = routes.Flag("dynamic", "Only show routes to non-static files").Bool()

	render     = app.Command("render", "Render a file or URL path")
	renderPath = render.Arg("PATH", "Path or URL").String()
)

func init() {
	build.Flag("dry-run", "Dry run").Short('n').BoolVar(&buildOptions.DryRun)
}

// This is the longest label. Pull it out here so we can both use it, and measure it for alignment.
const configurationFileLabel = "Configuration file:"

func printSetting(label string, value string) {
	fmt.Printf("%s %s\n", helpers.LeftPad(label, len(configurationFileLabel)), value)
}

func printPathSetting(label string, name string) {
	name, err := filepath.Abs(name)
	if err != nil {
		panic("Couldn't convert to absolute path")
	}
	printSetting(label, name)
}

func main() {
	app.HelpFlag.Short('h')
	cmd := kingpin.MustParse(app.Parse(os.Args[1:]))
	if err := run(cmd); err != nil {
		app.FatalIfError(err, "")
	}

}
func run(cmd string) error {
	site, err := loadSite(*source, configFlags)
	if err != nil {
		return err
	}

	switch cmd {
	case build.FullCommand():
		return buildCommand(site)
	case profile.FullCommand():
		return profileCommand(site)
	case render.FullCommand():
		return renderCommand(site)
	case routes.FullCommand():
		return routesCommand(site)
	case serve.FullCommand():
		return serveCommand(site)
	case variables.FullCommand():
		return varsCommand(site)
	}
	return nil
}

// Load the site, and print the common banner settings.
func loadSite(source string, flags config.Flags) (*sites.Site, error) {
	site, err := sites.NewSiteFromDirectory(source, flags)
	if err != nil {
		return nil, err
	}
	if site.ConfigFile != nil {
		printPathSetting(configurationFileLabel, *site.ConfigFile)
	} else {
		printSetting(configurationFileLabel, "none")
	}
	printPathSetting("Source:", site.SourceDir())
	err = site.Load()
	return site, err
}
