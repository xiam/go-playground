docker: Dockerfile
	docker build -t playground/sandbox .

test: docker
	go test
	docker run --rm playground/sandbox test

run:
	docker run -d -p 8080:8080 -t playground/sandbox
