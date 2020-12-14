VERSION=v1.0
IMAGE_NAME=$(shell pwd | xargs basename)
GOPATH:=$(shell go env GOPATH)

.PHONY: docker
docker:
	DOCKER_BUILDKIT=1 docker build -t ${IMAGE_NAME}:${VERSION} .
	docker system prune -f

.PHONY: push
push:
	docker tag ${IMAGE_NAME}:${VERSION} swr.cn-east-3.myhuaweicloud.com/wsxcc/${IMAGE_NAME}:${VERSION}
	docker push swr.cn-east-3.myhuaweicloud.com/wsxcc/${IMAGE_NAME}:${VERSION}