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
	printPathSetting("Destination:", site.DestDir())
	printSetting("Generating...", "")
	count, err := site.Build(buildOptions)
	if err != nil {
		return err
	}
	elapsed := time.Since(commandStartTime)
	printSetting("", fmt.Sprintf("wrote %d files in %.2fs.", count, elapsed.Seconds()))
	return nil
}

func cleanCommand(site *sites.Site) error {
	printPathSetting("Cleaner:", fmt.Sprintf("Removing %s...", site.DestDir()))
	return site.Clean(buildOptions)
}

func benchmarkCommand(_ *sites.Site) error {
	t0 := time.Now()
	for i := 0; time.Since(t0) < 10*time.Second; i++ {
		site, err := loadSite(*source, configFlags)
		if err != nil {
			return err
		}
		_, err = site.Build(buildOptions)
		if err != nil {
			return err
		}
		printSetting("", fmt.Sprintf("Run #%d; %.1fs elapsed", i+1, time.Since(t0).Seconds()))
	}
	return nil
}

func serveCommand(site *sites.Site) error {
	server := server.Server{Site: site}
	return server.Run(*open, printSetting)
}

func routesCommand(site *sites.Site) error {
	printSetting("Routes:", "")
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
	printPathSetting("Render:", filepath.Join(site.SourceDir(), p.SiteRelPath()))
	printSetting("URL:", p.Permalink())
	printSetting("Content:", "")
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

func varsCommand(site *sites.Site) error {
	var data interface{}
	switch {
	case strings.HasPrefix(*variablePath, "site"):
		data = site
		for _, name := range strings.Split(*variablePath, ".")[1:] {
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
			return fmt.Errorf("no such property: %q", name)
		}
	case *variablePath != "":
		page, err := pageFromPathOrRoute(site, *variablePath)
		if err != nil {
			return err
		}
		data = page.(liquid.Drop).ToLiquid()
	default:
		data = site
	}
	if drop, ok := data.(liquid.Drop); ok {
		data = drop
	}
	b, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	printSetting("Variables:", "")
	fmt.Println(string(b))
	return nil
}
