docker-push:
	for PROJECT in webapp unsafebox sandbox; do \
		$(MAKE) -C $$PROJECT docker-push; \
	done
