package config

import (
	"path/filepath"
	"strings"

	"github.com/osteele/gojekyll/templates"
	"github.com/osteele/gojekyll/utils"
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
	Theme       string

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

	// Conversion
	Incremental bool

	// Serving
	Host        string
	Port        int
	AbsoluteURL string `yaml:"url"`
	BaseURL     string

	// Outputting
	Permalink string
	Timezone  string
	Verbose   bool
	Defaults  []struct {
		Scope struct {
			Path string
			Type string
		}
		Values map[string]interface{}
	}

	// CLI-only
	DryRun       bool `yaml:"-"`
	ForcePolling bool `yaml:"-"`
	Watch        bool `yaml:"-"`

	// Unstructured data for templates
	Variables map[string]interface{} `yaml:"-"`

	// Plugins
	RequireFrontMatter        bool            `yaml:"-"`
	RequireFrontMatterExclude map[string]bool `yaml:"-"`
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
	return utils.MustAbs(c.Source)
}

// GetFrontMatterDefaults implements https://jekyllrb.com/docs/configuration/#front-matter-defaults
func (c *Config) GetFrontMatterDefaults(typename, rel string) (m map[string]interface{}) {
	for _, entry := range c.Defaults {
		scope := &entry.Scope
		hasPrefix := strings.HasPrefix(rel, scope.Path)
		hasType := scope.Type == "" || scope.Type == typename
		if hasPrefix && hasType {
			m = templates.MergeVariableMaps(m, entry.Values)
		}
	}
	return
}

// RequiresFrontMatter returns a bool indicating whether the file requires front matter in order to recognize as a page.
func (c *Config) RequiresFrontMatter(rel string) bool {
	switch {
	case c.RequireFrontMatter:
		return true
	case !c.IsMarkdown(rel):
		return true
	case contains(c.Include, rel):
		return false
	case c.RequireFrontMatterExclude[strings.ToUpper(utils.TrimExt(filepath.Base(rel)))]:
		return true
	default:
		return false
	}
}

func contains(array []string, s string) bool {
	for _, item := range array {
		if item == s {
			return true
		}
	}
	return false
}

// Unmarshal updates site from a YAML configuration file.
func Unmarshal(bytes []byte, c *Config) error {
	var (
		compat configCompat
		cList  collectionsList
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
