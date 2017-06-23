# Go Jekyll

When I grow up, I want to be a [Go](https://golang.org) implementation of [Jekyll](https://jekyllrb.com).

## Status
[![Build Status](https://travis-ci.org/osteele/gojekyll.svg?branch=master)](https://travis-ci.org/osteele/gojekyll)
[![Go Report Card](https://goreportcard.com/badge/github.com/osteele/gojekyll)](https://goreportcard.com/report/github.com/osteele/gojekyll)

This project is missing more functionality than it implements. It may accidentally work on tiny or simple sites.

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
    - [x] link tag
    - [x] include tag
    - [ ] Remaining Jekyll Liquid tags
    - [ ] Jekyll Liquid filters
  - [x] Includes
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

You get slightly better Liquid template parsing from some unmerged pull requests to the **acstech/liquid** library. If you want to use [my fork](https://github.com/osteele/liquid) instead:

```bash
cd $(go env GOPATH)/src/github.com/acstech/liquid
git remote add osteele https://github.com/osteele/liquid.git
git pull -f osteele
```

(See articles by [Shlomi Noach](http://code.openark.org/blog/development/forking-golang-repositories-on-github-and-managing-the-import-path) and [Francesc Campoy](http://blog.campoy.cat/2014/03/github-and-go-forking-pull-requests-and.html) for how this works and why it is necessary.)

## Usage

```bash
gojekyll -s path/to/site build                # builds into ./_site
gojekyll -s path/to/site serve                # serves from memory, w/ live reload
gojekyll help
gojekyll help build
```

### Liquid Template Server

The embedded Liquid server isn't very compliant with Shopfiy Liquid syntax.

You can run a "Liquid Template Server" on the same machine, and tell `gojekyll` to use this instead.
This is currently about 10x slower than using the embedded engine, but still 5x faster than Ruby `jekyll`.

1. Download and run (liquid-template-server)[https://github.com/osteele/liquid-template-server].
2. Invoke `gojekyll` with the `--use-liquid-server` option; e.g.:

  ```bash
  gojekyll --use-liquid-server build
  gojekyll --use-liquid-server serve
  ```

Neither the embedded Liquid server nor the Liquid Template Server implements very many Jekyll Liquid filters or tags. I'm adding to these as necessary to support my own sites.

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

For rendering Liquid templates: ACS Technologies's fork [acstech/liquid](https://github.com/acstech/liquid) of Karl Seguin's [karlseguin/liquid](https://github.com/karlseguin/liquid) Go implementation; or, Jun Yang's JavaScript implementation [harttle/shopify-liquid](https://github.com/harttle/shopify-liquid/) via JSON-RPC.

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
