exe = ./cmd/mpdbot

VERSION = `git describe --tags`
BUILD_DATE=`date +%FT%T%z`
LDFLAGS=-ldflags "-w -s -X main.Version=${VERSION} -X main.BuildDate=${BUILD_DATE}"

.PHONY: all build build-release install test deps 

all: test clean build

deps:
	dep ensure	

clean:
	rm -rf dist/*

build:
	go build -v -o dist/mpdbot $(exe)

build-release:
	GOOS=linux CC=arm-linux-gnueabihf-gcc GOOS=linux GOARCH=arm CGO_ENABLED=1 go build -v ${LDFLAGS} -o dist/linux_arm_mpdbot/mpdbot ./cmd/mpdbot
	CGO_ENABLED=1 go build -v ${LDFLAGS} -o dist/linux_amd64_mpdbot/mpdbot ./cmd/mpdbot

install:
	go install $(exe)

test:
	go test ./...
