.PHONY: all test install-deps

BINARY_NAME=gocam

all: build build-arm

build: main.go Makefile
	packr build -o dist/${BINARY_NAME} main.go

build-arm: main.go Makefile
	GOOS=linux GOARCH=arm packr build -o dist/${BINARY_NAME}_arm main.go

test:
	go test -cover -timeout 30s ./...

install-deps:
	go get
