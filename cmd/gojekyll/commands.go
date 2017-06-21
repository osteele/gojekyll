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

	"github.com/osteele/gojekyll"
	"github.com/osteele/gojekyll/helpers"
)

// main sets this
var commandStartTime = time.Now()

func buildCommand(site *gojekyll.Site) error {
	printPathSetting("Destination:", site.Destination)
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

func profileCommand(_ *gojekyll.Site) error {
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

func serveCommand(site *gojekyll.Site) error {
	server := gojekyll.Server{Site: site}
	return server.Run(printSetting)
}

func varsCommand(site *gojekyll.Site) error {
	printSetting("Variables:", "")
	siteData := site.Variables
	// The YAML representation including collections is impractically large for debugging.
	// (Actually it's circular, which the yaml package can't handle.)
	// Neuter it. This destroys it as Liquid data, but that's okay in this context.
	for _, c := range site.Collections {
		siteData[c.Name] = fmt.Sprintf("<elided page data for %d items>", len(siteData[c.Name].([]gojekyll.VariableMap)))
	}
	var data interface{} //gojekyll.VariableMap
	switch {
	case *siteVariable:
		data = siteData
	case *dataVariable:
		data = siteData["data"].(gojekyll.VariableMap)
		if *variablePath != "" {
			data = data.(gojekyll.VariableMap)[*variablePath]
		}
	default:
		page, err := cliPage(site, *variablePath)
		if err != nil {
			return err
		}
		data = page.Variables()
	}
	b, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func routesCommand(site *gojekyll.Site) error {
	printSetting("Routes:", "")
	urls := []string{}
	for u, p := range site.Paths {
		if !(*dynamicRoutes && p.Static()) {
			urls = append(urls, u)
		}
	}
	sort.Strings(urls)
	for _, u := range urls {
		fmt.Printf("  %s -> %s\n", u, site.Paths[u].Path())
	}
	return nil
}

func renderCommand(site *gojekyll.Site) error {
	page, err := cliPage(site, *renderPath)
	if err != nil {
		return err
	}
	printPathSetting("Render:", filepath.Join(site.Source, page.Path()))
	printSetting("URL:", page.Permalink())
	printSetting("Content:", "")
	if err := site.CollectionVariable(); err != nil {
		return err
	}
	return page.Write(os.Stdout)
}

// If path starts with /, it's a URL path. Else it's a file path relative
// to the site source directory.
func cliPage(site *gojekyll.Site, path string) (page gojekyll.Page, err error) {
	arg := "/"
	if path != "" {
		arg = path
	}
	if strings.HasPrefix(arg, "/") {
		page, _ = site.PageForURL(arg)
		if page == nil {
			err = helpers.NewPathError("render", arg, "the site does not include a file with this URL path")
		}
	} else {
		page = site.FindPageByFilePath(arg)
		if page == nil {
			err = helpers.NewPathError("render", arg, "no such file")
		}
	}
	return
}
