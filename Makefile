OS=linux
ARCH=amd64

dockerx:
	docker buildx build \
	--progress=plain \
	--tag ghcr.io/tubopo/tennis-tg-bot:main --tag tubopo/tennis-tg-bot:main .

.PHONY: dockerx
