FROM golang:1.9-alpine AS builder

RUN apk add --no-cache make
COPY . /go/src/github.com/blippar/alpine-package-browser
WORKDIR /go/src/github.com/blippar/alpine-package-browser

RUN make bin/apk

FROM alpine:3.4 AS runtime

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /go/src/github.com/blippar/alpine-package-browser/bin/apk /app/bin/apk
COPY config.example.json   /app/config.json
COPY templates/*.html.tmpl /app/templates/

EXPOSE 8000
ENTRYPOINT ["bin/apk"]
CMD ["-config=config.json"]
