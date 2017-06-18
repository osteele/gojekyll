package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/pprof"

	"github.com/osteele/gojekyll"
	"github.com/osteele/gojekyll/helpers"
	"gopkg.in/urfave/cli.v1"
)

// Command-line options
var buildOptions gojekyll.BuildOptions

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
	}

	withSite := func(cmd func(*cli.Context, *gojekyll.Site) error) func(*cli.Context) error {
		siteLoader := func() (*gojekyll.Site, error) { return loadSite(source, destination) }
		return loadSiteAndRun(siteLoader, cmd)
	}

	app.Commands = []cli.Command{
		{
			Name:    "server",
			Aliases: []string{"s", "serve"},
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
			Name:    "benchmark",
			Aliases: []string{"b"},
			Usage:   "Benchmark",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:        "dry-run, n",
					Usage:       "Dry run",
					Destination: &buildOptions.DryRun,
				},
			},
			Action: withSite(benchmarkCommand),
		}, {
			Name:    "data",
			Aliases: []string{"b"},
			Action:  withSite(dataCommand),
		},
		{
			Name: "routes",
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
			Action: withSite(renderCommand),
		},
	}

	if true {
		f, err := os.Create("gojekyll.prof")
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	_ = app.Run(os.Args)
}
