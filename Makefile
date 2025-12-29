# Go Requests Library Makefile
# ============================

.PHONY: help build test test-coverage test-race bench bench-save bench-compare bench-ci lint fmt clean

# Default target
.DEFAULT_GOAL := help

# Variables
BENCH_DIR := .benchmarks
BENCH_FILE := $(BENCH_DIR)/current.txt
BASELINE_FILE := $(BENCH_DIR)/baseline.txt
COVERAGE_FILE := coverage.out
COVERAGE_HTML := coverage.html

# ============================================================================
# Help
# ============================================================================

help: ## Show this help message
	@echo "Go Requests Library - Available Commands"
	@echo "========================================="
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "Examples:"
	@echo "  make build          # Build the project"
	@echo "  make test           # Run all tests"
	@echo "  make bench          # Run benchmarks"
	@echo "  make bench-compare  # Compare with baseline"

# ============================================================================
# Build
# ============================================================================

build: ## Build the project
	@echo "Building..."
	@go build ./...
	@echo "Build complete."

# ============================================================================
# Testing
# ============================================================================

test: ## Run all tests
	@echo "Running tests..."
	@go test ./...

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	@go test -coverprofile=$(COVERAGE_FILE) ./...
	@go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "Coverage report generated: $(COVERAGE_HTML)"
	@go tool cover -func=$(COVERAGE_FILE) | tail -1

test-race: ## Run tests with race detector
	@echo "Running tests with race detector..."
	@go test -race ./...

# ============================================================================
# Benchmarks
# ============================================================================

bench: ## Run benchmarks and save results
	@echo "Running benchmarks..."
	@mkdir -p $(BENCH_DIR)
	@go test -bench=. -benchmem -count=5 ./test/... 2>&1 | tee $(BENCH_FILE)
	@echo ""
	@echo "Results saved to $(BENCH_FILE)"

bench-save: ## Save current benchmark results as baseline
	@if [ ! -f $(BENCH_FILE) ]; then \
	  echo "No current benchmark results found. Run 'make bench' first."; \
	  exit 1; \
	fi
	@mkdir -p $(BENCH_DIR)
	@cp $(BENCH_FILE) $(BASELINE_FILE)
	@echo "Baseline saved to $(BASELINE_FILE)"

bench-compare: ## Compare current benchmarks with baseline
	@if [ ! -f $(BASELINE_FILE) ]; then \
	  echo "No baseline found. Run 'make bench-save' first."; \
	  exit 1; \
	fi
	@if [ ! -f $(BENCH_FILE) ]; then \
	  echo "No current benchmark results found. Run 'make bench' first."; \
	  exit 1; \
	fi
	@if ! command -v benchstat > /dev/null 2>&1; then \
		echo "benchstat not found. Install with:"; \
		echo "  go install golang.org/x/perf/cmd/benchstat@latest"; \
		exit 1; \
	fi
	@echo "Comparing benchmarks..."
	@echo ""
	@benchstat $(BASELINE_FILE) $(BENCH_FILE)

bench-version: ## Compare benchmarks with a specific git version (usage: make bench-version VERSION=v1.0.0)
	@if [ -z "$(VERSION)" ]; then \
	  echo "Usage: make bench-version VERSION=<tag|branch|commit>"; \
	  echo ""; \
	  echo "Example: make bench-version VERSION=v1.0.0"; \
	  echo ""; \
	  echo "Available tags:"; \
	  git tag -l --sort=-v:refname | head -10; \
	  exit 1; \
	fi
	@./scripts/bench-compare.sh $(VERSION)

bench-latest: ## Compare benchmarks with the latest git tag
	@./scripts/bench-compare.sh

bench-ci: ## Run benchmarks for CI (JSON output)
	@echo "Running benchmarks for CI..."
	@mkdir -p $(BENCH_DIR)
	@go test -bench=. -benchmem -count=10 ./test/... -json > $(BENCH_DIR)/bench-ci.json
	@echo "Results saved to $(BENCH_DIR)/bench-ci.json"

bench-profile: ## Run benchmarks with CPU and memory profiling
	@echo "Running benchmarks with profiling..."
	@mkdir -p $(BENCH_DIR)
	@go test -bench=. -benchmem -cpuprofile=$(BENCH_DIR)/cpu.prof -memprofile=$(BENCH_DIR)/mem.prof ./test/...
	@echo ""
	@echo "Profiles saved to:"
	@echo "  - $(BENCH_DIR)/cpu.prof (CPU profile)"
	@echo "  - $(BENCH_DIR)/mem.prof (Memory profile)"
	@echo ""
	@echo "View with: go tool pprof $(BENCH_DIR)/cpu.prof"

# ============================================================================
# Code Quality
# ============================================================================

lint: ## Run linters
	@echo "Running linters..."
	@if command -v golangci-lint > /dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not found. Install with:"; \
		echo "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		echo ""; \
		echo "Running go vet instead..."; \
		go vet ./...; \
	fi

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@if command -v gofumpt > /dev/null 2>&1; then \
		gofumpt -extra -w .; \
	else \
		echo "gofumpt not found. Install with:"; \
		echo "  go install mvdan.cc/gofumpt@latest"; \
	fi
	@echo "Formatting code complete."

# ============================================================================
# Cleanup
# ============================================================================

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf $(COVERAGE_FILE) $(COVERAGE_HTML)
	@rm -rf $(BENCH_DIR)
	@go clean -testcache
	@echo "Clean complete."

# ============================================================================
# Development
# ============================================================================

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies updated."

verify: ## Verify dependencies
	@echo "Verifying dependencies..."
	@go mod verify
	@echo "Dependencies verified."

all: fmt lint test build ## Run all checks (format, lint, test, build)
	@echo "All checks passed!"
