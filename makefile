# Makefile for go-bn128 library
# Author: Blockchain Engineer
# Description: Build, test, and benchmark automation for BN128 elliptic curve library

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Test parameters
TEST_FLAGS=-v -race -coverprofile=coverage.out -covermode=atomic
BENCH_FLAGS=-bench=. -benchmem -benchtime=3s
BENCH_COMPARE_FLAGS=-bench=. -benchmem -benchtime=1s -count=5

# Linter
GOLINT=golangci-lint

# Colors for output
COLOR_RESET=\033[0m
COLOR_BOLD=\033[1m
COLOR_GREEN=\033[32m
COLOR_YELLOW=\033[33m
COLOR_BLUE=\033[34m
COLOR_CYAN=\033[36m

.PHONY: all build test bench clean fmt vet lint coverage help install deps update check benchmark-all benchmark-report

## all: Run all checks (fmt, vet, lint, test)
all: check test

## help: Display this help message
help:
	@echo "$(COLOR_BOLD)go-bn128 - BN128 Elliptic Curve Library$(COLOR_RESET)"
	@echo ""
	@echo "$(COLOR_CYAN)Available targets:$(COLOR_RESET)"
	@grep -E '^## ' Makefile | sed 's/## /  $(COLOR_GREEN)/' | sed 's/:/$(COLOR_RESET):/'
	@echo ""

## deps: Download and install dependencies
deps:
	@echo "$(COLOR_BLUE)→ Installing dependencies...$(COLOR_RESET)"
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "$(COLOR_GREEN)✓ Dependencies installed$(COLOR_RESET)"

## update: Update dependencies to latest versions
update:
	@echo "$(COLOR_BLUE)→ Updating dependencies...$(COLOR_RESET)"
	$(GOGET) -u ./...
	$(GOMOD) tidy
	@echo "$(COLOR_GREEN)✓ Dependencies updated$(COLOR_RESET)"

## fmt: Format Go source files
fmt:
	@echo "$(COLOR_BLUE)→ Formatting code...$(COLOR_RESET)"
	$(GOFMT) ./...
	@echo "$(COLOR_GREEN)✓ Code formatted$(COLOR_RESET)"

## vet: Run go vet
vet:
	@echo "$(COLOR_BLUE)→ Running go vet...$(COLOR_RESET)"
	$(GOVET) ./...
	@echo "$(COLOR_GREEN)✓ Vet passed$(COLOR_RESET)"

## lint: Run golangci-lint (requires golangci-lint to be installed)
lint:
	@echo "$(COLOR_BLUE)→ Running linter...$(COLOR_RESET)"
	@if command -v $(GOLINT) > /dev/null; then \
		$(GOLINT) run ./...; \
		echo "$(COLOR_GREEN)✓ Lint passed$(COLOR_RESET)"; \
	else \
		echo "$(COLOR_YELLOW)⚠ golangci-lint not installed. Run: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin$(COLOR_RESET)"; \
	fi

## check: Run format, vet, and lint checks
check: fmt vet lint
	@echo "$(COLOR_GREEN)✓ All checks passed$(COLOR_RESET)"

## test: Run all tests with race detector and coverage
test:
	@echo "$(COLOR_BLUE)→ Running tests...$(COLOR_RESET)"
	$(GOTEST) $(TEST_FLAGS) ./...
	@echo "$(COLOR_GREEN)✓ Tests passed$(COLOR_RESET)"

## test-short: Run tests without race detector (faster)
test-short:
	@echo "$(COLOR_BLUE)→ Running tests (short mode)...$(COLOR_RESET)"
	$(GOTEST) -v -short ./...
	@echo "$(COLOR_GREEN)✓ Tests passed$(COLOR_RESET)"

## test-verbose: Run tests with verbose output
test-verbose:
	@echo "$(COLOR_BLUE)→ Running tests (verbose)...$(COLOR_RESET)"
	$(GOTEST) -v -race ./...
	@echo "$(COLOR_GREEN)✓ Tests passed$(COLOR_RESET)"

## coverage: Generate and display test coverage report
coverage: test
	@echo "$(COLOR_BLUE)→ Generating coverage report...$(COLOR_RESET)"
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "$(COLOR_GREEN)✓ Coverage report generated: coverage.html$(COLOR_RESET)"
	@echo "$(COLOR_CYAN)→ Opening coverage report in browser...$(COLOR_RESET)"
	@if [ "$$(uname)" = "Darwin" ]; then \
		open coverage.html; \
	elif [ "$$(uname)" = "Linux" ]; then \
		xdg-open coverage.html 2>/dev/null || echo "Please open coverage.html manually"; \
	else \
		echo "Please open coverage.html manually"; \
	fi

## coverage-text: Display coverage in terminal
coverage-text: test
	@echo "$(COLOR_BLUE)→ Coverage summary:$(COLOR_RESET)"
	@$(GOCMD) tool cover -func=coverage.out

## bench: Run all benchmarks
bench:
	@echo "$(COLOR_BLUE)→ Running benchmarks...$(COLOR_RESET)"
	$(GOTEST) $(BENCH_FLAGS) ./...
	@echo "$(COLOR_GREEN)✓ Benchmarks complete$(COLOR_RESET)"

## bench-fp: Benchmark field arithmetic operations
bench-fp:
	@echo "$(COLOR_BLUE)→ Benchmarking Fp operations...$(COLOR_RESET)"
	$(GOTEST) -bench=BenchmarkFp -benchmem -benchtime=3s ./...

## bench-fp2: Benchmark Fp2 operations
bench-fp2:
	@echo "$(COLOR_BLUE)→ Benchmarking Fp2 operations...$(COLOR_RESET)"
	$(GOTEST) -bench=BenchmarkFp2 -benchmem -benchtime=3s ./...

## bench-g1: Benchmark G1 operations
bench-g1:
	@echo "$(COLOR_BLUE)→ Benchmarking G1 operations...$(COLOR_RESET)"
	$(GOTEST) -bench=BenchmarkG1 -benchmem -benchtime=3s ./...

## bench-g2: Benchmark G2 operations
bench-g2:
	@echo "$(COLOR_BLUE)→ Benchmarking G2 operations...$(COLOR_RESET)"
	$(GOTEST) -bench=BenchmarkG2 -benchmem -benchtime=3s ./...

## bench-pairing: Benchmark pairing operations
bench-pairing:
	@echo "$(COLOR_BLUE)→ Benchmarking pairing operations...$(COLOR_RESET)"
	$(GOTEST) -bench=BenchmarkPairing -benchmem -benchtime=3s ./...

## bench-gt: Benchmark GT operations
bench-gt:
	@echo "$(COLOR_BLUE)→ Benchmarking GT operations...$(COLOR_RESET)"
	$(GOTEST) -bench=BenchmarkGT -benchmem -benchtime=3s ./...

## bench-apps: Benchmark application scenarios (zkSNARK, BLS)
bench-apps:
	@echo "$(COLOR_BLUE)→ Benchmarking application scenarios...$(COLOR_RESET)"
	$(GOTEST) -bench="BenchmarkZKSNARK|BenchmarkBLS|BenchmarkMulti" -benchmem -benchtime=3s ./...

## benchmark-all: Run comprehensive benchmarks with detailed output
benchmark-all:
	@echo "$(COLOR_BOLD)$(COLOR_BLUE)═══════════════════════════════════════════════════$(COLOR_RESET)"
	@echo "$(COLOR_BOLD)$(COLOR_CYAN)    BN128 Comprehensive Benchmark Suite$(COLOR_RESET)"
	@echo "$(COLOR_BOLD)$(COLOR_BLUE)═══════════════════════════════════════════════════$(COLOR_RESET)"
	@echo ""
	@$(MAKE) bench-fp
	@echo ""
	@$(MAKE) bench-fp2
	@echo ""
	@$(MAKE) bench-g1
	@echo ""
	@$(MAKE) bench-g2
	@echo ""
	@$(MAKE) bench-pairing
	@echo ""
	@$(MAKE) bench-gt
	@echo ""
	@$(MAKE) bench-apps
	@echo ""
	@echo "$(COLOR_BOLD)$(COLOR_GREEN)✓ All benchmarks complete$(COLOR_RESET)"

## benchmark-compare: Run benchmarks multiple times for comparison (saved to bench.txt)
benchmark-compare:
	@echo "$(COLOR_BLUE)→ Running benchmarks for comparison...$(COLOR_RESET)"
	$(GOTEST) $(BENCH_COMPARE_FLAGS) ./... | tee bench.txt
	@echo "$(COLOR_GREEN)✓ Benchmark results saved to bench.txt$(COLOR_RESET)"
	@echo "$(COLOR_CYAN)Tip: Use 'benchstat' to compare results$(COLOR_RESET)"

## benchmark-report: Generate benchmark report with memory profiling
benchmark-report:
	@echo "$(COLOR_BLUE)→ Generating benchmark report with profiling...$(COLOR_RESET)"
	$(GOTEST) -bench=. -benchmem -memprofile=mem.prof -cpuprofile=cpu.prof ./...
	@echo "$(COLOR_GREEN)✓ Profiles generated: mem.prof, cpu.prof$(COLOR_RESET)"
	@echo "$(COLOR_CYAN)View with: go tool pprof mem.prof$(COLOR_RESET)"

## bench-mem: Run benchmarks with memory allocation reporting
bench-mem:
	@echo "$(COLOR_BLUE)→ Running benchmarks with allocation tracking...$(COLOR_RESET)"
	$(GOTEST) -bench=Allocs -benchmem ./...

## build: Build example programs (if any exist in examples/)
build:
	@echo "$(COLOR_BLUE)→ Building library...$(COLOR_RESET)"
	$(GOBUILD) ./...
	@echo "$(COLOR_GREEN)✓ Build successful$(COLOR_RESET)"

## clean: Remove generated files and caches
clean:
	@echo "$(COLOR_BLUE)→ Cleaning generated files...$(COLOR_RESET)"
	$(GOCLEAN)
	rm -f coverage.out coverage.html bench.txt
	rm -f *.prof
	rm -f *.test
	@echo "$(COLOR_GREEN)✓ Clean complete$(COLOR_RESET)"

## install: Install the library (go install)
install:
	@echo "$(COLOR_BLUE)→ Installing library...$(COLOR_RESET)"
	$(GOCMD) install ./...
	@echo "$(COLOR_GREEN)✓ Library installed$(COLOR_RESET)"

## verify: Full verification pipeline (deps, check, test, bench)
verify: deps check test bench
	@echo ""
	@echo "$(COLOR_BOLD)$(COLOR_GREEN)═══════════════════════════════════════════════════$(COLOR_RESET)"
	@echo "$(COLOR_BOLD)$(COLOR_GREEN)    ✓ Full verification complete!$(COLOR_RESET)"
	@echo "$(COLOR_BOLD)$(COLOR_GREEN)═══════════════════════════════════════════════════$(COLOR_RESET)"

## ci: Run CI pipeline (for continuous integration)
ci: deps check test coverage-text
	@echo "$(COLOR_GREEN)✓ CI pipeline complete$(COLOR_RESET)"

## quick: Quick check (fmt, test-short)
quick: fmt test-short
	@echo "$(COLOR_GREEN)✓ Quick check complete$(COLOR_RESET)"

## profile-cpu: Generate CPU profile
profile-cpu:
	@echo "$(COLOR_BLUE)→ Generating CPU profile...$(COLOR_RESET)"
	$(GOTEST) -bench=BenchmarkPairing -cpuprofile=cpu.prof -benchtime=10s
	@echo "$(COLOR_GREEN)✓ CPU profile generated$(COLOR_RESET)"
	@echo "$(COLOR_CYAN)View with: go tool pprof cpu.prof$(COLOR_RESET)"

## profile-mem: Generate memory profile
profile-mem:
	@echo "$(COLOR_BLUE)→ Generating memory profile...$(COLOR_RESET)"
	$(GOTEST) -bench=BenchmarkPairing -memprofile=mem.prof -benchtime=10s
	@echo "$(COLOR_GREEN)✓ Memory profile generated$(COLOR_RESET)"
	@echo "$(COLOR_CYAN)View with: go tool pprof mem.prof$(COLOR_RESET)"

## stats: Show project statistics
stats:
	@echo "$(COLOR_BOLD)$(COLOR_CYAN)Project Statistics:$(COLOR_RESET)"
	@echo "$(COLOR_YELLOW)Lines of code:$(COLOR_RESET)"
	@find . -name '*.go' -not -path "./vendor/*" | xargs wc -l | tail -1
	@echo ""
	@echo "$(COLOR_YELLOW)Number of functions:$(COLOR_RESET)"
	@grep -r "^func " --include="*.go" --exclude-dir=vendor | wc -l
	@echo ""
	@echo "$(COLOR_YELLOW)Number of tests:$(COLOR_RESET)"
	@grep -r "^func Test" --include="*_test.go" | wc -l
	@echo ""
	@echo "$(COLOR_YELLOW)Number of benchmarks:$(COLOR_RESET)"
	@grep -r "^func Benchmark" --include="*_test.go" | wc -l

## watch: Watch for changes and run tests (requires entr)
watch:
	@if command -v entr > /dev/null; then \
		echo "$(COLOR_BLUE)→ Watching for changes...$(COLOR_RESET)"; \
		find . -name '*.go' | entr -c make test-short; \
	else \
		echo "$(COLOR_YELLOW)⚠ entr not installed. Install with: brew install entr (macOS) or apt-get install entr (Linux)$(COLOR_RESET)"; \
	fi

## doc: Generate and serve documentation
doc:
	@echo "$(COLOR_BLUE)→ Starting documentation server...$(COLOR_RESET)"
	@echo "$(COLOR_CYAN)Open http://localhost:6060/pkg/$(COLOR_RESET)"
	godoc -http=:6060

## mod-graph: Show dependency graph
mod-graph:
	@echo "$(COLOR_BLUE)→ Dependency graph:$(COLOR_RESET)"
	$(GOMOD) graph

## mod-why: Show why packages are needed
mod-why:
	@echo "$(COLOR_BLUE)→ Analyzing dependencies...$(COLOR_RESET)"
	$(GOMOD) why -m all

## run-examples: Run basic usage examples
run-examples:
	@if [ -f "examples/main.go" ]; then \
		echo "$(COLOR_BLUE)→ Running basic examples...$(COLOR_RESET)"; \
		cd examples && go run main.go; \
		echo "$(COLOR_GREEN)✓ Examples complete$(COLOR_RESET)"; \
	else \
		echo "$(COLOR_YELLOW)⚠ examples/main.go not found$(COLOR_RESET)"; \
	fi

## run-apps: Run real-world blockchain applications
run-apps:
	@if [ -f "examples/applications/main.go" ]; then \
		echo "$(COLOR_BLUE)→ Running blockchain applications...$(COLOR_RESET)"; \
		cd examples/applications && go run main.go; \
		echo "$(COLOR_GREEN)✓ All applications complete$(COLOR_RESET)"; \
	else \
		echo "$(COLOR_YELLOW)⚠ examples/applications/main.go not found$(COLOR_RESET)"; \
		echo "$(COLOR_CYAN)Create the file first with the 5 applications$(COLOR_RESET)"; \
	fi

## demo: Run complete demonstration (tests + benchmarks + examples + apps)
demo:
	@echo "$(COLOR_BOLD)$(COLOR_CYAN)═══════════════════════════════════════════════════$(COLOR_RESET)"
	@echo "$(COLOR_BOLD)$(COLOR_CYAN)    go-bn128 Complete Demonstration$(COLOR_RESET)"
	@echo "$(COLOR_BOLD)$(COLOR_CYAN)═══════════════════════════════════════════════════$(COLOR_RESET)"
	@echo ""
	@echo "$(COLOR_BOLD)Part 1: Library Tests$(COLOR_RESET)"
	@$(MAKE) test-short
	@echo ""
	@echo "$(COLOR_BOLD)Part 2: Performance Benchmarks$(COLOR_RESET)"
	@go test -bench=BenchmarkPairing -benchtime=1s
	@echo ""
	@echo "$(COLOR_BOLD)Part 3: Real-World Applications$(COLOR_RESET)"
	@$(MAKE) run-apps
	@echo ""
	@echo "$(COLOR_BOLD)$(COLOR_GREEN)═══════════════════════════════════════════════════$(COLOR_RESET)"
	@echo "$(COLOR_BOLD)$(COLOR_GREEN)    ✓ Demonstration Complete!$(COLOR_RESET)"
	@echo "$(COLOR_BOLD)$(COLOR_GREEN)═══════════════════════════════════════════════════$(COLOR_RESET)"

.DEFAULT_GOAL := help