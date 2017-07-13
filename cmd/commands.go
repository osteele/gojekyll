package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"

	"github.com/montanaflynn/stats"
	"github.com/osteele/gojekyll/pages"
	"github.com/osteele/gojekyll/server"
	"github.com/osteele/gojekyll/site"
	"github.com/osteele/gojekyll/utils"
	"github.com/osteele/liquid"
)

// main sets this
var commandStartTime = time.Now()

func buildCommand(site *site.Site) error {
	logger.path("Destination:", site.DestDir())
	logger.label("Generating...", "")
	count, err := site.Build(buildOptions)
	if err != nil {
		return err
	}
	elapsed := time.Since(commandStartTime)
	logger.label("", "wrote %d files in %.2fs.", count, elapsed.Seconds())
	return nil
}

func cleanCommand(site *site.Site) error {
	logger.label("Cleaner:", "Removing %s...", site.DestDir())
	return site.Clean(buildOptions)
}

func benchmarkCommand() (err error) {
	startTime := time.Now()
	samples := []float64{}
	for i := 0; time.Since(startTime) < 10*time.Second; i++ {
		sampleStart := time.Now()
		site, err := loadSite(*source, configFlags)
		if err != nil {
			return err
		}
		_, err = site.Build(buildOptions)
		if err != nil {
			return err
		}
		dur := time.Since(sampleStart).Seconds()
		samples = append(samples, dur)
		quiet = true
		fmt.Printf("Run #%d; %.1fs elapsed\n", i+1, time.Since(commandStartTime).Seconds())
	}
	median, _ := stats.Median(samples)
	stddev, _ := stats.StandardDeviationSample(samples)
	fmt.Printf("%d samples @ %.2fs ± %.2fs\n", len(samples), median, stddev)
	return nil
}

func serveCommand(site *site.Site) error {
	server := server.Server{Site: site}
	return server.Run(*open, func(label, value string) { logger.label(label, value) })
}

func routesCommand(site *site.Site) error {
	logger.label("Routes:", "")
	urls := []string{}
	for u, p := range site.Routes {
		if !(*dynamicRoutes && p.Static()) {
			urls = append(urls, u)
		}
	}
	sort.Strings(urls)
	for _, u := range urls {
		filename := site.Routes[u].SourcePath()
		fmt.Printf("  %s -> %s\n", u, filename)
	}
	return nil
}

func renderCommand(site *site.Site) error {
	p, err := pageFromPathOrRoute(site, *renderPath)
	if err != nil {
		return err
	}
	logger.path("Render:", filepath.Join(site.SourceDir(), p.SourcePath()))
	logger.label("URL:", p.Permalink())
	logger.label("Content:", "")
	return site.WriteDocument(os.Stdout, p)
}

// If path starts with /, it's a URL path. Else it's a file path relative
// to the site source directory.
func pageFromPathOrRoute(s *site.Site, path string) (pages.Document, error) {
	if path == "" {
		path = "/"
	}
	switch {
	case strings.HasPrefix(path, "/"):
		page, found := s.URLPage(path)
		if !found {
			return nil, utils.NewPathError("render", path, "the site does not include a file with this URL path")
		}
		return page, nil
	default:
		page, found := s.FilePathPage(path)
		if !found {
			return nil, utils.NewPathError("render", path, "no such file")
		}
		return page, nil
	}
}

func variablesCommand(site *site.Site) (err error) {
	var data interface{}
	switch {
	case strings.HasPrefix(*variablePath, "site"):
		data, err = utils.FollowDots(site, strings.Split(*variablePath, ".")[1:])
		if err != nil {
			return
		}
	case *variablePath != "":
		data, err = pageFromPathOrRoute(site, *variablePath)
		if err != nil {
			return
		}
	default:
		data = site
	}
	b, err := yaml.Marshal(liquid.FromDrop(data))
	if err != nil {
		return err
	}
	logger.label("Variables:", "")
	fmt.Println(string(b))
	return nil
}

func versionCommand() error {
	fmt.Printf("gojekyll %s\n", Version)
	return nil
}
