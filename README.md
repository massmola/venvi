# Venvi

EU Hackathon Aggregator and Suggestion Platform.

## Quick Start

This project uses **Nix** to provide a reproducible development environment.

1.  **Enter the Development Environment**:
    ```bash
    nix develop
    ```
2.  **Initialize the Database**:
    ```bash
    bash scripts/init_db.sh
    ```
3.  **Run the Application**:
    ```bash
    poetry run uvicorn venvi.main:app --reload
    ```
4.  **Run Tests**:
    ```bash
    poetry run pytest
    ```

## Documentation

Full project documentation is available in the `docs/` directory and can be served locally:

```bash
poetry run mkdocs serve
```

See [Guidelines](docs/guidelines.md) for coding standards and CI/CD requirements.
