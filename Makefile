NAME   := jonathanwthom/hosts
TAG    := $$(git rev-parse --short HEAD)
IMG    := ${NAME}:${TAG}
LATEST := ${NAME}:latest

build:
	@docker build -t ${IMG} .
	@docker tag ${IMG} ${LATEST}

push:
	@docker push ${NAME}

