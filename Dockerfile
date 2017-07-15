FROM golang:1.8.3-alpine

ADD . /go/src/github.com/osteele/gojekyll

RUN \
  apk add --update gcc g++ git make python py-pip && \
  pip install Pygments && \
  cd $GOPATH/src/github.com/osteele/gojekyll && \
  go get -v && \
  make install && \
  rm -rf /var/cache/apk/* && \
  rm -rf $GOPATH/src/*

EXPOSE 4000

ENTRYPOINT ["/go/bin/gojekyll"]

CMD [ "--help" ]
