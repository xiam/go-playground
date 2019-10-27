FROM node:8.15 AS node-builder

RUN apt-get update && \
  apt-get install -y make

RUN npm install uglify-js -g

COPY static /app/static

RUN cd /app/static && make

FROM golang:1.11 AS go-builder

WORKDIR /go/src/github.com/xiam/go-playground/webapp/

COPY . .

RUN go install github.com/xiam/go-playground/webapp

FROM debian:jessie

RUN apt-get update && \
  apt-get install -y ca-certificates

COPY --from=go-builder /go/bin/webapp /usr/local/bin/webapp
COPY --from=node-builder /app/static /var/app/static

WORKDIR /var/app

CMD [ "webapp" ]
