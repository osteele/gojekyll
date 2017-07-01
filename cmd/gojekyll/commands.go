package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/pprof"
	"sort"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"

	"github.com/osteele/gojekyll/helpers"
	"github.com/osteele/gojekyll/pages"
	"github.com/osteele/gojekyll/server"
	"github.com/osteele/gojekyll/sites"
)

// main sets this
var commandStartTime = time.Now()

func buildCommand(site *sites.Site) error {
	printPathSetting("Destination:", site.DestDir())
	printSetting("Generating...", "")
	if buildOptions.DryRun {
		buildOptions.Verbose = true
	}
	count, err := site.Build(buildOptions)
	if err != nil {
		return err
	}
	elapsed := time.Since(commandStartTime)
	printSetting("", fmt.Sprintf("wrote %d files in %.2fs.", count, elapsed.Seconds()))
	return nil
}

func profileCommand(_ *sites.Site) error {
	printSetting("Profiling...", "")
	var profilePath = "gojekyll.prof"
	f, err := os.Create(profilePath)
	if err != nil {
		return err
	}
	if err = pprof.StartCPUProfile(f); err != nil {
		return err
	}
	t0 := time.Now()
	for i := 0; time.Since(t0) < 10*time.Second; i++ {
		site, err := loadSite(*source, *destination)
		if err != nil {
			return err
		}
		_, err = site.Build(buildOptions)
		if err != nil {
			return err
		}
		printSetting("", fmt.Sprintf("Run #%d; %.1fs elapsed", i+1, time.Since(t0).Seconds()))
	}
	pprof.StopCPUProfile()
	if err := f.Close(); err != nil {
		return err
	}
	fmt.Println("Wrote", profilePath)
	return nil
}

func serveCommand(site *sites.Site) error {
	server := server.Server{Site: site}
	return server.Run(*open, printSetting)
}

func varsCommand(site *sites.Site) error {
	printSetting("Variables:", "")
	siteData := site.SiteVariables()
	// The YAML representation including collections is impractically large for debugging.
	// (Actually it's circular, which the yaml package can't handle.)
	// Neuter it. This destroys it as Liquid data, but that's okay in this context.
	for _, c := range site.Collections {
		siteData[c.Name] = fmt.Sprintf("<elided page data for %d items>", len(siteData[c.Name].([]interface{})))
	}
	var data map[string]interface{}
	switch {
	case *siteVariable:
		data = siteData
	case *dataVariable:
		data = siteData["data"].(map[string]interface{})
		if *variablePath != "" {
			data = data[*variablePath].(map[string]interface{})
		}
	default:
		page, err := pageFromPathOrRoute(site, *variablePath)
		if err != nil {
			return err
		}
		data = page.PageVariables()
	}
	b, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
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
		page, found := s.RelPathPage(path)
		if !found {
			return nil, helpers.NewPathError("render", path, "no such file")
		}
		return page, nil
	}
}
