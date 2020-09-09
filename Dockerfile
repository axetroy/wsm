# Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.
# builder for backend
FROM golang:1.15.2-alpine AS builder

WORKDIR /app

COPY main.go go.mod go.sum ./
COPY ./vendor ./vendor
COPY ./internal ./internal
COPY ./cmd/user/main.go ./cmd/user/main.go

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-s -w" -o ./bin/wsm ./cmd/user/main.go

# target
FROM alpine:3.12

WORKDIR /app
COPY --from=builder /app/bin .

ENV PORT=80

EXPOSE 80

CMD ["./wsm", "start"]