.PHONY: help
.DEFAULT_GOAL := help

test: ## Quick tests, without cache
	go test -short -count=1 ./...

test-all: ## All tests, without cache
	go test -count=1 ./...

bench: ## Quick benchmarks
	go test -short -bench .

bench-all: ## All benchmarks
	go test -bench .

pprof: ## Quick cpu & mem profile
	go test -short -cpuprofile cpu.prof -memprofile mem.prof -bench .

pprof-all: ## Complete cpu & mem profile
	go test -cpuprofile cpu.prof -memprofile mem.prof -bench .

lint: ## Lint all files
	golangci-lint run

fmt: ## Show autoformat diff
	gofumpt -d .

fmt-apply: ## Apply autoformatter
	gofumpt -l -w .

# Adapted from https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'