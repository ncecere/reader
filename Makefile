# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=reader
DOCKER_IMAGE=reader

# Build flags
LDFLAGS=-ldflags "-s -w"

.PHONY: all build clean test coverage deps docker-build docker-run help

all: clean build test ## Build and run tests

build: ## Build the application
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) ./cmd/reader

clean: ## Clean build files
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -rf screenshots/*

test: ## Run tests
	$(GOTEST) -v ./...

coverage: ## Run tests with coverage
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

deps: ## Download dependencies
	$(GOMOD) download
	$(GOMOD) tidy

docker-build: ## Build Docker image
	docker build -t $(DOCKER_IMAGE) .

docker-run: ## Run Docker container
	docker run -p 8080:8080 $(DOCKER_IMAGE)

lint: ## Run linters
	golangci-lint run

fmt: ## Format code
	gofmt -s -w .
	goimports -w .

run: ## Run the application
	$(GOCMD) run ./cmd/reader

dev: ## Run with hot reload
	air -c .air.toml

check: lint test ## Run linters and tests

# Help target
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

# Default target
.DEFAULT_GOAL := help
