# Go Jekyll

When I grow up, I want to be a [Go](https://golang.org) implementation of [Jekyll](https://jekyllrb.com).

## Status

This project is missing more functionality than it implements. It may accidentally work on tiny or simple sites, but I'd be surprised. Most egregious are an insufficiency of template variables, and limitations in the **liquid** library.

I'm writing this to learn my way around Go. It's not good for anytihng yet, and it may never come to anything.

## Install

```bash
go get
```

Sometimes this relies on unmerged improvements to the **acstech/liquid** library. If you want this branch instead:

```bash
cd $(go env GOPATH)/src/github.com/acstech/liquid
git remote set-url origin https://github.com/osteele/liquid.git
git fetch
git reset --hard origin/master
```

## Run

```bash
./scripts/gojekyll --source test build
./scripts/gojekyll --source test serve
./scripts/gojekyll --source test render index.md
./scripts/gojekyll --source test render /
```

`--source DIR` is optional.

`build` needn't be run before `server`. It serves from memory, and doesn't currently rebuild.

`render` renders a single file, identified by permalink if it starts with `/` and by pathname (relative to the source directory) if it doesn't.

## Credits

The [acstech/liquid](https://github.com/acstech/liquid) fork of [karlseguin/liquid](https://github.com/karlseguin/liquid).

## Related

[Hugo](https://gohugo.io) isn't Jekyll-compatible (-), but actually works (+++).

## License

MIT

## Alternate Naming Possibilities

* "Gekyll". (Hard or soft "g"? See [gif](https://en.wikipedia.org/wiki/GIF#Pronunciation_of_GIF).)
* "Gekko"
