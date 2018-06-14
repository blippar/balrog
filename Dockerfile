FROM golang:1.9-alpine AS builder

RUN apk add --no-cache make

ARG VERSION="unknown"
ENV BALROG_VERSION="$VERSION"

COPY . /go/src/github.com/blippar/balrog
WORKDIR /go/src/github.com/blippar/balrog

RUN make VERSION="${BALROG_VERSION}" bin/balrog

FROM alpine:3.4 AS runtime

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /go/src/github.com/blippar/balrog/bin/balrog /app/bin/balrog
COPY config.example.json   /app/config.json
COPY templates/*.html.tmpl /app/templates/
COPY dist/                 /app/dist/

EXPOSE 8000
ENTRYPOINT ["bin/balrog"]
CMD ["-config=config.json"]
