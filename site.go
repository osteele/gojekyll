package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/acstech/liquid"

	yaml "gopkg.in/yaml.v2"
)

// Site is a Jekyll site.
type Site struct {
	ConfigFile  *string
	Source      string
	Destination string

	Collections []*Collection
	Variables   VariableMap
	Paths       map[string]Page // URL path -> Page

	config SiteConfig
}

// For now (and maybe always?), there's just one site.
var site = NewSite()

// SiteConfig is the Jekyll site configuration, typically read from _config.yml.
// See https://jekyllrb.com/docs/configuration/#default-configuration
type SiteConfig struct {
	// Where things are:
	Source      string
	Destination string
	Collections map[string]VariableMap

	// Handling Reading
	Include     []string
	Exclude     []string
	MarkdownExt string `yaml:"markdown_ext"`

	// Outputting
	Permalink string
}

// From https://jekyllrb.com/docs/configuration/#default-configuration
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

// NewSite creates a new site record, initialized with the site defaults.
func NewSite() *Site {
	s := new(Site)
	if err := s.readConfigBytes([]byte(siteConfigDefaults)); err != nil {
		panic(err)
	}
	return s
}

// ReadConfiguration reads the configuration file, if it exists.
func (s *Site) ReadConfiguration(source, dest string) error {
	configPath := filepath.Join(source, "_config.yml")
	bytes, err := ioutil.ReadFile(configPath)
	switch {
	case err == nil:
		if err = site.readConfigBytes(bytes); err != nil {
			return err
		}
		s.Source = filepath.Join(source, s.config.Source)
		s.Destination = filepath.Join(s.Source, s.config.Destination)
		s.ConfigFile = &configPath
		if dest != "" {
			site.Destination = dest
		}
		return nil
	case os.IsNotExist(err):
		return nil
	default:
		return err
	}
}

func (s *Site) readConfigBytes(bytes []byte) error {
	configVariables := VariableMap{}
	if err := yaml.Unmarshal(bytes, &s.config); err != nil {
		return err
	}
	if err := yaml.Unmarshal(bytes, &configVariables); err != nil {
		return err
	}
	s.Variables = mergeVariableMaps(s.Variables, configVariables)
	return nil
}

// FindLayout returns a template for the named layout.
func (s *Site) FindLayout(name string, fm *VariableMap) (t *liquid.Template, err error) {
	exts := []string{"", ".html"}
	for _, ext := range strings.SplitN(s.config.MarkdownExt, `,`, -1) {
		exts = append(exts, "."+ext)
	}
	var (
		path    string
		content []byte
		found   bool
	)
	for _, ext := range exts {
		// TODO respect layout config
		path = filepath.Join(s.Source, "_layouts", name+ext)
		content, err = ioutil.ReadFile(path)
		if err == nil {
			found = true
			break
		}
		if !os.IsNotExist(err) {
			return nil, err
		}
	}
	if !found {
		panic(fmt.Errorf("no template for %s", name))
	}
	*fm, err = readFrontMatter(&content)
	if err != nil {
		return
	}
	return liquid.Parse(content, nil)
}

// KeepFile returns a boolean indicating that clean should leave the file in the destination directory.
func (s *Site) KeepFile(path string) bool {
	// TODO
	return false
}

// MarkdownExtensions returns a set of markdown extension, without the final dots.
func (s *Site) MarkdownExtensions() map[string]bool {
	extns := strings.SplitN(s.config.MarkdownExt, `,`, -1)
	return stringArrayToMap(extns)
}

// GetFileURL returns the URL path given a file path, relative to the site source directory.
func (s *Site) GetFileURL(path string) (string, bool) {
	for _, p := range s.Paths {
		if p.Path() == path {
			return p.Permalink(), true
		}
	}
	return "", false
}

// Exclude returns a boolean indicating that the site excludes a file.
func (s *Site) Exclude(path string) bool {
	// TODO exclude based on glob, not exact match
	inclusionMap := stringArrayToMap(s.config.Include)
	exclusionMap := stringArrayToMap(s.config.Exclude)
	base := filepath.Base(path)
	switch {
	case inclusionMap[path]:
		return false
	case path == ".":
		return false
	case exclusionMap[path]:
		return true
	case strings.HasPrefix(base, "."), strings.HasPrefix(base, "_"):
		return true
	default:
		return false
	}
}

// ReadFiles scans the source directory and creates pages and collections.
func (s *Site) ReadFiles() error {
	s.Paths = make(map[string]Page)
	defaults := VariableMap{}

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
		p, err := ReadPage(rel, defaults)
		if err != nil {
			return err
		}
		if p.Published() {
			s.Paths[p.Permalink()] = p
		}
		return nil
	}

	if err := filepath.Walk(s.Source, walkFn); err != nil {
		return err
	}
	if err := s.readCollections(); err != nil {
		return err
	}
	s.initTemplateAttributes()
	return nil
}

// readCollections scans the file system for collections. It adds each collection's
// pages to the site map, and creates a template site variable for each collection.
func (s *Site) readCollections() error {
	for name, d := range s.config.Collections {
		c := makeCollection(s, name, d)
		s.Collections = append(s.Collections, c)
		if c.Output { // TODO always read the pages; just don't build them / include them in routes
			if err := c.ReadPages(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Site) initTemplateAttributes() {
	// TODO site: {pages, posts, related_posts, static_files, html_pages, html_files, collections, data, documents, categories.CATEGORY, tags.TAG}
	s.Variables = mergeVariableMaps(s.Variables, VariableMap{
		"time": time.Now(),
	})
	for _, c := range s.Collections {
		s.Variables[c.Name] = c.PageTemplateObjects()
	}
}
