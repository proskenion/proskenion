#!/bin/bash -e
PROJECT_ROOT=$(git rev-parse --show-toplevel)
LOCAL_HOST_IP=`ifconfig en0 | grep inet | grep -v inet6 | sed -E "s/inet ([0-9]{1,3}.[0-9]{1,3}.[0-9].{1,3}.[0-9]{1,3}) .*$/\1/" | tr -d "\t"`
echo $LOCAL_HOST_IP
SRC=$PROJECT_ROOT/example
sed -i -e "s/[0-9]*\.[0-9]*\.[0-9]*\.[0-9]*/$LOCAL_HOST_IP/g" ${SRC}/*.yaml
rm -rf ${SRC}/*-e
