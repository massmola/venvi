---
trigger: always_on
description: Enforce Go coding standards for the project.
---

# Go Coding Standards

All Go code in this project must follow these standards.

## Formatting & Style
- **Formatter**: All code MUST be formatted with `gofmt` or `goimports`.
- **Naming**: Follow Go conventions (MixedCaps for exported, mixedCaps for unexported).
- **Imports**: Group imports in standard order: stdlib, external, internal.

## Documentation
- **Package Comments**: Every package must have a doc comment starting with "Package <name>".
- **Exported Identifiers**: All exported functions, types, and constants must have doc comments.
- **Comments Style**: Use complete sentences starting with the identifier name.

Example:
```go
// FetchEvents retrieves raw event data from the external API.
func (p *ODHProvider) FetchEvents(ctx context.Context) ([]RawEvent, error) {
```

## Logic & correctness
- **Edge Cases**: Always check for edge cases (empty strings, nil slices, zero values) at the beginning of utility functions. explicitly return `false`, `nil`, or an error as appropriate.

## Error Handling
- **Wrap Errors**: Use `fmt.Errorf("context: %w", err)` for error wrapping.
- **No Panic**: Never use `panic` in library code; return errors instead.
- **Handle Errors**: Never ignore errors with `_`.
- **Resource Cleanup**: Always handle or explicitly ignore errors from `Close()` calls (e.g., `_ = resp.Body.Close()`) to satisfy `errcheck`.

## Testing
- **Mandatory Validation**: Always run `./scripts/validate.sh` before submitting changes. This is REQUIRED for every commit.
- **Write Tests**: Write unit tests for all new logic. Aim for high test coverage.
- **Table-Driven Tests**: Prefer table-driven tests for multiple cases.
- **Testify**: Use `github.com/stretchr/testify` for assertions.
38:     - Use `require` for checks that must pass to avoid panics (e.g., `require.NotNil(t, obj)` before accessing `obj.Field`).
39:     - Use `assert` for checks that should not stop test execution on failure.
- **Naming**: Test functions must be named `Test<FunctionName>_<Scenario>`.

## Project-Specific
- **Providers**: New data sources must implement the `EventProvider` interface.
- **Routes**: Register routes in `routes/web.go` or `routes/api.go`.
- **Templates**: Use Go `html/template` in `views/` directory.
