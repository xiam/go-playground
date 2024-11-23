ARG ALPINE_VERSION=3.20
ARG GO_VERSION=1.23.3
ARG NODE_VERSION=23

# Compile playground and webapp
FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder

RUN apk add --no-cache \
  ca-certificates

WORKDIR /go/src/github.com/xiam/go-playground

COPY ./ ./

RUN go build -o /go/bin/go-playground-runner github.com/xiam/go-playground/runner

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

RUN apk add --no-cache \
  ca-certificates git

WORKDIR /app/

RUN mkdir -p ./bin

COPY --from=builder /go/bin/go-playground-runner ./bin/
COPY --from=builder /go/bin/go-playground-webapp ./bin/

COPY --from=node-builder /app/static ./static

ENV PATH=/app/bin:$PATH

EXPOSE 3000
EXPOSE 3020

CMD ["go-playground-webapp"]
