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

.SILENT:
tag-release:
	if [[ $(TAG) == v?.?.? ]]; then echo "Tagging $(TAG)"; elif [[ $(TAG) == v?.?.?? ]]; then echo "Tagging $(TAG)"; else echo "Bad Tag Format: $(TAG)"; exit 1; fi && git tag -a $(TAG) -m "Releasing $(TAG)" ; read -p "Push tag: $(TAG)? " push_tag ; if [ "${push_tag}"="yes" ]; then git push origin $(TAG); fi

.SILENT:
create-release:
	if [[ $(TAG) == v?.?.? ]]; then echo "Cutting release from $(TAG)"; elif [[ $(TAG) == v?.?.?? ]]; then echo "Cutting release from $(TAG)"; else echo "Bad Tag Format, cannot cut release: $(TAG)"; exit 1; fi && git tag -a $(TAG) -m "Releasing $(TAG)" ; read -p "Cut release from tag: $(TAG)? " push_tag ; if [ "${push_tag}"="yes" ]; then TAG=$(TAG) ./make-release.sh; fi

.PHONY: build build-linux build-macos clean
