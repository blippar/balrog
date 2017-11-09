FROM golang:1.9-alpine AS builder

COPY . /go/src/github.com/0rax/apk
WORKDIR /go/src/github.com/0rax/apk

RUN go build -v -o bin/apk ./cmd/apk/

FROM alpine:3.4 AS runtime

WORKDIR /app
COPY --from=builder /go/src/github.com/0rax/apk/bin/apk /app/bin/apk
COPY config.example.json   /app/config.json
COPY templates/*.html.tmpl /app/bin/templates/

EXPOSE 8000
ENTRYPOINT ["bin/apk"]
CMD ["-config=config.json"]
