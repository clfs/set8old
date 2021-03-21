.PHONY: help
.DEFAULT_GOAL := help

test: ## Run short tests (no cache)
	go test -short -count=1 ./...

test-all: ## Run all tests (no cache)
	go test -count=1 ./...

bench: ## Run short benchmarks
	go test -short -bench .

bench-all: ## Run all benchmarks
	go test -bench .

pprof: ## Profile cpu/mem for short tests/benches
	go test -short -cpuprofile cpu.prof -memprofile mem.prof -bench .

pprof-all: ## Profile cpu/mem for all tests/benches
	go test -cpuprofile cpu.prof -memprofile mem.prof -bench .

lint: ## Lint all files
	golangci-lint run

fmt: ## Show autoformat diff
	gofumpt -d .

fmt-apply: ## Apply autoformat
	gofumpt -l -w .

perf: ## Create a new 

clean: ## Delete generated files
	rm mem.prof cpu.prof set8.test

# Adapted from https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'