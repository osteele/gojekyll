package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/osteele/gojekyll"
	"github.com/osteele/gojekyll/helpers"
	"gopkg.in/urfave/cli.v1"
)

// Command-line options
var buildOptions gojekyll.BuildOptions
var useRemoteLiquidEngine bool

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
	var source, destination string

	app := cli.NewApp()
	app.Name = "gojekyll"
	app.Usage = "a (maybe someday) Jekyll-compatible blog generator in Go"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "source",
			Value:       ".",
			Usage:       "Source directory",
			Destination: &source,
		},
		cli.StringFlag{
			Name:        "destination",
			Value:       "",
			Usage:       "Destination directory",
			Destination: &destination,
		},
		cli.BoolFlag{
			Name:        "use-liquid-server",
			Usage:       "Use Liquid JSON-RPC server",
			Destination: &useRemoteLiquidEngine,
		},
	}

	withSite := func(cmd func(*cli.Context, *gojekyll.Site) error) func(*cli.Context) error {
		siteLoader := func() (*gojekyll.Site, error) { return loadSite(source, destination) }
		return loadSiteAndRun(siteLoader, cmd)
	}

	app.Commands = []cli.Command{
		{
			Name:    "serve",
			Aliases: []string{"server", "s"},
			Usage:   "Serve your site locally",
			Action:  withSite(serveCommand),
		},
		{
			Name:    "build",
			Aliases: []string{"b"},
			Usage:   "Build your site",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:        "dry-run, n",
					Usage:       "Dry run",
					Destination: &buildOptions.DryRun,
				},
			},
			Action: withSite(buildCommand),
		},
		{
			Name:   "profile",
			Usage:  "Build several times, and write a profile file",
			Action: withSite(profileCommand),
		}, {
			Name:    "variables",
			Aliases: []string{"v", "var", "vars"},
			Usage:   "Print a file or URL path's variables",
			Action:  withSite(dataCommand),
		},
		{
			Name:  "routes",
			Usage: "Display site permalinks and associated files",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "dynamic",
					Usage: "Only show routes to non-static files",
				},
			},
			Action: withSite(routesCommand),
		},
		{
			Name:   "render",
			Usage:  "Render a file or URL path",
			Action: withSite(renderCommand),
		},
	}

	_ = app.Run(os.Args)
}
