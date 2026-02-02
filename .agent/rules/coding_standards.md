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
## 5. Asynchronous Sessions (SQLModel/SQLAlchemy)
- When using `AsyncSession`, always use `await session.execute(query)` followed by `result.scalars().all()` (or `.first()`, etc.).
- The `.exec()` method is for synchronous sessions and must not be used with `AsyncSession`.

## 6. Web & HTMX Integrity
- **HTMX Triggers**: Dynamic UI components using HTMX must have explicit triggers (`hx-trigger`) and swap strategies (`hx-swap`).
- **Partial Templates**: Every HTMX load target must have a corresponding partial template in `src/venvi/templates/partials/`.
- **Integration Testing**: All web routes and partials must be covered by integration tests in `src/tests/test_web.py` to catch `TemplateNotFound` or rendering errors.
