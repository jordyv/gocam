BINARY_NAME=gocam

all: build build-arm

build: main.go Makefile
	go build -o dist/${BINARY_NAME} main.go

build-arm: main.go Makefile
	GOOS=linux GOARCH=arm go build -o dist/${BINARY_NAME}_arm main.go
