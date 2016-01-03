docker: Dockerfile
	docker build -t playground/unsafebox .

test: docker
	go test
	docker run --rm playground/unsafebox test

run:
	docker run -d -p 8080:8080 --name playground-unsafebox -t playground/unsafebox
