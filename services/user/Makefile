
GOPATH:=$(shell go env GOPATH)
MODIFY=Mproto/imports/api.proto=github.com/micro/go-micro/v2/api/proto

.PHONY: proto
proto:
	protoc -I=${GOPATH}/src  --proto_path=../../ --micro_out=${MODIFY}:. --gofast_out=${MODIFY}:.  proto/user/user.proto

    

.PHONY: build
build: proto

	CGO_ENABLED=0 GOOS=linux go build  -o  user-service *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker: build
	docker build . -t scg130/user:latest

.PHONY: push
push: docker
	docker push scg130/user