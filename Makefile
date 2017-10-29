exe = ./cmd/mpdbot

.PHONY: all build build-arm install test deps 

all: install

deps:
	dep ensure	

build:
	go build -v -o build/mpdbot $(exe)

build-arm:
	env GOOS=android GOARCH=arm64 go build -v -o build/arm64/mpdbot $(exe)

install:
	go install $(exe)

test:
	go test
