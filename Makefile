SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY = gojekyll
PACKAGE = github.com/osteele/gojekyll
COMMIT_HASH = `git rev-parse --short HEAD 2>/dev/null`

LDFLAGS=-ldflags "-X ${PACKAGE}.commands.Version=${COMMIT_HASH}"

.DEFAULT_GOAL: build
.PHONY: build clean deps setup install lint test help

$(BINARY): $(SOURCES)
	go build ${LDFLAGS} -o ${BINARY} ${PACKAGE}

build: $(BINARY) ## compile the package

clean: ## remove binary files
	rm -fI ${BINARY}

imports: ## list imports
	go list -f '{{join .Imports "\n"}}' ./... | grep -v `go list -f '{{.ImportPath}}'` | grep '\.' | sort | uniq

deps: ## list dependencies
	go list -f '{{join .Deps "\n"}}' ./... | grep -v `go list -f '{{.ImportPath}}'` | grep '\.' | sort | uniq

race: ## build a binary with race detection
	go build -race ${LDFLAGS} -o ${BINARY}-race ${PACKAGE}

setup: ## install dependencies and development tools
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install

install: ## compile and install the executable
	go install ${LDFLAGS} ${PACKAGE}/cmd/gojekyll

lint: ## Run all the linters
	gometalinter ./... --disable=gotype --disable aligncheck

test: ## test the package
	go test ./...

# Source: https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
