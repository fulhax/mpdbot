exe = ./cmd/mpdbot

.PHONY: all build install test deps 

all: test clean build

deps:
	dep ensure	

clean:
	rm -rf dist/*

build:
	go build -v -o dist/mpdbot $(exe)

install:
	go install $(exe)

test:
	go test ./...
