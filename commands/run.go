package commands

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime/pprof"

	"github.com/osteele/gojekyll/config"
	"github.com/osteele/gojekyll/site"
	"github.com/osteele/gojekyll/version"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

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

func run(cmd string) error { // nolint: gocyclo
	// dispatcher gets to ignore cyclo threshold ^
	if profile || cmd == benchmark.FullCommand() {
		defer setupProfiling()()
	}
	// These commands run *without* loading the site
	switch cmd {
	case benchmark.FullCommand():
		return benchmarkCommand()
	case pluginsApp.FullCommand():
		return pluginsCommand()
	case versionCmd.FullCommand():
		return versionCommand()
	}

	site, err := loadSite(*source, options)
	// Print the version at an awkward place, so its
	// labels will line up. And print it even if
	// loading the site produced an error.
	if *versionFlag {
		logger.label("Version:", version.Version)
	}
	if err != nil {
		return err
	}

	// These commands run *after* the site is loaded
	switch cmd {
	case build.FullCommand():
		return buildCommand(site)
	case clean.FullCommand():
		return cleanCommand(site)
	case render.FullCommand():
		return renderCommand(site)
	case routes.FullCommand():
		return routesCommand(site)
	case serve.FullCommand():
		return serveCommand(site)
	case variables.FullCommand():
		return variablesCommand(site)
	default:
		// kingpin should have provided help and exited before here
		panic("exhaustive switch")
	}
}

// Load the site, and print the common banner settings.
func loadSite(source string, flags config.Flags) (*site.Site, error) {
	site, err := site.FromDirectory(source, flags)
	if err != nil {
		return nil, err
	}
	const configurationFileLabel = "Configuration file:"
	if cf := site.Config().ConfigFile; cf != "" {
		logger.path(configurationFileLabel, cf)
	} else {
		logger.label(configurationFileLabel, "none")
	}
	logger.path("Source:", site.SourceDir())
	err = site.Read()
	return site, err
}

func setupProfiling() func() {
	profilePath := "gojekyll.prof"
	logger.label("Profiling...", "")
	f, err := os.Create(profilePath)
	app.FatalIfError(err, "")
	err = pprof.StartCPUProfile(f)
	app.FatalIfError(err, "")
	return func() {
		pprof.StopCPUProfile()
		err = f.Close()
		app.FatalIfError(err, "")
		logger.Info("Wrote", profilePath)
	}
}
