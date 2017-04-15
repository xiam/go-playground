CONTAINER_IMAGE ?= xiam/playgo-web
CONTAINER_NAME  ?= my-go-playground

require-glide:
	@if [ -z "$$(which glide)" ]; then \
		echo 'Missing "glide" command. See https://github.com/Masterminds/glide' && \
		exit 1; \
	fi

docker-build: require-glide
	glide install && \
	GOOS=linux GOARCH=amd64 go build -o app_linux_amd64 && \
	docker build -t $(CONTAINER_IMAGE) .

docker-run:
	(docker stop $(CONTAINER_NAME) || exit 0) && \
	(docker stop $(CONTAINER_NAME) || exit 0) && \
	docker run -d -p 127.0.0.1:3000:3000 --name $(CONTAINER_NAME) -t $(CONTAINER_IMAGE)

docker-push:
	docker push $(CONTAINER_IMAGE)
