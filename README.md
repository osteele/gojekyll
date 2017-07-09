# Gojekyll

 [![][travis-svg]][travis-url] [![][coveralls-svg]][coveralls-url] [![][go-report-card-svg]][go-report-card-url] [![][license-svg]][license-url]

Gojekyll is a clone of the [Jekyll](https://jekyllrb.com) static site generator, in the [Go](https://golang.org) programming language.

<!-- TOC -->

- [Gojekyll](#gojekyll)
    - [Installation](#installation)
    - [Usage](#usage)
    - [Current Limitations](#current-limitations)
    - [Other Differences](#other-differences)
    - [Timings](#timings)
    - [Feature Status](#feature-status)
    - [Contributing](#contributing)
    - [Credits](#credits)
    - [Related](#related)
    - [License](#license)

<!-- /TOC -->

## Installation

First-time install:

1. [Install go](https://golang.org/doc/install#install). On macOS running Homebrew, `brew install go` is easier than the linked instructions.
2. `go get osteele/gojekyll/cmd/gojekyll`
3. To use the `{% highlight %}` tag, you need Pygments. `pip install Pygments`.

Update to the latest version:

* `go get -u github.com/osteele/liquid github.com/osteele/gojekyll/cmd/gojekyll`

## Usage

```bash
gojekyll build       # builds the site in the current directory into _site
gojekyll serve       # serve the app at http://localhost:4000
gojekyll help
gojekyll help build
```

## Current Limitations

Missing features:
- Themes
- Excerpts
- Pagination
- Math
- CSV and JSON data files
- Plugins (except `jekyll-avatar`, `jekyll-gist`, and `jekyll-redirect-from`, which are simulated)
- `site-static_files`, `site.html_files`, and `site.tags`
- Jekyll liquid filters: `group_by_exp`, `cgi_escape`, `uri_escape`, `scssify`, and `smartify`
- Windows compatibility

For more detailed status:

- The [feature parity board](https://github.com/osteele/gojekyll/projects/1) board  lists differences between Jekyll and gojekyll in more detail.
- The [plugin board](https://github.com/osteele/gojekyll/projects/2) lists the implementation status of common plugins. (Gojekyll lacks an extensible plugin mechanism. The goal is to be able to use it to build Jekyll sites that use the most popular plugins.)
- The [Go Liquid feature parity board](https://github.com/osteele/liquid/projects/1) to see differences between the real Liquid library and the one that is used in gojekyll.

## Other Differences

These will probably not change.

- `uniq` doesn't work the same way as in Jekyll / Shopify Liquid. See the [Go Liquid differences](https://github.com/osteele/liquid#differences) for more on this.
- Real Jekyll provides an (undocumented) `jekyll.version` variable to templates. Copying this didn't seem right.
- `serve` generates pages on the fly; it doesn't write to the file system.
- Files are cached to `/tmp/gojekyll-${USER}`, not `./.sass-cache`
- Server live reload is always on.
- The server reloads the `_config.yml` too.
- An extensible plugin mechanism: support for plugins that aren't compiled into the executable.
- Incremental build isn't planned. The emphasis is on optimizing the non-incremental case.

## Timings

`[go]jekyll -s jekyll/docs build` on a late-2015 MacBook Pro, running current versions of everything as of 2017-07-01.

| Executable | Options                     | Time          |
|------------|-----------------------------|---------------|
| jekyll     |                             | 18.53s        |
| gojekyll   | single-threaded; cold cache | 3.31s ± 0.07s |
| gojekyll   | single-threaded; warm cache | 2.50s ± 0.04s |
| gojekyll   | multi-threaded              | 1.51s ± 0.01s |
| gojekyll   | multi-threaded              | 0.80s ± 0.09s |

This isn't an apples-to-ranges comparison. Gojekyll doesn't use all the plugins that Jekyll does. In particular, `jekyll-mentions` parses each page's HTML. This could slow Gojekyll down once it's added.

There's currently no way to disable the cache. It was disabled off by re-building the executable.

The cache is for calls to Pygments (via the `highlight` tag). For another site, SASS is greater overhead. This is another candidate for caching, but with multi-threading it may not matter.

## Feature Status

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
    - [ ] `--config`, `--baseurl`, `--lsi`, `--watch`, etc.
    - [ ] not planned: `--force-polling`, `--limit-posts`, `--incremental`, `JEKYLL_ENV=production`
  - [x] `clean`
  - [ ] `doctor`
  - [x] `help`
  - [ ] `import`
  - [ ] `new`
  - [ ] `new-theme`
  - [x] `serve`
    - [x] `--open-uri`
    - [ ] `--detach`, `--host`, `--port`, `--baseurl`
    - [ ] not planned: `--incremental`, `--ssl`-*
- [ ] Windows

## Contributing

Bug reports, test cases, and code contributions are more than welcome.
Please refer to the [contribution guidelines](./CONTRIBUTING.md).

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

The text for `gojekyll help` was taken from the output of `jekyll help`, and is used under the terms of the MIT license.

Many of the filter test cases are taken directly from the Jekyll documentation, and are used under the terms of the MIT license.

The template for `jekyll-redirect-from` page redirects was adapted from the template in <https://github.com/jekyll/jekyll-redirect-from>, and is used under the terms of the MIT license.

The gopher image in the `testdata` directory is from [Wikimedia Commons](https://commons.wikimedia.org/wiki/File:Gophercolor.jpg). It is used under the [Creative Commons Attribution-Share Alike 3.0 Unported license](https://creativecommons.org/licenses/by-sa/3.0/deed.en).

In addition to being totally and obviously inspired by Jekyll and its plugins, Jekyll's *documentation* was solid and indispensible. During development the [Jekyll docs](https://jekyllrb.com/docs/home/) were always open in at least one tab.

## Related

[Hugo](https://gohugo.io) is *the* pre-eminent Go static site generator. It isn't Jekyll-compatible (-), but it's extraordinarily polished, performant, and productized (+++).

[Liquid](https://github.com/osteele/liquid) is a Go implementation of Liquid templates. I wrote it for gojekyll, but it's implemented as a standalone library.

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
