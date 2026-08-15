package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime/pprof"

	kingpin "github.com/alecthomas/kingpin/v2"
	"github.com/osteele/gojekyll/config"
	"github.com/osteele/gojekyll/site"
	"github.com/osteele/gojekyll/version"
)

// commandEntry describes how to run one CLI subcommand.
type commandEntry struct {
	needsSite bool
	run       func(*site.Site) error
}

var commandTable = map[string]commandEntry{
	benchmark.FullCommand():  {needsSite: false, run: func(*site.Site) error { return benchmarkCommand() }},
	pluginsApp.FullCommand(): {needsSite: false, run: func(*site.Site) error { pluginsCommand(); return nil }},
	versionCmd.FullCommand(): {needsSite: false, run: func(*site.Site) error { return versionCommand() }},
	build.FullCommand():      {needsSite: true, run: buildCommand},
	clean.FullCommand():      {needsSite: true, run: cleanCommand},
	render.FullCommand():     {needsSite: true, run: renderCommand},
	routes.FullCommand():     {needsSite: true, run: func(s *site.Site) error { routesCommand(s); return nil }},
	serve.FullCommand():      {needsSite: true, run: serveCommand},
	variables.FullCommand():  {needsSite: true, run: variablesCommand},
}

// ParseAndRun parses and executes the command-line arguments.
func ParseAndRun(args []string) error {
	if reflect.DeepEqual(args, []string{"--version"}) {
		return versionCommand()
	}
	cmd := kingpin.MustParse(app.Parse(args))
	if options.Destination != nil {
		dest, err := filepath.Abs(*options.Destination)
		app.FatalIfError(err, "")
		options.Destination = &dest
	}
	if options.DryRun {
		verbose := true
		options.Verbose = &verbose
	}
	return run(cmd)
}

func run(cmd string) error {
	// Set quiet mode on logger
	log.SetQuiet(quiet)

	if profile || cmd == benchmark.FullCommand() {
		defer setupProfiling()()
	}

	entry, ok := commandTable[cmd]
	if !ok {
		// kingpin should have provided help and exited before here
		panic(fmt.Sprintf("unknown command: %s", cmd))
	}

	if !entry.needsSite {
		return entry.run(nil)
	}

	site, err := loadSite(*source, options)
	// Print the version at an awkward place, so its
	// labels will line up. And print it even if
	// loading the site produced an error.
	if *versionFlag {
		bannerLog.label("Version:", "%s", version.Version)
	}
	if err != nil {
		return err
	}
	return entry.run(site)
}

// Load the site, and print the common banner settings.
func loadSite(source string, flags config.Flags) (*site.Site, error) {
	site, err := site.FromDirectory(source, flags)
	if err != nil {
		return nil, err
	}
	const configurationFileLabel = "Configuration file:"
	if cf := site.Config().ConfigFile; cf != "" {
		bannerLog.path(configurationFileLabel, cf)
	} else {
		bannerLog.label(configurationFileLabel, "none")
	}
	bannerLog.path("Source:", site.SourceDir())
	err = site.Read()
	return site, err
}

func setupProfiling() func() {
	profilePath := "gojekyll.prof"
	bannerLog.label("Profiling...", "")
	f, err := os.Create(profilePath)
	app.FatalIfError(err, "")
	err = pprof.StartCPUProfile(f)
	app.FatalIfError(err, "")
	return func() {
		pprof.StopCPUProfile()
		err = f.Close()
		app.FatalIfError(err, "")
		bannerLog.Info("Wrote", profilePath)
	}
}
