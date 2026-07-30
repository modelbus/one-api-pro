FROM --platform=$BUILDPLATFORM node:20 AS builder

WORKDIR /web
COPY ./web /web

RUN npm install && npm run build

FROM golang:alpine AS builder2

RUN apk add --no-cache \
    gcc \
    musl-dev \
    sqlite-dev \
    build-base

ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux

WORKDIR /build

ADD go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=builder /web/build ./web/build

RUN go build -trimpath -ldflags "-s -w -X 'github.com/Leon-PanPan/one-api-pro/common.Version=$(cat VERSION)' -linkmode external -extldflags '-static'" -o one-api-pro

FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder2 /build/one-api-pro /

EXPOSE 3000
WORKDIR /data
ENTRYPOINT ["/one-api-pro"]