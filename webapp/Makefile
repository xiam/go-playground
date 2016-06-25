IMAGE_NAME ?= xiam/playgo-web

docker:
	go get -d ./... && \
	GOOS=linux GOARCH=amd64 go build -o app_linux_amd64 && \
	docker build -t $(IMAGE_NAME) .

run:
	docker run -d -p 3000:3000 --name go-playground -t $(IMAGE_NAME)

push:
	docker push $(IMAGE_NAME)
