# Project Guidelines and Standards

## Core Mentality

1.  **"Handled Errors are the Truth."**
    -   Never ignore errors with `_`.
    -   Wrap errors with context: `fmt.Errorf("context: %w", err)`.
    -   Panic is strictly forbidden in library code.

2.  **"Documentation is Alive."**
    -   Documentation must be updated alongside code.
    -   We use Go doc comments for all exported identifiers.
    -   Project-wide information lives in `README.md` and `GUIDELINES.md`.

3.  **"Green Build or Bust."**
    -   The CI pipeline is the ultimate source of truth.
    -   If the build fails (formatting, linting, tests), the code cannot be merged.
    -   Run `./scripts/validate.sh` before pushing.

4.  **"Keep it Simple."**
    -   Leverage PocketBase features before writing custom logic.
    -   Follow standard Go idiomatic patterns.

## Technology Stack

-   **Language**: Go 1.24+
-   **Framework**: PocketBase (Go-extended)
-   **Database**: SQLite (embedded in PocketBase)
-   **Testing**: Go `testing` package + Testify
-   **Linting/Formatting**: `gofmt` and `golangci-lint`

## Development Workflow

1.  **Dependency Management**:
    -   We use Go Modules (`go.mod`, `go.sum`).
    -   Add dependencies: `go get <package>`
    -   Clean up: `go mod tidy`

2.  **Code Style**:
    -   Format code: `go fmt ./...`
    -   Run linter: `golangci-lint run ./...`

3.  **Testing**:
    -   Run tests: `go test -v ./...`
    -   Run with coverage: `go test -cover ./...`

4.  **Database Migrations**:
    -   Create new migration: `go run main.go migrate create <name>`
    -   Migrations are automatically applied on server start if triggered via `migratecmd`.

## CI/CD Requirements

All Pull Requests must pass the validation script:
-   `go fmt` check.
-   `golangci-lint` check.
-   All Go tests passing.
-   Successful build of the `venvi` binary.
