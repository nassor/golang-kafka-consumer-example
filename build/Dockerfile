# BUILD
FROM golang:1.13-alpine AS build-env
RUN apk upgrade --update-cache --available
RUN apk add -U --repository http://dl-cdn.alpinelinux.org/alpine/edge/community \
    gcc \
    librdkafka-dev \
    musl-dev
RUN mkdir -p /source
WORKDIR /source
ADD . /source
ARG CGO_ENABLED=1
ARG GOOS=linux
ARG GOARCH=amd64
RUN go build -mod=vendor -ldflags '-s -w' -o /device-manager ./cmd/device-manager

# PRODUCTION IMAGE
FROM alpine:3
RUN apk upgrade --update-cache --available
RUN apk add -U --repository http://dl-cdn.alpinelinux.org/alpine/edge/community \
    ca-certificates \
    librdkafka \
    jq \
    bash
RUN update-ca-certificates
COPY --from=build-env /device-manager /bin/
ENTRYPOINT ["/bin/device-manager"]
