package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"

	"github.com/acstech/liquid"
)

// Command-line options
var options struct {
	useHardLinks bool
	dryRun       bool
}

// This is the longest label. Pull it out here so we can both use it, and measure it for alignment.
const configurationFileLabel = "Configuration file:"

func printSetting(label string, value string) {
	fmt.Printf("%s %s\n",
		leftPad(label, len(configurationFileLabel)), value)
}

func printPathSetting(label string, path string) {
	path, err := filepath.Abs(path)
	if err != nil {
		panic("Couldn't convert to absolute path")
	}
	printSetting(label, path)
}

func main() {
	liquid.Tags["link"] = LinkFactory

	// general options
	source := flag.String("source", ".", "Source directory")
	dest := flag.String("destination", "./_site", "Destination directory")

	// maybe add flags for these
	flag.BoolVar(&options.dryRun, "dry-run", false, "Dry run")
	// flag.BoolVar(&options.useHardLinks, "-n", false, "Dry run")

	// routes subcommand
	dynamic := flag.Bool("dynamic", false, "Dynamic routes only")

	flag.Parse()
	if len(flag.Args()) < 1 {
		fmt.Println("A subcommand is required.")
		return
	}

	configPath := filepath.Join(*source, "_config.yml")
	_, err := os.Stat(configPath)
	switch {
	case err == nil:
		if err = site.ReadConfig(configPath); err != nil {
			fmt.Println(err)
			return
		}
		site.Config.SourceDir = *source
		site.Config.DestinationDir = *dest
		printPathSetting(configurationFileLabel, configPath)
	case os.IsNotExist(err):
		site.Initialize()
		printSetting(configurationFileLabel, "none")
	default:
		fmt.Println(err)
		return
	}
	printPathSetting("Source:", site.Config.SourceDir)

	start := time.Now()
	if err = site.ReadFiles(); err != nil {
		fmt.Println(err)
		return
	}

	switch flag.Arg(0) {
	case "s", "serve", "server":
		if err = server(); err != nil {
			fmt.Println(err)
		}
	case "b", "build":
		printPathSetting("Destination:", site.Config.DestinationDir)
		printSetting("Generating...", "")
		if err = site.Build(); err != nil {
			fmt.Println(err)
			break
		}
		elapsed := time.Since(start)
		printSetting("", fmt.Sprintf("done in %.2fs.", elapsed.Seconds()))
	case "data":
		path := flag.Arg(1)
		page := findPageForCLIArg(path)
		if page == nil {
			fmt.Println("No page at", path)
			return
		}

		// The YAML representation including collections is impractically large for debugging.
		// (Actually it's circular, which the yaml package can't handle.)
		// Neuter it. This destroys it as Liquid data, but that's okay in this context.
		for c := range site.Config.Collections {
			site.Data[c] = fmt.Sprintf("<elided page data for %d items>", len(site.Data[c].([]interface{}))) //"..."
		}
		b, _ := yaml.Marshal(stringMap(page.Data()))
		fmt.Println(string(b))
	default:
		fmt.Println("A subcommand is required.")
	case "routes":
		fmt.Printf("\nRoutes:\n")
		urls := []string{}
		for u, p := range site.Paths {
			if !(*dynamic && p.Static) {
				urls = append(urls, u)
			}
		}
		sort.Strings(urls)
		for _, u := range urls {
			fmt.Printf("  %s -> %s\n", u, site.Paths[u].Path)
		}
	case "render":
		path := flag.Arg(1)
		page := findPageForCLIArg(path)
		if page == nil {
			fmt.Println("No page at", path)
			return
		}
		printPathSetting("Render:", filepath.Join(site.Config.SourceDir, page.Path))
		printSetting("URL:", page.Permalink)
		if err := page.Render(os.Stdout); err != nil {
			fmt.Println(err)
			break
		}
	}
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
