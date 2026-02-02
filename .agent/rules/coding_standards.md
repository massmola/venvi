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
- A "Green Build" means all checks in the `pre-push` hook pass:
    1. `poetry run ruff format .` (Formatting)
    2. `poetry run ruff check . --fix` (Linting)
    3. `poetry run mypy .` (Type Checking - Strict Mode)
    4. `poetry run pytest` (Testing - Must maintain 100% coverage)

## 3. Documentation
- **Docstrings**: Public modules, classes, and functions must have Google Style docstrings.
- **Autogen**: This project uses `mkdocstrings`. Any new module must be added to `docs/api_reference.md` using the `::: module.name` directive.
- **Guides**: Updates to code logic must include updates to relevant documentation (`README.md`, `docs/`, `walkthrough.md`).

## 4. Testing
- **100% Coverage**: Every line of code must be covered by tests. Verify with `pytest-cov`.
- **Exhaustive Cases**: Include negative tests (error paths), filter edge cases, and core utility verification.
- **No Mocks for DB**: Do not mock database models if you can use a test database fixture (preferred for SQLModel).
## 5. Asynchronous Sessions (SQLModel/SQLAlchemy)
- When using `AsyncSession`, always use `await session.execute(query)` followed by `result.scalars().all()` (or `.first()`, etc.).
- The `.exec()` method is for synchronous sessions and must not be used with `AsyncSession`.

## 6. Web & HTMX Integrity
- **HTMX Triggers**: Dynamic UI components using HTMX must have explicit triggers (`hx-trigger`) and swap strategies (`hx-swap`).
- **Partial Templates**: Every HTMX load target must have a corresponding partial template in `src/venvi/templates/partials/`.
- **Integration Testing**: All web routes and partials must be covered by integration tests in `src/tests/test_web.py` to catch `TemplateNotFound` or rendering errors.
