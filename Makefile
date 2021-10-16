NAME   := jonathanwthom/hosts
TAG    := $$(git rev-parse --short HEAD)
IMG    := ${NAME}:${TAG}
LATEST := ${NAME}:latest

build:
	@docker build -t ${IMG} .
	@docker tag ${IMG} ${LATEST}

push:
	@docker push ${NAME}

start:
	go run ./...

pop:
	go run ./... -p

popc:
	go run ./... -p -h $(h)

