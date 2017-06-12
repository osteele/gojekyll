package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/acstech/liquid"
)

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

	flag.StringVar(&siteConfig.DestinationDir, "destination", siteConfig.DestinationDir, "Destination directory")
	flag.StringVar(&siteConfig.SourceDir, "source", siteConfig.SourceDir, "Source directory")

	// routes subcommand
	dynamic := flag.Bool("dynamic", false, "Dynamic routes only")

	flag.Parse()
	if len(flag.Args()) < 1 {
		fmt.Println("A subcommand is required.")
		return
	}

	configPath := filepath.Join(siteConfig.SourceDir, "_config.yml")
	// TODO error if file is e.g. unreadable
	if _, err := os.Stat(configPath); err == nil {
		if err := siteConfig.read(configPath); err != nil {
			fmt.Println(err)
			return
		}
		printPathSetting(configurationFileLabel, configPath)
	} else {
		printSetting(configurationFileLabel, "none")
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
		if path == "" {
			path = "index.md"
		}
		if strings.HasPrefix(path, "/") {
			if page, found := siteMap[path]; found {
				path = page.Path
			} else {
				fmt.Println("No page at", path)
				return
			}
		}
		page, err := readPage(path, siteData)
		if err != nil {
			fmt.Println(err)
			break
		}
		body, err := page.Render()
		if err != nil {
			fmt.Println(err)
			break
		}
		printPathSetting("Render:", filepath.Join(siteConfig.SourceDir, path))
		printSetting("URL:", page.Permalink)
		fmt.Println(string(body))
	default:
		fmt.Println("A subcommand is required.")
	}
}
