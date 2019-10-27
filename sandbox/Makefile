IMAGE_NAME        ?= xiam/go-playground-sandbox
IMAGE_TAG         ?= latest
IMAGE_VERSION     ?=

GIT_SHORTHASH     ?= $(shell git rev-parse --short HEAD)

CONTAINER_NAME    ?= go-playground-sandbox

.PHONY: vendor

test: docker-build
	go test && \
	docker run --rm $(IMAGE_NAME):$(GIT_SHORTHASH) test

fmt:
	for i in $$(find -name \*.go | grep -v vendor); do \
		gofmt -w $$i && \
		goimports -w $$i; \
	done

docker-build: fmt
	docker build -t $(IMAGE_NAME):$(GIT_SHORTHASH) .

docker-run:
	(docker rm -f $(CONTAINER_NAME) || exit 0) && \
	docker run \
		-p 0.0.0.0:8080:8080 \
		--name $(CONTAINER_NAME) \

docker-push: docker-build
	docker tag $(IMAGE_NAME):$(GIT_SHORTHASH) $(IMAGE_NAME):$(IMAGE_TAG) && \
	docker push $(IMAGE_NAME):$(GIT_SHORTHASH) && \
	docker push $(IMAGE_NAME):$(IMAGE_TAG)
