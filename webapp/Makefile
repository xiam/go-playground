CONTAINER_IMAGE   ?= xiam/go-playground
CONTAINER_NAME    ?= go-playground

TAG               ?= $(shell date +'%Y.%-m.%-d')

docker-build:
	docker build -t $(CONTAINER_IMAGE) .

docker-run:
	(docker rm -f $(CONTAINER_NAME) || exit 0) && \
	docker run \
		-p 0.0.0.0:3000:3000 \
		--name $(CONTAINER_NAME) \
		-t $(CONTAINER_IMAGE)

docker-push: docker-build
	docker tag $(CONTAINER_IMAGE) $(CONTAINER_IMAGE):$(TAG) && \
	docker tag $(CONTAINER_IMAGE):$(TAG) $(CONTAINER_IMAGE):latest && \
	docker push $(CONTAINER_IMAGE):$(TAG) && \
	docker push $(CONTAINER_IMAGE):latest

clean:
	rm -rf bin/ && \
	mkdir -p bin

build: clean
	go build -o bin/webapp *.go

run:
	go run *.go
