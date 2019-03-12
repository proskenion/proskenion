#!/bin/bash -e
PROJECT_ROOT=$(git rev-parse --show-toplevel)
SRC=$PROJECT_ROOT/proto
protoc --doc_out=markdown,index.md:docs --proto_path=${SRC} --go_out=plugins=grpc:${SRC} ${SRC}/*.proto
sed -i -e "s/github.com\/golang/github.com\/satellitex/g" ${SRC}/*.pb.go
rm -rf ${SRC}/*-e
sed -i -e "s/\<p\>/<pre>/g" ${PROJECT_ROOT}/docs/*
sed -i -e "s/\<\/p\>/<\/pre>/g" ${PROJECT_ROOT}/docs/*
rm -rf ${PROJECT_ROOT}/docs/*-e
