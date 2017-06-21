package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/osteele/gojekyll"
	"github.com/osteele/gojekyll/helpers"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Command-line options
var (
	buildOptions          gojekyll.BuildOptions
	useRemoteLiquidEngine bool
	site                  *gojekyll.Site
)

var (
	app         = kingpin.New("gojekyll", "a (maybe someday) Jekyll-compatible blog generator in Go")
	source      = app.Flag("source", "Source directory").Default(".").String()
	destination = app.Flag("destination", "Destination directory").Default("").String()

	serve = app.Command("serve", "Serve your site locally").Alias("server").Alias("s")
	open  = app.Flag("open-url", "Launch your site in a browser").Short('o').Bool()

	build = app.Command("build", "Build your site").Alias("b")

	profile = app.Command("profile", "Build several times, and write a profile file")

	variables    = app.Command("variables", "Display a file or URL path's variables").Alias("v").Alias("var").Alias("vars")
	dataVariable = variables.Flag("data", "Display site.data").Bool()
	siteVariable = variables.Flag("site", "Display site variables instead of page variables").Bool()
	variablePath = variables.Arg("PATH", "Path or URL").String()

	routes        = app.Command("routes", "Display site permalinks and associated files")
	dynamicRoutes = routes.Flag("dynamic", "Only show routes to non-static files").Bool()

	render     = app.Command("render", "Render a file or URL path")
	renderPath = variables.Arg("PATH2", "Path or URL").String()
)

func init() {
	app.Flag("use-liquid-server", "Use Liquid JSON-RPC server").BoolVar(&useRemoteLiquidEngine)
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
	cmd := kingpin.MustParse(app.Parse(os.Args[1:]))
	if err := run(cmd); err != nil {
		app.FatalIfError(err, "")
	}
}
func run(cmd string) error {
	site, err := loadSite(*source, *destination)
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

// Load the site specified at destination into the site global, and print the common banner settings.
func loadSite(source, destination string) (*gojekyll.Site, error) {
	site, err := gojekyll.NewSiteFromDirectory(source)
	if err != nil {
		return nil, err
	}
	site.UseRemoteLiquidEngine = useRemoteLiquidEngine
	if destination != "" {
		site.Destination = destination
	}
	if site.ConfigFile != nil {
		printPathSetting(configurationFileLabel, *site.ConfigFile)
	} else {
		printSetting(configurationFileLabel, "none")
	}
	printPathSetting("Source:", site.Source)
	err = site.Load()
	return site, err
}
