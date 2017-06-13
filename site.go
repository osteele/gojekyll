package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// Site is a Jekyll site.
type Site struct {
	ConfigFile *string
	Source     string
	Dest       string

	Collections []*Collection
	Data        map[interface{}]interface{}
	Paths       map[string]*Page // URL path -> *Page

	config SiteConfig
}

// For now (and maybe always?), there's just one site.
var site Site

// SiteConfig is the Jekyll site configuration, typically read from _config.yml.
// See https://jekyllrb.com/docs/configuration/#default-configuration
type SiteConfig struct {
	// Where things are:
	Source      string
	Destination string
	Collections map[string]interface{}

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

//TODO permalink:      "/:categories/:year/:month/:day/:title.html",

// For unit tests
func init() {
	site.Initialize()
}

// Initialize sets the defaults
func (s *Site) Initialize() {
	y := []byte(siteConfigDefaults)
	if err := yaml.Unmarshal(y, &s.config); err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(y, &s.Data); err != nil {
		panic(err)
	}
	s.Paths = make(map[string]*Page)
}

// ReadConfiguration reads the configuration file, if it's present
func ReadConfiguration(source, dest string) error {
	site.Initialize()
	configPath := filepath.Join(source, "_config.yml")
	_, err := os.Stat(configPath)
	switch {
	case err == nil:
		if err = site.ReadConfig(configPath); err != nil {
			return err
		}
		site.Source = filepath.Join(source, site.config.Source)
		site.Dest = filepath.Join(site.Source, site.config.Destination)
		if dest != "" {
			site.Dest = dest
		}
		site.ConfigFile = &configPath
		return nil
	case os.IsNotExist(err):
		return nil
	default:
		return err
	}
}

func (s *Site) ReadConfig(path string) error {
	switch configBytes, err := ioutil.ReadFile(path); {
	case err != nil && !os.IsNotExist(err):
		return nil
	case err != nil:
		return err
	default:
		if err := yaml.Unmarshal(configBytes, s.config); err != nil {
			return err
		}
		return yaml.Unmarshal(configBytes, s.Data)
	}
}

func (s *Site) KeepFile(p string) bool {
	// TODO
	return false
}

// MarkdownExtensions returns a set of markdown extension.
func (s *Site) MarkdownExtensions() map[string]bool {
	extns := strings.SplitN(s.config.MarkdownExt, `,`, -1)
	return stringArrayToMap(extns)
}

// GetFileURL returns the URL path given a file path, relative to the site source directory.
func (s *Site) GetFileURL(path string) (string, bool) {
	for _, v := range s.Paths {
		if v.Path == path {
			return v.Permalink, true
		}
	}
	return "", false
}

// Exclude returns true iff a site excludes a file.
func (s *Site) Exclude(path string) bool {
	// TODO exclude based on glob, not exact match
	exclusionMap := stringArrayToMap(s.config.Exclude)
	base := filepath.Base(path)
	switch {
	case path == ".":
		return false
	case exclusionMap[path]:
		return true
		// TODO check Include
	case strings.HasPrefix(base, "."), strings.HasPrefix(base, "_"):
		return true
	default:
		return false
	}
}

// ReadFiles scans the source directory and creates pages and collections.
func (s *Site) ReadFiles() error {
	d := map[interface{}]interface{}{}

	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(s.Source, path)
		if err != nil {
			return err
		}
		switch {
		case info.IsDir() && s.Exclude(rel):
			return filepath.SkipDir
		case info.IsDir(), s.Exclude(rel):
			return nil
		}
		p, err := ReadPage(rel, d)
		if err != nil {
			return err
		}
		if p.Published {
			s.Paths[p.Permalink] = p
		}
		return nil
	}

	if err := filepath.Walk(s.Source, walkFn); err != nil {
		return err
	}
	return s.ReadCollections()
}

// ReadCollections scans the file system for collections. It adds each collection's
// pages to the site map, and creates a template site variable for each collection.
func (s *Site) ReadCollections() error {
	for name, d := range s.config.Collections {
		data, ok := d.(map[interface{}]interface{})
		if !ok {
			panic("expected collection value to be a map")
		}
		c := makeCollection(s, name, data)
		s.Collections = append(s.Collections, c)
		if c.Output { // TODO always read the pages; just don't build them / include them in routes
			if err := c.ReadPages(); err != nil {
				return err
			}
		}
		s.Data[c.Name] = c.PageData()
	}
	return nil
}
