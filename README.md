# Go Jekyll

When I grow up, I want to be a [Go](https://golang.org) implementation of [Jekyll](https://jekyllrb.com).

## Status
[![Build Status](https://travis-ci.org/osteele/gojekyll.svg?branch=master)](https://travis-ci.org/osteele/gojekyll)
[![Go Report Card](https://goreportcard.com/badge/github.com/osteele/gojekyll)](https://goreportcard.com/report/github.com/osteele/gojekyll)

See the significant missing functionality. This currently works on some simple sites that don't use drafts, templates, future posts, or various other features listed below.

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
      - [ ] Sass caching
- [ ] Customization
  - [x] Templates
    - [ ] Jekyll filters (partial)
    - [ ] Jekyll tags (partial)
    - [ ] `markdown=1`
  - [x] Includes
    - [ ] `include_relative`
    - [x] parameters
      - [ ] variable parameters
    - [ ] variables `{% include {{ expr }} %}`
  - [x] Permalinks
  - [ ] Pagination
  - [ ] Themes
  - [x] Layouts
- [x] Server
  - [x] Directory watch
  - [x] Live reload
- [ ] Windows -- not tested

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

Jascha Ephraim's [jaschaephraim/lrserver](https://github.com/jaschaephraim/lrserver) Live Reload server.

<https://github.com/pkg/browser> to open the URL in a browser.

The gopher image in the test directory is from [Wikimedia Commons](https://commons.wikimedia.org/wiki/File:Gophercolor.jpg). It is used under the [Creative Commons Attribution-Share Alike 3.0 Unported license](https://creativecommons.org/licenses/by-sa/3.0/deed.en).

## Related

[Hugo](https://gohugo.io) isn't Jekyll-compatible (-), but actually works (+++).

[Jekyll](https://jekyllrb.com), of course.

This project is a clean-room implementation of Jekyll, based solely on Jekyll's documentation and testing it against a few sites. Hopefully this can pay off in contributing towards Jekyll's documentation.

## License

MIT

## Alternate Naming Possibilities

* "Gekyll". (Hard or soft "g"? See [gif](https://en.wikipedia.org/wiki/GIF#Pronunciation_of_GIF).)
* "Gekko"
