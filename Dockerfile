FROM docker.io/golang:1.22-alpine3.19 AS builder

WORKDIR /go

COPY . .

RUN go build -o main main.go

FROM docker.io/alpine:3.19

WORKDIR /go

COPY --from=builder /go/main main

RUN apk add --no-cache dumb-init

ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["/go/main"]