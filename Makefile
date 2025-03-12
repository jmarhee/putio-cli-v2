# Makefile for building the Go program for different platforms and architectures
SHELL=/bin/bash
BINARY_NAME=putio-cli

build:
	@echo "Building for all platforms..."
	$(MAKE) build-linux
	$(MAKE) build-macos

build-linux:
	@echo "Building for Linux..."
	GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY_NAME)-linux-amd64 ./...
	GOOS=linux GOARCH=arm64 go build -o bin/$(BINARY_NAME)-linux-arm64 ./...

build-macos:
	@echo "Building for macOS..."
	GOOS=darwin GOARCH=amd64 go build -o bin/$(BINARY_NAME)-darwin-amd64 ./...
	GOOS=darwin GOARCH=arm64 go build -o bin/$(BINARY_NAME)-darwin-arm64 ./...

clean:
	@echo "Cleaning up..."
	rm -rf bin/*

.PHONY: build build-linux build-macos clean
