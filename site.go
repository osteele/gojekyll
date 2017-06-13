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
	// Where things are:
	SourceDir      string `yaml:"source"`
	DestinationDir string `yaml:"destination"`
	Collections    map[string]interface{}

	// Handling Reading
	Include     []string
	Exclude     []string
	MarkdownExt string `yaml:"markdown_ext"`

	// Outputting
	Permalink string
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

var siteData = map[interface{}]interface{}{}

// For unit tests
func init() {
	siteConfig.setDefaults()
}

func (c *SiteConfig) setDefaults() {
	if err := yaml.Unmarshal([]byte(siteConfigDefaults), &siteData); err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal([]byte(siteConfigDefaults), &siteConfig); err != nil {
		panic(err)
	}
}

func (c *SiteConfig) read(path string) error {
	c.setDefaults()
	switch configBytes, err := ioutil.ReadFile(path); {
	case err != nil && !os.IsNotExist(err):
		return nil
	case err != nil:
		return err
	default:
		if err := yaml.Unmarshal(configBytes, siteData); err != nil {
			return err
		}
		return yaml.Unmarshal(configBytes, c)
	}
}

// MarkdownExtensions returns a set of markdown extension.
func (c *SiteConfig) MarkdownExtensions() map[string]bool {
	extns := strings.SplitN(siteConfig.MarkdownExt, `,`, -1)
	return stringArrayToMap(extns)
}

func getFileURL(path string) (string, bool) {
	for _, v := range siteMap {
		if v.Path == path {
			return v.Permalink, true
		}
	}
	return "", false
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

		rel, err := filepath.Rel(basePath, path)
		if err != nil {
			return err
		}
		base := filepath.Base(rel)
		// TODO exclude based on glob, not exact match
		_, exclude := exclusionMap[rel]
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
		p, err := ReadPage(rel, defaultPageData)
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
	if err := ReadCollections(fileMap); err != nil {
		return nil, err
	}
	return fileMap, nil
}

// ReadCollections scans the file system for collections. It adds each collection's
// pages to the site map, and creates a template site variable for each collection.
func ReadCollections(fileMap map[string]*Page) error {
	for s, d := range siteConfig.Collections {
		data, ok := d.(map[interface{}]interface{})
		if !ok {
			panic("expected collection value to be a map")
		}
		c := makeCollection(s, data)
		if c.Output { // TODO always read the pages; just don't build them
			if err := c.ReadPages(fileMap); err != nil {
				return err
			}
		}
		siteData[c.Name] = c.PageData()
	}
	return nil
}

// Collection is a Jekyll collection.
type Collection struct {
	Name   string
	Data   map[interface{}]interface{}
	Output bool
	Pages  []*Page
}

func makeCollection(name string, data map[interface{}]interface{}) *Collection {
	return &Collection{
		Name:   name,
		Data:   data,
		Output: getBool(data, "output", false),
	}
}

// PageData returns an array of a page data, for use as the template variable
// value of the collection.
func (c *Collection) PageData() (d []interface{}) {
	for _, p := range c.Pages {
		d = append(d, p.PageData())
	}
	return
}

// Posts returns true if the collection is the special "posts" collection.
func (c *Collection) Posts() bool {
	return c.Name == "posts"
}

// SourceDir returns the source directory for pages in the collection.
func (c *Collection) SourceDir() string {
	return filepath.Join(siteConfig.SourceDir, "_"+c.Name)
}

// ReadPages scans the file system for collection pages, and adds them to c.Pages.
func (c *Collection) ReadPages(fileMap map[string]*Page) error {
	basePath := siteConfig.SourceDir
	d := map[interface{}]interface{}{
		"site":       siteData,
		"collection": c.Name,
	}
	d = mergeMaps(c.Data, d)

	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// if the issue is simply that the directory doesn't exist, ignore the error
			if os.IsNotExist(err) {
				if !c.Posts() {
					fmt.Println("Missing directory for collection", c.Name)
				}
				return nil
			}
			return err
		}
		rel, err := filepath.Rel(basePath, path)
		switch {
		case err != nil:
			return err
		case info.IsDir():
			return nil
		}
		p, err := ReadPage(rel, d)
		switch {
		case err != nil:
			return err
		case p.Static:
			fmt.Printf("skipping static file inside collection: %s\n", path)
		case p.Published:
			fileMap[p.Permalink] = p
			c.Pages = append(c.Pages, p)
		}
		return nil
	}
	return filepath.Walk(c.SourceDir(), walkFn)
}
