package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/acstech/liquid"
)

// This is the longest label. Pull it out here so we can both use it, and measure it for alignment.
const configurationFileLabel = "Configuration file:"

func printSetting(label string, value string) {
	fmt.Printf("%s %s\n",
		leftPad(label, len(configurationFileLabel)-len(label)), value)
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
	flag.Parse()

	configPath := filepath.Join(siteConfig.SourceDir, "_config.yml")
	// TODO error if file is e.g. unreadable
	if _, err := os.Stat(configPath); err == nil {
		err := siteConfig.read(configPath)
		if err != nil {
			fmt.Println(err)
			return
		}
		printPathSetting(configurationFileLabel, configPath)
	} else {
		printSetting(configurationFileLabel, "none")
	}
	printPathSetting("Source:", siteConfig.SourceDir)
	printPathSetting("Destination:", siteConfig.DestinationDir)

	fileMap, err := buildSiteMap()
	if err != nil {
		fmt.Println(err)
		return
	}
	siteMap = fileMap

	switch flag.Arg(0) {
	case "s", "serve", "server":
		err = server()
	case "b", "build":
		printSetting("Generating...", "")
		start := time.Now()
		err = build()
		elapsed := time.Since(start)
		printSetting("", fmt.Sprintf("done in %.2fs.", elapsed.Seconds()))
	case "routes":
		fmt.Printf("\nRoutes:\n")
		for url, p := range siteMap {
			fmt.Printf("  %s -> %s\n", url, p.Path)
		}
	case "render":
		// build a single page, and print it to stdout; for testing
		page, err2 := readFile("index.md", siteData, true)
		if err2 != nil {
			err = err2
			break
		}
		fmt.Println(string(page.Body))
	default:
		fmt.Println("A subcommand is required.")
	}
	if err != nil {
		fmt.Println(err)
	}
}
