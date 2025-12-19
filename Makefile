# goscaffold Makefile

BINARY_NAME=goscaffold
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

.PHONY: all build clean test lint run install help

all: lint test build

## build: Build the binary
build:
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/$(BINARY_NAME)

## install: Install to GOPATH/bin
install:
	go install $(LDFLAGS) ./cmd/$(BINARY_NAME)

## clean: Clean build artifacts
clean:
	rm -rf bin/
	rm -rf dist/
	rm -f coverage.out

## test: Run tests
test:
	go test -v -race -coverprofile=coverage.out ./...

## lint: Run linter
lint:
	golangci-lint run ./...

## run: Run the application
run:
	go run ./cmd/$(BINARY_NAME)

## demo: Create a demo project
demo:
	rm -rf demo-output
	go run ./cmd/$(BINARY_NAME) new demo-output -t api -g demo -D -Q --no-interactive
	@echo "\nDemo project created in demo-output/"

## release-dry: Test goreleaser locally
release-dry:
	goreleaser release --snapshot --clean

## help: Show this help
help:
	@echo "goscaffold - Go Project Scaffolding Tool\n"
	@echo "Available targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'
