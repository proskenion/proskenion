SHELL := /bin/bash

.PHONY: proto
proto:
	script/proto.sh

.PHONY: build
build:
	go build -o ./bin/proskenion main.go
	go build -o ./bin/keygen ./script/keygen.go
	go build -o ./bin/example ./example/sender.go

.PHONY: build-osx
build-osx:
	GOOS=darwin GOARCH=amd64 go build -o ./bin/darwin64/proskenion main.go

# https://github.com/mattn/go-sqlite3/issues/372#issuecomment-396863368
.PHONY: build-linux
build-linux:
	CC=x86_64-linux-musl-gcc CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -o ./bin/linux/proskenion main.go

.PHONY: test
test:
	go test -cover -v $(shell glide novendor)

.PHONY: test-ci
test-ci:
	go test -v $(shell glide novendor)

.PHONY: dockerup
dockerup:
	docker build . -t proskenion

.PHONY: build-docker
build-docker: build-linux dockerup

.PHONY: example
example: build build-docker