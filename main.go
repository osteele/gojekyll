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
	flag.StringVar(&siteConfig.DestinationDir, "destination", siteConfig.DestinationDir, "Destination directory")
	source := flag.String("source", ".", "Source directory")

	// maybe add flags for these
	// options.useHardLinks = true

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
		if err = siteConfig.read(configPath); err != nil {
			fmt.Println(err)
			return
		}
		siteConfig.SourceDir = *source
		printPathSetting(configurationFileLabel, configPath)
	case os.IsNotExist(err):
		printSetting(configurationFileLabel, "none")
	default:
		fmt.Println(err)
		return
	}
	printPathSetting("Source:", siteConfig.SourceDir)

	start := time.Now()
	fileMap, err := buildSiteMap()
	if err != nil {
		fmt.Println(err)
		return
	}
	siteMap = fileMap

	switch flag.Arg(0) {
	case "s", "serve", "server":
		if err = server(); err != nil {
			fmt.Println(err)
		}
	case "b", "build":
		printPathSetting("Destination:", siteConfig.DestinationDir)
		printSetting("Generating...", "")
		if err = build(); err != nil {
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
		for c := range siteConfig.Collections {
			siteData[c] = fmt.Sprintf("<elided page data for %d items>", len(siteData[c].([]interface{}))) //"..."
		}
		b, _ := yaml.Marshal(stringMap(page.Data()))
		fmt.Println(string(b))
	default:
		fmt.Println("A subcommand is required.")
	case "routes":
		fmt.Printf("\nRoutes:\n")
		urls := []string{}
		for u, p := range siteMap {
			if !(*dynamic && p.Static) {
				urls = append(urls, u)
			}
		}
		sort.Strings(urls)
		for _, u := range urls {
			fmt.Printf("  %s -> %s\n", u, siteMap[u].Path)
		}
	case "render":
		path := flag.Arg(1)
		page := findPageForCLIArg(path)
		if page == nil {
			fmt.Println("No page at", path)
			return
		}
		printPathSetting("Render:", filepath.Join(siteConfig.SourceDir, page.Path))
		printSetting("URL:", page.Permalink)
		if err := page.Render(os.Stdout); err != nil {
			fmt.Println(err)
			break
		}
	}
}

func findPageForCLIArg(path string) *Page {
	if path == "" {
		path = "/"
	}
	if strings.HasPrefix(path, "/") {
		return siteMap[path]
	}
	return findPageByFilePath(path)
}

func findPageByFilePath(path string) *Page {
	for _, p := range siteMap {
		if p.Path == path {
			return p
		}
	}
	return nil
}
