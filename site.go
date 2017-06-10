package main

import (
	"fmt"
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
	Safe           bool
	Exclude        []string
	Include        []string
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

var siteData = map[interface{}]interface{}{}

func (config *SiteConfig) read(path string) error {
	configBytes, err := ioutil.ReadFile(path)
	if err == nil {
		err = yaml.Unmarshal(configBytes, &config)
	} else if os.IsNotExist(err) {
		err = nil
	}
	return err
}

func buildSiteMap() (map[string]*Page, error) {
	basePath := siteConfig.SourceDir
	fileMap := map[string]*Page{}
	exclusionMap := stringArrayToMap(siteConfig.Exclude)

	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == basePath {
			return nil
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
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if info.IsDir() {
			return nil
		}
		p, err := readFile(relPath, siteData, false)
		if err != nil {
			return err
		}
		fileMap[p.Permalink] = p
		return nil
	}
	err := filepath.Walk(basePath, walkFn)
	if err != nil {
		return nil, err
	}

	for name, colVal := range siteConfig.Collections {
		data := colVal.(map[interface{}]interface{})
		output := false
		if val, found := data["output"]; found {
			output = val.(bool)
		}
		if output {
			err = addCollectionFiles(fileMap, name, data)
			if err != nil {
				return nil, err
			}
		}
	}
	return fileMap, err
}

func addCollectionFiles(fileMap map[string]*Page, name string, data map[interface{}]interface{}) error {
	basePath := siteConfig.SourceDir
	collData := mergeMaps(siteData, data)
	collData["collection"] = name

	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(basePath, path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		p, err := readFile(relPath, collData, false)
		if err != nil {
			return err
		}
		if p.Static {
			fmt.Printf("skipping static file inside collection: %s\n", path)
		} else {
			fileMap[p.Permalink] = p
		}
		return nil
	}
	err := filepath.Walk(filepath.Join(basePath, "_"+name), walkFn)
	return err
}

func getFileURL(path string) (string, bool) {
	for _, v := range siteMap {
		if v.Path == path {
			return v.Permalink, true
		}
	}
	return "", false
}
