package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/urfave/cli"

	_ "gopkg.in/urfave/cli.v1"
	yaml "gopkg.in/yaml.v2"
)

// Command-line options
var options struct {
	useHardLinks bool
	dryRun       bool
}

// This is the longest label. Pull it out here so we can both use it, and measure it for alignment.
const configurationFileLabel = "Configuration file:"

func printSetting(label string, value string) {
	fmt.Printf("%s %s\n", LeftPad(label, len(configurationFileLabel)), value)
}

func printPathSetting(label string, path string) {
	path, err := filepath.Abs(path)
	if err != nil {
		panic("Couldn't convert to absolute path")
	}
	printSetting(label, path)
}

var commandStartTime = time.Now()

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
			Destination: &source,
		},
	}

	// Load the site specified at destination into the site global, and print the common banner settings.
	loadSite := func() error {
		if err := site.ReadConfiguration(source, destination); err != nil {
			return err
		}
		if site.ConfigFile != nil {
			printPathSetting(configurationFileLabel, *site.ConfigFile)
		} else {
			printSetting(configurationFileLabel, "none")

		}
		printPathSetting("Source:", site.Source)
		return site.ReadFiles()
	}

	// Given a subcommand function, load the site and then call the subcommand.
	withSite := func(cmd func(c *cli.Context) error) func(c *cli.Context) error {
		return func(c *cli.Context) error {
			if err := loadSite(); err != nil {
				return cli.NewExitError(err, 1)
			}
			if err := cmd(c); err != nil {
				return cli.NewExitError(err, 1)
			}
			return nil
		}
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
					Destination: &options.dryRun,
				},
			},
			Action: withSite(buildCommand),
		},
		{
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
			Action: withSite(routeCommand),
		},
		{
			Name:   "render",
			Action: withSite(renderCommand),
		},
	}

	_ = app.Run(os.Args)
}

func buildCommand(c *cli.Context) error {
	printPathSetting("Destination:", site.Dest)
	printSetting("Generating...", "")
	count, err := site.Build()
	if err != nil {
		return err
	}
	elapsed := time.Since(commandStartTime)
	printSetting("", fmt.Sprintf("created %d files in %.2fs.", count, elapsed.Seconds()))
	return nil
}

func serveCommand(c *cli.Context) error {
	return server()
}

func dataCommand(c *cli.Context) error {
	page, err := cliPage(c)
	if err != nil {
		return err
	}

	printSetting("Data:", "")
	// The YAML representation including collections is impractically large for debugging.
	// (Actually it's circular, which the yaml package can't handle.)
	// Neuter it. This destroys it as Liquid data, but that's okay in this context.
	for _, c := range site.Collections {
		site.Data[c.Name] = fmt.Sprintf("<elided page data for %d items>", len(site.Data[c.Name].([]interface{})))
	}
	b, _ := yaml.Marshal(stringMap(page.Data()))
	fmt.Println(string(b))
	return nil
}

func routeCommand(c *cli.Context) error {
	printSetting("Routes:", "")
	urls := []string{}
	for u, p := range site.Paths {
		if !(c.Bool("dynamic") && p.Static) {
			urls = append(urls, u)
		}
	}
	sort.Strings(urls)
	for _, u := range urls {
		fmt.Printf("  %s -> %s\n", u, site.Paths[u].Path)
	}
	return nil
}

func renderCommand(c *cli.Context) error {
	page, err := cliPage(c)
	if err != nil {
		return err
	}
	printPathSetting("Render:", filepath.Join(site.Source, page.Path))
	printSetting("URL:", page.Permalink)
	printSetting("Content:", "")
	return page.Render(os.Stdout)
}

// If path starts with /, it's a URL path. Else it's a file path relative
// to the site source directory.
func cliPage(c *cli.Context) (page *Page, err error) {
	path := "/"
	if c.NArg() > 0 {
		path = c.Args().Get(0)
	}
	if strings.HasPrefix(path, "/") {
		page = site.Paths[path]
		if page == nil {
			err = &os.PathError{Op: "render", Path: path, Err: errors.New("the site does not include a file with this URL path")}
		}
	} else {
		page = site.FindPageByFilePath(path)
		if page == nil {
			err = &os.PathError{Op: "render", Path: path, Err: errors.New("no such file")}
		}
	}
	return
}

// FindPageByFilePath returns a Page or nil, referenced by relative path.
func (s *Site) FindPageByFilePath(path string) *Page {
	for _, p := range s.Paths {
		if p.Path == path {
			return p
		}
	}
	return nil
}
