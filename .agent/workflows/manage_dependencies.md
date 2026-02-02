---
description: Safely add or remove dependencies using Poetry.
---
# Manage Dependencies Workflow

## Adding a Dependency

1. **Identify Package**: Determine the package name and version (if specific).

2. **Add Package**:
   ```bash
   poetry add <package_name>
   ```
   *For dev tools (like linters), use:* `poetry add -D <package_name>`

3. **Verify Installation**:
   Check `pyproject.toml` to see if it was added.

4. **Update Lockfile**:
   This happens automatically with `poetry add`.

## Removing a Dependency

1. **Remove Package**:
   ```bash
   poetry remove <package_name>
   ```

2. **Clean Environment**:
   ```bash
   poetry install --sync
   ```
