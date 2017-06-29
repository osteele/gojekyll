# Go Jekyll

When I grow up, I want to be a [Go](https://golang.org) implementation of [Jekyll](https://jekyllrb.com).

## Status
[![Build Status](https://travis-ci.org/osteele/gojekyll.svg?branch=master)](https://travis-ci.org/osteele/gojekyll)
[![Go Report Card](https://goreportcard.com/badge/github.com/osteele/gojekyll)](https://goreportcard.com/report/github.com/osteele/gojekyll)

See the significant missing functionality below. This tool currently works on some simple sites that don't use drafts, templates, future posts, or various other features listed below.

We are currently ~5x faster than Jekyll. Some obvious improvements would include caching SASS, caching templates, and using goroutines to render pages.

- [ ] Content
  - [x] Front Matter
  - [ ] Posts
    - [ ] Categories
    - [ ] Tags
    - [ ] Drafts
    - [ ] Future
    - [ ] Related
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
      - [ ] Sass cache
- [ ] Customization
  - [x] Templates
    - [ ] Jekyll filters
      - [ ] array filters: group_by group_by_exp sample pop shift
      - [ ] string filters: cgi_escape uri_escape scssify smartify slugify normalize_whitespace
      - [x] everything else
    - [ ] Jekyll tags
      - [x] `include`
      - [ ] `include_relative`
      - [x] `link`
      - [ ] `post_url`
      - [ ] `gist`
      - [ ] `highlight`
    - [ ] `markdown=1`
  - [x] Includes
      - [x] `include` parameters
      - [ ] `include` variables (e.g. `{% include {{ expr }} %}`)
  - [x] Permalinks
  - [ ] Pagination
  - [ ] Themes
  - [x] Layouts
- [x] Server
  - [x] Directory watch
  - [x] Live reload
- [ ] Windows

Intentional differences from Jekyll:

- `serve` doesn't write to the file system
- No `.sass-cache`. (When caching is added, it will go to a temporary directory.)
- Server live reload is always on.

## Install

```bash
go get -u osteele/gojekyll/cmd/gojekyll
```

## Usage

```bash
gojekyll -s path/to/site build                # builds into ./_site
gojekyll -s path/to/site serve                # serves from memory, w/ live reload
gojekyll help
gojekyll help build
```

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

| Package | Author | Description |
| --- | --- | --- |
| [github.com/jaschaephraim/lrserver](https://github.com/jaschaephraim/lrserver) | Jascha Ephraim | Live Reload server |
| [github.com/osteele/liquid](https://github.com/osteele/liquid) | Oliver Steele | Liquid processor |
| [github.com/pkg/browser](https://github.com/pkg/browser) | [pkg](https://github.com/pkg) | The `serve -o` option to open the site in the browser |
| [github.com/russross/blackfriday](https://github.com/russross/blackfriday) | Russ Ross | Markdown processor |
| [github.com/sass/libsass](https://github.com/sass/libsass) | Listed [here](https://https://github.com/sass/libsass) | C port of the Ruby SASS compiler |
| [github.com/wellington/go-libsass](https://github.com/wellington/go-libsass) | Drew Wells | Go bindings to libsass |
| [gopkg.in/alecthomas/kingpin.v2](https://github.com/alecthomas/kingpin)  | Alec Thomas | command line and flag parser |
| [gopkg.in/yaml.v2](https://github.com/go-yaml) | Canonical | YAML support |

In addition to being totally and obviously inspired by the Jekyll Ruby implementation, Jekyll's solid documentation was indispensible. The [Jekyll docs](https://jekyllrb.com/docs/home/) were always open in at least one tab.

The gopher image in the test directory is from [Wikimedia Commons](https://commons.wikimedia.org/wiki/File:Gophercolor.jpg). It is used under the [Creative Commons Attribution-Share Alike 3.0 Unported license](https://creativecommons.org/licenses/by-sa/3.0/deed.en).

## Related

[Hugo](https://gohugo.io) isn't Jekyll-compatible (-), but actually works (+++).

[Jekyll](https://jekyllrb.com), of course.

## License

MIT

## Alternate Naming Possibilities

* "Gekyll". (Hard or soft "g"? See [gif](https://en.wikipedia.org/wiki/GIF#Pronunciation_of_GIF).)
* "Gekko"
