# Lessons Learned

This document tracks significant failures and the adjustments made to infrastructure, rules, or workflows to prevent them in the future.

## 2026-02-09: Validation and Commit Failures

### 1. Nix Environment Awareness
- **Issue**: Running `./scripts/validate.sh` directly failed because `go` was not in the base PATH.
- **Lesson**: The agent must always be aware of the project's development environment (Nix).
- **Fix**: Updated `validate.sh` to detect missing binaries and suggest `nix develop`. Updated `tech_stack.md` to mandate Nix for all commands.

### 2. Gitleaks Failure on `server.log`
- **Issue**: A commit failed because `server.log` contained sensitive data (JWTs).
- **Lesson**: Log files and temporary artifacts should never be staged. `.gitignore` was too specific.
- **Fix**: Updated `.gitignore` to use `*.log` and `run_logs_*.txt`. Explicitly removed current log files from the index.

### 3. Unhandled `Close()` Errors
- **Issue**: `golangci-lint` (errcheck) failed because `resp.Body.Close()` errors were ignored.
- **Lesson**: Library code must be explicit about error ignoring to pass strict linting.
- **Fix**: Updated `coding_standards.md` to specify the pattern for resource cleanup (`_ = Close()`).

## Meta-Reflection
- Granular verification after every significant change (the "Ralph Wiggum" approach) is essential to catch environment-specific issues early.
- Pre-push hooks are a vital last line of defense but should be supplemented by frequent manual validation.
