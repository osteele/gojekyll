FROM bufbuild/buf AS buf
FROM dart:stable AS sass

COPY --from=buf /usr/local/bin/buf /usr/local/bin/

RUN git clone https://github.com/sass/dart-sass.git /dart-sass && \
    cd /dart-sass && \
    dart pub get && \
    dart run grinder protobuf && \
    dart compile exe bin/sass.dart


FROM golangci/golangci-lint:latest AS golangci-lint
FROM golang:latest AS gojekyll

ADD . /gojekyll

COPY --from=golangci-lint /usr/bin/golangci-lint /usr/bin/golangci-lint
COPY --from=sass /dart-sass/bin/sass.exe /usr/bin/sass

WORKDIR /gojekyll

RUN go test ./...
RUN golangci-lint run
RUN go build main.go

FROM debian:stable-slim

COPY --from=gojekyll /gojekyll/main /usr/bin/gojekyll
COPY --from=sass /dart-sass/bin/sass.exe /usr/bin/sass

WORKDIR /app

ENTRYPOINT [ "/usr/bin/gojekyll" ]

CMD [ "--help" ]
