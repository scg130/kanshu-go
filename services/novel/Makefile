
GOPATH:=$(shell go env GOPATH)
MODIFY=Mproto/imports/api.proto=github.com/micro/go-micro/v2/api/proto

.PHONY: proto
proto:
    
	protoc -I=${GOPATH}/src  --proto_path=../../ --micro_out=${MODIFY}:. --gofast_out=${MODIFY}:.  proto/novel/novel.proto

    

.PHONY: build
build: proto

	CGO_ENABLED=0 GOOS=linux go build -o novel-service *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t scg130/novel:latest

.PHONY: push
push:
	docker push scg130/novel