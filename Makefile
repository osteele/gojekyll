BINARY = gojekyll
PACKAGE = github.com/osteele/gojekyll

SOURCEDIR = .
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

COMMIT_HASH = `git rev-parse --short HEAD 2>/dev/null`
BUILD_DATE = `date +%FT%T%z`
VERSION := $(COMMIT_HASH)
OS := $(shell uname -o)

LDFLAGS=-ldflags "-X ${PACKAGE}/version.Version=${VERSION} -X ${PACKAGE}/version.BuildDate=${BUILD_DATE}"

.DEFAULT_GOAL: build
.PHONY: build clean deps setup install lint release test help

$(BINARY): $(SOURCES)
	go build ${LDFLAGS} -o ${BINARY} ${PACKAGE}

build: $(BINARY) ## compile the package

clean:
	rm -f ${BINARY}

imports:
	go list -f '{{join .Imports "\n"}}' ./... | grep -v `go list -f '{{.ImportPath}}'` | grep '\.' | sort | uniq

deps:
	go list -f '{{join .Deps "\n"}}' ./... | grep -v `go list -f '{{.ImportPath}}'` | grep '\.' | sort | uniq

race:
	go build -race ${LDFLAGS} -o ${BINARY}-race ${PACKAGE}

release: build
	mkdir -p dist
	tar -cvzf dist/$(BINARY)_$(VERSION)_$(OS:GNU/%=%)_$(shell uname -m).tar.gz $(BINARY) LICENSE README.md

setup:
	go get -t -u ./...
	go get github.com/alecthomas/gometalinter
	gometalinter --install

install:
	go install ${LDFLAGS} ${PACKAGE}

lint:
	gometalinter ./... --tests --deadline=5m --disable=gotype --disable=aligncheck

test:
	go test ./...
