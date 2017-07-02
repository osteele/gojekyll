# Gojekyll

[![Build Status](https://travis-ci.org/osteele/gojekyll.svg?branch=master)](https://travis-ci.org/osteele/gojekyll)
[![Go Report Card](https://goreportcard.com/badge/github.com/osteele/gojekyll)](https://goreportcard.com/report/github.com/osteele/gojekyll)

Gojekyll is an incomplete implementation of the [Jekyll](https://jekyllrb.com) static site generator, in the [Go](https://golang.org) programming language.

This was my “weekend project” for learning Go. It took on a life of its own (and more than a weekend).

<!-- TOC -->

- [Gojekyll](#gojekyll)
    - [Installation](#installation)
    - [Usage](#usage)
    - [Status](#status)
        - [Major Omissions](#major-omissions)
        - [Other Caveats](#other-caveats)
        - [Intentional Differences](#intentional-differences)
        - [Timings](#timings)
        - [Features](#features)
    - [Contributing](#contributing)
        - [Testing](#testing)
        - [Profiling](#profiling)
    - [Credits](#credits)
    - [Related](#related)
    - [License](#license)

<!-- /TOC -->

## Installation

1. [Install go](https://golang.org/doc/install#install). On macOS running Homebrew, `brew install go` is easier than the linked instructions.
2. `go get -u osteele/gojekyll/cmd/gojekyll`
3. To use the `{% highlight %}` tag, you need to install Pygments: `pip install Pygments`.

## Usage

```bash
gojekyll -s path/to/site build                # builds into ./_site
gojekyll -s path/to/site serve                # serves from memory, w/ live reload
gojekyll help
gojekyll help build
```

## Status

### Major Omissions

- Major features: themes, page tags, excerpts, plugins (except for a few listed below), pagination, math, warning mode.
- Site variables: `pages`, `static_files`, `html_pages`, `html_files`, `documents`, and `tags`
- Jekyll filters: `group_by_exp`, `pop`, `shift`, `cgi_escape`, `uri_escape`, `scssify`, and `smartify`.
- Jekyll's `include_relative` tag
- The Go Liquid engine is also missing some tags and a few filters. See its [README](https://github.com/osteele/gojekyll/#status) for status.
- Data files must be YAML. CSV and JSON are not (yet) supported.
- `{% highlight %}` uses Pygments. There's no way to tell it to use Rouge. Also, I don't know what will happen if Pygments isn't installed.
- `<div markdown=1>` doesn't work. I think this is a limitation of the Blackfriday Markdown processor.

### Other Caveats

- This Liquid engine is probably much stricter than the original.
- Both the Liquid and Jekyll implementations are much less robust than the originals. They probably panic or otherwise fails on a lot of legitimate constructs.
- Liquid errors aren't reported very nicely.

### Intentional Differences

- `serve` generates pages on the fly; it doesn't write to the file system.
- Files are cached to `/tmp/gojekyll-${USER}`, not `./.sass-cache`
- Server live reload is always on.
- The server reloads the `_config.yml` (and the rest of the site) when that file changes.
- `build` with no `-d` option resolves the destination relative to the source directory, not the current directory.

### Timings

`[go]jekyll -s jekyll/docs build`

| Executable | Options                              | Time   |
|------------|--------------------------------------|--------|
| jekyll     |                                      | 18.53s |
| gojekyll   | single-threaded; cold cache          | 6.85s  |
| gojekyll   | single-threaded; warm cache          | 0.61s  |
| gojekyll   | multi-threaded; cache doesn't matter | 0.34s  |

There's currently no way to disable concurrency or the cache. They were switched off by editing the code for these timings.

The cache is for calls to Pygments (via the `highlight` tag).

### Features

- [ ] Content
  - [x] Front Matter
  - [ ] Posts
    - [x] Categories
    - [ ] Tags
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
      - [ ] `group_by_exp` `pop` `shift` `cgi_escape` `uri_escape` `scssify` `smartify`
      - [x] everything else
    - [ ] Jekyll tags
      - [x] `include`
      - [ ] `include_relative`
      - [x] `link`
      - [x] `post_url`
      - [ ] `gist`
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
    - [ ] `--config`, `--baseurl`, `--lsi`, `--watch`, etc.
    - [ ] won't implement: `--force-polling`, `--limit-posts`, `--incremental`, `JEKYLL_ENV=production`
  - [x] `clean`
  - [ ] `doctor`
  - [x] `help`
  - [ ] `import`
  - [ ] `new`
  - [ ] `new-theme`
  - [x] `serve`
    - [x] `--open-uri`
    - [ ] `--detach`, `--host`, `--port`, `--baseurl`
    - [ ] won't implement: `--incremental`, `--ssl-*`
- [ ] Windows

## Contributing

Install package dependencies and development tools:

```bash
make setup
```

### Testing

```bash
make test
make lint
gojekyll  -s path/to/site render index.md      # render a file to stdout
gojekyll  -s path/to/site render /             # render a URL to stdout
gojekyll  -s path/to/site variables /          # print a file or URL's variables
./scripts/coverage && go tool cover -html=coverage.out
```

`./scripts/gojekyll` is an alternative to the `gojekyll` executable, that uses `go run` each time it's invoked.

### Profiling

```bash
gojekyll -s path/to/site profile
go tool pprof gojekyll gojekyll.prof
```

## Credits

Gojekyll uses these libraries:

| Package                                                                        | Author(s)                                              | Description                                           | License                           |
|--------------------------------------------------------------------------------|--------------------------------------------------------|-------------------------------------------------------|-----------------------------------|
| [github.com/jaschaephraim/lrserver](https://github.com/jaschaephraim/lrserver) | Jascha Ephraim                                         | Live Reload server                                    | MIT                               |
| [github.com/osteele/liquid](https://github.com/osteele/liquid)                 | yours truly                                            | Liquid processor                                      | MIT                               |
| [github.com/pkg/browser](https://github.com/pkg/browser)                       | [pkg](https://github.com/pkg)                          | The `serve -o` option to open the site in the browser | BSD 2-clause "Simplified" License |
| [github.com/russross/blackfriday](https://github.com/russross/blackfriday)     | Russ Ross                                              | Markdown processor                                    | Simplified BSD License            |
| [github.com/sass/libsass](https://github.com/sass/libsass)                     | Listed [here](https://https://github.com/sass/libsass) | C port of the Ruby SASS compiler                      | MIT                               |
| [github.com/wellington/go-libsass](https://github.com/wellington/go-libsass)   | Drew Wells                                             | Go bindings for **libsass**                           | ???                               |
| [gopkg.in/alecthomas/kingpin.v2](https://github.com/alecthomas/kingpin)        | Alec Thomas                                            | command line and flag parser                          | MIT                               |
| [gopkg.in/yaml.v2](https://github.com/go-yaml/yaml)                            | Canonical                                              | YAML support                                          | Apache License 2.0                |

In addition to being totally and obviously inspired by Jekyll, Jekyll's *documentation* was solid and indispensible. Many of the filter test cases are taken directly from the Jekyll documentation, and during development the [Jekyll docs](https://jekyllrb.com/docs/home/) were always open in at least one tab.

The text for `gojekyll help` was taken from the output of `jekyll help`.

The gopher image in the `testdata` directory is from [Wikimedia Commons](https://commons.wikimedia.org/wiki/File:Gophercolor.jpg). It is used under the [Creative Commons Attribution-Share Alike 3.0 Unported license](https://creativecommons.org/licenses/by-sa/3.0/deed.en).

## Related

[Hugo](https://gohugo.io) is *the* pre-eminent Go static site generator. It isn't Jekyll-compatible (-), but it's extraordinarily polished, performant, and productized (+++).

[Liquid](https://github.com/osteele/liquid) is a Go implementation of Liquid templates. I wrote it for gojekyll, but it's implemented as a standalone library.

[Jekyll](https://jekyllrb.com), of course.

## License

MIT
