docker: Dockerfile
	docker build -t xiam/go-playground/unsafebox .

test: docker
	go test
	docker run --rm xiam/go-playground/unsafebox test

run:
	docker run -d -p 8080:8080 -t xiam/go-playground/unsafebox
