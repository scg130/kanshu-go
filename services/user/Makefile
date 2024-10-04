
GOPATH:=$(shell go env GOPATH)
MODIFY=Mproto/imports/api.proto=github.com/micro/go-micro/v2/api/proto

.PHONY: build
build:

	CGO_ENABLED=0 GOOS=linux go build  -o  user-service *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker: build
	docker build . -t scg130/kanshu-user:latest

.PHONY: push
push: docker
	docker push scg130/kanshu-user