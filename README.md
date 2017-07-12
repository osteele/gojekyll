# Gojekyll

 [![][travis-svg]][travis-url] [![][coveralls-svg]][coveralls-url] [![][go-report-card-svg]][go-report-card-url] [![][license-svg]][license-url]

> “It is easier to write an incorrect program than understand a correct one.” - Alan Perlis

Gojekyll is a clone of the [Jekyll](https://jekyllrb.com) static site generator, written in the [Go](https://golang.org) programming language. It provides `build` and `serve` commands, with directory watch and live reload for the latter.

| &nbsp;                  | Gojekyll                                  | Jekyll | Hugo         |
|-------------------------|-------------------------------------------|--------|--------------|
| Stable                  |                                           | ✓      | ✓            |
| Fast                    | ✓<br>[[~20×Jekyll](./docs/benchmarks.md)] |        | ✓            |
| Template language       | Liquid                                    | Liquid | Go templates |
| Jekyll compatibility    | [partial](#current-limitations)           | ✓      |              |
| Plugins                 | [some](./docs/plugins.md)                 | yes    | ?            |
| Runs on Windows         |                                           | ✓      | ✓            |
| Implementation language | Go                                        | Ruby   | Go           |



<!-- TOC -->

- [Gojekyll](#gojekyll)
    - [Installation](#installation)
    - [Usage](#usage)
    - [Status](#status)
        - [Current Limitations](#current-limitations)
        - [Other Differences](#other-differences)
        - [Feature Checklist](#feature-checklist)
    - [Contributing](#contributing)
    - [Attribution](#attribution)
    - [Related](#related)
    - [License](#license)

<!-- /TOC -->

## Installation

First-time install:

1. [Install go](https://golang.org/doc/install#install). (On macOS running [Homebrew](https://brew.sh), `brew install go` is easier than the Install instructions on the Go site.)
2. `go get osteele/gojekyll/cmd/gojekyll`
3. To use the `{% highlight %}` tag, you also need [Pygments](http://pygments.org). `pip install Pygments`.

Update to the latest version:

* `go get -u github.com/osteele/liquid github.com/osteele/gojekyll/cmd/gojekyll`

## Usage

```bash
gojekyll build       # builds the site in the current directory into _site
gojekyll serve       # serve the app at http://localhost:4000; reload on changes
gojekyll help
gojekyll help build
```

## Status

This project is at an early stage of development.

It works on the Google Pages sites that I care about, and it looks credible on a spot-check of other Jekyll sites.

It doesn't run on Windows, and it may not work for you.

In addition to the limitations listed below, this software isn't robust. Jekyll, Hugo, and other mature projects have lots of test coverage, and have had lots of testing by lots of people. I've this used only in limited ways, in the month since I started writing it.

### Current Limitations

Missing features:

- Themes
- Excerpts
- Pagination
- Math
- CSV and JSON data files
- Plugins. (Some plugins are [emulated](./docs/plugins.md).)
- Template filters `group_by_exp` and `scssify`
- More Liquid tags and filters, listed [here](https://github.com/osteele/liquid#differences-from-liquid).
- Markdown features: [attribute lists](https://kramdown.gettalong.org/syntax.html#attribute-list-definitions), [automatic ids](https://kramdown.gettalong.org/converter/html.html#auto-ids), [`markdown=1`](https://kramdown.gettalong.org/syntax.html#html-blocks).
- `site.data` is not sorted.
- Windows compatibility

Also see the [detailed status](#feature-status) below.

### Other Differences

These will probably not change:

By design:

- Having the wrong type in a `_config.yml` is an error.
- Plugins must be listed in the config file, not a Gemfile.
- `serve` generates pages on the fly; it doesn't write to the file system.
- Files are cached to `/tmp/gojekyll-${USER}`, not `./.sass-cache`
- Server live reload is always on.
- `serve --watch` (the default) reloads `_config.yml` too.
- Jekyll provides an (undocumented) `jekyll.version` variable to templates. Copying this didn't seem right.
- Incremental build. The emphasis is on optimizing the non-incremental case. This is easier with fewer code paths.

Muzukashii:

- An extensible plugin mechanism – support for plugins that aren't compiled into the executable.

### Feature Checklist

- [ ] Content
  - [x] Front Matter
  - [x] Posts
    - [x] Categories
    - [x] Tags
    - [x] Drafts
    - [x] Future
    - [x] Related
  - [x] Static Files
  - [x] Variables
  - [x] Collections
  - [ ] Data Files
    - [ ] CSV
    - [ ] JSON
    - [x] YAML
  - [ ] Assets
    - [ ] Coffeescript
    - [x] Sass/SCSS
- [ ] Customization
  - [x] Templates
    - [ ] Jekyll filters
      - [ ] `group_by_exp` and `scssify`
      - [x] everything else
    - [x] Jekyll tags
      - [x] `include`
      - [x] `include_relative`
      - [x] `link`
      - [x] `post_url`
      - [x] `highlight`
  - [x] Includes
      - [x] `include` parameters
      - [x] `include` variables (e.g. `{% include {{ expr }} %}`)
  - [x] Permalinks
  - [ ] Pagination
  - [ ] Plugins
    - [x] `jekyll-avatar`
    - [ ] `jekyll-coffeescript`
    - [x] `jekyll-gist` (ignores `noscript: false`)
    - [x] `jekyll-live-reload` (always on)
    - [ ] `jekyll-paginate`
  - [ ] Themes
  - [x] Layouts
- [x] Server
  - [x] Directory watch
- [ ] Commands
  - [x] `build`
    - [x] `--source`, `--destination`, `--drafts`, `--future`, `--unpublished`
    - [ ] `--config`, `--baseurl`, `--lsi`, `--no-watch`
    - [ ] not planned: `--force-polling`, `--limit-posts`, `--incremental`, `JEKYLL_ENV=production`
  - [x] `clean`
  - [x] `help`
  - [x] `serve`
    - [x] `--open-uri`, `--host`, `--port`
    - [ ] `--detach`, `--baseurl`
    - [ ] not planned: `--incremental`, `--ssl`-*
  - [ ] not planned: `doctor`, `import`, `new`, `new-theme`
- [ ] Windows

Also see:

- The [feature parity board](https://github.com/osteele/gojekyll/projects/1) board  lists differences between Jekyll and gojekyll in more detail.
- The [plugin board](https://github.com/osteele/gojekyll/projects/2) lists the implementation status of common plugins. (Gojekyll lacks an extensible plugin mechanism. The goal is to be able to use it to build Jekyll sites that use the most popular plugins.)
- The [Go Liquid feature parity board](https://github.com/osteele/liquid/projects/1) to see differences between the real Liquid library and the one that is used in gojekyll.

## Contributing

Bug reports, test cases, and code contributions are more than welcome.
Please refer to the [contribution guidelines](./CONTRIBUTING.md).

## Attribution

Gojekyll uses these libraries:

| Package                                                                        | Author(s)                                              | Usage                            | License                                 |
|--------------------------------------------------------------------------------|--------------------------------------------------------|----------------------------------|-----------------------------------------|
| [github.com/jaschaephraim/lrserver](https://github.com/jaschaephraim/lrserver) | Jascha Ephraim                                         | Live Reload                      | MIT License                             |
| [github.com/dchest/cssmin](https://github.com/dchest/cssmin)                   | Dmitry Chestnykh                                       | CSS minimization                 | BSD 3-clause "New" or "Revised" License |
| [github.com/kyokomi/emoji](https://github.com/kyokomi/emoji)                   | kyokomi                                                | `jemoji` plugin emulation        | MIT License                             |
| [github.com/osteele/liquid](https://github.com/osteele/liquid)                 | yours truly                                            | Liquid processor                 | MIT License                             |
| [github.com/pkg/browser](https://github.com/pkg/browser)                       | [pkg](https://github.com/pkg)                          | `serve --open-url` option        | BSD 2-clause "Simplified" License       |
| [github.com/russross/blackfriday](https://github.com/russross/blackfriday)     | Russ Ross                                              | Markdown processing              | Simplified BSD License                  |
| [github.com/sass/libsass](https://github.com/sass/libsass)                     | Listed [here](https://https://github.com/sass/libsass) | C port of the Ruby SASS compiler | MIT License                             |
| [github.com/wellington/go-libsass](https://github.com/wellington/go-libsass)   | Drew Wells                                             | Go bindings for **libsass**      | ???                                     |
| [gopkg.in/alecthomas/kingpin.v2](https://github.com/alecthomas/kingpin)        | Alec Thomas                                            | command-line arguments           | MIT License                             |
| [gopkg.in/yaml.v2](https://github.com/go-yaml/yaml)                            | Canonical                                              | YAML support                     | Apache License 2.0                      |

In addition, the following pieces of text were taken from Jekyll and its plugins.
They are used under the terms of the MIT License.

| Source                                                                          | Use                  | Description            |
|---------------------------------------------------------------------------------|----------------------|------------------------|
| `jekyll help` command                                                           | `gojekyll help` text | help text              |
| [Jekyll template documentation](https://jekyllrb.com/docs/templates/)           | test cases           | filter examples        |
| [`jekyll-feed` plugin](https://github.com/jekyll/jekyll-feed)                   | plugin emulation     | `feed.xml` template    |
| [`jekyll-redirect-from` plugin](https://github.com/jekyll/jekyll-redirect-from) | plugin emulation     | redirect page template |
| [`jekyll-sitemap` plugin](https://github.com/jekyll/jekyll-redirect-from)       | plugin emulation     | sitemap template       |
| [`jekyll-seo-tag` plugin](https://github.com/jekyll/jekyll-redirect-from)       | plugin emulation     | feed template          |

The gopher image in the `testdata` directory is from [Wikimedia Commons](https://commons.wikimedia.org/wiki/File:Gophercolor.jpg). It is used under the [Creative Commons Attribution-Share Alike 3.0 Unported license](https://creativecommons.org/licenses/by-sa/3.0/deed.en).

In addition to being totally and obviously inspired by Jekyll and its plugins, Jekyll's  solid *documentation* was indispensible --- especially since I wanted to implement Jekyll as documented, not port its source code. The [Jekyll docs](https://jekyllrb.com/docs/home/) were always open in at least one tab during development.

## Related

[Hugo](https://gohugo.io) is the pre-eminent Go static site generator. It isn't Jekyll-compatible (-), but it's extraordinarily polished, performant, and productized (+++).

[jkl](https://github.com/drone/jkl) is another Go clone of Jekyll. If I'd found it sooner I might have started this project by forking that one. It's got a better name, too.

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
