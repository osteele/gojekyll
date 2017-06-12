# Go Jekyll

When I grow up, I want to be a [Go](https://golang.org) implementation of [Jekyll](https://jekyllrb.com).

## Status
[![Go Report Card](https://goreportcard.com/badge/github.com/osteele/gojekyll)](https://goreportcard.com/report/github.com/osteele/gojekyll)

This project is missing more functionality than it implements. It may accidentally work on tiny or simple sites, but I'd be surprised. Most egregious are an insufficiency of template variables, and limitations in the **liquid** library.

I'm writing this to learn my way around Go. It's not good for anything yet, and it may never come to anything.

## Install

```bash
go get -t
```

Sometimes this package benefits from my unmerged improvements to the **acstech/liquid** library. If you want to use [my fork](https://github.com/osteele/liquid) instead:

```bash
cd $(go env GOPATH)/src/github.com/acstech/liquid
git remote add osteele https://github.com/osteele/liquid.git
git pull -f osteele
```

(See articles by [Shlomi Noach](http://code.openark.org/blog/development/forking-golang-repositories-on-github-and-managing-the-import-path) and [Francesc Campoy](http://blog.campoy.cat/2014/03/github-and-go-forking-pull-requests-and.html) for how this works and why it is necessary.)

## Run

```bash
./scripts/gojekyll --source test build
./scripts/gojekyll --source test serve
./scripts/gojekyll --source test render index.md
./scripts/gojekyll --source test render /
```

`--source DIR` is optional.

`build` needn't be run before `server`. The latter serves from memory.

`server` only rebuilds individual changed pages, doesn't rebuild collections, and doesn't detect new pages.

`render` renders a single file, identified by permalink if it starts with `/`, and by pathname (relative to the source directory) if it doesn't.

`./scripts/gojekyll` uses `go run` each time it's invoked. Alternatives to it are: `go build && ./gojekyll ...`; or `go install && gojekyll ...` (if `$GOPATH/bin` is on your `$PATH`). These would be nicer for actual use (where the **gojekyll** sources don't change between invocations), but they aren't as handy during development.

## Credits

For rendering Liquid templates: the [acstech/liquid](https://github.com/acstech/liquid) fork of [karlseguin/liquid](https://github.com/karlseguin/liquid).

## Related

[Hugo](https://gohugo.io) isn't Jekyll-compatible (-), but actually works (+++).

[Jekyll](https://jekyllrb.com), of course.

## License

MIT

## Alternate Naming Possibilities

* "Gekyll". (Hard or soft "g"? See [gif](https://en.wikipedia.org/wiki/GIF#Pronunciation_of_GIF).)
* "Gekko"
