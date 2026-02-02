# Project Guidelines and Standards

## Core Mentality

1.  **"If it's not typed, it doesn't exist."**
    -   Strict static type checking is mandatory.
    -   `Any` is forbidden unless absolutely necessary and documented.
    -   We use **Mypy** in strict mode.

2.  **"Documentation is code."**
    -   Documentation must be updated alongside code.
    -   We use **MkDocs** with Material theme.
    -   API documentation is auto-generated via FastAPI.

3.  **"Green Build or Bust."**
    -   The CI pipeline is the ultimate source of truth.
    -   If the build fails (linting, types, tests), the code cannot be merged.

4.  **"Keep it Simple."**
    -   Avoid over-engineering.
    -   Use standard library and established patterns where possible.

## Technology Stack

-   **Language**: Python 3.12+
-   **Package Manager**: Poetry
-   **Web Framework**: FastAPI
-   **ORM**: SQLModel (SQLAlchemy + Pydantic)
-   **Database**: PostgreSQL
-   **Testing**: Pytest + Pytest-Cov
-   **Linting/Formatting**: Ruff
-   **Type Checking**: Mypy

## Development Workflow

1.  **Dependency Management**:
    -   Add dependencies: `poetry add <package>`
    -   Add dev dependencies: `poetry add -D <package>`

2.  **Code Style**:
    -   Run formatters: `poetry run ruff format .`
    -   Run linters: `poetry run ruff check . --fix`

3.  **Type Checking**:
    -   Run Mypy: `poetry run mypy .`

4.  **Testing**:
    -   Run tests: `poetry run pytest`

## CI/CD Requirements

All Pull Requests must pass:
-   Ruff formatting and linting.
-   Mypy strict type checking.
-   Pytest suite passing with acceptable coverage.
