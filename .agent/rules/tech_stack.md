---
trigger: always_on
description: Enforce the use of the defined technology stack.
---

# Technology Stack Rules

All code changes must adhere to the following technology choices. Do not introduce new libraries or frameworks without explicit user approval.

## Backend
- **Framework**: FastAPI
- **ORM**: SQLModel (SQLAlchemy + Pydantic)
- **Database**: PostgreSQL (via SQLModel)
- **Asynchronous**: All I/O bound operations (DB, HTTP) MUST be `async`.

## Development Tools
- **Dependency Manager**: Poetry
- **Linter/Formatter**: Ruff
- **Type Checker**: Mypy (Strict Mode)
- **Testing**: Pytest + Pytest-Cov

## Documentation
- **Format**: Markdown
- **Engine**: MkDocs (Material Theme)