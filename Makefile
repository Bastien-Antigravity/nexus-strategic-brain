.PHONY: version build test run clean help

VERSION := $(shell cat VERSION.txt 2>/dev/null || echo "0.0.1")

version:
	@echo $(VERSION)

build:
	@echo "Building strategic-nexus binary (v$(VERSION))..."
	@go build -o bin/strategic-nexus ./cmd/strategic-nexus

test:
	@echo "Running tests in 01-Strategic-Nexus..."
	@go test -v ./...

run: build
	@echo "Running strategic-nexus daemon..."
	@./bin/strategic-nexus

clean:
	@echo "Cleaning binaries and build artifacts..."
	@rm -rf bin/

help:
	@echo "01-Strategic-Nexus Makefile (v$(VERSION))"
	@echo "  make version - Print current version"
	@echo "  make build   - Build binary"
	@echo "  make test    - Run unit tests"
	@echo "  make run     - Build and run service"
	@echo "  make clean   - Clean build artifacts"
