FROM alpine:latest

RUN apk add libstdc++ libgcc

COPY gojekyll /usr/local/bin/gojekyll
  
EXPOSE 4000

ENTRYPOINT ["/usr/local/bin/gojekyll"]

CMD [ "--help" ]
