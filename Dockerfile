# step1
#ARG BASE=golang:1.18-alpine
#ARG BASE=registry.shdocker.tuya-inc.top/tedge/golang:1.18-alpine
ARG BASE=registry-shdocker-registry.cn-shanghai.cr.aliyuncs.com/tedgedriver/golang:1.18-alpine
FROM ${BASE} AS builder

ARG MAKE='make build'

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk add --update --no-cache make git

WORKDIR /driver-example
COPY . .
RUN source script/env.sh

RUN --mount=type=cache,target=/go,id=example_cache,sharing=shared
RUN --mount=type=cache,target=/root/.cache,id=example_build_cache,sharing=shared ${MAKE}

# step2
#FROM alpine:3.12
#FROM registry.shdocker.tuya-inc.top/tedge/alpine:3.12
FROM registry-shdocker-registry.cn-shanghai.cr.aliyuncs.com/tedgedriver/alpine:3.12

ENV TZ Asia/Shanghai
RUN apk add --no-cache --upgrade tzdata \
    && ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone \
    && mkdir -p /home/driver/

COPY --from=builder /driver-example/driver-example /home/driver/driver-example

ENTRYPOINT ["/home/driver/driver-example"]

