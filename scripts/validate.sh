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

echo "[2/4] Running golangci-lint..."
if command -v golangci-lint &> /dev/null; then
    LINT_VERSION=$(golangci-lint --version)
    if echo "$LINT_VERSION" | grep -q "2\.8"; then
        echo "Detected golangci-lint v2 (CI), using .golangci-v2.yml"
        golangci-lint run -c .golangci-v2.yml ./...
    else
        golangci-lint run ./...
    fi
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
