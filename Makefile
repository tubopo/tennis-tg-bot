OS=linux
ARCH=amd64
GIT_SHA=$(shell git rev-parse --short HEAD)

dockerx:
	docker buildx build \
	--progress=plain \
	--build-arg GIT_SHA=$(GIT_SHA) \
	--tag ghcr.io/tubopo/tennis-tg-bot:main --tag tubopo/tennis-tg-bot:main .

.PHONY: dockerx
