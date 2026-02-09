---
trigger: always_on
description: Enforce the use of the defined technology stack.
---

# Technology Stack Rules

All code changes must adhere to the following technology choices. Do not introduce new libraries or frameworks without explicit user approval.

## Backend
- **Framework**: PocketBase (Go-extended)
- **Database**: SQLite (embedded in PocketBase)
- **Concurrency**: Go goroutines and channels

## Development Tools
- **Dependency Manager**: Go Modules
- **Linter/Formatter**: golangci-lint, gofmt
- **Testing**: go test + testify
- **Development Env**: nix flake (use it and install tools trough it)

## Documentation
- **Format**: Markdown
- **In-Code**: Go doc comments (godoc format)

## Project Structure
- **Providers**: `providers/` - Event data source implementations
- **Routes**: `routes/` - HTTP routes for web and API
- **Views**: `views/` - Go html/template files
- **Static**: `pb_public/` - Static assets served by PocketBase
- **Migrations**: `pb_migrations/` - PocketBase schema migrations