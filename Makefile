docker-push:
	parallel $(MAKE) -C {} docker-push ::: webapp unsafebox sandbox
