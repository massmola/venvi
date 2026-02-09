#!/usr/bin/env bash
# Venvi Validation Script
# Runs all quality checks for the Go/PocketBase project

set -e

echo "=========================================="
echo "  Venvi Validation Script"
echo "=========================================="

cd "$(dirname "$0")/.."

# Check for required binaries
if ! command -v go &> /dev/null; then
    echo "❌ Error: 'go' is not found in your PATH."
    echo "💡 Hint: This project uses Nix. Try running: nix develop --command ./scripts/validate.sh"
    exit 1
fi

echo ""
echo "[1/4] Running go fmt..."
go fmt ./...
echo "✓ Format check passed"

echo "[2/4] Running golangci-lint..."
if command -v golangci-lint &> /dev/null; then
    golangci-lint run ./...
    echo "✓ Lint check passed"
else
    echo "⚠ Warning: 'golangci-lint' not found. Linting skipped."
    echo "💡 Hint: Use 'nix develop' to access all required tools."
fi

echo ""
echo "[3/4] Running tests..."
go test -v -cover ./...
echo "✓ Tests passed"

echo ""
echo "[4/4] Building executable..."
go build -o venvi .
echo "✓ Build successful"

echo ""
echo "=========================================="
echo "  All checks passed! ✓"
echo "=========================================="
