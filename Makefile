# Go parameters
BINARY_NAME=reader
MAIN_PATH=cmd/reader/main.go
BUILD_DIR=build

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOLINT=golangci-lint

# Docker parameters
DOCKER_IMAGE=reader
DOCKER_TAG=latest
DOCKER_FILE=Dockerfile

# Git parameters
GIT_HASH=$(shell git rev-parse --short HEAD)
BUILD_TIME=$(shell date +%FT%T%z)

# Build flags for optimization
LDFLAGS=-ldflags "-w -s -X main.version=$(GIT_HASH) -X main.buildTime=$(BUILD_TIME)"
GCFLAGS=-gcflags=all="-N -l"
BUILDTAGS=-tags 'netgo osusergo static_build'

# Performance test parameters
BENCH_TIME=5s
BENCH_COUNT=5
PROF_DIR=profiles

.PHONY: all build clean test coverage lint fmt mod-tidy docker-build docker-run help bench profile-cpu profile-mem profile-heap optimize

all: lint test build

build: ## Build the binary with optimizations
	mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) $(BUILDTAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

build-debug: ## Build with debug information
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(GCFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

clean: ## Clean build files
	$(GOCLEAN)
	rm -rf $(BUILD_DIR) $(PROF_DIR)

test: ## Run tests
	$(GOTEST) -v ./...

coverage: ## Run tests with coverage
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

lint: ## Run linter
	$(GOLINT) run

fmt: ## Format code
	$(GOFMT) -s -w .

mod-tidy: ## Tidy and verify Go modules
	$(GOMOD) tidy
	$(GOMOD) verify

docker-build: ## Build Docker image
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) -f $(DOCKER_FILE) .

docker-run: ## Run Docker container
	docker run -p 4444:4444 $(DOCKER_IMAGE):$(DOCKER_TAG)

bench: ## Run benchmarks
	mkdir -p $(PROF_DIR)
	$(GOTEST) -bench=. -benchtime=$(BENCH_TIME) -count=$(BENCH_COUNT) -benchmem ./... | tee $(PROF_DIR)/benchmark.txt

profile-cpu: ## Run CPU profiling
	mkdir -p $(PROF_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	$(BUILD_DIR)/$(BINARY_NAME) -cpuprofile=$(PROF_DIR)/cpu.prof

profile-mem: ## Run memory profiling
	mkdir -p $(PROF_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	$(BUILD_DIR)/$(BINARY_NAME) -memprofile=$(PROF_DIR)/mem.prof

profile-heap: ## Run heap profiling
	mkdir -p $(PROF_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	$(BUILD_DIR)/$(BINARY_NAME) -memprofile=$(PROF_DIR)/heap.prof

analyze-profile: ## Analyze profile data
	$(GOCMD) tool pprof -http=:8080 $(PROF_DIR)/cpu.prof

optimize: lint test bench profile-cpu profile-mem ## Run full optimization suite
	@echo "Optimization complete. Check $(PROF_DIR) for results"

dev: ## Run development server with hot reload
	air

monitor: ## Start monitoring stack (Prometheus + Grafana)
	docker-compose up -d prometheus grafana

help: ## Display this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# Default target
.DEFAULT_GOAL := help
