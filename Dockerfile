FROM golang:latest AS gojekyll

ADD . /gojekyll

WORKDIR /gojekyll

RUN go build main.go

FROM bufbuild/buf AS buf
FROM dart:stable AS sass

COPY --from=buf /usr/local/bin/buf /usr/local/bin/

RUN git clone https://github.com/sass/dart-sass.git /dart-sass && \
    cd /dart-sass && \
    dart pub get && \
    dart run grinder protobuf && \
    dart compile exe bin/sass.dart

FROM cgr.dev/chainguard/glibc-dynamic:latest

COPY --from=gojekyll /gojekyll/main /usr/bin/gojekyll
COPY --from=sass /dart-sass/bin/sass.exe /usr/bin/sass

ENTRYPOINT [ "/usr/bin/gojekyll" ]

CMD [ "--help" ]