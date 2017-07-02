SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY = gojekyll
PACKAGE = github.com/osteele/gojekyll
COMMIT_HASH = `git rev-parse --short HEAD 2>/dev/null`
BUILD_TIME = `date +%FT%T%z`

VERSION=0.0.0

LDFLAGS=-ldflags "-X ${PACKAGE}.Version=${VERSION} -X ${PACKAGE}.BuildTime=${BUILD_TIME}"

.DEFAULT_GOAL: $(BINARY)
.PHONY: build clean dependencies setup install lint test help

$(BINARY): $(SOURCES)
	go build ${LDFLAGS} -o ${BINARY} ${PACKAGE}/cmd/gojekyll

build: $(BINARY) ## compile the package

clean: ## remove binary files
	rm -fI ${BINARY}

deps: ## list dependencies
	go list -f '{{join .Imports "\n"}}' ./... | grep -v ${PACKAGE} | grep '\.' | sort | uniq

setup: ## install dependencies and development tools
	go get -t ./...
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install

install: ## compile and install the executable
	go install ${LDFLAGS} ${PACKAGE}/cmd/gojekyll

lint: ## lint the package
	gometalinter ./...

test: ## test the package
	go test ./...

# Source: https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
