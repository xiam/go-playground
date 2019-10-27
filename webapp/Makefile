IMAGE_NAME        ?= xiam/go-playground
IMAGE_TAG         ?= latest
IMAGE_VERSION     ?=

GIT_SHORTHASH     ?= $(shell git rev-parse --short HEAD)

CONTAINER_NAME    ?= go-playground

.PHONY: vendor static

build: clean vendor fmt
	go build -o bin/webapp *.go

run: static
	go run *.go

static:
	$(MAKE) -C static

vendor:
	go mod vendor

fmt:
	for i in $$(find -name \*.go | grep -v vendor); do \
		gofmt -w $$i && \
		goimports -w $$i; \
	done

clean:
	rm -rf bin/ && \
	mkdir -p bin

docker-build: vendor
	docker build -t $(IMAGE_NAME):$(GIT_SHORTHASH) .

docker-run:
	(docker rm -f $(CONTAINER_NAME) || exit 0) && \
	docker run \
		-p 0.0.0.0:3000:3000 \
		--name $(CONTAINER_NAME) \
		-t $(IMAGE_NAME)

docker-push: docker-build
	docker tag $(IMAGE_NAME):$(GIT_SHORTHASH) $(IMAGE_NAME):$(IMAGE_TAG) && \
	docker push $(IMAGE_NAME):$(GIT_SHORTHASH) && \
	docker push $(IMAGE_NAME):$(IMAGE_TAG)
