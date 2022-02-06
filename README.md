# Gojekyll
<!-- ALL-CONTRIBUTORS-BADGE:START - Do not remove or modify this section -->
[![All Contributors](https://img.shields.io/badge/all_contributors-3-orange.svg?style=flat-square)](#contributors-)
<!-- ALL-CONTRIBUTORS-BADGE:END -->

 [![Travis badge][travis-svg]][travis-url]
 [![Golangci-lint badge][golangci-lint-svg]][golangci-lint-url]
 [![Coveralls badge][coveralls-svg]][coveralls-url]
 [![Go Report Card badge][go-report-card-svg]][go-report-card-url]
 [![MIT License][license-svg]][license-url]

Gojekyll is a partially-compatible clone of the [Jekyll](https://jekyllrb.com)
static site generator, written in the [Go](https://golang.org) programming
language. It provides `build` and `serve` commands, with directory watch and
live reload.

| &nbsp;                  | Gojekyll                                  | Jekyll | Hugo                         |
| ----------------------- | ----------------------------------------- | ------ | ---------------------------- |
| Stable                  |                                           | ✓      | ✓                            |
| Fast                    | ✓<br>([~20×Jekyll](./docs/benchmarks.md)) |        | ✓                            |
| Template language       | Liquid                                    | Liquid | Go, Ace and Amber templates  |
| SASS                    | ✓                                         | ✓      | ✓                            |
| Jekyll compatibility    | [partial](#current-limitations)           | ✓      |                              |
| Plugins                 | [some](./docs/plugins.md)                 | yes    | shortcodes, theme components |
| Windows support         |                                           | ✓      | ✓                            |
| Implementation language | Go                                        | Ruby   | Go                           |

<!-- TOC -->

- [Usage](#usage)
- [Installation](#installation)
  - [Binary Downloads](#binary-downloads)
  - [From Source](#from-source)
- [[Optional] Install command-line autocompletion](#optional-install-command-line-autocompletion)
- [Status](#status)
  - [Current Limitations](#current-limitations)
  - [Other Differences](#other-differences)
  - [Feature Checklist](#feature-checklist)
- [Troubleshooting](#troubleshooting)
- [Contributors](#contributors)
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

1. Ubuntu (64-bit) and macOS binaries are available from the [releases
   page](https://github.com/osteele/gojekyll/releases).
2. [Optional] **Themes**. To use a theme, you need to install Ruby and
   [bundler](http://bundler.io/). Create a `Gemfile` that lists the theme., and
   run `bundle install`. The [Jekyll theme
   instructions](https://jekyllrb.com/docs/themes/) provide more detail, and
   should work for Gojekyll too.

### From Source

Pre-requisites:

1. **Install go** (1) via [Homebrew](https://brew.sh): `brew install go`; or (2)
   [download](https://golang.org/doc/install#tarball).
2. See items (2-3) under **Binary Downloads**, above, for optional installations.

First-time install:

```bash
go get github.com/osteele/gojekyll
```

Update to the latest version:

```bash
go get -u github.com/osteele/liquid github.com/osteele/gojekyll
```

## [Optional] Install command-line autocompletion

Add this to your `.bashrc` or `.zshrc`:

```bash
# Bash:
eval "$(gojekyll --completion-script-bash)"
# Zsh:
eval "$(gojekyll --completion-script-zsh)"
```

## Status

This project works on the GitHub Pages sites that I and other contributors care
about. It looks credible on a spot-check of other Jekyll sites.

### Current Limitations

Missing features:

- Pagination
- Windows compatibility
- Math
- Plugin system. ([Some individual plugins](./docs/plugins.md) are emulated.)
- Liquid filter `sassify` is not implemented
- Liquid is run in strict mode: undefined filters and variables are errors.
- Missing markdown features:
  - [attribute lists](https://kramdown.gettalong.org/syntax.html#attribute-list-definitions)
  - [`markdown="span"`, `markdown="block"`](https://kramdown.gettalong.org/syntax.html#html-blocks)
  - Markdown configuration options

Also see the [detailed status](#feature-status) below.

### Other Differences

These will probably not change:

By design:

- Plugins must be listed in the config file, not a Gemfile.
- The wrong type in a `_config.yml` file – for example, a list where a string is
  expected, or vice versa – is generally an error.
- Server live reload is always on.
- `serve --watch` (the default) reloads the `_config.yml` and data files too.
- `serve` generates pages on the fly; it doesn't write to the file system.
- Files are cached in `/tmp/gojekyll-${USER}`, not `./.sass-cache`

Upstream:

- Markdown:
  - `<` and `>` inside markdown is interpreted as HTML. For example, `This is
    <b>bold</b>` renders as <b>bold</b>. This behavior matches the [Markdown
    spec](https://daringfireball.net/projects/markdown/syntax#html), but differs
    from Jekyll's default Kramdown processor.
  - The autogenerated id of a header that includes HTML is computed from the
    text of the title, ignoring its attributes. For example, the id of `## Title
    (<a href="https://example.com/path/to/details">ref</a>))` is `#title-ref`,
    not `#title-https-example-path-to-details-ref`.
  - Autogenerated header ids replace punctuation by the hyphens, rather than the
    empty string. For example, the id of `## Either/or` is `#either-or` not
    `#eitheror`; the id of `## I'm Lucky` is `#i-m-lucky` not `#im-lucky`.

Muzukashii:

- An extensible plugin mechanism – support for plugins that aren't compiled into
  the executable.

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
      - [ ] `scssify`
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
    - [x] `--incremental`, `--watch`, `--force_polling`, `JEKYLL_ENV=production`
    - [ ] `--baseurl`, `--config`, `--lsi`
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

## Troubleshooting

If the error is "403 API rate limit exceeded", you are probably building a
repository that uses the `jekyll-github-metadata` gem. Try setting the
`JEKYLL_GITHUB_TOKEN`, `JEKYLL_GITHUB_TOKEN`, or `OCTOKIT_ACCESS_TOKEN`
environment variable to the value of a [GitHub personal access
token][personal-access-token] and trying again.

[personal-access-token]: https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token

## Contributors

Thanks goes to these wonderful people ([emoji key](https://allcontributors.org/docs/en/emoji-key)):

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tr>
    <td align="center"><a href="https://code.osteele.com/"><img src="https://avatars.githubusercontent.com/u/674?v=4?s=100" width="100px;" alt=""/><br /><sub><b>Oliver Steele</b></sub></a><br /><a href="https://github.com/osteele/gojekyll/commits?author=osteele" title="Code">💻</a> <a href="#design-osteele" title="Design">🎨</a> <a href="https://github.com/osteele/gojekyll/commits?author=osteele" title="Documentation">📖</a> <a href="#ideas-osteele" title="Ideas, Planning, & Feedback">🤔</a> <a href="#infra-osteele" title="Infrastructure (Hosting, Build-Tools, etc)">🚇</a> <a href="#maintenance-osteele" title="Maintenance">🚧</a> <a href="#projectManagement-osteele" title="Project Management">📆</a> <a href="https://github.com/osteele/gojekyll/pulls?q=is%3Apr+reviewed-by%3Aosteele" title="Reviewed Pull Requests">👀</a> <a href="https://github.com/osteele/gojekyll/commits?author=osteele" title="Tests">⚠️</a></td>
    <td align="center"><a href="https://bep.is/"><img src="https://avatars.githubusercontent.com/u/394382?v=4?s=100" width="100px;" alt=""/><br /><sub><b>Bjørn Erik Pedersen</b></sub></a><br /><a href="https://github.com/osteele/gojekyll/commits?author=bep" title="Documentation">📖</a></td>
    <td align="center"><a href="https://tqdev.com/"><img src="https://avatars.githubusercontent.com/u/1288217?v=4?s=100" width="100px;" alt=""/><br /><sub><b>Maurits van der Schee</b></sub></a><br /><a href="https://github.com/osteele/gojekyll/commits?author=mevdschee" title="Code">💻</a></td>
  </tr>
</table>

<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->

<!-- ALL-CONTRIBUTORS-LIST:END -->

This project follows the
[all-contributors](https://github.com/all-contributors/all-contributors)
specification. [Contributions of any kind welcome](./CONTRIBUTING.md)!

## Attribution

Gojekyll uses these libraries:

| Package                                                                        | Author(s)                                              | Usage                                  | License                                 |
| ------------------------------------------------------------------------------ | ------------------------------------------------------ | -------------------------------------- | --------------------------------------- |
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
| ------------------------------------------------------------------------------- | -------------------- | ---------------------- |
| [Jekyll template documentation](https://jekyllrb.com/docs/templates/)           | test cases           | filter examples        |
| `jekyll help` command                                                           | `gojekyll help` text | help text              |
| [`jekyll-feed` plugin](https://github.com/jekyll/jekyll-feed)                   | plugin emulation     | `feed.xml` template    |
| [`jekyll-redirect-from` plugin](https://github.com/jekyll/jekyll-redirect-from) | plugin emulation     | redirect page template |
| [`jekyll-sitemap` plugin](https://github.com/jekyll/jekyll-redirect-from)       | plugin emulation     | sitemap template       |
| [`jekyll-seo-tag` plugin](https://github.com/jekyll/jekyll-redirect-from)       | plugin emulation     | feed template          |

The theme for in-browser error reporting was adapted from facebookincubator/create-react-app.

The gopher image in the `testdata` directory is from [Wikimedia
Commons](https://commons.wikimedia.org/wiki/File:Gophercolor.jpg). It is used
under the [Creative Commons Attribution-Share Alike 3.0 Unported
license](https://creativecommons.org/licenses/by-sa/3.0/deed.en).

In addition to being totally and obviously inspired by Jekyll and its plugins,
Jekyll's  solid *documentation* was indispensible --- especially since I wanted
to implement Jekyll as documented, not port its source code. The [Jekyll
docs](https://jekyllrb.com/docs/home/) were always open in at least one tab
during development.

## Related

[Hugo](https://gohugo.io) is the pre-eminent Go static site generator. It isn't
Jekyll-compatible (-), but it's highly polished, performant, and productized
(+++).

[jkl](https://github.com/drone/jkl) is another Go clone of Jekyll. If I'd found
it sooner I might have started this project by forking that one. It's got a
better name.

[Liquid](https://github.com/osteele/liquid) is a pure Go implementation of
Liquid templates, that I finally caved and wrote in order to use in this
project.

[Jekyll](https://jekyllrb.com), of course.

## License

MIT

[coveralls-url]: https://coveralls.io/r/osteele/gojekyll
[coveralls-svg]: https://img.shields.io/coveralls/osteele/gojekyll.svg?branch=master

[license-url]: https://github.com/osteele/gojekyll/blob/master/LICENSE
[license-svg]: https://img.shields.io/badge/license-MIT-blue.svg

[golangci-lint-url]: https://github.com/osteele/gojekyll/actions?query=workflow%3Agolangci-lint
[golangci-lint-svg]: https://github.com/osteele/gojekyll/actions/workflows/golangci-lint.yml/badge.svg

[go-report-card-url]: https://goreportcard.com/report/github.com/osteele/gojekyll
[go-report-card-svg]:  https://goreportcard.com/badge/github.com/osteele/gojekyll

[travis-url]: https://travis-ci.com/osteele/gojekyll
[travis-svg]: https://img.shields.io/travis/osteele/gojekyll.svg?branch=master
