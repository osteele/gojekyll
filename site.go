package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"

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

var siteData = map[interface{}]interface{}{
	"site": map[string]interface{}{},
}

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
		if p.Published {
			fileMap[p.Permalink] = p
		}
		return nil
	}

	if err := filepath.Walk(basePath, walkFn); err != nil {
		return nil, err
	}

	for name, colVal := range siteConfig.Collections {
		data, ok := colVal.(map[interface{}]interface{})
		if !ok {
			panic("expected collection value to be a map")
		}
		output := getBool(data, "output", false)
		if output {
			if err := addCollectionFiles(fileMap, name, data); err != nil {
				return nil, err
			}
		}
	}

	return fileMap, nil
}

func addCollectionFiles(fileMap map[string]*Page, name string, data map[interface{}]interface{}) error {
	basePath := siteConfig.SourceDir
	collData := mergeMaps(siteData, data)
	collData["collection"] = name
	pages := []*Page{}

	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// if the issue is simply that the directory doesn't exist, ignore the error
			if pathErr, ok := err.(*os.PathError); ok {
				if pathErr.Err == syscall.ENOENT {
					fmt.Println("Missing directory for collection", name)
					return nil
				}
			}
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
		} else if p.Published {
			fileMap[p.Permalink] = p
			pages = append(pages, p)
		}
		return nil
	}
	if err := filepath.Walk(filepath.Join(basePath, "_"+name), walkFn); err != nil {
		return err
	}
	pageData := []interface{}{}
	for _, p := range pages {
		pageData = append(pageData, p.CollectionItemData())
	}
	siteData["site"].(map[string]interface{})[name] = pageData
	return nil
}

func getFileURL(path string) (string, bool) {
	for _, v := range siteMap {
		if v.Path == path {
			return v.Permalink, true
		}
	}
	return "", false
}
