# Gojekyll

 [![][travis-svg]][travis-url] [![][coveralls-svg]][coveralls-url] [![][go-report-card-svg]][go-report-card-url] [![][license-svg]][license-url]

Gojekyll is a partially-compatible clone of the [Jekyll](https://jekyllrb.com) static site generator, written in the [Go](https://golang.org) programming language. It provides `build` and `serve` commands, with directory watch and live reload.

| &nbsp;                  | Gojekyll                                  | Jekyll | Hugo         |
|-------------------------|-------------------------------------------|--------|--------------|
| Stable                  |                                           | ✓      | ✓            |
| Fast                    | ✓<br>([~20×Jekyll](./docs/benchmarks.md)) |        | ✓            |
| Template language       | Liquid                                    | Liquid | Go templates |
| SASS                    | ✓                                         | ✓      |              |
| Jekyll compatibility    | [partial](#current-limitations)           | ✓      |              |
| Plugins                 | [some](./docs/plugins.md)                 | yes    | ?            |
| Windows support         |                                           | ✓      | ✓            |
| Implementation language | Go                                        | Ruby   | Go           |

<!-- TOC -->

- [Gojekyll](#gojekyll)
    - [Usage](#usage)
    - [Installation](#installation)
        - [Binary Downloads](#binary-downloads)
        - [From Source](#from-source)
    - [Status](#status)
        - [Current Limitations](#current-limitations)
        - [Other Differences](#other-differences)
        - [Feature Checklist](#feature-checklist)
    - [Contributing](#contributing)
    - [Attribution](#attribution)
    - [Related](#related)
    - [License](#license)

<!-- /TOC -->

## Usage

```bash
gojekyll build       # builds the site in the current directory into _site
gojekyll serve       # serve the app at http://localhost:4000; reload on changes
gojekyll help
gojekyll help build
```

## Installation

### Binary Downloads

1. Ubuntu (64-bit) and macOS binaries are available from the [releases page](https://github.com/osteele/gojekyll/releases).
2. [Optional] **Highlight**. To use the `{% highlight %}` tag, you also need [Pygments](http://pygments.org): `pip install Pygments`.
3. [Optional] **Themes**. To use a theme, you need to create a `Gemfile` that lists the theme, and run (the Ruby command) `bundle install`. The [Jekyll theme instructions](https://jekyllrb.com/docs/themes/) provides more detail, and should work for Gojekyll too.

### From Source

Pre-requisites:

1. **Install go** (1) via [Homebrew](https://brew.sh): `brew install go`; or (2) [download](https://golang.org/doc/install#tarball).
2. See items (2-3) under **Binary Downloads**, above, for optional installations.

First-time install:

```bash
go get github.com/osteele/gojekyll
```

[Later] Update to the latest version:

```bash
go get -u github.com/osteele/liquid github.com/osteele/gojekyll
```

[Optional] Install command-line autocompletion:

```bash
# Bash:
eval "$(gojekyll --completion-script-bash)"
# Zsh:
eval "$(gojekyll --completion-script-zsh)"
```

## Status

This project is at an early stage of development.

It works on the GitHub Pages sites that I care about, and it looks credible on a spot-check of other Jekyll sites.

### Current Limitations

Missing features:

- Pagination
- Windows compatibility
- Math
- Plugin system. ([Some individual plugins](./docs/plugins.md) are emulated.)
- Liquid filter `sassify` is not implemented
- Liquid is run in strict mode; undefined filters are errors.
- Markdown features
  - [Attribute lists](https://kramdown.gettalong.org/syntax.html#attribute-list-definitions)
  - [`markdown="span"`, `markdown="block"`](https://kramdown.gettalong.org/syntax.html#html-blocks)
- Markdown configuration options

Differences:

- The order of YAML maps, in `_config` and `site.data`, is not preserved.

Also see the [detailed status](#feature-status) below.

### Other Differences

These will probably not change:

By design:

- Plugins must be listed in the config file, not a Gemfile.
- The wrong type in a `_config.yml` is an error.
- Server live reload is always on.
- `serve --watch` (the default) reloads the `_config.yml` and data files too.
- `serve` generates pages on the fly; it doesn't write to the file system.
- Files are cached in `/tmp/gojekyll-${USER}`, not `./.sass-cache`
- `markdown=1` is only processed in Markdown files. This matches the Jekyll documentation, but not its implementation (which also expands markdown inside of `markdown=1` elements inside \*.html files).
- `<` and `>` inside markdown is interpreted as HTML. For example, `This is <b>bold</b>` renders as <b>bold</b>. This behavior matches the [Markdown spec](https://daringfireball.net/projects/markdown/syntax#html), but differs from Jekyll's default Kramdown processor.

Muzukashii:

- An extensible plugin mechanism – support for plugins that aren't compiled into the executable.

### Feature Checklist

- [ ] Content
  - [x] Front Matter
  - [x] Posts
  - [x] Static Files
  - [x] Variables
  - [x] Collections
  - [x] Data Files
  - [ ] Assets
    - [ ] Coffeescript
    - [x] Sass/SCSS
- [ ] Customization
  - [x] Templates
    - [ ] Jekyll filters
      - [ ] `group_by_exp` and `scssify`
      - [x] everything else
    - [x] Jekyll tags
  - [x] Includes
  - [x] Permalinks
  - [ ] Pagination
  - [ ] Plugins – partial; see [here](./docs/plugins.md)
  - [x] Themes
  - [x] Layouts
- [x] Server
  - [x] Directory watch
- [ ] Commands
  - [x] `build`
    - [x] `--source`, `--destination`, `--drafts`, `--future`, `--unpublished`
    - [x] `--incremental`, `--watch`, `--force_polling`
    - [ ] `--baseurl`, `--config`, `--lsi`, `JEKYLL_ENV=production`
    - [ ] `--limit-posts`
  - [x] `clean`
  - [x] `help`
  - [x] `serve`
    - [x] `--open-uri`, `--host`, `--port`
    - [x] `--incremental`, `–watch`, `--force_polling`
    - [ ] `--baseurl`, `--config`
    - [ ] `--detach`, `--ssl`-* – not planned
  - [ ] `doctor`, `import`, `new`, `new-theme` – not planned
- [ ] Windows

## Contributing

Bug reports, test cases, and code contributions are [more than welcome](./CONTRIBUTING.md).

## Attribution

Gojekyll uses these libraries:

| Package                                                                        | Author(s)                                              | Usage                                  | License                                 |
|--------------------------------------------------------------------------------|--------------------------------------------------------|----------------------------------------|-----------------------------------------|
| [github.com/jaschaephraim/lrserver](https://github.com/jaschaephraim/lrserver) | Jascha Ephraim                                         | Live Reload                            | MIT License                             |
| [github.com/kyokomi/emoji](https://github.com/kyokomi/emoji)                   | kyokomi                                                | `jemoji` plugin emulation              | MIT License                             |
| [github.com/osteele/liquid](https://github.com/osteele/liquid)                 | yours truly                                            | Liquid processor                       | MIT License                             |
| [github.com/pkg/browser](https://github.com/pkg/browser)                       | [pkg](https://github.com/pkg)                          | `serve --open-url` option              | BSD 2-clause "Simplified" License       |
| [github.com/radovskyb/watcher](https://github.com/radovskyb/watcher)           | Benjamin Radovsky                                      | Polling file watch (`--force_polling`) | BSD 3-clause "New" or "Revised" License |
| [github.com/russross/blackfriday](https://github.com/russross/blackfriday)     | Russ Ross                                              | Markdown processing                    | Simplified BSD License                  |
| [github.com/sass/libsass](https://github.com/sass/libsass)                     | Listed [here](https://https://github.com/sass/libsass) | C port of the Ruby SASS compiler       | MIT License                             |
| [github.com/tdewolff/minify](https://github.com/tdewolff/minify)               | Taco de Wolff                                          | CSS minimization                       | MIT License                             |
| [github.com/wellington/go-libsass](https://github.com/wellington/go-libsass)   | Drew Wells                                             | Go bindings for **libsass**            | ???                                     |
| [gopkg.in/alecthomas/kingpin.v2](https://github.com/alecthomas/kingpin)        | Alec Thomas                                            | command-line arguments                 | MIT License                             |
| [gopkg.in/yaml.v2](https://github.com/go-yaml/yaml)                            | Canonical                                              | YAML support                           | Apache License 2.0                      |

In addition, the following pieces of text were taken from Jekyll and its plugins.
They are used under the terms of the MIT License.

| Source                                                                          | Use                  | Description            |
|---------------------------------------------------------------------------------|----------------------|------------------------|
| [Jekyll template documentation](https://jekyllrb.com/docs/templates/)           | test cases           | filter examples        |
| `jekyll help` command                                                           | `gojekyll help` text | help text              |
| [`jekyll-feed` plugin](https://github.com/jekyll/jekyll-feed)                   | plugin emulation     | `feed.xml` template    |
| [`jekyll-redirect-from` plugin](https://github.com/jekyll/jekyll-redirect-from) | plugin emulation     | redirect page template |
| [`jekyll-sitemap` plugin](https://github.com/jekyll/jekyll-redirect-from)       | plugin emulation     | sitemap template       |
| [`jekyll-seo-tag` plugin](https://github.com/jekyll/jekyll-redirect-from)       | plugin emulation     | feed template          |

The theme for in-browser error reporting was adapted from facebookincubator/create-react-app.

The gopher image in the `testdata` directory is from [Wikimedia Commons](https://commons.wikimedia.org/wiki/File:Gophercolor.jpg). It is used under the [Creative Commons Attribution-Share Alike 3.0 Unported license](https://creativecommons.org/licenses/by-sa/3.0/deed.en).

In addition to being totally and obviously inspired by Jekyll and its plugins, Jekyll's  solid *documentation* was indispensible --- especially since I wanted to implement Jekyll as documented, not port its source code. The [Jekyll docs](https://jekyllrb.com/docs/home/) were always open in at least one tab during development.

## Related

[Hugo](https://gohugo.io) is the pre-eminent Go static site generator. It isn't Jekyll-compatible (-), but it's highly polished, performant, and productized (+++).

[jkl](https://github.com/drone/jkl) is another Go clone of Jekyll. If I'd found it sooner I might have started this project by forking that one. It's got a better name.

[Liquid](https://github.com/osteele/liquid) is a pure Go implementation of Liquid templates, that I finally caved and wrote in order to use in this project.

[Jekyll](https://jekyllrb.com), of course.

## License

MIT

[coveralls-url]: https://coveralls.io/r/osteele/gojekyll
[coveralls-svg]: https://img.shields.io/coveralls/osteele/gojekyll.svg?branch=master

[license-url]: https://github.com/osteele/gojekyll/blob/master/LICENSE
[license-svg]: https://img.shields.io/badge/license-MIT-blue.svg

[go-report-card-url]: https://goreportcard.com/report/github.com/osteele/gojekyll
[go-report-card-svg]:  https://goreportcard.com/badge/github.com/osteele/gojekyll

[travis-url]: https://travis-ci.org/osteele/gojekyll
[travis-svg]: https://img.shields.io/travis/osteele/gojekyll.svg?branch=master
