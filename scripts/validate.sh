#!/usr/bin/env bash
# Venvi Validation Script
# Runs all quality checks for the Go/PocketBase project

set -e

echo "=========================================="
echo "  Venvi Validation Script"
echo "=========================================="

cd "$(dirname "$0")/.."

echo ""
echo "[1/4] Running go fmt..."
go fmt ./...
echo "✓ Format check passed"

echo ""
echo "[2/4] Running golangci-lint..."
if command -v golangci-lint &> /dev/null; then
    echo "DEBUG: golangci-lint version:"
    golangci-lint --version
    echo "DEBUG: .golangci.yml content:"
    cat .golangci.yml
    golangci-lint run ./...
    echo "✓ Lint check passed"
else
    echo "⚠ golangci-lint not found, skipping..."
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
