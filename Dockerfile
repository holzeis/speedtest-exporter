FROM golang:1.17 AS build

ARG ARCH=arm64

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=${ARCH}
ENV GO111MODULE=on

# Magic line, notice in use that the lib name is different!
# RUN apt-get update && apt-get install -y gcc-aarch64-linux-gnu

WORKDIR /go/src/github.com/holzeis/speedtest

COPY go.mod go.sum ./

RUN go mod download 

COPY . .

# Build application
RUN go build -o speedtest .

FROM ubuntu:21.10

WORKDIR /app

RUN apt-get update && \
    apt-get install golang-go -y && \
    apt-get install python -y && \
    apt-get install wget -y && \
    wget -O speedtest-cli https://raw.githubusercontent.com/sivel/speedtest-cli/master/speedtest.py &&  \
    chmod +x speedtest-cli

LABEL org.opencontainers.image.authors="richard@holzeis.me"
LABEL org.opencontainers.image.source="https://github.com/holzeis/speedtest-exporter"

COPY --from=build /go/src/github.com/holzeis/speedtest/speedtest /app/speedtest

ENV PORT 9112
CMD ["./speedtest"]

EXPOSE 9112
