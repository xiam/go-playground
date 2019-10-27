# Copyright 2017 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
FROM debian:buster as builder
LABEL maintainer "golang-dev@googlegroups.com"

ENV GOPATH /go
ENV GOCACHE /go/.cache
ENV PATH /usr/local/go/bin:$GOPATH/bin:$PATH
ENV GOROOT_BOOTSTRAP /usr/local/gobootstrap
ENV GO_VERSION 1.13.3

RUN apt-get update && apt-get install -y ca-certificates curl

# Get the Go binary.
RUN curl -sSL https://dl.google.com/go/go$GO_VERSION.linux-amd64.tar.gz -o /tmp/go.tar.gz
RUN curl -sSL https://dl.google.com/go/go$GO_VERSION.linux-amd64.tar.gz.sha256 -o /tmp/go.tar.gz.sha256
RUN echo "$(cat /tmp/go.tar.gz.sha256) /tmp/go.tar.gz" | sha256sum -c -
RUN tar -C /usr/local/ -vxzf /tmp/go.tar.gz

# Add and compile playground daemon
COPY . /go/src/playground/
RUN go install playground

FROM debian:buster

RUN apt-get update && apt-get install -y ca-certificates git --no-install-recommends

COPY --from=builder /usr/local/go /usr/local/go

ENV GOPATH /go
ENV GOCACHE /tmp/.cache

ENV PATH /usr/local/go/bin:$GOPATH/bin:$PATH

# Add and compile tour packages
RUN go get \
    golang.org/x/tour/pic \
    golang.org/x/tour/reader \
    golang.org/x/tour/tree \
    golang.org/x/tour/wc \
    golang.org/x/talks/content/2016/applicative/google && \
    rm -rf $GOPATH/src/golang.org/x/tour/.git && \
    rm -rf $GOPATH/src/golang.org/x/talks/.git

# Add tour packages under their old import paths (so old snippets still work)
RUN mkdir -p $GOPATH/src/code.google.com/p/go-tour && \
    cp -R $GOPATH/src/golang.org/x/tour/* $GOPATH/src/code.google.com/p/go-tour/ && \
    sed -i 's_// import_// public import_' $(find $GOPATH/src/code.google.com/p/go-tour/ -name *.go) && \
    go install \
    code.google.com/p/go-tour/pic \
    code.google.com/p/go-tour/reader \
    code.google.com/p/go-tour/tree \
    code.google.com/p/go-tour/wc

RUN mkdir /app

COPY --from=builder /go/bin/playground /app
WORKDIR /app

# Run tests
RUN /app/playground test

EXPOSE 8080
ENTRYPOINT ["/app/playground"]
