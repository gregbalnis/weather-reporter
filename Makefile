# Makefile for weather-reporter

BINARY_NAME=weather-reporter
BUILD_DIR=bin
SRC_DIR=./src/cmd/weather-reporter
VERSION ?= dev
COMMIT ?= $(shell git rev-parse --short HEAD)
DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.Date=$(DATE)"

.PHONY: all build clean test lint snapshot

all: clean lint test build

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(SRC_DIR)

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -rf dist/

test:
	@echo "Running tests..."
	go test -race -coverprofile=coverage.out ./src/...
	go tool cover -func=coverage.out

lint:
	@echo "Running linter..."
	golangci-lint run ./src/...

# snapshot:
# 	@echo "Creating snapshot release..."
# 	goreleaser release --snapshot --clean
