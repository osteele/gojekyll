# Contributing

Refer to the [Liquid contribution guidelines](https://github.com/Shopify/liquid/blob/master/CONTRIBUTING.md).

In addition to those checklists, I also won't merge:

- [ ] Performance improvements that don't include a benchmark.
- [ ] Meager (<3%) performance improvements that increase code verbosity or complexity.

## Developer Cookbook

### Set up your machine

Fork and clone the repo.

[Install go](https://golang.org/doc/install#install). On macOS running Homebrew, `brew install go` is easier than the linked instructions.

Install package dependencies and development tools:

```bash
make install_dev_tools
go get -t ./...
```

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
gojekyll -s path/to/site benchmark && go tool pprof gojekyll gojekyll.prof
```
