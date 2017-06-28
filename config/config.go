package config

import (
	"strings"

	"github.com/osteele/gojekyll/templates"
	yaml "gopkg.in/yaml.v2"
)

// Config is the Jekyll site configuration, typically read from _config.yml.
// See https://jekyllrb.com/docs/configuration/#default-configuration
type Config struct {
	// Where things are:
	Source      string
	Destination string
	LayoutsDir  string `yaml:"layouts_dir"`
	DataDir     string `yaml:"data_dir"`
	IncludesDir string `yaml:"includes_dir"`
	Collections map[string]templates.VariableMap

	// Handling Reading
	Include     []string
	Exclude     []string
	MarkdownExt string `yaml:"markdown_ext"`

	// Serving
	AbsoluteURL string `yaml:"url"`
	BaseURL     string

	// Outputting
	Permalink string

	Defaults []struct {
		Scope struct {
			Path string
			Type string
		}
		Values templates.VariableMap
	}

	Variables templates.VariableMap `yaml:"-"`
}

// Default returns a default site configuration.
// This is a function instead of a global variable, and returns a new value each time,
// since the caller may overwrite it.
func Default() Config {
	var c Config
	err := Unmarshal([]byte(defaultSiteConfig), &c)
	if err != nil {
		panic(err)
	}
	return c
}

// Unmarshal reads a YAML configuration.
func Unmarshal(bytes []byte, c *Config) error {
	if err := yaml.Unmarshal(bytes, &c); err != nil {
		return err
	}
	if err := yaml.Unmarshal(bytes, &c.Variables); err != nil {
		return err
	}
	return nil
}

// GetFrontMatterDefaults implements https://jekyllrb.com/docs/configuration/#front-matter-defaults
func (c *Config) GetFrontMatterDefaults(relpath, typename string) (m templates.VariableMap) {
	for _, entry := range c.Defaults {
		scope := &entry.Scope
		hasPrefix := strings.HasPrefix(relpath, scope.Path)
		hasType := scope.Type == "" || scope.Type == typename
		if hasPrefix && hasType {
			m = templates.MergeVariableMaps(m, entry.Values)
		}
	}
	return
}

// From https://jekyllrb.com/docs/configuration/#default-configuration
const defaultSiteConfig = `
# Where things are
source:       .
destination:  ./_site
layouts_dir:  _layouts
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
