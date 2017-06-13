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
	Config SiteConfig
	Source string
	Dest   string
	Data   map[interface{}]interface{}
	Paths  map[string]*Page // URL path -> *Page
}

// For now (and maybe always?), there's just one site.
var site Site

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

//TODO permalink:      "/:categories/:year/:month/:day/:title.html",

// For unit tests
func init() {
	site.Initialize()
}

func (s *Site) Initialize() {
	y := []byte(siteConfigDefaults)
	if err := yaml.Unmarshal(y, &s.Config); err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(y, &s.Data); err != nil {
		panic(err)
	}
	s.Paths = make(map[string]*Page)
}

func (s *Site) ReadConfig(path string) error {
	s.Initialize()
	switch configBytes, err := ioutil.ReadFile(path); {
	case err != nil && !os.IsNotExist(err):
		return nil
	case err != nil:
		return err
	default:
		if err := yaml.Unmarshal(configBytes, s.Config); err != nil {
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
	extns := strings.SplitN(s.Config.MarkdownExt, `,`, -1)
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
	exclusionMap := stringArrayToMap(s.Config.Exclude)
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

		rel, err := filepath.Rel(s.Config.SourceDir, path)
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

	if err := filepath.Walk(s.Config.SourceDir, walkFn); err != nil {
		return err
	}
	return s.ReadCollections()
}

// ReadCollections scans the file system for collections. It adds each collection's
// pages to the site map, and creates a template site variable for each collection.
func (s *Site) ReadCollections() error {
	for name, d := range s.Config.Collections {
		data, ok := d.(map[interface{}]interface{})
		if !ok {
			panic("expected collection value to be a map")
		}
		c := makeCollection(s, name, data)
		if c.Output { // TODO always read the pages; just don't build them / include them in routes
			if err := c.ReadPages(); err != nil {
				return err
			}
		}
		s.Data[c.Name] = c.PageData()
	}
	return nil
}
