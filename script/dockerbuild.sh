#!/bin/sh

build() {
    v=$1
    if [ "$1" = '' ]; then
      echo "please input version, example: v1.0.0"
      exit 1
    fi

    make clean
    go mod tidy && go mod vendor

    platform=linux/amd64,linux/arm/v7
    docker buildx build --platform "$platform"\
      --cache-from=registry-shdocker-registry.cn-shanghai.cr.aliyuncs.com/tedgedriver/alpine:3.12 \
      --cache-from=registry-shdocker-registry.cn-shanghai.cr.aliyuncs.com/tedgedriver/golang:1.18-alpine \
      -t registry-shdocker-registry.cn-shanghai.cr.aliyuncs.com/tedgedriver/driver-example:"$v"\
      -f Dockerfile --push .

    if [ $? != 0 ];then
       exit 1
    fi
}

build "$1";
