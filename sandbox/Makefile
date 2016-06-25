IMAGE_NAME ?= xiam/playgo-sandbox

docker: Dockerfile
	docker build -t $(IMAGE_NAME) .

test: docker
	go test
	docker run --rm $(IMAGE_NAME) test

run:
	docker run -d -p 8080:8080 -t $(IMAGE_NAME)

push:
	docker push $(IMAGE_NAME)
