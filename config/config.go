package config

import (
	"strings"

	"github.com/osteele/gojekyll/helpers"
	"github.com/osteele/gojekyll/templates"
	yaml "gopkg.in/yaml.v2"
)

// Config is the Jekyll site configuration, typically read from _config.yml.
// See https://jekyllrb.com/docs/configuration/#default-configuration
type Config struct {
	// Where things are:
	Source      string
	Destination string
	LayoutsDir  string                            `yaml:"layouts_dir"`
	DataDir     string                            `yaml:"data_dir"`
	IncludesDir string                            `yaml:"includes_dir"`
	Collections map[string]map[string]interface{} `yaml:"-"`

	// Handling Reading
	Include     []string
	Exclude     []string
	KeepFiles   []string `yaml:"keep_files"`
	MarkdownExt string   `yaml:"markdown_ext"`

	// Filtering Content
	Drafts      bool `yaml:"show_drafts"`
	Future      bool
	Unpublished bool

	// Plugins
	Plugins []string

	// Plugins
	ExcerptSeparator string `yaml:"excerpt_separator"`

	// Serving
	Host        string
	Port        int
	AbsoluteURL string `yaml:"url"`
	BaseURL     string

	// Outputting
	Permalink string

	Defaults []struct {
		Scope struct {
			Path string
			Type string
		}
		Values map[string]interface{}
	}

	Variables map[string]interface{} `yaml:"-"`
}

type configCompat struct {
	Gems []string
}

type collectionsList struct {
	Collections []string
}

type collectionsMap struct {
	Collections map[string]map[string]interface{}
}

// SourceDir returns the source directory as an absolute path.
func (c *Config) SourceDir() string {
	return helpers.MustAbs(c.Source)
}

// GetFrontMatterDefaults implements https://jekyllrb.com/docs/configuration/#front-matter-defaults
func (c *Config) GetFrontMatterDefaults(typename, relpath string) (m map[string]interface{}) {
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

// Unmarshal updates site from a YAML configuration file.
func Unmarshal(bytes []byte, c *Config) error {
	var (
		compat configCompat
		cList  collectionsList
		// cMap   collectionsMap
	)
	if err := yaml.Unmarshal(bytes, &c); err != nil {
		return err
	}
	if err := yaml.Unmarshal(bytes, &c.Variables); err != nil {
		return err
	}
	if err := yaml.Unmarshal(bytes, &cList); err == nil {
		if len(c.Collections) == 0 {
			c.Collections = make(map[string]map[string]interface{})
		}
		for _, name := range cList.Collections {
			c.Collections[name] = map[string]interface{}{}
		}
	}
	cMap := collectionsMap{c.Collections}
	if err := yaml.Unmarshal(bytes, &cMap); err == nil {
		c.Collections = cMap.Collections
	}
	if err := yaml.Unmarshal(bytes, &compat); err != nil {
		return err
	}
	if len(c.Plugins) == 0 {
		c.Plugins = compat.Gems
	}
	return nil
}
