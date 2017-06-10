package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// SiteConfig is the Jekyll site configuration, typically read from _config.yml.
type SiteConfig struct {
	Permalink      string
	SourceDir      string
	DestinationDir string
	// Safe           bool
	Exclude []string
	Include []string
	// KeepFiles      []string
	// TimeZone       string
	// Encoding       string
	Collections map[string]interface{}
}

// Initialize with defaults.
var siteConfig = SiteConfig{
	SourceDir:      "./",
	DestinationDir: "./_site",
	Permalink:      "/:categories/:year/:month/:day/:title.html",
}

// A map from URL path -> *Page
var siteMap map[string]*Page

func (config *SiteConfig) readFromDirectory(path string) error {
	configBytes, err := ioutil.ReadFile(path)
	if err == nil {
		err = yaml.Unmarshal(configBytes, &config)
	} else if os.IsNotExist(err) {
		err = nil
	}
	return err
}

func buildFileMap() (map[string]*Page, error) {
	basePath := siteConfig.SourceDir
	fileMap := map[string]*Page{}
	exclusionMap := stringArrayToMap(siteConfig.Exclude)

	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == siteConfig.SourceDir {
			return nil
		}
		// TODO replace by info.IsDir
		stat, err := os.Stat(path)
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(basePath, path)
		if err != nil {
			return err
		}
		base := filepath.Base(relPath)
		// TODO exclude based on glob, not exact match
		_, exclude := exclusionMap[relPath]
		exclude = exclude || strings.HasPrefix(base, ".") || strings.HasPrefix(base, "_")
		if exclude {
			if stat.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if !stat.IsDir() {
			page, err := readFile(relPath, false)
			if err != nil {
				return err
			}
			fileMap[page.Permalink] = page
		}
		return nil
	}
	err := filepath.Walk(basePath, walkFn)
	return fileMap, err
}

func getFilePermalink(path string) (string, bool) {
	for _, v := range siteMap {
		if v.Path == path {
			return v.Permalink, true
		}
	}
	return "", false
}
