CONTAINER_IMAGE ?= xiam/go-playground
CONTAINER_NAME  ?= go-playground
TAG ?= latest

require-glide:
	@if [ -z "$$(which glide)" ]; then \
		echo 'Missing "glide" command. See https://github.com/Masterminds/glide' && \
		exit 1; \
	fi

require-uglifyjs:
	@if [ -z "$$(which uglifyjs)" ]; then \
		echo 'Missing "uglifyjs" command. See https://github.com/mishoo/UglifyJS' && \
		exit 1; \
	fi

docker-build: require-glide require-uglifyjs
	glide install && \
	make -C static && \
	GOOS=linux GOARCH=amd64 go build -o app_linux_amd64 && \
	docker build -t $(CONTAINER_IMAGE) .

docker-run:
	(docker stop $(CONTAINER_NAME) || exit 0) && \
	(docker rm $(CONTAINER_NAME) || exit 0) && \
	docker run -d -p 127.0.0.1:3000:3000 --name $(CONTAINER_NAME) -t $(CONTAINER_IMAGE)

docker-push: docker-build
	docker tag $(CONTAINER_IMAGE) $(CONTAINER_IMAGE):$(TAG) && \
	docker push $(CONTAINER_IMAGE):$(TAG)
