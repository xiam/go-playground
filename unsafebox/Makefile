CONTAINER_IMAGE   ?= xiam/go-playground-unsafebox
CONTAINER_NAME    ?= go-playground-unsafebox

TAG               ?= $(shell date +'%Y.%-m.%-d')

docker-build: Dockerfile
	docker build -t $(CONTAINER_IMAGE) .

docker-push: docker-build
	docker tag $(CONTAINER_IMAGE) $(CONTAINER_IMAGE):$(TAG) && \
	docker tag $(CONTAINER_IMAGE):$(TAG) $(CONTAINER_IMAGE):latest && \
	docker push $(CONTAINER_IMAGE):$(TAG) && \
	docker push $(CONTAINER_IMAGE):latest

test: docker-build
	go test && \
	docker run --rm $(CONTAINER_IMAGE) test

docker-run: docker-build
	(docker rm -f $(CONTAINER_NAME) || exit 0) && \
	docker run \
		-p 0.0.0.0:8080:8080 \
		--name $(CONTAINER_NAME) \
		-t $(CONTAINER_IMAGE)
