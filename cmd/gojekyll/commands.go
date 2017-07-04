package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"

	"github.com/osteele/gojekyll/helpers"
	"github.com/osteele/gojekyll/pages"
	"github.com/osteele/gojekyll/server"
	"github.com/osteele/gojekyll/sites"
	"github.com/osteele/liquid"
)

// main sets this
var commandStartTime = time.Now()

func buildCommand(site *sites.Site) error {
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

func cleanCommand(site *sites.Site) error {
	logger.label("Cleaner:", "Removing %s...", site.DestDir())
	return site.Clean(buildOptions)
}

func benchmarkCommand(site *sites.Site) (err error) {
	for i := 0; time.Since(commandStartTime) < 10*time.Second; i++ {
		// skip this the first time, since the caller has already
		// started the profile and loaded the site once.
		if i > 0 {
			site, err = loadSite(*source, configFlags)
			if err != nil {
				return err
			}
		}
		_, err = site.Build(buildOptions)
		if err != nil {
			return err
		}
		logger.label("", "Run #%d; %.1fs elapsed", i+1, time.Since(commandStartTime).Seconds())
	}
	return nil
}

func serveCommand(site *sites.Site) error {
	server := server.Server{Site: site}
	return server.Run(*open, func(label, value string) { logger.label(label, value) })
}

func routesCommand(site *sites.Site) error {
	logger.label("Routes:", "")
	urls := []string{}
	for u, p := range site.Routes {
		if !(*dynamicRoutes && p.Static()) {
			urls = append(urls, u)
		}
	}
	sort.Strings(urls)
	for _, u := range urls {
		filename := site.Routes[u].SiteRelPath()
		fmt.Printf("  %s -> %s\n", u, filename)
	}
	return nil
}

func renderCommand(site *sites.Site) error {
	p, err := pageFromPathOrRoute(site, *renderPath)
	if err != nil {
		return err
	}
	logger.path("Render:", filepath.Join(site.SourceDir(), p.SiteRelPath()))
	logger.label("URL:", p.Permalink())
	logger.label("Content:", "")
	return site.WriteDocument(p, os.Stdout)
}

// If path starts with /, it's a URL path. Else it's a file path relative
// to the site source directory.
func pageFromPathOrRoute(s *sites.Site, path string) (pages.Document, error) {
	if path == "" {
		path = "/"
	}
	switch {
	case strings.HasPrefix(path, "/"):
		page, found := s.URLPage(path)
		if !found {
			return nil, helpers.NewPathError("render", path, "the site does not include a file with this URL path")
		}
		return page, nil
	default:
		page, found := s.FilePathPage(path)
		if !found {
			return nil, helpers.NewPathError("render", path, "no such file")
		}
		return page, nil
	}
}

func varsCommand(site *sites.Site) (err error) {
	var data interface{}
	switch {
	case strings.HasPrefix(*variablePath, "site"):
		data, err = followDots(site, strings.Split(*variablePath, ".")[1:])
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
	b, err := yaml.Marshal(toLiquid(data))
	if err != nil {
		return err
	}
	logger.label("Variables:", "")
	fmt.Println(string(b))
	return nil
}

func followDots(data interface{}, props []string) (interface{}, error) {
	for _, name := range props {
		if drop, ok := data.(liquid.Drop); ok {
			data = drop.ToLiquid()
		}
		if reflect.TypeOf(data).Kind() == reflect.Map {
			item := reflect.ValueOf(data).MapIndex(reflect.ValueOf(name))
			if item.CanInterface() && !item.IsNil() {
				data = item.Interface()
				continue
			}
		}
		return nil, fmt.Errorf("no such property: %q", name)
	}
	return data, nil
}

func toLiquid(value interface{}) interface{} {
	switch value := value.(type) {
	case liquid.Drop:
		return value.ToLiquid()
	default:
		return value
	}
}
