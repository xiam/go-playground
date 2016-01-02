docker:
	GOOS=linux GOARCH=amd64 go build -o app_linux_amd64 && \
	docker build -t xiam/go-playground .

run:
	docker run -d -p 3000:3000 --name go-playground -t xiam/go-playground
