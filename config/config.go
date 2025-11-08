package config

import (
	"os"
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

	// Conversion
	ExcerptSeparator string `yaml:"excerpt_separator"`
	Incremental      bool
	Sass             struct {
		Dir string `yaml:"sass_dir"`
		// TODO Style string // compressed
	}

	// Serving
	Host        string
	Port        int
	AbsoluteURL string `yaml:"url"`
	BaseURL     string

	// Outputting
	Permalink         string
	PermalinkTimezone string `yaml:"permalink_timezone,omitempty"`
	Timezone          string
	Verbose           bool
	Defaults          []struct {
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

	// Meta
	ConfigFile string                 `yaml:"-"`
	m          map[string]interface{} `yaml:"-"` // config file, as map
	ms         yaml.MapSlice          `yaml:"-"` // config file, as MapSlice

	// Plugins
	RequireFrontMatter        bool            `yaml:"-"`
	RequireFrontMatterExclude map[string]bool `yaml:"-"`
}

// FromDirectory updates the config from the config file in
// the directory, if such a file exists.
func (c *Config) FromDirectory(dir string) error {
	path := filepath.Join(dir, "_config.yml")
	bytes, err := os.ReadFile(path)
	switch {
	case os.IsNotExist(err):
		// break
	case err != nil:
		return err
	default:
		if err = Unmarshal(bytes, c); err != nil {
			return utils.WrapPathError(err, path)
		}
		c.ConfigFile = path
	}
	c.Source = dir
	return nil
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

// IsConfigPath returns true if its arguments is a site configuration file.
func (c *Config) IsConfigPath(rel string) bool {
	return rel == "_config.yml"
}

// SassDir returns the relative path of the SASS directory.
func (c *Config) SassDir() string {
	return "_sass"
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
	case utils.StringArrayContains(c.Include, rel):
		return false
	case c.RequireFrontMatterExclude[strings.ToUpper(utils.TrimExt(filepath.Base(rel)))]:
		return true
	default:
		return false
	}
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
	if err := yaml.Unmarshal(bytes, &c.ms); err != nil {
		return err
	}
	if err := yaml.Unmarshal(bytes, &c.m); err != nil {
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

// Variables returns the configuration as a Liquid variable map.
func (c *Config) Variables() map[string]interface{} {
	m := map[string]interface{}{}
	for _, item := range c.ms {
		if s, ok := item.Key.(string); ok {
			m[s] = item.Value
		}
	}
	return m
}

// Set sets a value in the Liquid variable map.
// This does not update the corresponding value in the Config struct.
//
// Note: Iterates by index rather than value to modify c.ms in place.
// Range-over-value creates copies, so `item.Value = val` would modify
// the copy instead of the original slice element.
//
// Thread safety: This method is called during site initialization and reload.
// Reload is protected by Server.m mutex, so concurrent calls don't occur.
func (c *Config) Set(key string, val interface{}) {
	c.m[key] = val
	for i := range c.ms {
		if c.ms[i].Key == key {
			c.ms[i].Value = val
			return
		}
	}
	c.ms = append(c.ms, yaml.MapItem{Key: key, Value: val})
}

// Map returns the config indexed by key, if it's a map.
func (c *Config) Map(key string) (map[string]interface{}, bool) {
	if m, ok := c.m[key]; ok {
		if m, ok := m.(map[string]interface{}); ok {
			return m, ok
		}
	}
	return nil, false
}

// String returns the config indexed by key, if it's a string.
func (c *Config) String(key string) (string, bool) {
	if m, ok := c.m[key]; ok {
		if m, ok := m.(string); ok {
			return m, ok
		}
	}
	return "", false
}
