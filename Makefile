SHELL := /bin/bash

.PHONY: proto
proto:
	script/proto.sh

.PHONY: build
build:
	go build -o ./bin/proskenion main.go

.PHONY: build-osx
build-osx:
	GOOS=darwin GOARCH=amd64 go build -o ./bin/darwin64/proskenion main.go

.PHONY: build-linux
build-linux:
	GOOS=linux GOARCH=amd64 go build -o ./bin/linux/proskenion main.go

.PHONY: test
test:
	go test -cover -v $(shell glide novendor | grep -v grpc_test)

.PHONY: test-ci
test-ci:
	go test -v $(shell glide novendor nogrpc_test | grep -v grpc_test)

.PHONY: test-grpc
test-grpc:
	./bin/proskenion ../grpc_test/config.yaml &
	sleep 1
	go test -cover -v ./grpc_test