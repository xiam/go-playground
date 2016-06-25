# Copyright 2014 The Go Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

FROM golang:1.6.2

# add and compile tour packages
RUN go get \
	golang.org/x/tour/pic \
	golang.org/x/tour/reader \
	golang.org/x/tour/tree \
	golang.org/x/tour/wc \
	golang.org/x/talks/2016/applicative/google

# add tour packages under their old import paths (so old snippets still work)
RUN mkdir -p $GOPATH/src/code.google.com/p/go-tour && \
	cp -R $GOPATH/src/golang.org/x/tour/* $GOPATH/src/code.google.com/p/go-tour/ && \
	sed -i 's_// import_// public import_' $(find $GOPATH/src/code.google.com/p/go-tour/ -name *.go) && \
	go install \
		code.google.com/p/go-tour/pic \
		code.google.com/p/go-tour/reader \
		code.google.com/p/go-tour/tree \
		code.google.com/p/go-tour/wc

# add and compile sandbox daemon
ADD . $GOPATH/src/sandbox/
RUN go install sandbox

# make sure it works
RUN $GOPATH/bin/sandbox test

EXPOSE 8080
ENTRYPOINT ["/go/bin/sandbox"]
