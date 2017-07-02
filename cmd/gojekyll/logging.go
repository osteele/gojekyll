package main

import (
	"fmt"
	"path/filepath"

	"github.com/osteele/gojekyll/helpers"
)

// This is the longest label. Pull it out here so we can both use it, and measure it for alignment.
const configurationFileLabel = "Configuration file:"

func printSetting(label string, value string) {
	if !quiet {
		fmt.Printf("%s %s\n", helpers.LeftPad(label, len(configurationFileLabel)), value)
	}
}

func printPathSetting(label string, name string) {
	name, err := filepath.Abs(name)
	if err != nil {
		panic("Couldn't convert to absolute path")
	}
	if !quiet {
		printSetting(label, name)
	}
}
