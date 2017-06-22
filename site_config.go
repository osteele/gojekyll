package gojekyll

import (
	"github.com/osteele/gojekyll/templates"
	yaml "gopkg.in/yaml.v2"
)

// SiteConfig is the Jekyll site configuration, typically read from _config.yml.
// See https://jekyllrb.com/docs/configuration/#default-configuration
type SiteConfig struct {
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

	// Outputting
	Permalink string

	Defaults []struct {
		Scope struct {
			Path string
			Type string
		}
		Values templates.VariableMap
	}
}

func (s *Site) readConfigBytes(bytes []byte) error {
	configVariables := templates.VariableMap{}
	if err := yaml.Unmarshal(bytes, &s.config); err != nil {
		return err
	}
	if err := yaml.Unmarshal(bytes, &configVariables); err != nil {
		return err
	}
	s.Variables = templates.MergeVariableMaps(s.Variables, configVariables)
	return nil
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
