#!/bin/bash -e

docker buildx build --platform linux/arm64,linux/amd64 . \
	-f Dockerfile \
	-t danog/gojekyll:latest \
	--push

goreleaser release