#!/bin/bash -e
PROJECT_ROOT=$(git rev-parse --show-toplevel)
SRC=$PROJECT_ROOT/proto
protoc --proto_path=${SRC} --go_out=plugins=grpc:${SRC} ${SRC}/*.proto
