package config

// Default returns a default site configuration.
//
// This is a function instead of a global variable, and returns a new value each time,
// since the caller may overwrite it.
func Default() Config {
	return FromString(defaultSiteConfig)
}

// FromString returns a new configuration initialized from a string
func FromString(src string) Config {
	var c Config
	// TODO this doesn't set c.Variables. Should it? If so,
	// config.Unmarshal needs to merge them instead of overwriting them (unless yaml.Unmarshal already does this)
	err := Unmarshal([]byte(src), &c)
	if err != nil {
		panic(err)
	}
	return c
}

// From https://jekyllrb.com/docs/configuration/#default-configuration
// The following includes only those keys that are currently implemented.
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

# Plugins
plugins:   []

# Conversion
excerpt_separator: "\n\n"
incremental: false

# Serving
detach:  false
port:    4000
host:    127.0.0.1
baseurl: "" # does not include hostname

# Outputting
permalink:     date
paginate_path: /page:num
timezone:      null

verbose:  false
`
