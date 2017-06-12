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
// See https://jekyllrb.com/docs/configuration/#default-configuration
type SiteConfig struct {
	//Where things are:
	SourceDir      string // `source`
	DestinationDir string `yaml:"destination"`
	Collections    map[string]interface{}

	Permalink string
	Safe      bool
	Exclude   []string
	Include   []string
}

const siteConfigDefaults = `
# Where things are
source:       .
destination:  ./_site
include: [".htaccess"]
data_dir:     _data
includes_dir: _includes
collections:
  posts:
    output:   true

# Handling Reading
include:              [".htaccess"]
exclude:              ["Gemfile", "Gemfile.lock", "node_modules", "vendor/bundle/", "vendor/cache/", "vendor/gems/", "vendor/ruby/"]
keep_files:           [".git", ".svn"]
encoding:             "utf-8"
markdown_ext:         "markdown,mkdown,mkdn,mkd,md"
strict_front_matter: false

# Outputting
permalink:     date
paginate_path: /page:num
timezone:      null
`

//permalink:      "/:categories/:year/:month/:day/:title.html",

var siteConfig SiteConfig

// A map from URL path -> *Page
var siteMap map[string]*Page

var siteData = map[interface{}]interface{}{
	"site": map[string]interface{}{},
}

func (c *SiteConfig) read(path string) error {
	if err := yaml.Unmarshal([]byte(siteConfigDefaults), c); err != nil {
		return err
	}
	switch configBytes, err := ioutil.ReadFile(path); {
	case err != nil && !os.IsNotExist(err):
		return nil
	case err != nil:
		return err
	default:
		return yaml.Unmarshal(configBytes, c)
	}
}

func buildSiteMap() (map[string]*Page, error) {
	basePath := siteConfig.SourceDir
	fileMap := map[string]*Page{}
	exclusionMap := stringArrayToMap(siteConfig.Exclude)

	defaultPageData := map[interface{}]interface{}{}

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
		p, err := readPage(relPath, defaultPageData)
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

func addCollectionFiles(fileMap map[string]*Page, collectionName string, data map[interface{}]interface{}) error {
	basePath := siteConfig.SourceDir
	pages := []*Page{}
	defaultPageData := map[interface{}]interface{}{
		"site":       siteData,
		"collection": collectionName,
	}

	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// if the issue is simply that the directory doesn't exist, ignore the error
			if os.IsNotExist(err) {
				if collectionName != "posts" {
					fmt.Println("Missing directory for collection", collectionName)
				}
				return nil
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
		p, err := readPage(relPath, defaultPageData)
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
	if err := filepath.Walk(filepath.Join(basePath, "_"+collectionName), walkFn); err != nil {
		return err
	}
	collectionPageData := []interface{}{}
	for _, p := range pages {
		collectionPageData = append(collectionPageData, p.PageData())
	}
	siteData[collectionName] = collectionPageData
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
