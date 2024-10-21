GOPATH:=$(shell go env GOPATH)
MODIFY=Mproto/imports/api.proto=github.com/micro/go-micro/v2/api/proto


.PHONY: proto
proto:
	protoc -I=${GOPATH}/src  --proto_path=. --micro_out=${MODIFY}:. --gofast_out=${MODIFY}:.  proto/user/*.proto
	cp -rf  proto/user/*.go services/user/proto/user
	protoc -I=${GOPATH}/src  --proto_path=. --micro_out=${MODIFY}:. --gofast_out=${MODIFY}:.  proto/wallet/*.proto
	cp -rf  proto/wallet/*.go services/wallet/proto/wallet
	protoc -I=${GOPATH}/src  --proto_path=. --micro_out=${MODIFY}:. --gofast_out=${MODIFY}:.  proto/novel/*.proto
	cp -rf  proto/novel/*.go services/novel/proto/novel
	protoc -I=${GOPATH}/src  --proto_path=. --micro_out=${MODIFY}:. --gofast_out=${MODIFY}:.  proto/charge/*.proto
	cp -rf  proto/charge/*.go services/charge/proto/charge
	protoc -I=${GOPATH}/src  --proto_path=. --micro_out=${MODIFY}:. --gofast_out=${MODIFY}:.  proto/admin/*.proto
	cp -rf  proto/admin/*.go services/admin/proto/admin
	protoc -I=${GOPATH}/src  --proto_path=. --micro_out=${MODIFY}:. --gofast_out=${MODIFY}:.  proto/comment/*.proto
	cp -rf  proto/comment/*.go services/comment/proto/comment

#.PHONY: api
#api:
#	micro --registry=etcd --registry_address=192.168.18.14:2379  api >> /dev/null 2>&1 &
#
#.PHONY: web
#web:
#	micro --registry=etcd --registry_address=192.168.18.14:2379  --enable_stats  web >> /dev/null 2>&1 &

.PHONY: run
run:
	swag init --parseDependency
	go run -tags "doc" .

.PHONY: build
build: 
	#swag init --parseDependency
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o  runapp .

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t scg130/kanshu-runapp:latest

.PHONY: push
push:
	docker push scg130/kanshu-runapp
