# Contributing

Here's some ways to help:

* Try using gojekyll on your site. Use this as fodder for test cases.
* Choose an item to work on from the [issues list](https://github.com/osteele/gojekyll/issues).
* Search the sources for FIXME and TODO comments.
* Improve the [code coverage](https://coveralls.io/github/osteele/gojekyll?branch=master).

If you choose to contribute code, please review the [pull request template](https://github.com/osteele/gojekyll/blob/master/.github/PULL_REQUEST_TEMPLATE.md) before you get too far along.

## Developer Cookbook

### Set up your machine

Fork and clone the repo.

[Install go](https://golang.org/doc/install#install). On macOS running Homebrew,
`brew install go` is easier than the linked instructions.

Install package dependencies and development tools:

```bash
make setup
go get -t ./...
```

[Install golangci-lint](https://golangci-lint.run/usage/install/#local-installation).
On macOS: `brew install golangci-lint`

### Test and Lint

```bash
make test
make lint
```

### Debugging tools

```bash
gojekyll -s path/to/site render index.md              # render a file to stdout
gojekyll -s path/to/site render page.md               # render a file to stdout
gojekyll -s path/to/site render /                     # render a URL to stdout
gojekyll -s path/to/site variables /                  # print a file or URL's variables
gojekyll -s path/to/site variables site               # print the site variables
gojekyll -s path/to/site variables site.twitter.name  # print a specific site variable
```

`./scripts/gojekyll` is an alternative to the `gojekyll` executable, that uses `go run` each time it's invoked.

### Coverage

```bash
./scripts/coverage && go tool cover -html=coverage.out
```

### Profiling

```bash
gojekyll -s path/to/site benchmark
go tool pprof --web gojekyll gojekyll.prof
```
