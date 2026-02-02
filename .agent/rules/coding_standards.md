---
description: Enforce coding standards, type safety, and build integrity.
---
# Coding Standards & Integrity

## 1. Type Safety ("If it's not typed, it doesn't exist")
- All functions and methods must have type hints for arguments and return values.
- `Any` is strictly forbidden unless absolutely necessary and accompanied by a comment explaining why.
- Mypy **strict mode** must pass without errors.

## 2. Green Build Policy ("Green Build or Bust")
- You may not consider a task complete until the build passes locally.
- A "Green Build" means:
    1. `poetry run ruff format .` (Formatting)
    2. `poetry run ruff check .` (Linting)
    3. `poetry run mypy .` (Type Checking)
    4. `poetry run pytest` (Testing)

## 3. Documentation
- Public modules, classes, and functions must have docstrings.
- Docstrings should follow the Google Style Guide.
- Updates to code logic must include updates to relevant documentation (README, guides).

## 4. Testing
- All new logic must be satisfied by a unit test.
- Do not mock database models if you can use a test database fixture (preferred for SQLModel).
