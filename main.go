package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/urfave/cli"

	yaml "gopkg.in/yaml.v2"

	_ "gopkg.in/urfave/cli.v1"
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

	app.Before = func(c *cli.Context) error {
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

	start := time.Now()

	app.Commands = []cli.Command{
		{
			Name:    "server",
			Aliases: []string{"s", "serve"},
			Usage:   "Serve your site locally",
			Action: func(c *cli.Context) error {
				if err := server(); err != nil {
					fmt.Println(err)
					return err
				}
				return nil
			},
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

			Action: func(c *cli.Context) error {
				printPathSetting("Destination:", site.Dest)
				printSetting("Generating...", "")
				count, err := site.Build()
				if err != nil {
					fmt.Println(err)
					return nil
				}
				elapsed := time.Since(start)
				printSetting("", fmt.Sprintf("created %d files in %.2fs.", count, elapsed.Seconds()))
				return nil
			},
		},
		{
			Name:    "data",
			Aliases: []string{"b"},
			Action: func(c *cli.Context) error {
				path := "/"
				if c.NArg() > 0 {
					path = c.Args().Get(0)
				}
				page := findPageForCLIArg(path)
				if page == nil {
					fmt.Println("No page at", path)
					return nil
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
			},
		},
		{
			Name: "routes",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "dynamic",
					Usage: "Dynamic routes only",
				},
			},
			Action: func(c *cli.Context) error {
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
			},
		},
		{
			Name: "render",
			Action: func(c *cli.Context) error {
				path := "/"
				if c.NArg() > 0 {
					path = c.Args().Get(0)
				}
				page := findPageForCLIArg(path)
				if page == nil {
					fmt.Println("No page at", path)
					return nil
				}
				printPathSetting("Render:", filepath.Join(site.Source, page.Path))
				printSetting("URL:", page.Permalink)
				printSetting("Content:", "")
				if err := page.Render(os.Stdout); err != nil {
					fmt.Println(err)
					return nil
				}
				return nil
			},
		},
	}

	_ = app.Run(os.Args)
}

// If path starts with /, it's a URL path. Else it's a file path relative
// to the site source directory. Either way, return the Page or nil.
func findPageForCLIArg(path string) *Page {
	if path == "" {
		path = "/"
	}
	if strings.HasPrefix(path, "/") {
		return site.Paths[path]
	}
	return site.FindPageByFilePath(path)
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
