ARG ALPINE_VERSION=3.20
ARG GO_VERSION=1.23.3
ARG NODE_VERSION=23

ARG CGO_ENABLED=0

# Compile playground and webapp
FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder

RUN apk add --no-cache \
  ca-certificates

WORKDIR /go/src/github.com/xiam/go-playground

COPY ./ ./

RUN go build -o /go/bin/go-playground-executor github.com/xiam/go-playground/executor

RUN go build -o /go/bin/go-playground-webapp github.com/xiam/go-playground/webapp

# Build web assets
FROM node:${NODE_VERSION}-alpine AS node-builder

RUN apk update && \
  apk add make

RUN npm install uglify-js -g

WORKDIR /app
COPY ./webapp/static ./static

RUN cd ./static && \
  make

# Compose final image
FROM alpine:${ALPINE_VERSION}

ARG GO_VERSION
ARG CGO_ENABLED

ENV GOCACHE=/tmp/.gocache

RUN apk add --no-cache \
  ca-certificates \
  curl

ENV GOLANG_URL=https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz

RUN curl -sSL ${GOLANG_URL} | tar -C /usr/local -xz

ENV GOARCH=amd64
ENV GOOS=linux
ENV GOPATH=/go
ENV CGO_ENABLED=${CGO_ENABLED:-0}

ENV PATH=/usr/local/go/bin:$PATH

WORKDIR /app/

RUN mkdir -p ./bin

COPY --from=builder /go/bin/go-playground-executor ./bin/
COPY --from=builder /go/bin/go-playground-webapp ./bin/

COPY --from=node-builder /app/static ./static

ENV PATH=/app/bin:$PATH

EXPOSE 3000
EXPOSE 3003
