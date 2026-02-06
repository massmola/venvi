---
description: Safely add or remove dependencies using Go Modules.
---
# Manage Dependencies Workflow

// turbo-all

## Adding a Dependency

1. **Enter Nix Environment**:
   ```bash
   nix develop
   ```

2. **Add the dependency**:
   ```bash
   go get github.com/example/package@latest
   ```

3. **Tidy modules**:
   ```bash
   go mod tidy
   ```

4. **Verify build**:
   ```bash
   go build ./...
   ```

## Removing a Dependency

1. **Remove imports** from all Go files that use the package.

2. **Tidy modules** to remove unused dependencies:
   ```bash
   go mod tidy
   ```

3. **Verify build**:
   ```bash
   go build ./...
   ```

## Updating Dependencies

1. **Update a specific package**:
   ```bash
   go get github.com/example/package@latest
   ```

2. **Update all dependencies**:
   ```bash
   go get -u ./...
   go mod tidy
   ```
