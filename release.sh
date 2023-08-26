#!/bin/bash -e

for arch in amd64 arm64; do
	docker buildx build --platform linux/$arch . \
		-f Dockerfile.base \
		-t gojekyll_base:next-$arch \
		--load &
done

wait

docker buildx imagetools create -t gojekyll_base:latest gojekyll_base:next-{arm64,amd64}

goreleaser release