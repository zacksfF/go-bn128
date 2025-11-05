# Usage Guide

## Building from Source

```bash
# Clone the repository
git clone https://github.com/zacksfF/go-bn128.git
cd go-bn128

# Initialize module
go mod download
go mod tidy
```

## Quick Start

```bash
# Run quick validation
make quick
```

## Testing

```bash
# Run all tests with coverage
make test
```

Test coverage: **91%**

```bash
# View coverage report
make coverage
```

## Benchmarking

```bash
# Run all benchmarks
make bench
```

Performance on modern hardware:
- G1 scalar multiplication: ~3ms
- G2 scalar multiplication: ~6ms
- Pairing operation: ~15ms

```bash
# Specific benchmarks
make bench-g1        # G1 operations only
make bench-pairing   # Pairing operations only
```

## Running Examples

```bash
# Run 5 real-world applications
make run-apps
```

Applications included:
1. zkSNARK Proof Verification
2. BLS Multi-Signature
3. Identity-Based Encryption
4. Verifiable Random Function
5. Anonymous Voting

## Code Quality

```bash
# Format and check code
make check
```

## Complete Verification

```bash
# Run everything: tests, benchmarks, quality checks
make verify
```

Time: 5-10 minutes

## Profiling

```bash
# CPU profiling
make profile-cpu

# Memory profiling
make profile-mem
```

## Cleanup

```bash
# Remove generated files
make clean
```

## Help

```bash
# Show all available commands
make help
```

## Prerequisites

- Go 1.21+
- Make
