exe = ./cmd/mpdbot

.PHONY: all build build-arm install test deps 

all: test build

deps:
	dep ensure	

build:
	go build -v -o dist/mpdbot $(exe)

build-arm:
	env GOOS=linux GOARCH=arm64 go build -v -o dist/arm64/mpdbot $(exe)

install:
	go install $(exe)

test:
	go test ./...
