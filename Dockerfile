# syntax = docker/dockerfile:experimental

FROM swr.cn-east-3.myhuaweicloud.com/wsxcc/golang:latest as builder
WORKDIR /build
copy . .
RUN --mount=type=cache,target=/go,id=gomod,sharing=locked \
    GOPROXY=https://goproxy.cn CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -v -tags timetzdata -o app

FROM swr.cn-east-3.myhuaweicloud.com/wsxcc/alpine:latest
COPY --from=builder /build/app /
ENTRYPOINT [ "/app" ]