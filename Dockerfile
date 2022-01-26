FROM golang:1.11 AS builder

ARG ARCH=arm64

# Magic line, notice in use that the lib name is different!
RUN apt-get update && apt-get install -y gcc-aarch64-linux-gnu

# Add your app and do what you need to for dependencies
ADD . /go/src/speedtest
WORKDIR /go/bin/speedtest
RUN CGO_ENABLED=1 CC=aarch64-linux-gnu-gcc GOOS=linux GOARCH=${ARCH} go build -o speedtest /go/src/speedtest

FROM ubuntu:21.10

RUN apt-get update && \
    apt-get install golang-go && \
    apt-get install python && \
    apt-get install wget && \
    wget -O speedtest-cli https://raw.githubusercontent.com/sivel/speedtest-cli/master/speedtest.py


LABEL org.opencontainers.image.authors="richard@holzeis.me"
LABEL org.opencontainers.image.source="https://github.com/holzeis/speedtest"

WORKDIR /root/

COPY --from=builder /go/bin/speedtest .

ENV PORT 8080
CMD ["./speedtest"]
