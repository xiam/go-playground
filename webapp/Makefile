IMAGE_NAME        ?= xiam/go-playground
IMAGE_TAG         ?= latest
IMAGE_VERSION     ?=

GIT_SHORTHASH     ?= $(shell git rev-parse --short HEAD)

CONTAINER_NAME    ?= go-playground

.PHONY: vendor

build: clean vendor fmt
	go build -o bin/webapp *.go

run:
	go run main.go

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

docker-build: fmt vendor
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

require-version:
	@test $(IMAGE_VERSION) || (echo "missing IMAGE_VERSION" && exit 1)

release: require-version docker-build
	git fetch --tags && \
	git tag -a v$(IMAGE_VERSION) -m 'release v$(IMAGE_VERSION)' && \
	git push origin v$(IMAGE_VERSION) && \
	IMAGE_TAG=$(IMAGE_VERSION) $(MAKE) docker-push && \
	IMAGE_TAG=latest $(MAKE) docker-push
