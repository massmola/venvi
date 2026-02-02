---
description: Verify that code changes meet the project's quality standards.
---
# Verify Changes Workflow

1. **Format Code**:
   ```bash
   poetry run ruff format .
   ```

2. **Lint Code**:
   ```bash
   poetry run ruff check . --fix
   ```

3. **Type Check (Strict)**:
   ```bash
   poetry run mypy .
   ```

4. **Run Tests**:
   ```bash
   poetry run pytest
   ```

// turbo-all
5. **Report**:
   If any command fails, fix the issue and restart the workflow. If all pass, the code is ready for review.
