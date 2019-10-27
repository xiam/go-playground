IMAGE_VERSION     ?=

docker-push:
	parallel $(MAKE) -C {} docker-push ::: webapp unsafebox sandbox

require-version:
	@test $(IMAGE_VERSION) || (echo "missing IMAGE_VERSION" && exit 1)

release: require-version docker-push
	git fetch --tags && \
	git tag -a v$(IMAGE_VERSION) -m 'release v$(IMAGE_VERSION)' && \
	git push origin v$(IMAGE_VERSION) && \
	IMAGE_TAG=$(IMAGE_VERSION) $(MAKE) docker-push && \
	IMAGE_TAG=latest $(MAKE) docker-push
