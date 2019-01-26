CONTAINER_IMAGE ?= xiam/playgo-sandbox
CONTAINER_NAME 	?= go-playground-sandbox

docker: Dockerfile
	docker build -t $(CONTAINER_IMAGE) .

test: docker
	go test
	docker run --rm $(CONTAINER_IMAGE) test

run:
	(docker stop $(CONTAINER_NAME) || exit 0) && \
    (docker rm $(CONTAINER_NAME) || exit 0) && \
	docker run -d -p 8080:8080 --name $(CONTAINER_NAME) -t $(CONTAINER_IMAGE)

push:
	docker push $(CONTAINER_IMAGE)
